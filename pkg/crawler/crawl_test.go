package crawler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unicornultrafoundation/go-u2u/libs/common"
	"github.com/unicornultrafoundation/go-u2u/libs/crypto"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p/discover"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p/enode"
	"testing"
	"time"
)

func TestDialNode(t *testing.T) {
	p2pNode, _ := parseNode("enode://a06d1b232c30abc73200373bebf6f28a8c4133edb133e0c8d520c04a1576ca4daaa41de465ac3635cc10c6e1c5fa21c3b89bffb03e6eeee3d4d6f30504a1a4e7@171.242.12.202:0?discport=30301")
	genesisHash := common.HexToHash("0xb89e32765f34b70eaa3e0760d112690431f182112f345e77f065cfca83e4524a")
	info, err := getClientInfo(genesisHash, 4339, "http://localhost:8545/", p2pNode)
	assert.Nil(t, err, "err: %s", err.Error())
	fmt.Printf("Info: %+v \n", info)
}

func TestDiscovery_DiscV4(t *testing.T) {
	n, _ := parseNode("enode://ccc5ae45b01c3cbdcdd15275af187144416718a208d66d41d2ecf5cd874b13f834a022dd589163bb759f7a8907ee7987ba5a29f003347f57178d720a2998575b@13.215.34.27:5050")

	disc, _ := startV4()
	defer disc.Close()
	start := time.Now()
	fmt.Println(disc.Resolve(n).String())
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//fmt.Println(info)
	fmt.Printf("node responded to ping (RTT %v).\n", time.Since(start))

}

func TestDiscovery_DiscV5(t *testing.T) {
	n, _ := parseNode("enode://ccc5ae45b01c3cbdcdd15275af187144416718a208d66d41d2ecf5cd874b13f834a022dd589163bb759f7a8907ee7987ba5a29f003347f57178d720a2998575b@13.215.34.27:5050")

	disc, _ := startV5()
	defer disc.Close()
	start := time.Now()
	nodeInfo, err := disc.RequestENR(n)
	if err != nil {
		fmt.Println("ping err", err.Error())
	}
	fmt.Println("node info", nodeInfo)

	fmt.Printf("node responded to ping (RTT %v).\n", time.Since(start))

}
func startV5() (*discover.UDPv5, discover.Config) {
	ln, config := makeDiscoveryConfig()
	socket := listen(ln, "0.0.0.0:0")
	disc, err := discover.ListenV5(socket, ln, config)
	if err != nil {
		fmt.Println("cannot start v5", err.Error())
	}
	return disc, config
}

func startV4() (*discover.UDPv4, discover.Config) {
	ln, config := makeDiscoveryConfig()
	socket := listen(ln, "0.0.0.0:0")
	disc, err := discover.ListenV4(socket, ln, config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return disc, config
}

func makeDiscoveryConfig() (*enode.LocalNode, discover.Config) {
	var cfg discover.Config
	key, err := crypto.HexToECDSA("507b50d7f75308f5bef0dc49e776cec963827bad8a0df73accac855192f68b8b")
	cfg.PrivateKey = key

	bn, err := parseBootnodes()
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg.Bootnodes = bn

	db, err := enode.OpenDB("./test-db.db")
	if err != nil {
		fmt.Println("cannot open db", err.Error())
	}
	ln := enode.NewLocalNode(db, cfg.PrivateKey)
	return ln, cfg
}

func parseBootnodes() ([]*enode.Node, error) {
	s := U2UMainnetBootnodes

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
