package u2u

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc/rpc"
)

func TestGetBlockByNumber(t *testing.T) {
	c, err := rpc.NewClient("https://rpc-mainnet.uniultra.xyz/", "")
	if err != nil {
		t.Fatal(err)
	}
	u2uNode := NewU2U(c)
	blockNumber, err := u2uNode.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	block, err := u2uNode.GetBlocByNumber(big.NewInt(int64(blockNumber)), true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("block hash %s has %v txs\n", block.Hash(), len(block.Transactions()))
}

func TestPollBlock(t *testing.T) {
	c, err := rpc.NewClient("https://rpc-mainnet.uniultra.xyz/", "")
	if err != nil {
		t.Fatal(err)
	}
	u2uNode := NewU2U(c)
	for {
		blockNumber, err := u2uNode.GetBlockNumber()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("get block %v\n", blockNumber)
		block, err := u2uNode.GetBlocByNumber(big.NewInt(int64(blockNumber)), true)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("block hash %s has %v txs\n", block.Hash(), len(block.Transactions()))
		time.Sleep(time.Duration(5) * time.Second)
	}

}
