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

package discovery

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	ethCommon "github.com/unicornultrafoundation/go-u2u/libs/common"

	"github.com/oschwald/geoip2-golang"
	"github.com/unicornultrafoundation/go-u2u/libs/log"
	discoverProtocol "github.com/unicornultrafoundation/go-u2u/libs/p2p/discover"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p/enode"
	"github.com/unicornultrafoundation/network-stats/pkg/common"
	"github.com/unicornultrafoundation/network-stats/pkg/crawlerdb"
)

type Discover struct {
	// These are probably from flags
	GenesisHash string
	NetworkID   uint64
	NodeURL     string
	ListenAddr  string
	NodeKey     string
	Bootnodes   []string
	Timeout     time.Duration
	Workers     uint64

	NodeDB *enode.DB
}

type discover struct {
	output common.NodeSet

	genesisHash ethCommon.Hash
	networkID   uint64
	nodeURL     string

	disc resolver

	inputIter enode.Iterator
	iters     []enode.Iterator

	ch     chan *enode.Node
	closed chan struct{}

	// settings
	revalidateInterval time.Duration

	reqCh   chan *enode.Node
	workers uint64

	sync.WaitGroup
	sync.RWMutex
}

type resolver interface {
	RequestENR(*enode.Node) (*enode.Node, error)
	RandomNodes() enode.Iterator
}

func NewDiscover(
	genesisHash string,
	networkID uint64,
	nodeURL string,
	input common.NodeSet,
	workers uint64,
	disc resolver,
	iters ...enode.Iterator,
) *discover {
	c := &discover{
		genesisHash: ethCommon.HexToHash(genesisHash),
		output:      make(common.NodeSet, len(input)),
		networkID:   networkID,
		nodeURL:     nodeURL,
		disc:        disc,
		iters:       iters,
		inputIter:   enode.IterNodes(input.Nodes()),
		ch:          make(chan *enode.Node),
		reqCh:       make(chan *enode.Node, 512), // TODO: define this in config
		workers:     workers,
		closed:      make(chan struct{}),
	}
	c.iters = append(c.iters, c.inputIter)
	// Copy input to output initially. Any nodes that fail validation
	// will be dropped from output during the run.
	for id, n := range input {
		c.output[id] = n
	}
	return c
}

func (c *discover) Run(timeout time.Duration) common.NodeSet {
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

	for i := c.workers; i > 0; i-- {
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

func (c *discover) runIterator(done chan<- enode.Iterator, it enode.Iterator) {
	defer func() { done <- it }()
	for it.Next() {
		select {
		case c.ch <- it.Node():
		case <-c.closed:
			return
		}
	}
}

func (c *discover) getClientInfoLoop() {
	defer func() { c.Done() }()
	for n := range c.reqCh {
		if n == nil {
			return
		}

		var tooManyPeers bool
		var scoreInc int

		info, err := getClientInfo(c.genesisHash, c.networkID, c.nodeURL, n)
		if err != nil {
			log.Warn("GetClientInfo failed", "error", err, "nodeID", n.ID())
			if strings.Contains(err.Error(), "too many peers") {
				tooManyPeers = true
			}
		} else {
			scoreInc = 10
		}

		if info != nil {
			log.Info(
				"Updating node info",
				"client_type", info.ClientType,
				"version", info.SoftwareVersion,
				"network_id", info.NetworkID,
				"caps", info.Capabilities,
				"fork_id", info.ForkID,
				"height", info.Blockheight,
				"td", info.TotalDifficulty,
				"head", info.HeadHash,
			)
		}

		c.Lock()
		node := c.output[n.ID()]
		node.N = n
		node.Seq = n.Seq()
		if info != nil {
			node.Info = info
		}
		node.TooManyPeers = tooManyPeers
		node.Score += scoreInc
		c.output[n.ID()] = node
		c.Unlock()
	}
}

func (c *discover) updateNode(n *enode.Node) {
	c.Lock()
	defer c.Unlock()

	node, ok := c.output[n.ID()]

	// Skip validation of recently-seen nodes.
	if ok && !node.TooManyPeers && time.Since(node.LastCheck) < c.revalidateInterval {
		return
	}

	node.LastCheck = time.Now().UTC().Truncate(time.Second)

	// Request the node record.
	nn, err := c.disc.RequestENR(n)
	if err != nil {
		log.Info("request enr error", "err", err.Error(), "node", node)
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

func (c Discover) CrawlRound(
	inputSet common.NodeSet,
	db *sql.DB,
	geoipDB *geoip2.Reader,
) common.NodeSet {
	var v4, v5 common.NodeSet
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		v5 = c.discv5(inputSet)
		log.Info("DiscV5", "nodes", len(v5.Nodes()))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		v4 = c.discv4(inputSet)
		log.Info("DiscV4", "nodes", len(v4.Nodes()))
	}()

	wg.Wait()

	output := make(common.NodeSet, len(v4)+len(v5))
	for _, n := range v5 {
		output[n.N.ID()] = n
	}
	for _, n := range v4 {
		output[n.N.ID()] = n
	}

	var nodes []common.NodeJSON
	for _, node := range output {
		nodes = append(nodes, node)
	}

	// Write the node info to database
	if db != nil {
		if err := crawlerdb.UpdateNodes(db, geoipDB, nodes); err != nil {
			panic(err)
		}
	}
	return output
}

func (c Discover) discv5(inputSet common.NodeSet) common.NodeSet {
	ln, config := c.makeDiscoveryConfig()
	log.Info("disc v5", "config", config, "addr", c.ListenAddr)

	socket := listen(ln, c.ListenAddr)

	disc, err := discoverProtocol.ListenV5(socket, ln, config)
	if err != nil {
		panic(err)
	}
	defer disc.Close()

	return c.runDiscover(disc, inputSet)
}

func (c Discover) discv4(inputSet common.NodeSet) common.NodeSet {
	ln, config := c.makeDiscoveryConfig()

	log.Info("disc v4", "config", config, "addr", c.ListenAddr)

	socket := listen(ln, c.ListenAddr)

	disc, err := discoverProtocol.ListenV4(socket, ln, config)
	if err != nil {
		panic(err)
	}
	defer disc.Close()

	return c.runDiscover(disc, inputSet)
}

func (c Discover) runDiscover(disc resolver, inputSet common.NodeSet) common.NodeSet {
	log.Info("New discover with node url", "url", c)
	crawler := NewDiscover(c.GenesisHash, c.NetworkID, c.NodeURL, inputSet, c.Workers, disc, disc.RandomNodes())
	crawler.revalidateInterval = 1 * time.Minute
	return crawler.Run(c.Timeout)
}
