package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"gopkg.in/urfave/cli.v1"
)

var MainnetBootnodes = []string{
	"enode://21dfee41ddd127ebbd68fb14b39945f6e993ad9eb35c57e5e2e17ec1740960400d6d174f6c119fb9940072eec2d468ee5d767752bf9a44900ac8ac6d6de61330@18.143.208.170:5050",
	"enode://a1e1999ab32c7ea71b3fb4fd4e2143beadc3f71365e2a5a0e54e15780d28e5a80576a387406d9b60eee7c31289618c6a5ef93bfe295215518cecbf23bc50211e@3.1.11.147:5050",
}

func makeDiscoveryConfig(ctx *cli.Context, db *enode.DB) (*enode.LocalNode, discover.Config) {
	var cfg discover.Config
	var err error

	if ctx.IsSet(nodekeyFlag.Name) {
		key, err := crypto.HexToECDSA(ctx.String(nodekeyFlag.Name))
		if err != nil {
			panic(fmt.Errorf("-%s: %v", nodekeyFlag.Name, err))
		}
		cfg.PrivateKey = key
	} else {
		cfg.PrivateKey, _ = crypto.GenerateKey()
	}

	cfg.Bootnodes, err = parseBootnodes(ctx)
	if err != nil {
		panic(err)
	}

	return enode.NewLocalNode(db, cfg.PrivateKey), cfg
}

func listen(ln *enode.LocalNode, addr string) *net.UDPConn {
	if addr == "" {
		addr = "0.0.0.0:0"
	}
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

func parseBootnodes(ctx *cli.Context) ([]*enode.Node, error) {
	s := MainnetBootnodes
	if ctx.IsSet(bootnodesFlag.Name) {
		input := ctx.String(bootnodesFlag.Name)
		if input == "" {
			return nil, nil
		}
		s = strings.Split(input, ",")
	}
	nodes := make([]*enode.Node, len(s))
	var err error
	for i, record := range s {
		nodes[i], err = parseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %v", err)
		}
	}
	return nodes, nil
}
