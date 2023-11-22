// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchainconnector

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

// OnchainVerifierMetaData contains all meta data concerning the OnchainVerifier contract.
var OnchainVerifierMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"_preimage\",\"type\":\"bytes32\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\"}],\"name\":\"submitPreimage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"submitHash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stateChecks\",\"outputs\":[{\"name\":\"submittedHash\",\"type\":\"bytes32\"},{\"name\":\"submittedPreimage\",\"type\":\"bytes32\"},{\"name\":\"verified\",\"type\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getStateCheckByIndex\",\"outputs\":[{\"name\":\"submittedHash\",\"type\":\"bytes32\"},{\"name\":\"submittedPreimage\",\"type\":\"bytes32\"},{\"name\":\"verified\",\"type\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"}],\"name\":\"getPendingSubmissions\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\"}],\"name\":\"getStateCheckByBlockNumber\",\"outputs\":[{\"name\":\"submittedHash\",\"type\":\"bytes32\"},{\"name\":\"submittedPreimage\",\"type\":\"bytes32\"},{\"name\":\"verified\",\"type\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]",
}

// OnchainVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use OnchainVerifierMetaData.ABI instead.
var OnchainVerifierABI = OnchainVerifierMetaData.ABI

// OnchainVerifier is an auto generated Go binding around an Ethereum contract.
type OnchainVerifier struct {
	OnchainVerifierCaller     // Read-only binding to the contract
	OnchainVerifierTransactor // Write-only binding to the contract
	OnchainVerifierFilterer   // Log filterer for contract events
}

// OnchainVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type OnchainVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OnchainVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OnchainVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OnchainVerifierSession struct {
	Contract     *OnchainVerifier  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OnchainVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OnchainVerifierCallerSession struct {
	Contract *OnchainVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// OnchainVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OnchainVerifierTransactorSession struct {
	Contract     *OnchainVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// OnchainVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type OnchainVerifierRaw struct {
	Contract *OnchainVerifier // Generic contract binding to access the raw methods on
}

// OnchainVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OnchainVerifierCallerRaw struct {
	Contract *OnchainVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// OnchainVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OnchainVerifierTransactorRaw struct {
	Contract *OnchainVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOnchainVerifier creates a new instance of OnchainVerifier, bound to a specific deployed contract.
func NewOnchainVerifier(address common.Address, backend bind.ContractBackend) (*OnchainVerifier, error) {
	contract, err := bindOnchainVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnchainVerifier{OnchainVerifierCaller: OnchainVerifierCaller{contract: contract}, OnchainVerifierTransactor: OnchainVerifierTransactor{contract: contract}, OnchainVerifierFilterer: OnchainVerifierFilterer{contract: contract}}, nil
}

// NewOnchainVerifierCaller creates a new read-only instance of OnchainVerifier, bound to a specific deployed contract.
func NewOnchainVerifierCaller(address common.Address, caller bind.ContractCaller) (*OnchainVerifierCaller, error) {
	contract, err := bindOnchainVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnchainVerifierCaller{contract: contract}, nil
}

// NewOnchainVerifierTransactor creates a new write-only instance of OnchainVerifier, bound to a specific deployed contract.
func NewOnchainVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*OnchainVerifierTransactor, error) {
	contract, err := bindOnchainVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnchainVerifierTransactor{contract: contract}, nil
}

// NewOnchainVerifierFilterer creates a new log filterer instance of OnchainVerifier, bound to a specific deployed contract.
func NewOnchainVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*OnchainVerifierFilterer, error) {
	contract, err := bindOnchainVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnchainVerifierFilterer{contract: contract}, nil
}

// bindOnchainVerifier binds a generic wrapper to an already deployed contract.
func bindOnchainVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OnchainVerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnchainVerifier *OnchainVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnchainVerifier.Contract.OnchainVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnchainVerifier *OnchainVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.OnchainVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnchainVerifier *OnchainVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.OnchainVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnchainVerifier *OnchainVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnchainVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnchainVerifier *OnchainVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnchainVerifier *OnchainVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.contract.Transact(opts, method, params...)
}

