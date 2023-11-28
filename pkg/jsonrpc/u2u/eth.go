package u2u

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"math/big"
	"sort"

	"github.com/unicornultrafoundation/go-u2u/libs/accounts"
	"github.com/unicornultrafoundation/go-u2u/libs/common"
	"github.com/unicornultrafoundation/go-u2u/libs/common/hexutil"
	eTypes "github.com/unicornultrafoundation/go-u2u/libs/core/types"
	"github.com/unicornultrafoundation/go-u2u/libs/crypto"

	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc/rpc"
	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc/types"
	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc/utils"
)

const (
	PRIORITY_FEE_INCREASE_BOUNDARY = 200
)

// U2U is the eth namespace
type U2U struct {
	c             *rpc.Client
	privateKey    *ecdsa.PrivateKey
	address       common.Address
	chainId       *big.Int
	txPollTimeout int
	utils         *utils.Utils
}

// NewU2U Create a u2u instance
func NewU2U(c *rpc.Client) *U2U {

	return &U2U{
		c:     c,
		utils: &utils.Utils{},
	}
}

// SetAccount Setup default ethereum account with privateKey (hex format)
func (e *U2U) SetAccount(privateKey string) error {
	if len(privateKey) == 0 {
		return fmt.Errorf("private key is empty")
	}
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	e.privateKey = privKey

	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	copy(e.address[:], addr[:])

	return nil
}

func (e *U2U) GetPrivateKey() *ecdsa.PrivateKey {
	return e.privateKey
}

func (e *U2U) GetChainId() *big.Int {
	return e.chainId
}

// SetChainId Setup current network chainId
func (e *U2U) SetChainId(chainId int64) {
	e.chainId = big.NewInt(chainId)
}

// SetTxPollTimeout  Setup timeout for polling confirmation from txs (unit second)
func (e *U2U) SetTxPollTimeout(timeout int) {
	if timeout == 0 {
		// default tx poll timeout is 720s
		e.txPollTimeout = 720
		return
	}
	e.txPollTimeout = timeout
}

