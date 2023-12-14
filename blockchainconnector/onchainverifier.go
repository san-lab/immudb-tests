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

// OnChainVerifierMetaData contains all meta data concerning the OnChainVerifier contract.
var OnChainVerifierMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"submitHash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"_preimage\",\"type\":\"bytes32\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\"}],\"name\":\"submitPreimage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"HashSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"\",\"type\":\"bool\"}],\"name\":\"PreimageSubmitted\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"}],\"name\":\"getPendingSubmissions\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\"}],\"name\":\"getStateCheckByBlockNumber\",\"outputs\":[{\"name\":\"submittedHash\",\"type\":\"bytes32\"},{\"name\":\"submittedPreimage\",\"type\":\"bytes32\"},{\"name\":\"verified\",\"type\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_originatorBank\",\"type\":\"address\"},{\"name\":\"_recipientBank\",\"type\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getStateCheckByIndex\",\"outputs\":[{\"name\":\"submittedHash\",\"type\":\"bytes32\"},{\"name\":\"submittedPreimage\",\"type\":\"bytes32\"},{\"name\":\"verified\",\"type\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stateChecks\",\"outputs\":[{\"name\":\"submittedHash\",\"type\":\"bytes32\"},{\"name\":\"submittedPreimage\",\"type\":\"bytes32\"},{\"name\":\"verified\",\"type\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// OnChainVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use OnChainVerifierMetaData.ABI instead.
var OnChainVerifierABI = OnChainVerifierMetaData.ABI

// OnChainVerifier is an auto generated Go binding around an Ethereum contract.
type OnChainVerifier struct {
	OnChainVerifierCaller     // Read-only binding to the contract
	OnChainVerifierTransactor // Write-only binding to the contract
	OnChainVerifierFilterer   // Log filterer for contract events
}

// OnChainVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type OnChainVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnChainVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OnChainVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnChainVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OnChainVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnChainVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OnChainVerifierSession struct {
	Contract     *OnChainVerifier  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OnChainVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OnChainVerifierCallerSession struct {
	Contract *OnChainVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// OnChainVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OnChainVerifierTransactorSession struct {
	Contract     *OnChainVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// OnChainVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type OnChainVerifierRaw struct {
	Contract *OnChainVerifier // Generic contract binding to access the raw methods on
}

// OnChainVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OnChainVerifierCallerRaw struct {
	Contract *OnChainVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// OnChainVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OnChainVerifierTransactorRaw struct {
	Contract *OnChainVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOnChainVerifier creates a new instance of OnChainVerifier, bound to a specific deployed contract.
func NewOnChainVerifier(address common.Address, backend bind.ContractBackend) (*OnChainVerifier, error) {
	contract, err := bindOnChainVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnChainVerifier{OnChainVerifierCaller: OnChainVerifierCaller{contract: contract}, OnChainVerifierTransactor: OnChainVerifierTransactor{contract: contract}, OnChainVerifierFilterer: OnChainVerifierFilterer{contract: contract}}, nil
}

// NewOnChainVerifierCaller creates a new read-only instance of OnChainVerifier, bound to a specific deployed contract.
func NewOnChainVerifierCaller(address common.Address, caller bind.ContractCaller) (*OnChainVerifierCaller, error) {
	contract, err := bindOnChainVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnChainVerifierCaller{contract: contract}, nil
}

// NewOnChainVerifierTransactor creates a new write-only instance of OnChainVerifier, bound to a specific deployed contract.
func NewOnChainVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*OnChainVerifierTransactor, error) {
	contract, err := bindOnChainVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnChainVerifierTransactor{contract: contract}, nil
}

// NewOnChainVerifierFilterer creates a new log filterer instance of OnChainVerifier, bound to a specific deployed contract.
func NewOnChainVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*OnChainVerifierFilterer, error) {
	contract, err := bindOnChainVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnChainVerifierFilterer{contract: contract}, nil
}

