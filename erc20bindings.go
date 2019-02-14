// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package beth

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// CompatibleERC20ABI is the input ABI used to generate the binding from.
const CompatibleERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// CompatibleERC20Bin is the compiled bytecode used for deploying new contracts.
const CompatibleERC20Bin = `0x`

// DeployCompatibleERC20 deploys a new Ethereum contract, binding an instance of CompatibleERC20 to it.
func DeployCompatibleERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CompatibleERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(CompatibleERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CompatibleERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CompatibleERC20{CompatibleERC20Caller: CompatibleERC20Caller{contract: contract}, CompatibleERC20Transactor: CompatibleERC20Transactor{contract: contract}, CompatibleERC20Filterer: CompatibleERC20Filterer{contract: contract}}, nil
}

// CompatibleERC20 is an auto generated Go binding around an Ethereum contract.
type CompatibleERC20 struct {
	CompatibleERC20Caller     // Read-only binding to the contract
	CompatibleERC20Transactor // Write-only binding to the contract
	CompatibleERC20Filterer   // Log filterer for contract events
}

// CompatibleERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type CompatibleERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type CompatibleERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CompatibleERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CompatibleERC20Session struct {
	Contract     *CompatibleERC20  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CompatibleERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CompatibleERC20CallerSession struct {
	Contract *CompatibleERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CompatibleERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CompatibleERC20TransactorSession struct {
	Contract     *CompatibleERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CompatibleERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type CompatibleERC20Raw struct {
	Contract *CompatibleERC20 // Generic contract binding to access the raw methods on
}

// CompatibleERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CompatibleERC20CallerRaw struct {
	Contract *CompatibleERC20Caller // Generic read-only contract binding to access the raw methods on
}

// CompatibleERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CompatibleERC20TransactorRaw struct {
	Contract *CompatibleERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewCompatibleERC20 creates a new instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20(address common.Address, backend bind.ContractBackend) (*CompatibleERC20, error) {
	contract, err := bindCompatibleERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20{CompatibleERC20Caller: CompatibleERC20Caller{contract: contract}, CompatibleERC20Transactor: CompatibleERC20Transactor{contract: contract}, CompatibleERC20Filterer: CompatibleERC20Filterer{contract: contract}}, nil
}

// NewCompatibleERC20Caller creates a new read-only instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20Caller(address common.Address, caller bind.ContractCaller) (*CompatibleERC20Caller, error) {
	contract, err := bindCompatibleERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20Caller{contract: contract}, nil
}

// NewCompatibleERC20Transactor creates a new write-only instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*CompatibleERC20Transactor, error) {
	contract, err := bindCompatibleERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20Transactor{contract: contract}, nil
}

// NewCompatibleERC20Filterer creates a new log filterer instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*CompatibleERC20Filterer, error) {
	contract, err := bindCompatibleERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20Filterer{contract: contract}, nil
}

// bindCompatibleERC20 binds a generic wrapper to an already deployed contract.
func bindCompatibleERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CompatibleERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompatibleERC20 *CompatibleERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CompatibleERC20.Contract.CompatibleERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompatibleERC20 *CompatibleERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.CompatibleERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompatibleERC20 *CompatibleERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.CompatibleERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompatibleERC20 *CompatibleERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CompatibleERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompatibleERC20 *CompatibleERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompatibleERC20 *CompatibleERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CompatibleERC20.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.Allowance(&_CompatibleERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.Allowance(&_CompatibleERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Caller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CompatibleERC20.contract.Call(opts, out, "balanceOf", who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Session) BalanceOf(who common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.BalanceOf(&_CompatibleERC20.CallOpts, who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20CallerSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.BalanceOf(&_CompatibleERC20.CallOpts, who)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CompatibleERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Session) TotalSupply() (*big.Int, error) {
	return _CompatibleERC20.Contract.TotalSupply(&_CompatibleERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _CompatibleERC20.Contract.TotalSupply(&_CompatibleERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Approve(&_CompatibleERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Approve(&_CompatibleERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Transfer(&_CompatibleERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Transfer(&_CompatibleERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.TransferFrom(&_CompatibleERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.TransferFrom(&_CompatibleERC20.TransactOpts, from, to, value)
}

// CompatibleERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the CompatibleERC20 contract.
type CompatibleERC20ApprovalIterator struct {
	Event *CompatibleERC20Approval // Event containing the contract specifics and raw log

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
func (it *CompatibleERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompatibleERC20Approval)
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
		it.Event = new(CompatibleERC20Approval)
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
func (it *CompatibleERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompatibleERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompatibleERC20Approval represents a Approval event raised by the CompatibleERC20 contract.
type CompatibleERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*CompatibleERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CompatibleERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20ApprovalIterator{contract: _CompatibleERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CompatibleERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CompatibleERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompatibleERC20Approval)
				if err := _CompatibleERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// CompatibleERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the CompatibleERC20 contract.
type CompatibleERC20TransferIterator struct {
	Event *CompatibleERC20Transfer // Event containing the contract specifics and raw log

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
func (it *CompatibleERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompatibleERC20Transfer)
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
		it.Event = new(CompatibleERC20Transfer)
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
func (it *CompatibleERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompatibleERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompatibleERC20Transfer represents a Transfer event raised by the CompatibleERC20 contract.
type CompatibleERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CompatibleERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CompatibleERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20TransferIterator{contract: _CompatibleERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CompatibleERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CompatibleERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompatibleERC20Transfer)
				if err := _CompatibleERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