// Accounts Get accounts from rpc providers
func (e *U2U) Accounts() ([]common.Address, error) {
	var out []common.Address
	if err := e.c.Call("eth_accounts", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddressGet current default account address
func (e *U2U) Address() common.Address {
	return e.address
}

// GetBlockNumber Get current block height
func (e *U2U) GetBlockNumber() (uint64, error) {
	var out string
	if err := e.c.Call("eth_blockNumber", &out); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

// GetBlockHeaderByNumber Get block header by block number
func (e *U2U) GetBlockHeaderByNumber(number *big.Int, full bool) (*eTypes.Header, error) {
	var head *eTypes.Header
	if err := e.c.Call("eth_getBlockByNumber", &head, utils.ToBlockNumArg(number), full); err != nil {
		return nil, err
	}
	return head, nil
}

// GetBlockByNumber Get block by block number
func (e *U2U) GetBlockByNumber(number *big.Int, full bool) (*eTypes.Block, error) {
	return e.getBlock("eth_getBlockByNumber", utils.ToBlockNumArg(number), full)
}

// GetBlockByHash Get block by block hash
func (e *U2U) GetBlockByHash(hash common.Hash, full bool) (*eTypes.Block, error) {
	var b *eTypes.Block
	if err := e.c.Call("eth_getBlockByHash", &b, hash, full); err != nil {
		return nil, err
	}
	return b, nil
}

// SendTransaction Send transaction
func (e *U2U) SendTransaction(txn *eTypes.Transaction) (common.Hash, error) {
	var hash common.Hash
	err := e.c.Call("eth_sendTransaction", &hash, txn)
	return hash, err
}

// GetTransactionByHash Get transaction by transaction hash
func (e *U2U) GetTransactionByHash(hash common.Hash) (*eTypes.Transaction, error) {
	var tx *eTypes.Transaction
	err := e.c.Call("eth_getTransactionByHash", &tx, hash)
	return tx, err
}

// GetTransactionReceipt Get transaction receipt by transaction hash
func (e *U2U) GetTransactionReceipt(hash common.Hash) (*eTypes.Receipt, error) {
	var receipt *eTypes.Receipt
	err := e.c.Call("eth_getTransactionReceipt", &receipt, hash)
	return receipt, err
}

// GetNonce  Get nonce of account
func (e *U2U) GetNonce(addr common.Address, blockNumber *big.Int) (uint64, error) {
	var nonce string
	if err := e.c.Call("eth_getTransactionCount", &nonce, addr, utils.ToBlockNumArg(blockNumber)); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(nonce)
}

// Get ether balance of account
func (e *U2U) GetBalance(addr common.Address, blockNumber *big.Int) (*big.Int, error) {
	var out string
	if err := e.c.Call("eth_getBalance", &out, addr, utils.ToBlockNumArg(blockNumber)); err != nil {
		return nil, err
	}
	b, ok := new(big.Int).SetString(out[2:], 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert to big.int")
	}
	return b, nil
}

// Get gas price for Non-EIP1559 tx
func (e *U2U) GasPrice() (uint64, error) {
	var out string
	if err := e.c.Call("eth_gasPrice", &out); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

// Get fee history for EIP1559 blocks
func (e *U2U) FeeHistory(historicalBlocks int, blockNumber *big.Int, feeHistoryPercentile []float64) (*types.FeeHistory, error) {
	var out *types.FeeHistory
	if err := e.c.Call("eth_feeHistory", &out, historicalBlocks, utils.ToBlockNumArg(blockNumber), feeHistoryPercentile); err != nil {
		return nil, err
	}
	return out, nil
}

// Call Do Call functions
func (e *U2U) Call(msg *types.CallMsg, block *big.Int) (string, error) {
	var out string
	if err := e.c.Call("eth_call", &out, msg, utils.ToBlockNumArg(block)); err != nil {
		return "", err
	}
	return out, nil
}

// EstimateGasContract Estimate gas for deploying contract
func (e *U2U) EstimateGasContract(bin []byte) (uint64, error) {
	var out string
	msg := map[string]interface{}{
		"data": "0x" + hex.EncodeToString(bin),
	}
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

// EstimateGas Estimate gas for excuting transaction
func (e *U2U) EstimateGas(msg *types.CallMsg) (uint64, error) {
	var out string
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

// ChainID Get currnet network chainId from provider
func (e *U2U) ChainID() (*big.Int, error) {
	if e.chainId != nil {
		return e.chainId, nil
	}
	var out string
	if err := e.c.Call("eth_chainId", &out); err != nil {
		return nil, err
	}
	return utils.ParseBigInt(out), nil
}

// GetLogs Get past event logs with filter
func (e *U2U) GetLogs(filter *types.Filter) ([]*types.Event, error) {
	out := make([]*types.Event, 0)
	if err := e.c.Call("eth_getLogs", &out, filter); err != nil {
		return nil, err
	}
	return out, nil
}

func (e *U2U) SuggestGasTipCap() (*big.Int, error) {
	var hex hexutil.Big
	if err := e.c.Call("eth_maxPriorityFeePerGas", &hex); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// EstimatePriorityFee request ETA priority fee
func (e *U2U) EstimatePriorityFee(historicalBlocks int, blockNumber *big.Int, feeHistoryPercentile []float64) (*big.Int, error) {
	feeHistory, err := e.FeeHistory(historicalBlocks, blockNumber, feeHistoryPercentile)
	if err != nil {
		return nil, err
	}

	rewards := make(types.Bigs, 0)
	for _, item := range feeHistory.Reward {
		if len(item) == 0 {
			continue
		}
		if item[0].ToInt().Int64() == 0 {
			continue
		}
		rewards = append(rewards, item[0])
	}

	sort.Sort(rewards)
	if len(rewards) == 0 {
		return nil, fmt.Errorf("reward is empty")
	}

	highestIncrease := float64(0)
	highestIncreaseIndex := 0
	for i := range rewards {
		if i == len(rewards)-1 {
			break
		}
		cur := rewards[i].ToInt()
		next := rewards[i+1].ToInt()

		curF := big.NewFloat(0).SetInt(cur)
		v := big.NewFloat(0).Sub(big.NewFloat(0).SetInt(next), curF)
		v = big.NewFloat(0).Quo(v, curF)
		v = big.NewFloat(0).Mul(v, big.NewFloat(100))
		vf, _ := v.Float64()

		if vf > highestIncrease {
			highestIncrease = vf
			highestIncreaseIndex = i
		}
	}
	midIndex := len(rewards) / 2
	if highestIncrease >= PRIORITY_FEE_INCREASE_BOUNDARY && highestIncreaseIndex >= midIndex {

		newRewards := make(types.Bigs, 0)
		for i, item := range rewards {
			if i < highestIncreaseIndex {
				continue
			}
			newRewards = append(newRewards, item)
		}

		return newRewards[midIndex].ToInt(), nil
	}
	return rewards[midIndex].ToInt(), nil
}

// EstimateFee request ETA tx fee
func (e *U2U) EstimateFee() (*EstimateFee, error) {
	header, err := e.GetBlockHeaderByNumber(nil, false)
	if err != nil {
		return nil, err
	}
	priorityFee, err := e.SuggestGasTipCap()
	if err != nil {
		return nil, err
	}
	potentialMaxFee := big.NewInt(1).Mul(header.BaseFee, getBaseFeeMultiplier(header.BaseFee))
	potentialMaxFee = big.NewInt(1).Div(potentialMaxFee, big.NewInt(10))

	maxFeePerGas := big.NewInt(0)

	if priorityFee.Cmp(potentialMaxFee) > 0 {
		maxFeePerGas = big.NewInt(1).Add(potentialMaxFee, priorityFee)
	} else {
		maxFeePerGas = potentialMaxFee
	}

	fee := &EstimateFee{
		BaseFee:              header.BaseFee,
		MaxPriorityFeePerGas: priorityFee,
		MaxFeePerGas:         maxFeePerGas,
	}
	return fee, nil
}

// DecodeParameters decode input data of smc call
func (e *U2U) DecodeParameters(parameters []string, data []byte) ([]interface{}, error) {
	return e.utils.DecodeParameters(parameters, data)
}

// EncodeParameters encode input data of smc call
func (e *U2U) EncodeParameters(parameters []string, data []interface{}) ([]byte, error) {
	return e.utils.EncodeParameters(parameters, data)
}

// SignText signs raw text message
// keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
func (e *U2U) SignText(data []byte) ([]byte, error) {
	hashData := accounts.TextHash(data)
	signature, err := crypto.Sign(hashData, e.privateKey)
	if err != nil {
		return nil, err
	}
	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	return signature, nil
}

// EcSign signs raw hash message
func (e *U2U) EcSign(hashData []byte) ([]byte, error) {
	signature, err := crypto.Sign(hashData, e.privateKey)
	if err != nil {
		return nil, err
	}
	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	return signature, nil
}

//// SignTypedData signs EIP-712 conformant typed data
//// hash = keccak256("\x19${byteVersion}${domainSeparator}${hashStruct(message)}")
//// It returns
//// - the signature,
//// - and/or any error
//func (e *U2U) SignTypedData(data core.TypedData) ([]byte, error) {
//	if e.privateKey == nil {
//		return nil, fmt.Errorf("please setup private key before signing")
//	}
//	domainSeparator, err := data.HashStruct("EIP712Domain", data.Domain.Map())
//	if err != nil {
//		return nil, err
//	}
//
//	typedDataHash, err := data.HashStruct(data.PrimaryType, data.Message)
//	if err != nil {
//		return nil, err
//	}
//
//	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
//
//	sighash := crypto.Keccak256(rawData)
//
//	signature, err := crypto.Sign(sighash, e.privateKey)
//	if err != nil {
//		return nil, err
//	}
//
//	if signature[64] == 0 || signature[64] == 1 {
//		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
//	}
//
//	return signature, nil
//}

func getBaseFeeMultiplier(baseFee *big.Int) *big.Int {
	u := utils.Utils{}
	if baseFee.Cmp(u.ToGWei(40)) <= 0 {
		return big.NewInt(20)
	}
	if baseFee.Cmp(u.ToGWei(100)) <= 0 {
		return big.NewInt(16)
	}
	if baseFee.Cmp(u.ToGWei(200)) <= 0 {
		return big.NewInt(14)
	}
	return big.NewInt(12)
}