// bindOnChainVerifier binds a generic wrapper to an already deployed contract.
func bindOnChainVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OnChainVerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnChainVerifier *OnChainVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnChainVerifier.Contract.OnChainVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnChainVerifier *OnChainVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.OnChainVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnChainVerifier *OnChainVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.OnChainVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnChainVerifier *OnChainVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnChainVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnChainVerifier *OnChainVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnChainVerifier *OnChainVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.contract.Transact(opts, method, params...)
}

// GetPendingSubmissions is a free data retrieval call binding the contract method 0xc392db84.
//
// Solidity: function getPendingSubmissions(address _originatorBank, address _recipientBank) view returns(uint256[])
func (_OnChainVerifier *OnChainVerifierCaller) GetPendingSubmissions(opts *bind.CallOpts, _originatorBank common.Address, _recipientBank common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _OnChainVerifier.contract.Call(opts, &out, "getPendingSubmissions", _originatorBank, _recipientBank)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetPendingSubmissions is a free data retrieval call binding the contract method 0xc392db84.
//
// Solidity: function getPendingSubmissions(address _originatorBank, address _recipientBank) view returns(uint256[])
func (_OnChainVerifier *OnChainVerifierSession) GetPendingSubmissions(_originatorBank common.Address, _recipientBank common.Address) ([]*big.Int, error) {
	return _OnChainVerifier.Contract.GetPendingSubmissions(&_OnChainVerifier.CallOpts, _originatorBank, _recipientBank)
}

// GetPendingSubmissions is a free data retrieval call binding the contract method 0xc392db84.
//
// Solidity: function getPendingSubmissions(address _originatorBank, address _recipientBank) view returns(uint256[])
func (_OnChainVerifier *OnChainVerifierCallerSession) GetPendingSubmissions(_originatorBank common.Address, _recipientBank common.Address) ([]*big.Int, error) {
	return _OnChainVerifier.Contract.GetPendingSubmissions(&_OnChainVerifier.CallOpts, _originatorBank, _recipientBank)
}

// GetStateCheckByBlockNumber is a free data retrieval call binding the contract method 0xd80e96d5.
//
// Solidity: function getStateCheckByBlockNumber(address _originatorBank, address _recipientBank, uint256 _blockNumber) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnChainVerifier *OnChainVerifierCaller) GetStateCheckByBlockNumber(opts *bind.CallOpts, _originatorBank common.Address, _recipientBank common.Address, _blockNumber *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	var out []interface{}
	err := _OnChainVerifier.contract.Call(opts, &out, "getStateCheckByBlockNumber", _originatorBank, _recipientBank, _blockNumber)

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
func (_OnChainVerifier *OnChainVerifierSession) GetStateCheckByBlockNumber(_originatorBank common.Address, _recipientBank common.Address, _blockNumber *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnChainVerifier.Contract.GetStateCheckByBlockNumber(&_OnChainVerifier.CallOpts, _originatorBank, _recipientBank, _blockNumber)
}

// GetStateCheckByBlockNumber is a free data retrieval call binding the contract method 0xd80e96d5.
//
// Solidity: function getStateCheckByBlockNumber(address _originatorBank, address _recipientBank, uint256 _blockNumber) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnChainVerifier *OnChainVerifierCallerSession) GetStateCheckByBlockNumber(_originatorBank common.Address, _recipientBank common.Address, _blockNumber *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnChainVerifier.Contract.GetStateCheckByBlockNumber(&_OnChainVerifier.CallOpts, _originatorBank, _recipientBank, _blockNumber)
}

