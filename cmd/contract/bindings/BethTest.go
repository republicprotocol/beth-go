// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BethtestABI is the input ABI used to generate the binding from.
const BethtestABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_dataToDelete\",\"type\":\"uint256\"}],\"name\":\"remove\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"read\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"set\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"size\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"increment\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"append\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// BethtestBin is the compiled bytecode used for deploying new contracts.
const BethtestBin = `608060405234801561001057600080fd5b5061043f806100206000396000f300608060405260043610610083576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634cc822151461008857806357de26a4146100b557806360fe47b1146100e0578063949d225d1461010d5780639507d39a14610138578063d09de08a14610184578063e33b87071461019b575b600080fd5b34801561009457600080fd5b506100b3600480360381019080803590602001909291905050506101c8565b005b3480156100c157600080fd5b506100ca6102e7565b6040518082815260200191505060405180910390f35b3480156100ec57600080fd5b5061010b600480360381019080803590602001909291905050506102f1565b005b34801561011957600080fd5b506101226102fb565b6040518082815260200191505060405180910390f35b34801561014457600080fd5b5061016360048036038101908080359060200190929190505050610307565b60405180838152602001821515151581526020019250505060405180910390f35b34801561019057600080fd5b5061019961034b565b005b3480156101a757600080fd5b506101c66004803603810190808035906020019092919050505061035f565b005b60006001600083815260200190815260200160002054905060008114156101ee576102e3565b60016000805490501015156102cc57600060016000806001850381548110151561021457fe5b9060005260206000200154815260200190815260200160002081905550600060016000805490500381548110151561024857fe5b906000526020600020015460006001830381548110151561026557fe5b9060005260206000200181905550600060016000805490500381548110151561028a57fe5b9060005260206000200160009055806001600080600185038154811015156102ae57fe5b90600052602060002001548152602001908152602001600020819055505b60008054809190600190036102e191906103c2565b505b5050565b6000600254905090565b8060028190555050565b60008080549050905090565b600080600060016000858152602001908152602001600020549050600081141561033a5760008081915092509250610345565b600181036001925092505b50915091565b600260008154809291906001019190505550565b600061036a82610307565b9150508015156103be57600082908060018154018082558091505090600182039060005260206000200160009091929091909150555060008054905060016000848152602001908152602001600020819055505b5050565b8154818355818111156103e9578183600052602060002091820191016103e891906103ee565b5b505050565b61041091905b8082111561040c5760008160009055506001016103f4565b5090565b905600a165627a7a72305820277b4e2b2a4c40cf5756cfb58aa8648917ec198aeeee73f1f4a305c4f1eca4080029`