// GetPendingSubmissions is a free data retrieval call binding the contract method 0xc392db84.
//
// Solidity: function getPendingSubmissions(address _originatorBank, address _recipientBank) view returns(uint256[])
func (_OnchainVerifier *OnchainVerifierCaller) GetPendingSubmissions(opts *bind.CallOpts, _originatorBank common.Address, _recipientBank common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _OnchainVerifier.contract.Call(opts, &out, "getPendingSubmissions", _originatorBank, _recipientBank)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetPendingSubmissions is a free data retrieval call binding the contract method 0xc392db84.
//
// Solidity: function getPendingSubmissions(address _originatorBank, address _recipientBank) view returns(uint256[])
func (_OnchainVerifier *OnchainVerifierSession) GetPendingSubmissions(_originatorBank common.Address, _recipientBank common.Address) ([]*big.Int, error) {
	return _OnchainVerifier.Contract.GetPendingSubmissions(&_OnchainVerifier.CallOpts, _originatorBank, _recipientBank)
}

// GetPendingSubmissions is a free data retrieval call binding the contract method 0xc392db84.
//
// Solidity: function getPendingSubmissions(address _originatorBank, address _recipientBank) view returns(uint256[])
func (_OnchainVerifier *OnchainVerifierCallerSession) GetPendingSubmissions(_originatorBank common.Address, _recipientBank common.Address) ([]*big.Int, error) {
	return _OnchainVerifier.Contract.GetPendingSubmissions(&_OnchainVerifier.CallOpts, _originatorBank, _recipientBank)
}

// GetStateCheckByBlockNumber is a free data retrieval call binding the contract method 0xd80e96d5.
//
// Solidity: function getStateCheckByBlockNumber(address _originatorBank, address _recipientBank, uint256 _blockNumber) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierCaller) GetStateCheckByBlockNumber(opts *bind.CallOpts, _originatorBank common.Address, _recipientBank common.Address, _blockNumber *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	var out []interface{}
	err := _OnchainVerifier.contract.Call(opts, &out, "getStateCheckByBlockNumber", _originatorBank, _recipientBank, _blockNumber)

	outstruct := new(struct {
		SubmittedHash     [32]byte
		SubmittedPreimage [32]byte
		Verified          bool
		BlockNumber       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubmittedHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.SubmittedPreimage = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Verified = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.BlockNumber = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStateCheckByBlockNumber is a free data retrieval call binding the contract method 0xd80e96d5.
//
// Solidity: function getStateCheckByBlockNumber(address _originatorBank, address _recipientBank, uint256 _blockNumber) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierSession) GetStateCheckByBlockNumber(_originatorBank common.Address, _recipientBank common.Address, _blockNumber *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnchainVerifier.Contract.GetStateCheckByBlockNumber(&_OnchainVerifier.CallOpts, _originatorBank, _recipientBank, _blockNumber)
}

// GetStateCheckByBlockNumber is a free data retrieval call binding the contract method 0xd80e96d5.
//
// Solidity: function getStateCheckByBlockNumber(address _originatorBank, address _recipientBank, uint256 _blockNumber) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierCallerSession) GetStateCheckByBlockNumber(_originatorBank common.Address, _recipientBank common.Address, _blockNumber *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnchainVerifier.Contract.GetStateCheckByBlockNumber(&_OnchainVerifier.CallOpts, _originatorBank, _recipientBank, _blockNumber)
}

// GetStateCheckByIndex is a free data retrieval call binding the contract method 0xa7f0ab94.
//
// Solidity: function getStateCheckByIndex(address _originatorBank, address _recipientBank, uint256 index) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierCaller) GetStateCheckByIndex(opts *bind.CallOpts, _originatorBank common.Address, _recipientBank common.Address, index *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	var out []interface{}
	err := _OnchainVerifier.contract.Call(opts, &out, "getStateCheckByIndex", _originatorBank, _recipientBank, index)

	outstruct := new(struct {
		SubmittedHash     [32]byte
		SubmittedPreimage [32]byte
		Verified          bool
		BlockNumber       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubmittedHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.SubmittedPreimage = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Verified = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.BlockNumber = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStateCheckByIndex is a free data retrieval call binding the contract method 0xa7f0ab94.
//
// Solidity: function getStateCheckByIndex(address _originatorBank, address _recipientBank, uint256 index) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierSession) GetStateCheckByIndex(_originatorBank common.Address, _recipientBank common.Address, index *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnchainVerifier.Contract.GetStateCheckByIndex(&_OnchainVerifier.CallOpts, _originatorBank, _recipientBank, index)
}

// GetStateCheckByIndex is a free data retrieval call binding the contract method 0xa7f0ab94.
//
// Solidity: function getStateCheckByIndex(address _originatorBank, address _recipientBank, uint256 index) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierCallerSession) GetStateCheckByIndex(_originatorBank common.Address, _recipientBank common.Address, index *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnchainVerifier.Contract.GetStateCheckByIndex(&_OnchainVerifier.CallOpts, _originatorBank, _recipientBank, index)
}

// StateChecks is a free data retrieval call binding the contract method 0x96535fed.
//
// Solidity: function stateChecks(address , address , uint256 ) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierCaller) StateChecks(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	var out []interface{}
	err := _OnchainVerifier.contract.Call(opts, &out, "stateChecks", arg0, arg1, arg2)

	outstruct := new(struct {
		SubmittedHash     [32]byte
		SubmittedPreimage [32]byte
		Verified          bool
		BlockNumber       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubmittedHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.SubmittedPreimage = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Verified = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.BlockNumber = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// StateChecks is a free data retrieval call binding the contract method 0x96535fed.
//
// Solidity: function stateChecks(address , address , uint256 ) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierSession) StateChecks(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnchainVerifier.Contract.StateChecks(&_OnchainVerifier.CallOpts, arg0, arg1, arg2)
}

// StateChecks is a free data retrieval call binding the contract method 0x96535fed.
//
// Solidity: function stateChecks(address , address , uint256 ) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnchainVerifier *OnchainVerifierCallerSession) StateChecks(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnchainVerifier.Contract.StateChecks(&_OnchainVerifier.CallOpts, arg0, arg1, arg2)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OnchainVerifier *OnchainVerifierCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OnchainVerifier.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OnchainVerifier *OnchainVerifierSession) Version() (string, error) {
	return _OnchainVerifier.Contract.Version(&_OnchainVerifier.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OnchainVerifier *OnchainVerifierCallerSession) Version() (string, error) {
	return _OnchainVerifier.Contract.Version(&_OnchainVerifier.CallOpts)
}

// SubmitHash is a paid mutator transaction binding the contract method 0x542686b0.
//
// Solidity: function submitHash(address _originatorBank, address _recipientBank, bytes32 _hash) returns()
func (_OnchainVerifier *OnchainVerifierTransactor) SubmitHash(opts *bind.TransactOpts, _originatorBank common.Address, _recipientBank common.Address, _hash [32]byte) (*types.Transaction, error) {
	return _OnchainVerifier.contract.Transact(opts, "submitHash", _originatorBank, _recipientBank, _hash)
}

// SubmitHash is a paid mutator transaction binding the contract method 0x542686b0.
//
// Solidity: function submitHash(address _originatorBank, address _recipientBank, bytes32 _hash) returns()
func (_OnchainVerifier *OnchainVerifierSession) SubmitHash(_originatorBank common.Address, _recipientBank common.Address, _hash [32]byte) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.SubmitHash(&_OnchainVerifier.TransactOpts, _originatorBank, _recipientBank, _hash)
}

// SubmitHash is a paid mutator transaction binding the contract method 0x542686b0.
//
// Solidity: function submitHash(address _originatorBank, address _recipientBank, bytes32 _hash) returns()
func (_OnchainVerifier *OnchainVerifierTransactorSession) SubmitHash(_originatorBank common.Address, _recipientBank common.Address, _hash [32]byte) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.SubmitHash(&_OnchainVerifier.TransactOpts, _originatorBank, _recipientBank, _hash)
}

