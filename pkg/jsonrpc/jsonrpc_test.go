package jsonrpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unicornultrafoundation/go-u2u/libs/common"
	"github.com/unicornultrafoundation/network-stats/contracts"
	"math/big"
	"testing"
)

func TestDumpContract(t *testing.T) {

	web3, err := NewWeb3("https://rpc-nebulas-testnet.uniultra.xyz/")
	assert.Nil(t, err)
	sfcContract, err := contracts.NewSFC(common.HexToAddress("0xfc00face00000000000000000000000000000000"), web3.GetContractConn())

	currentEpoch, err := sfcContract.CurrentEpoch(nil)
	assert.Nil(t, err)
	fmt.Println("current epoch", currentEpoch)

	snapshot, err := sfcContract.GetEpochSnapshot(nil, new(big.Int).SetUint64(50310))
	assert.Nil(t, err)
	fmt.Printf("Snapshot: %+v \n", snapshot)
}