// GetStateCheckByIndex is a free data retrieval call binding the contract method 0xa7f0ab94.
//
// Solidity: function getStateCheckByIndex(address _originatorBank, address _recipientBank, uint256 index) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnChainVerifier *OnChainVerifierCaller) GetStateCheckByIndex(opts *bind.CallOpts, _originatorBank common.Address, _recipientBank common.Address, index *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	var out []interface{}
	err := _OnChainVerifier.contract.Call(opts, &out, "getStateCheckByIndex", _originatorBank, _recipientBank, index)

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
func (_OnChainVerifier *OnChainVerifierSession) GetStateCheckByIndex(_originatorBank common.Address, _recipientBank common.Address, index *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnChainVerifier.Contract.GetStateCheckByIndex(&_OnChainVerifier.CallOpts, _originatorBank, _recipientBank, index)
}

// GetStateCheckByIndex is a free data retrieval call binding the contract method 0xa7f0ab94.
//
// Solidity: function getStateCheckByIndex(address _originatorBank, address _recipientBank, uint256 index) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnChainVerifier *OnChainVerifierCallerSession) GetStateCheckByIndex(_originatorBank common.Address, _recipientBank common.Address, index *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnChainVerifier.Contract.GetStateCheckByIndex(&_OnChainVerifier.CallOpts, _originatorBank, _recipientBank, index)
}

// StateChecks is a free data retrieval call binding the contract method 0x96535fed.
//
// Solidity: function stateChecks(address , address , uint256 ) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnChainVerifier *OnChainVerifierCaller) StateChecks(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	var out []interface{}
	err := _OnChainVerifier.contract.Call(opts, &out, "stateChecks", arg0, arg1, arg2)

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
func (_OnChainVerifier *OnChainVerifierSession) StateChecks(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnChainVerifier.Contract.StateChecks(&_OnChainVerifier.CallOpts, arg0, arg1, arg2)
}

// StateChecks is a free data retrieval call binding the contract method 0x96535fed.
//
// Solidity: function stateChecks(address , address , uint256 ) view returns(bytes32 submittedHash, bytes32 submittedPreimage, bool verified, uint256 blockNumber)
func (_OnChainVerifier *OnChainVerifierCallerSession) StateChecks(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (struct {
	SubmittedHash     [32]byte
	SubmittedPreimage [32]byte
	Verified          bool
	BlockNumber       *big.Int
}, error) {
	return _OnChainVerifier.Contract.StateChecks(&_OnChainVerifier.CallOpts, arg0, arg1, arg2)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OnChainVerifier *OnChainVerifierCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OnChainVerifier.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OnChainVerifier *OnChainVerifierSession) Version() (string, error) {
	return _OnChainVerifier.Contract.Version(&_OnChainVerifier.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_OnChainVerifier *OnChainVerifierCallerSession) Version() (string, error) {
	return _OnChainVerifier.Contract.Version(&_OnChainVerifier.CallOpts)
}

// SubmitHash is a paid mutator transaction binding the contract method 0x542686b0.
//
// Solidity: function submitHash(address _originatorBank, address _recipientBank, bytes32 _hash) returns()
func (_OnChainVerifier *OnChainVerifierTransactor) SubmitHash(opts *bind.TransactOpts, _originatorBank common.Address, _recipientBank common.Address, _hash [32]byte) (*types.Transaction, error) {
	return _OnChainVerifier.contract.Transact(opts, "submitHash", _originatorBank, _recipientBank, _hash)
}

// SubmitHash is a paid mutator transaction binding the contract method 0x542686b0.
//
// Solidity: function submitHash(address _originatorBank, address _recipientBank, bytes32 _hash) returns()
func (_OnChainVerifier *OnChainVerifierSession) SubmitHash(_originatorBank common.Address, _recipientBank common.Address, _hash [32]byte) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.SubmitHash(&_OnChainVerifier.TransactOpts, _originatorBank, _recipientBank, _hash)
}

// SubmitHash is a paid mutator transaction binding the contract method 0x542686b0.
//
// Solidity: function submitHash(address _originatorBank, address _recipientBank, bytes32 _hash) returns()
func (_OnChainVerifier *OnChainVerifierTransactorSession) SubmitHash(_originatorBank common.Address, _recipientBank common.Address, _hash [32]byte) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.SubmitHash(&_OnChainVerifier.TransactOpts, _originatorBank, _recipientBank, _hash)
}

