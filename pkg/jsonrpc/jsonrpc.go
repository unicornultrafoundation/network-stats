package jsonrpc

import (
	"strings"

	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc/rpc"
	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc/u2u"
	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc/utils"

	"github.com/unicornultrafoundation/go-u2u/libs/ethclient"
)

type Web3 struct {
	U2U   *u2u.U2U
	Utils *utils.Utils
	c     *rpc.Client
	be    *ethclient.Client
}

func NewWeb3(provider string) (*Web3, error) {

	return NewWeb3WithProxy(provider, "")
}

func NewWeb3WithProxy(provider, proxy string) (*Web3, error) {
	c, err := rpc.NewClient(provider, proxy)
	if err != nil {
		return nil, err
	}
	be, err := ethclient.Dial(provider)
	if err != nil {
		panic(err)
	}

	e := u2u.NewU2U(c)

	providerLowerStr := strings.ToLower(provider)

	if strings.Contains(providerLowerStr, "ropsten") {
		e.SetChainId(3)
	} else if strings.Contains(providerLowerStr, "kovan") {
		e.SetChainId(42)
	} else if strings.Contains(providerLowerStr, "rinkeby") {
		e.SetChainId(4)
	} else if strings.Contains(providerLowerStr, "goerli") {
		e.SetChainId(5)
	} else {
		e.SetChainId(1)
	}

	u := utils.NewUtils()
	w := &Web3{
		U2U:   e,
		Utils: u,
		c:     c,
		be:    be,
	}

	// Default poll timeout 2 hours
	w.U2U.SetTxPollTimeout(7200)
	return w, nil
}

func (w *Web3) GetContractConn() *ethclient.Client {
	return w.be
}

func (w *Web3) Version() (string, error) {
	var out string
	err := w.c.Call("web3_clientVersion", &out)
	return out, err
}
