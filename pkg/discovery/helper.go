package discovery

import (
	"fmt"
	"net"

	"github.com/unicornultrafoundation/go-u2u/libs/crypto"
	u2uProtocol "github.com/unicornultrafoundation/go-u2u/libs/p2p/discover"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p/enode"
)

func (c Discover) makeDiscoveryConfig() (*enode.LocalNode, u2uProtocol.Config) {
	var cfg u2uProtocol.Config
	var err error

	if c.NodeKey != "" {
		key, err := crypto.HexToECDSA(c.NodeKey)
		if err != nil {
			panic(fmt.Errorf("-%s: %v", c.NodeKey, err))
		}
		cfg.PrivateKey = key
	} else {
		cfg.PrivateKey, _ = crypto.GenerateKey()
	}

	cfg.Bootnodes, err = c.parseBootnodes()
	if err != nil {
		panic(err)
	}

	return enode.NewLocalNode(c.NodeDB, cfg.PrivateKey), cfg
}

func listen(ln *enode.LocalNode, addr string) *net.UDPConn {
	socket, err := net.ListenPacket("udp4", addr)
	if err != nil {
		panic(err)
	}
	usocket := socket.(*net.UDPConn)
	uaddr := socket.LocalAddr().(*net.UDPAddr)
	if uaddr.IP.IsUnspecified() {
		ln.SetFallbackIP(net.IP{127, 0, 0, 1})
	} else {
		ln.SetFallbackIP(uaddr.IP)
	}
	ln.SetFallbackUDP(uaddr.Port)
	return usocket
}

func (c Discover) parseBootnodes() ([]*enode.Node, error) {
	bootnodes := CrawlerMainnetBootnodes
	if len(c.Bootnodes) != 0 {
		bootnodes = c.Bootnodes
	}

	nodes := make([]*enode.Node, len(bootnodes))
	var err error
	for i, record := range bootnodes {
		nodes[i], err = parseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %v", err)
		}
	}
	return nodes, nil
}