// SubmitPreimage is a paid mutator transaction binding the contract method 0x2eb05625.
//
// Solidity: function submitPreimage(address _originatorBank, address _recipientBank, bytes32 _preimage, uint256 _blockNumber) returns()
func (_OnChainVerifier *OnChainVerifierTransactor) SubmitPreimage(opts *bind.TransactOpts, _originatorBank common.Address, _recipientBank common.Address, _preimage [32]byte, _blockNumber *big.Int) (*types.Transaction, error) {
	return _OnChainVerifier.contract.Transact(opts, "submitPreimage", _originatorBank, _recipientBank, _preimage, _blockNumber)
}

// SubmitPreimage is a paid mutator transaction binding the contract method 0x2eb05625.
//
// Solidity: function submitPreimage(address _originatorBank, address _recipientBank, bytes32 _preimage, uint256 _blockNumber) returns()
func (_OnChainVerifier *OnChainVerifierSession) SubmitPreimage(_originatorBank common.Address, _recipientBank common.Address, _preimage [32]byte, _blockNumber *big.Int) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.SubmitPreimage(&_OnChainVerifier.TransactOpts, _originatorBank, _recipientBank, _preimage, _blockNumber)
}

// SubmitPreimage is a paid mutator transaction binding the contract method 0x2eb05625.
//
// Solidity: function submitPreimage(address _originatorBank, address _recipientBank, bytes32 _preimage, uint256 _blockNumber) returns()
func (_OnChainVerifier *OnChainVerifierTransactorSession) SubmitPreimage(_originatorBank common.Address, _recipientBank common.Address, _preimage [32]byte, _blockNumber *big.Int) (*types.Transaction, error) {
	return _OnChainVerifier.Contract.SubmitPreimage(&_OnChainVerifier.TransactOpts, _originatorBank, _recipientBank, _preimage, _blockNumber)
}

