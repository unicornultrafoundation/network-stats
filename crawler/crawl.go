// Copyright 2019 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type crawler struct {
	output nodeSet

	genesis   *core.Genesis
	networkID uint64
	nodeURL   string

	disc resolver

	inputIter enode.Iterator
	iters     []enode.Iterator

	ch     chan *enode.Node
	closed chan struct{}

	// settings
	revalidateInterval time.Duration

	reqCh   chan *enode.Node
	workers int

	// errors
	tooManyPeersError uint64

	sync.WaitGroup
	sync.RWMutex
}

type resolver interface {
	RequestENR(*enode.Node) (*enode.Node, error)
	RandomNodes() enode.Iterator
}

func newCrawler(genesis *core.Genesis, networkID uint64, nodeURL string, input nodeSet, disc resolver, iters ...enode.Iterator) *crawler {
	c := &crawler{
		output:    make(nodeSet, len(input)),
		genesis:   genesis,
		networkID: networkID,
		nodeURL:   nodeURL,
		disc:      disc,
		iters:     iters,
		inputIter: enode.IterNodes(input.nodes()),
		ch:        make(chan *enode.Node),
		reqCh:     make(chan *enode.Node, 1024), // TODO: define this in config
		workers:   32,                           // TODO: define this in config
		closed:    make(chan struct{}),
	}
	c.iters = append(c.iters, c.inputIter)
	// Copy input to output initially. Any nodes that fail validation
	// will be dropped from output during the run.
	for id, n := range input {
		c.output[id] = n
	}
	return c
}

func (c *crawler) run(timeout time.Duration) nodeSet {
	var (
		timeoutTimer = time.NewTimer(timeout)
		timeoutCh    <-chan time.Time
		doneCh       = make(chan enode.Iterator, len(c.iters))
		liveIters    = len(c.iters)
		inputSetLen  = len(c.output)
	)
	defer timeoutTimer.Stop()

	for _, it := range c.iters {
		go c.runIterator(doneCh, it)
	}

	for i := 0; i < c.workers; i++ {
		c.Add(1)
		go c.getClientInfoLoop()
	}

loop:
	for {
		select {
		case n := <-c.ch:
			c.updateNode(n)
		case it := <-doneCh:
			if it == c.inputIter {
				// Enable timeout when we're done revalidating the input nodes.
				log.Info("Revalidation of input set is done", "len", inputSetLen)
				if timeout > 0 {
					timeoutCh = timeoutTimer.C
				}
			}
			if liveIters--; liveIters <= 0 {
				break loop
			}
		case <-timeoutCh:
			break loop
		}
	}

	close(c.closed)
	close(c.reqCh)
	for _, it := range c.iters {
		it.Close()
	}
	for ; liveIters > 0; liveIters-- {
		<-doneCh
	}
	c.Wait()

	close(c.ch)

	return c.output
}

func (c *crawler) runIterator(done chan<- enode.Iterator, it enode.Iterator) {
	defer func() { done <- it }()
	for it.Next() {
		select {
		case c.ch <- it.Node():
		case <-c.closed:
			return
		}
	}
}

func (c *crawler) getClientInfoLoop() {
	defer func() { c.Done() }()
	for {
		select {
		case n, ok := <-c.reqCh:
			if !ok {
				return
			}

			errorReason := 0
			errorString := ""
			var scoreInc int

			info, err := getClientInfo(c.genesis, c.networkID, c.nodeURL, n)
			if err != nil {
				errStrings := strings.Split(err.Error(), ":")
				if len(errStrings) >= 2 {
					lastError := errStrings[len(errStrings)-1]
					if strings.Contains(lastError, "decoding into (main.Status).NetworkID") {
						lastError = " error decoding NetworkID"
					}
					clientType := fmt.Sprintf(errStrings[0] + lastError)
					clientType = strings.Replace(clientType, " ", "_", 10)
					info.ClientType = clientType
				}
				errorReason = -1
				errorString = err.Error()
				log.Warn("GetClientInfo failed", "error", err, "nodeID", n.ID())
			} else {
				scoreInc = 10
			}

			if info != nil {
				log.Info(
					"Updating node info",
					"client_type", info.ClientType,
					"clientVersion", info.ClientVersion,
					"clientDescription", info.ClientDesc,
					"osType", info.OsType,
					"goVersion", info.GoVersion,
					"version", info.SoftwareVersion,
					"network_id", info.NetworkID,
					"caps", info.Capabilities,
					"fork_id", info.ForkID,
					"height", info.Blockheight,
					"td", info.TotalDifficulty,
					"head", info.HeadHash,
					"url", n.URLv4(),
				)
			}

			c.Lock()
			node := c.output[n.ID()]
			node.N = n
			node.Seq = n.Seq()
			if info != nil {
				node.Info = info
			}
			node.ErrorReason = errorReason
			node.ErrorString = errorString
			node.Score += scoreInc
			c.output[n.ID()] = node
			c.Unlock()
		}
	}
}

func (c *crawler) updateNode(n *enode.Node) {
	c.Lock()
	defer c.Unlock()

	node, ok := c.output[n.ID()]

	// Skip validation of recently-seen nodes.
	if ok && time.Since(node.LastCheck) < c.revalidateInterval {
		return
	}

	node.LastCheck = time.Now().UTC().Truncate(time.Second)

	// Request the node record.
	nn, err := c.disc.RequestENR(n)
	if err != nil {
		if node.Score == 0 {
			// Node doesn't implement EIP-868.
			log.Debug("Skipping node", "id", n.ID())
			return
		}
		node.Score /= 2
	} else {
		node.N = nn
		node.Seq = nn.Seq()
		node.Score++
		if node.FirstResponse.IsZero() {
			node.FirstResponse = node.LastCheck
		}
		node.LastResponse = node.LastCheck
	}

	// Store/update node in output set.
	if node.Score <= 0 {
		log.Info("Removing node", "id", n.ID())
		delete(c.output, n.ID())
	} else {
		log.Info("Updating node", "id", n.ID(), "seq", n.Seq(), "score", node.Score)
		c.reqCh <- n
		c.output[n.ID()] = node
	}
}
