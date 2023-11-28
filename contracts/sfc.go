// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/unicornultrafoundation/go-u2u/libs"
	"github.com/unicornultrafoundation/go-u2u/libs/accounts/abi"
	"github.com/unicornultrafoundation/go-u2u/libs/accounts/abi/bind"
	"github.com/unicornultrafoundation/go-u2u/libs/common"
	"github.com/unicornultrafoundation/go-u2u/libs/core/types"
	"github.com/unicornultrafoundation/go-u2u/libs/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SFCMetaData contains all meta data concerning the SFC contract.
var SFCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"}],\"name\":\"getEpochSnapshot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBaseRewardWeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalTxRewardWeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseRewardPerSecond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SFCABI is the input ABI used to generate the binding from.
// Deprecated: Use SFCMetaData.ABI instead.
var SFCABI = SFCMetaData.ABI

// SFC is an auto generated Go binding around an Ethereum contract.
type SFC struct {
	SFCCaller     // Read-only binding to the contract
	SFCTransactor // Write-only binding to the contract
	SFCFilterer   // Log filterer for contract events
}

// SFCCaller is an auto generated read-only Go binding around an Ethereum contract.
type SFCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SFCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SFCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SFCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SFCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SFCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SFCSession struct {
	Contract     *SFC              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SFCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SFCCallerSession struct {
	Contract *SFCCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SFCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SFCTransactorSession struct {
	Contract     *SFCTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SFCRaw is an auto generated low-level Go binding around an Ethereum contract.
type SFCRaw struct {
	Contract *SFC // Generic contract binding to access the raw methods on
}

// SFCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SFCCallerRaw struct {
	Contract *SFCCaller // Generic read-only contract binding to access the raw methods on
}

// SFCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SFCTransactorRaw struct {
	Contract *SFCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSFC creates a new instance of SFC, bound to a specific deployed contract.
func NewSFC(address common.Address, backend bind.ContractBackend) (*SFC, error) {
	contract, err := bindSFC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SFC{SFCCaller: SFCCaller{contract: contract}, SFCTransactor: SFCTransactor{contract: contract}, SFCFilterer: SFCFilterer{contract: contract}}, nil
}

// NewSFCCaller creates a new read-only instance of SFC, bound to a specific deployed contract.
func NewSFCCaller(address common.Address, caller bind.ContractCaller) (*SFCCaller, error) {
	contract, err := bindSFC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SFCCaller{contract: contract}, nil
}

// NewSFCTransactor creates a new write-only instance of SFC, bound to a specific deployed contract.
func NewSFCTransactor(address common.Address, transactor bind.ContractTransactor) (*SFCTransactor, error) {
	contract, err := bindSFC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SFCTransactor{contract: contract}, nil
}

// NewSFCFilterer creates a new log filterer instance of SFC, bound to a specific deployed contract.
func NewSFCFilterer(address common.Address, filterer bind.ContractFilterer) (*SFCFilterer, error) {
	contract, err := bindSFC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SFCFilterer{contract: contract}, nil
}

// bindSFC binds a generic wrapper to an already deployed contract.
func bindSFC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SFCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SFC *SFCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SFC.Contract.SFCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SFC *SFCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SFC.Contract.SFCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SFC *SFCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SFC.Contract.SFCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SFC *SFCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SFC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SFC *SFCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SFC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SFC *SFCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SFC.Contract.contract.Transact(opts, method, params...)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_SFC *SFCCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SFC.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_SFC *SFCSession) CurrentEpoch() (*big.Int, error) {
	return _SFC.Contract.CurrentEpoch(&_SFC.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_SFC *SFCCallerSession) CurrentEpoch() (*big.Int, error) {
	return _SFC.Contract.CurrentEpoch(&_SFC.CallOpts)
}

// GetEpochSnapshot is a free data retrieval call binding the contract method 0x39b80c00.
//
// Solidity: function getEpochSnapshot(uint256 _epoch) view returns(uint256 endTime, uint256 epochFee, uint256 totalBaseRewardWeight, uint256 totalTxRewardWeight, uint256 baseRewardPerSecond, uint256 totalStake, uint256 totalSupply)
func (_SFC *SFCCaller) GetEpochSnapshot(opts *bind.CallOpts, _epoch *big.Int) (struct {
	EndTime               *big.Int
	EpochFee              *big.Int
	TotalBaseRewardWeight *big.Int
	TotalTxRewardWeight   *big.Int
	BaseRewardPerSecond   *big.Int
	TotalStake            *big.Int
	TotalSupply           *big.Int
}, error) {
	var out []interface{}
	err := _SFC.contract.Call(opts, &out, "getEpochSnapshot", _epoch)

	outstruct := new(struct {
		EndTime               *big.Int
		EpochFee              *big.Int
		TotalBaseRewardWeight *big.Int
		TotalTxRewardWeight   *big.Int
		BaseRewardPerSecond   *big.Int
		TotalStake            *big.Int
		TotalSupply           *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EndTime = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.EpochFee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalBaseRewardWeight = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalTxRewardWeight = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.BaseRewardPerSecond = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TotalStake = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.TotalSupply = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetEpochSnapshot is a free data retrieval call binding the contract method 0x39b80c00.
//
// Solidity: function getEpochSnapshot(uint256 _epoch) view returns(uint256 endTime, uint256 epochFee, uint256 totalBaseRewardWeight, uint256 totalTxRewardWeight, uint256 baseRewardPerSecond, uint256 totalStake, uint256 totalSupply)
func (_SFC *SFCSession) GetEpochSnapshot(_epoch *big.Int) (struct {
	EndTime               *big.Int
	EpochFee              *big.Int
	TotalBaseRewardWeight *big.Int
	TotalTxRewardWeight   *big.Int
	BaseRewardPerSecond   *big.Int
	TotalStake            *big.Int
	TotalSupply           *big.Int
}, error) {
	return _SFC.Contract.GetEpochSnapshot(&_SFC.CallOpts, _epoch)
}

// GetEpochSnapshot is a free data retrieval call binding the contract method 0x39b80c00.
//
// Solidity: function getEpochSnapshot(uint256 _epoch) view returns(uint256 endTime, uint256 epochFee, uint256 totalBaseRewardWeight, uint256 totalTxRewardWeight, uint256 baseRewardPerSecond, uint256 totalStake, uint256 totalSupply)
func (_SFC *SFCCallerSession) GetEpochSnapshot(_epoch *big.Int) (struct {
	EndTime               *big.Int
	EpochFee              *big.Int
	TotalBaseRewardWeight *big.Int
	TotalTxRewardWeight   *big.Int
	BaseRewardPerSecond   *big.Int
	TotalStake            *big.Int
	TotalSupply           *big.Int
}, error) {
	return _SFC.Contract.GetEpochSnapshot(&_SFC.CallOpts, _epoch)
}