// OnChainVerifierHashSubmittedIterator is returned from FilterHashSubmitted and is used to iterate over the raw logs and unpacked data for HashSubmitted events raised by the OnChainVerifier contract.
type OnChainVerifierHashSubmittedIterator struct {
	Event *OnChainVerifierHashSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnChainVerifierHashSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnChainVerifierHashSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnChainVerifierHashSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnChainVerifierHashSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnChainVerifierHashSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnChainVerifierHashSubmitted represents a HashSubmitted event raised by the OnChainVerifier contract.
type OnChainVerifierHashSubmitted struct {
	Arg0 common.Address
	Arg1 common.Address
	Arg2 [32]byte
	Arg3 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterHashSubmitted is a free log retrieval operation binding the contract event 0xbffd7b7f91e2d2d9779e3361f717fb43ef16c8b9239bfb9a2a983ba95cf7b7d6.
//
// Solidity: event HashSubmitted(address arg0, address arg1, bytes32 arg2, uint256 arg3)
func (_OnChainVerifier *OnChainVerifierFilterer) FilterHashSubmitted(opts *bind.FilterOpts) (*OnChainVerifierHashSubmittedIterator, error) {

	logs, sub, err := _OnChainVerifier.contract.FilterLogs(opts, "HashSubmitted")
	if err != nil {
		return nil, err
	}
	return &OnChainVerifierHashSubmittedIterator{contract: _OnChainVerifier.contract, event: "HashSubmitted", logs: logs, sub: sub}, nil
}

// WatchHashSubmitted is a free log subscription operation binding the contract event 0xbffd7b7f91e2d2d9779e3361f717fb43ef16c8b9239bfb9a2a983ba95cf7b7d6.
//
// Solidity: event HashSubmitted(address arg0, address arg1, bytes32 arg2, uint256 arg3)
func (_OnChainVerifier *OnChainVerifierFilterer) WatchHashSubmitted(opts *bind.WatchOpts, sink chan<- *OnChainVerifierHashSubmitted) (event.Subscription, error) {

	logs, sub, err := _OnChainVerifier.contract.WatchLogs(opts, "HashSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnChainVerifierHashSubmitted)
				if err := _OnChainVerifier.contract.UnpackLog(event, "HashSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHashSubmitted is a log parse operation binding the contract event 0xbffd7b7f91e2d2d9779e3361f717fb43ef16c8b9239bfb9a2a983ba95cf7b7d6.
//
// Solidity: event HashSubmitted(address arg0, address arg1, bytes32 arg2, uint256 arg3)
func (_OnChainVerifier *OnChainVerifierFilterer) ParseHashSubmitted(log types.Log) (*OnChainVerifierHashSubmitted, error) {
	event := new(OnChainVerifierHashSubmitted)
	if err := _OnChainVerifier.contract.UnpackLog(event, "HashSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnChainVerifierPreimageSubmittedIterator is returned from FilterPreimageSubmitted and is used to iterate over the raw logs and unpacked data for PreimageSubmitted events raised by the OnChainVerifier contract.
type OnChainVerifierPreimageSubmittedIterator struct {
	Event *OnChainVerifierPreimageSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnChainVerifierPreimageSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnChainVerifierPreimageSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnChainVerifierPreimageSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnChainVerifierPreimageSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnChainVerifierPreimageSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnChainVerifierPreimageSubmitted represents a PreimageSubmitted event raised by the OnChainVerifier contract.
type OnChainVerifierPreimageSubmitted struct {
	Arg0 common.Address
	Arg1 common.Address
	Arg2 [32]byte
	Arg3 *big.Int
	Arg4 bool
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterPreimageSubmitted is a free log retrieval operation binding the contract event 0x97f8c3a2688124b2be511340aa0247f2236bd4e9905dea96412dd8083b1c0b14.
//
// Solidity: event PreimageSubmitted(address arg0, address arg1, bytes32 arg2, uint256 arg3, bool arg4)
func (_OnChainVerifier *OnChainVerifierFilterer) FilterPreimageSubmitted(opts *bind.FilterOpts) (*OnChainVerifierPreimageSubmittedIterator, error) {

	logs, sub, err := _OnChainVerifier.contract.FilterLogs(opts, "PreimageSubmitted")
	if err != nil {
		return nil, err
	}
	return &OnChainVerifierPreimageSubmittedIterator{contract: _OnChainVerifier.contract, event: "PreimageSubmitted", logs: logs, sub: sub}, nil
}

// WatchPreimageSubmitted is a free log subscription operation binding the contract event 0x97f8c3a2688124b2be511340aa0247f2236bd4e9905dea96412dd8083b1c0b14.
//
// Solidity: event PreimageSubmitted(address arg0, address arg1, bytes32 arg2, uint256 arg3, bool arg4)
func (_OnChainVerifier *OnChainVerifierFilterer) WatchPreimageSubmitted(opts *bind.WatchOpts, sink chan<- *OnChainVerifierPreimageSubmitted) (event.Subscription, error) {

	logs, sub, err := _OnChainVerifier.contract.WatchLogs(opts, "PreimageSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnChainVerifierPreimageSubmitted)
				if err := _OnChainVerifier.contract.UnpackLog(event, "PreimageSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePreimageSubmitted is a log parse operation binding the contract event 0x97f8c3a2688124b2be511340aa0247f2236bd4e9905dea96412dd8083b1c0b14.
//
// Solidity: event PreimageSubmitted(address arg0, address arg1, bytes32 arg2, uint256 arg3, bool arg4)
func (_OnChainVerifier *OnChainVerifierFilterer) ParsePreimageSubmitted(log types.Log) (*OnChainVerifierPreimageSubmitted, error) {
	event := new(OnChainVerifierPreimageSubmitted)
	if err := _OnChainVerifier.contract.UnpackLog(event, "PreimageSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