// DeployBethtest deploys a new Ethereum contract, binding an instance of Bethtest to it.
func DeployBethtest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Bethtest, error) {
	parsed, err := abi.JSON(strings.NewReader(BethtestABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BethtestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bethtest{BethtestCaller: BethtestCaller{contract: contract}, BethtestTransactor: BethtestTransactor{contract: contract}, BethtestFilterer: BethtestFilterer{contract: contract}}, nil
}

// Bethtest is an auto generated Go binding around an Ethereum contract.
type Bethtest struct {
	BethtestCaller     // Read-only binding to the contract
	BethtestTransactor // Write-only binding to the contract
	BethtestFilterer   // Log filterer for contract events
}

// BethtestCaller is an auto generated read-only Go binding around an Ethereum contract.
type BethtestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BethtestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BethtestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BethtestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BethtestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BethtestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BethtestSession struct {
	Contract     *Bethtest         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BethtestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BethtestCallerSession struct {
	Contract *BethtestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BethtestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BethtestTransactorSession struct {
	Contract     *BethtestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BethtestRaw is an auto generated low-level Go binding around an Ethereum contract.
type BethtestRaw struct {
	Contract *Bethtest // Generic contract binding to access the raw methods on
}

// BethtestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BethtestCallerRaw struct {
	Contract *BethtestCaller // Generic read-only contract binding to access the raw methods on
}

// BethtestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BethtestTransactorRaw struct {
	Contract *BethtestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBethtest creates a new instance of Bethtest, bound to a specific deployed contract.
func NewBethtest(address common.Address, backend bind.ContractBackend) (*Bethtest, error) {
	contract, err := bindBethtest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bethtest{BethtestCaller: BethtestCaller{contract: contract}, BethtestTransactor: BethtestTransactor{contract: contract}, BethtestFilterer: BethtestFilterer{contract: contract}}, nil
}

// NewBethtestCaller creates a new read-only instance of Bethtest, bound to a specific deployed contract.
func NewBethtestCaller(address common.Address, caller bind.ContractCaller) (*BethtestCaller, error) {
	contract, err := bindBethtest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BethtestCaller{contract: contract}, nil
}

// NewBethtestTransactor creates a new write-only instance of Bethtest, bound to a specific deployed contract.
func NewBethtestTransactor(address common.Address, transactor bind.ContractTransactor) (*BethtestTransactor, error) {
	contract, err := bindBethtest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BethtestTransactor{contract: contract}, nil
}

// NewBethtestFilterer creates a new log filterer instance of Bethtest, bound to a specific deployed contract.
func NewBethtestFilterer(address common.Address, filterer bind.ContractFilterer) (*BethtestFilterer, error) {
	contract, err := bindBethtest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BethtestFilterer{contract: contract}, nil
}

// bindBethtest binds a generic wrapper to an already deployed contract.
func bindBethtest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BethtestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bethtest *BethtestRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bethtest.Contract.BethtestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bethtest *BethtestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bethtest.Contract.BethtestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bethtest *BethtestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bethtest.Contract.BethtestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bethtest *BethtestCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bethtest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bethtest *BethtestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bethtest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bethtest *BethtestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bethtest.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x9507d39a.
//
// Solidity: function get(data uint256) constant returns(uint256, bool)
func (_Bethtest *BethtestCaller) Get(opts *bind.CallOpts, data *big.Int) (*big.Int, bool, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(bool)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Bethtest.contract.Call(opts, out, "get", data)
	return *ret0, *ret1, err
}

// Get is a free data retrieval call binding the contract method 0x9507d39a.
//
// Solidity: function get(data uint256) constant returns(uint256, bool)
func (_Bethtest *BethtestSession) Get(data *big.Int) (*big.Int, bool, error) {
	return _Bethtest.Contract.Get(&_Bethtest.CallOpts, data)
}

// Get is a free data retrieval call binding the contract method 0x9507d39a.
//
// Solidity: function get(data uint256) constant returns(uint256, bool)
func (_Bethtest *BethtestCallerSession) Get(data *big.Int) (*big.Int, bool, error) {
	return _Bethtest.Contract.Get(&_Bethtest.CallOpts, data)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() constant returns(uint256)
func (_Bethtest *BethtestCaller) Read(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bethtest.contract.Call(opts, out, "read")
	return *ret0, err
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() constant returns(uint256)
func (_Bethtest *BethtestSession) Read() (*big.Int, error) {
	return _Bethtest.Contract.Read(&_Bethtest.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() constant returns(uint256)
func (_Bethtest *BethtestCallerSession) Read() (*big.Int, error) {
	return _Bethtest.Contract.Read(&_Bethtest.CallOpts)
}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() constant returns(uint256)
func (_Bethtest *BethtestCaller) Size(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bethtest.contract.Call(opts, out, "size")
	return *ret0, err
}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() constant returns(uint256)
func (_Bethtest *BethtestSession) Size() (*big.Int, error) {
	return _Bethtest.Contract.Size(&_Bethtest.CallOpts)
}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() constant returns(uint256)
func (_Bethtest *BethtestCallerSession) Size() (*big.Int, error) {
	return _Bethtest.Contract.Size(&_Bethtest.CallOpts)
}

// Append is a paid mutator transaction binding the contract method 0xe33b8707.
//
// Solidity: function append(x uint256) returns()
func (_Bethtest *BethtestTransactor) Append(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _Bethtest.contract.Transact(opts, "append", x)
}

// Append is a paid mutator transaction binding the contract method 0xe33b8707.
//
// Solidity: function append(x uint256) returns()
func (_Bethtest *BethtestSession) Append(x *big.Int) (*types.Transaction, error) {
	return _Bethtest.Contract.Append(&_Bethtest.TransactOpts, x)
}

// Append is a paid mutator transaction binding the contract method 0xe33b8707.
//
// Solidity: function append(x uint256) returns()
func (_Bethtest *BethtestTransactorSession) Append(x *big.Int) (*types.Transaction, error) {
	return _Bethtest.Contract.Append(&_Bethtest.TransactOpts, x)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Bethtest *BethtestTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bethtest.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Bethtest *BethtestSession) Increment() (*types.Transaction, error) {
	return _Bethtest.Contract.Increment(&_Bethtest.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Bethtest *BethtestTransactorSession) Increment() (*types.Transaction, error) {
	return _Bethtest.Contract.Increment(&_Bethtest.TransactOpts)
}

// Remove is a paid mutator transaction binding the contract method 0x4cc82215.
//
// Solidity: function remove(_dataToDelete uint256) returns()
func (_Bethtest *BethtestTransactor) Remove(opts *bind.TransactOpts, _dataToDelete *big.Int) (*types.Transaction, error) {
	return _Bethtest.contract.Transact(opts, "remove", _dataToDelete)
}

// Remove is a paid mutator transaction binding the contract method 0x4cc82215.
//
// Solidity: function remove(_dataToDelete uint256) returns()
func (_Bethtest *BethtestSession) Remove(_dataToDelete *big.Int) (*types.Transaction, error) {
	return _Bethtest.Contract.Remove(&_Bethtest.TransactOpts, _dataToDelete)
}

// Remove is a paid mutator transaction binding the contract method 0x4cc82215.
//
// Solidity: function remove(_dataToDelete uint256) returns()
func (_Bethtest *BethtestTransactorSession) Remove(_dataToDelete *big.Int) (*types.Transaction, error) {
	return _Bethtest.Contract.Remove(&_Bethtest.TransactOpts, _dataToDelete)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(x uint256) returns()
func (_Bethtest *BethtestTransactor) Set(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _Bethtest.contract.Transact(opts, "set", x)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(x uint256) returns()
func (_Bethtest *BethtestSession) Set(x *big.Int) (*types.Transaction, error) {
	return _Bethtest.Contract.Set(&_Bethtest.TransactOpts, x)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(x uint256) returns()
func (_Bethtest *BethtestTransactorSession) Set(x *big.Int) (*types.Transaction, error) {
	return _Bethtest.Contract.Set(&_Bethtest.TransactOpts, x)
}