// SubmitPreimage is a paid mutator transaction binding the contract method 0x2eb05625.
//
// Solidity: function submitPreimage(address _originatorBank, address _recipientBank, bytes32 _preimage, uint256 _blockNumber) returns()
func (_OnchainVerifier *OnchainVerifierTransactor) SubmitPreimage(opts *bind.TransactOpts, _originatorBank common.Address, _recipientBank common.Address, _preimage [32]byte, _blockNumber *big.Int) (*types.Transaction, error) {
	return _OnchainVerifier.contract.Transact(opts, "submitPreimage", _originatorBank, _recipientBank, _preimage, _blockNumber)
}

// SubmitPreimage is a paid mutator transaction binding the contract method 0x2eb05625.
//
// Solidity: function submitPreimage(address _originatorBank, address _recipientBank, bytes32 _preimage, uint256 _blockNumber) returns()
func (_OnchainVerifier *OnchainVerifierSession) SubmitPreimage(_originatorBank common.Address, _recipientBank common.Address, _preimage [32]byte, _blockNumber *big.Int) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.SubmitPreimage(&_OnchainVerifier.TransactOpts, _originatorBank, _recipientBank, _preimage, _blockNumber)
}

// SubmitPreimage is a paid mutator transaction binding the contract method 0x2eb05625.
//
// Solidity: function submitPreimage(address _originatorBank, address _recipientBank, bytes32 _preimage, uint256 _blockNumber) returns()
func (_OnchainVerifier *OnchainVerifierTransactorSession) SubmitPreimage(_originatorBank common.Address, _recipientBank common.Address, _preimage [32]byte, _blockNumber *big.Int) (*types.Transaction, error) {
	return _OnchainVerifier.Contract.SubmitPreimage(&_OnchainVerifier.TransactOpts, _originatorBank, _recipientBank, _preimage, _blockNumber)
}
