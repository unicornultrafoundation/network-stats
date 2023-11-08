package common

import (
	"math/big"

	"github.com/unicornultrafoundation/go-u2u/libs/common"
	"github.com/unicornultrafoundation/go-u2u/libs/core/forkid"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p"
)

type ClientInfo struct {
	ClientType      string
	SoftwareVersion uint64
	Capabilities    []p2p.Cap
	NetworkID       uint64
	ForkID          forkid.ID
	Blockheight     string
	TotalDifficulty *big.Int
	HeadHash        common.Hash
}
