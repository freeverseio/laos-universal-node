// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// CollectionMetaData contains all meta data concerning the Collection contract.
var CollectionMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_collectionAddress\",\"type\":\"address\"}],\"name\":\"NewCollection\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"createCollection\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CollectionABI is the input ABI used to generate the binding from.
// Deprecated: Use CollectionMetaData.ABI instead.
var CollectionABI = CollectionMetaData.ABI

// Collection is an auto generated Go binding around an Ethereum contract.
type Collection struct {
	CollectionCaller     // Read-only binding to the contract
	CollectionTransactor // Write-only binding to the contract
	CollectionFilterer   // Log filterer for contract events
}

// CollectionCaller is an auto generated read-only Go binding around an Ethereum contract.
type CollectionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollectionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CollectionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollectionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CollectionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollectionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CollectionSession struct {
	Contract     *Collection       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CollectionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CollectionCallerSession struct {
	Contract *CollectionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CollectionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CollectionTransactorSession struct {
	Contract     *CollectionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CollectionRaw is an auto generated low-level Go binding around an Ethereum contract.
type CollectionRaw struct {
	Contract *Collection // Generic contract binding to access the raw methods on
}

// CollectionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CollectionCallerRaw struct {
	Contract *CollectionCaller // Generic read-only contract binding to access the raw methods on
}

// CollectionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CollectionTransactorRaw struct {
	Contract *CollectionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCollection creates a new instance of Collection, bound to a specific deployed contract.
func NewCollection(address common.Address, backend bind.ContractBackend) (*Collection, error) {
	contract, err := bindCollection(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Collection{CollectionCaller: CollectionCaller{contract: contract}, CollectionTransactor: CollectionTransactor{contract: contract}, CollectionFilterer: CollectionFilterer{contract: contract}}, nil
}

// NewCollectionCaller creates a new read-only instance of Collection, bound to a specific deployed contract.
func NewCollectionCaller(address common.Address, caller bind.ContractCaller) (*CollectionCaller, error) {
	contract, err := bindCollection(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CollectionCaller{contract: contract}, nil
}

// NewCollectionTransactor creates a new write-only instance of Collection, bound to a specific deployed contract.
func NewCollectionTransactor(address common.Address, transactor bind.ContractTransactor) (*CollectionTransactor, error) {
	contract, err := bindCollection(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CollectionTransactor{contract: contract}, nil
}

// NewCollectionFilterer creates a new log filterer instance of Collection, bound to a specific deployed contract.
func NewCollectionFilterer(address common.Address, filterer bind.ContractFilterer) (*CollectionFilterer, error) {
	contract, err := bindCollection(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CollectionFilterer{contract: contract}, nil
}

// bindCollection binds a generic wrapper to an already deployed contract.
func bindCollection(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CollectionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Collection *CollectionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Collection.Contract.CollectionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Collection *CollectionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Collection.Contract.CollectionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Collection *CollectionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Collection.Contract.CollectionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Collection *CollectionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Collection.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Collection *CollectionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Collection.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Collection *CollectionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Collection.Contract.contract.Transact(opts, method, params...)
}

// CreateCollection is a paid mutator transaction binding the contract method 0x2069e953.
//
// Solidity: function createCollection(address _owner) returns(address)
func (_Collection *CollectionTransactor) CreateCollection(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Collection.contract.Transact(opts, "createCollection", _owner)
}

// CreateCollection is a paid mutator transaction binding the contract method 0x2069e953.
//
// Solidity: function createCollection(address _owner) returns(address)
func (_Collection *CollectionSession) CreateCollection(_owner common.Address) (*types.Transaction, error) {
	return _Collection.Contract.CreateCollection(&_Collection.TransactOpts, _owner)
}

// CreateCollection is a paid mutator transaction binding the contract method 0x2069e953.
//
// Solidity: function createCollection(address _owner) returns(address)
func (_Collection *CollectionTransactorSession) CreateCollection(_owner common.Address) (*types.Transaction, error) {
	return _Collection.Contract.CreateCollection(&_Collection.TransactOpts, _owner)
}

// CollectionNewCollectionIterator is returned from FilterNewCollection and is used to iterate over the raw logs and unpacked data for NewCollection events raised by the Collection contract.
type CollectionNewCollectionIterator struct {
	Event *CollectionNewCollection // Event containing the contract specifics and raw log

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
func (it *CollectionNewCollectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollectionNewCollection)
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
		it.Event = new(CollectionNewCollection)
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
func (it *CollectionNewCollectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollectionNewCollectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollectionNewCollection represents a NewCollection event raised by the Collection contract.
type CollectionNewCollection struct {
	Owner             common.Address
	CollectionAddress common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterNewCollection is a free log retrieval operation binding the contract event 0x5b84d9550adb7000df7bee717735ecd3af48ea3f66c6886d52e8227548fb228c.
//
// Solidity: event NewCollection(address indexed _owner, address _collectionAddress)
func (_Collection *CollectionFilterer) FilterNewCollection(opts *bind.FilterOpts, _owner []common.Address) (*CollectionNewCollectionIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}

	logs, sub, err := _Collection.contract.FilterLogs(opts, "NewCollection", _ownerRule)
	if err != nil {
		return nil, err
	}
	return &CollectionNewCollectionIterator{contract: _Collection.contract, event: "NewCollection", logs: logs, sub: sub}, nil
}

// WatchNewCollection is a free log subscription operation binding the contract event 0x5b84d9550adb7000df7bee717735ecd3af48ea3f66c6886d52e8227548fb228c.
//
// Solidity: event NewCollection(address indexed _owner, address _collectionAddress)
func (_Collection *CollectionFilterer) WatchNewCollection(opts *bind.WatchOpts, sink chan<- *CollectionNewCollection, _owner []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}

	logs, sub, err := _Collection.contract.WatchLogs(opts, "NewCollection", _ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollectionNewCollection)
				if err := _Collection.contract.UnpackLog(event, "NewCollection", log); err != nil {
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

// ParseNewCollection is a log parse operation binding the contract event 0x5b84d9550adb7000df7bee717735ecd3af48ea3f66c6886d52e8227548fb228c.
//
// Solidity: event NewCollection(address indexed _owner, address _collectionAddress)
func (_Collection *CollectionFilterer) ParseNewCollection(log types.Log) (*CollectionNewCollection, error) {
	event := new(CollectionNewCollection)
	if err := _Collection.contract.UnpackLog(event, "NewCollection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
