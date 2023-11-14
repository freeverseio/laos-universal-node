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

// EvolutionMetaData contains all meta data concerning the Evolution contract.
var EvolutionMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"collectionId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"EvolvedWithExternalURI\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"collectionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"slot\",\"type\":\"uint96\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"MintedWithExternalURI\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"collectionId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NewCollection\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"createCollection\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"collectionId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"evolveWithExternalURI\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"collectionId\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"slot\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mintWithExternalURI\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"collectionId\",\"type\":\"uint64\"}],\"name\":\"ownerOfCollection\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"collectionId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EvolutionABI is the input ABI used to generate the binding from.
// Deprecated: Use EvolutionMetaData.ABI instead.
var EvolutionABI = EvolutionMetaData.ABI

// Evolution is an auto generated Go binding around an Ethereum contract.
type Evolution struct {
	EvolutionCaller     // Read-only binding to the contract
	EvolutionTransactor // Write-only binding to the contract
	EvolutionFilterer   // Log filterer for contract events
}

// EvolutionCaller is an auto generated read-only Go binding around an Ethereum contract.
type EvolutionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EvolutionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EvolutionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EvolutionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EvolutionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EvolutionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EvolutionSession struct {
	Contract     *Evolution        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EvolutionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EvolutionCallerSession struct {
	Contract *EvolutionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// EvolutionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EvolutionTransactorSession struct {
	Contract     *EvolutionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// EvolutionRaw is an auto generated low-level Go binding around an Ethereum contract.
type EvolutionRaw struct {
	Contract *Evolution // Generic contract binding to access the raw methods on
}

// EvolutionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EvolutionCallerRaw struct {
	Contract *EvolutionCaller // Generic read-only contract binding to access the raw methods on
}

// EvolutionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EvolutionTransactorRaw struct {
	Contract *EvolutionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEvolution creates a new instance of Evolution, bound to a specific deployed contract.
func NewEvolution(address common.Address, backend bind.ContractBackend) (*Evolution, error) {
	contract, err := bindEvolution(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Evolution{EvolutionCaller: EvolutionCaller{contract: contract}, EvolutionTransactor: EvolutionTransactor{contract: contract}, EvolutionFilterer: EvolutionFilterer{contract: contract}}, nil
}

// NewEvolutionCaller creates a new read-only instance of Evolution, bound to a specific deployed contract.
func NewEvolutionCaller(address common.Address, caller bind.ContractCaller) (*EvolutionCaller, error) {
	contract, err := bindEvolution(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EvolutionCaller{contract: contract}, nil
}

// NewEvolutionTransactor creates a new write-only instance of Evolution, bound to a specific deployed contract.
func NewEvolutionTransactor(address common.Address, transactor bind.ContractTransactor) (*EvolutionTransactor, error) {
	contract, err := bindEvolution(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EvolutionTransactor{contract: contract}, nil
}

// NewEvolutionFilterer creates a new log filterer instance of Evolution, bound to a specific deployed contract.
func NewEvolutionFilterer(address common.Address, filterer bind.ContractFilterer) (*EvolutionFilterer, error) {
	contract, err := bindEvolution(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EvolutionFilterer{contract: contract}, nil
}

// bindEvolution binds a generic wrapper to an already deployed contract.
func bindEvolution(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EvolutionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Evolution *EvolutionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Evolution.Contract.EvolutionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Evolution *EvolutionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Evolution.Contract.EvolutionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Evolution *EvolutionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Evolution.Contract.EvolutionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Evolution *EvolutionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Evolution.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Evolution *EvolutionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Evolution.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Evolution *EvolutionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Evolution.Contract.contract.Transact(opts, method, params...)
}

// OwnerOfCollection is a free data retrieval call binding the contract method 0xfb34ae53.
//
// Solidity: function ownerOfCollection(uint64 collectionId) view returns(address)
func (_Evolution *EvolutionCaller) OwnerOfCollection(opts *bind.CallOpts, collectionId uint64) (common.Address, error) {
	var out []interface{}
	err := _Evolution.contract.Call(opts, &out, "ownerOfCollection", collectionId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOfCollection is a free data retrieval call binding the contract method 0xfb34ae53.
//
// Solidity: function ownerOfCollection(uint64 collectionId) view returns(address)
func (_Evolution *EvolutionSession) OwnerOfCollection(collectionId uint64) (common.Address, error) {
	return _Evolution.Contract.OwnerOfCollection(&_Evolution.CallOpts, collectionId)
}

// OwnerOfCollection is a free data retrieval call binding the contract method 0xfb34ae53.
//
// Solidity: function ownerOfCollection(uint64 collectionId) view returns(address)
func (_Evolution *EvolutionCallerSession) OwnerOfCollection(collectionId uint64) (common.Address, error) {
	return _Evolution.Contract.OwnerOfCollection(&_Evolution.CallOpts, collectionId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc8a3f102.
//
// Solidity: function tokenURI(uint64 collectionId, uint256 tokenId) view returns(string)
func (_Evolution *EvolutionCaller) TokenURI(opts *bind.CallOpts, collectionId uint64, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Evolution.contract.Call(opts, &out, "tokenURI", collectionId, tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc8a3f102.
//
// Solidity: function tokenURI(uint64 collectionId, uint256 tokenId) view returns(string)
func (_Evolution *EvolutionSession) TokenURI(collectionId uint64, tokenId *big.Int) (string, error) {
	return _Evolution.Contract.TokenURI(&_Evolution.CallOpts, collectionId, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc8a3f102.
//
// Solidity: function tokenURI(uint64 collectionId, uint256 tokenId) view returns(string)
func (_Evolution *EvolutionCallerSession) TokenURI(collectionId uint64, tokenId *big.Int) (string, error) {
	return _Evolution.Contract.TokenURI(&_Evolution.CallOpts, collectionId, tokenId)
}

// CreateCollection is a paid mutator transaction binding the contract method 0x2069e953.
//
// Solidity: function createCollection(address owner) returns(uint64)
func (_Evolution *EvolutionTransactor) CreateCollection(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _Evolution.contract.Transact(opts, "createCollection", owner)
}

// CreateCollection is a paid mutator transaction binding the contract method 0x2069e953.
//
// Solidity: function createCollection(address owner) returns(uint64)
func (_Evolution *EvolutionSession) CreateCollection(owner common.Address) (*types.Transaction, error) {
	return _Evolution.Contract.CreateCollection(&_Evolution.TransactOpts, owner)
}

// CreateCollection is a paid mutator transaction binding the contract method 0x2069e953.
//
// Solidity: function createCollection(address owner) returns(uint64)
func (_Evolution *EvolutionTransactorSession) CreateCollection(owner common.Address) (*types.Transaction, error) {
	return _Evolution.Contract.CreateCollection(&_Evolution.TransactOpts, owner)
}

// EvolveWithExternalURI is a paid mutator transaction binding the contract method 0x0ef2629f.
//
// Solidity: function evolveWithExternalURI(uint64 collectionId, uint256 tokenId, string tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactor) EvolveWithExternalURI(opts *bind.TransactOpts, collectionId uint64, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Evolution.contract.Transact(opts, "evolveWithExternalURI", collectionId, tokenId, tokenURI)
}

// EvolveWithExternalURI is a paid mutator transaction binding the contract method 0x0ef2629f.
//
// Solidity: function evolveWithExternalURI(uint64 collectionId, uint256 tokenId, string tokenURI) returns(uint256)
func (_Evolution *EvolutionSession) EvolveWithExternalURI(collectionId uint64, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.EvolveWithExternalURI(&_Evolution.TransactOpts, collectionId, tokenId, tokenURI)
}

// EvolveWithExternalURI is a paid mutator transaction binding the contract method 0x0ef2629f.
//
// Solidity: function evolveWithExternalURI(uint64 collectionId, uint256 tokenId, string tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactorSession) EvolveWithExternalURI(collectionId uint64, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.EvolveWithExternalURI(&_Evolution.TransactOpts, collectionId, tokenId, tokenURI)
}

// MintWithExternalURI is a paid mutator transaction binding the contract method 0xd4af5bbb.
//
// Solidity: function mintWithExternalURI(uint64 collectionId, uint96 slot, address to, string tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactor) MintWithExternalURI(opts *bind.TransactOpts, collectionId uint64, slot *big.Int, to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Evolution.contract.Transact(opts, "mintWithExternalURI", collectionId, slot, to, tokenURI)
}

// MintWithExternalURI is a paid mutator transaction binding the contract method 0xd4af5bbb.
//
// Solidity: function mintWithExternalURI(uint64 collectionId, uint96 slot, address to, string tokenURI) returns(uint256)
func (_Evolution *EvolutionSession) MintWithExternalURI(collectionId uint64, slot *big.Int, to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.MintWithExternalURI(&_Evolution.TransactOpts, collectionId, slot, to, tokenURI)
}

// MintWithExternalURI is a paid mutator transaction binding the contract method 0xd4af5bbb.
//
// Solidity: function mintWithExternalURI(uint64 collectionId, uint96 slot, address to, string tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactorSession) MintWithExternalURI(collectionId uint64, slot *big.Int, to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.MintWithExternalURI(&_Evolution.TransactOpts, collectionId, slot, to, tokenURI)
}

// EvolutionEvolvedWithExternalURIIterator is returned from FilterEvolvedWithExternalURI and is used to iterate over the raw logs and unpacked data for EvolvedWithExternalURI events raised by the Evolution contract.
type EvolutionEvolvedWithExternalURIIterator struct {
	Event *EvolutionEvolvedWithExternalURI // Event containing the contract specifics and raw log

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
func (it *EvolutionEvolvedWithExternalURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EvolutionEvolvedWithExternalURI)
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
		it.Event = new(EvolutionEvolvedWithExternalURI)
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
func (it *EvolutionEvolvedWithExternalURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EvolutionEvolvedWithExternalURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EvolutionEvolvedWithExternalURI represents a EvolvedWithExternalURI event raised by the Evolution contract.
type EvolutionEvolvedWithExternalURI struct {
	CollectionId uint64
	TokenId      *big.Int
	TokenURI     string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterEvolvedWithExternalURI is a free log retrieval operation binding the contract event 0x95c167d04a267f10e6b3f373c7a336dc65cf459caf048854dc32a2d37ab1607c.
//
// Solidity: event EvolvedWithExternalURI(uint64 collectionId, uint256 indexed tokenId, string tokenURI)
func (_Evolution *EvolutionFilterer) FilterEvolvedWithExternalURI(opts *bind.FilterOpts, tokenId []*big.Int) (*EvolutionEvolvedWithExternalURIIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Evolution.contract.FilterLogs(opts, "EvolvedWithExternalURI", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &EvolutionEvolvedWithExternalURIIterator{contract: _Evolution.contract, event: "EvolvedWithExternalURI", logs: logs, sub: sub}, nil
}

// WatchEvolvedWithExternalURI is a free log subscription operation binding the contract event 0x95c167d04a267f10e6b3f373c7a336dc65cf459caf048854dc32a2d37ab1607c.
//
// Solidity: event EvolvedWithExternalURI(uint64 collectionId, uint256 indexed tokenId, string tokenURI)
func (_Evolution *EvolutionFilterer) WatchEvolvedWithExternalURI(opts *bind.WatchOpts, sink chan<- *EvolutionEvolvedWithExternalURI, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Evolution.contract.WatchLogs(opts, "EvolvedWithExternalURI", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EvolutionEvolvedWithExternalURI)
				if err := _Evolution.contract.UnpackLog(event, "EvolvedWithExternalURI", log); err != nil {
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

// ParseEvolvedWithExternalURI is a log parse operation binding the contract event 0x95c167d04a267f10e6b3f373c7a336dc65cf459caf048854dc32a2d37ab1607c.
//
// Solidity: event EvolvedWithExternalURI(uint64 collectionId, uint256 indexed tokenId, string tokenURI)
func (_Evolution *EvolutionFilterer) ParseEvolvedWithExternalURI(log types.Log) (*EvolutionEvolvedWithExternalURI, error) {
	event := new(EvolutionEvolvedWithExternalURI)
	if err := _Evolution.contract.UnpackLog(event, "EvolvedWithExternalURI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EvolutionMintedWithExternalURIIterator is returned from FilterMintedWithExternalURI and is used to iterate over the raw logs and unpacked data for MintedWithExternalURI events raised by the Evolution contract.
type EvolutionMintedWithExternalURIIterator struct {
	Event *EvolutionMintedWithExternalURI // Event containing the contract specifics and raw log

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
func (it *EvolutionMintedWithExternalURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EvolutionMintedWithExternalURI)
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
		it.Event = new(EvolutionMintedWithExternalURI)
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
func (it *EvolutionMintedWithExternalURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EvolutionMintedWithExternalURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EvolutionMintedWithExternalURI represents a MintedWithExternalURI event raised by the Evolution contract.
type EvolutionMintedWithExternalURI struct {
	CollectionId uint64
	Slot         *big.Int
	To           common.Address
	TokenURI     string
	TokenId      *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMintedWithExternalURI is a free log retrieval operation binding the contract event 0x4b3b5da28a351f8bb73b960d7c80b2cef3e3570cb03448234dee173942c74786.
//
// Solidity: event MintedWithExternalURI(uint64 collectionId, uint96 slot, address indexed to, string tokenURI, uint256 tokenId)
func (_Evolution *EvolutionFilterer) FilterMintedWithExternalURI(opts *bind.FilterOpts, to []common.Address) (*EvolutionMintedWithExternalURIIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Evolution.contract.FilterLogs(opts, "MintedWithExternalURI", toRule)
	if err != nil {
		return nil, err
	}
	return &EvolutionMintedWithExternalURIIterator{contract: _Evolution.contract, event: "MintedWithExternalURI", logs: logs, sub: sub}, nil
}

// WatchMintedWithExternalURI is a free log subscription operation binding the contract event 0x4b3b5da28a351f8bb73b960d7c80b2cef3e3570cb03448234dee173942c74786.
//
// Solidity: event MintedWithExternalURI(uint64 collectionId, uint96 slot, address indexed to, string tokenURI, uint256 tokenId)
func (_Evolution *EvolutionFilterer) WatchMintedWithExternalURI(opts *bind.WatchOpts, sink chan<- *EvolutionMintedWithExternalURI, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Evolution.contract.WatchLogs(opts, "MintedWithExternalURI", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EvolutionMintedWithExternalURI)
				if err := _Evolution.contract.UnpackLog(event, "MintedWithExternalURI", log); err != nil {
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

// ParseMintedWithExternalURI is a log parse operation binding the contract event 0x4b3b5da28a351f8bb73b960d7c80b2cef3e3570cb03448234dee173942c74786.
//
// Solidity: event MintedWithExternalURI(uint64 collectionId, uint96 slot, address indexed to, string tokenURI, uint256 tokenId)
func (_Evolution *EvolutionFilterer) ParseMintedWithExternalURI(log types.Log) (*EvolutionMintedWithExternalURI, error) {
	event := new(EvolutionMintedWithExternalURI)
	if err := _Evolution.contract.UnpackLog(event, "MintedWithExternalURI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EvolutionNewCollectionIterator is returned from FilterNewCollection and is used to iterate over the raw logs and unpacked data for NewCollection events raised by the Evolution contract.
type EvolutionNewCollectionIterator struct {
	Event *EvolutionNewCollection // Event containing the contract specifics and raw log

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
func (it *EvolutionNewCollectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EvolutionNewCollection)
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
		it.Event = new(EvolutionNewCollection)
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
func (it *EvolutionNewCollectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EvolutionNewCollectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EvolutionNewCollection represents a NewCollection event raised by the Evolution contract.
type EvolutionNewCollection struct {
	CollectionId uint64
	Owner        common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterNewCollection is a free log retrieval operation binding the contract event 0x6eb24fd767a7bcfa417f3fe25a2cb245d2ae52293d3c4a8f8c6450a09795d289.
//
// Solidity: event NewCollection(uint64 collectionId, address indexed owner)
func (_Evolution *EvolutionFilterer) FilterNewCollection(opts *bind.FilterOpts, owner []common.Address) (*EvolutionNewCollectionIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Evolution.contract.FilterLogs(opts, "NewCollection", ownerRule)
	if err != nil {
		return nil, err
	}
	return &EvolutionNewCollectionIterator{contract: _Evolution.contract, event: "NewCollection", logs: logs, sub: sub}, nil
}

// WatchNewCollection is a free log subscription operation binding the contract event 0x6eb24fd767a7bcfa417f3fe25a2cb245d2ae52293d3c4a8f8c6450a09795d289.
//
// Solidity: event NewCollection(uint64 collectionId, address indexed owner)
func (_Evolution *EvolutionFilterer) WatchNewCollection(opts *bind.WatchOpts, sink chan<- *EvolutionNewCollection, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Evolution.contract.WatchLogs(opts, "NewCollection", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EvolutionNewCollection)
				if err := _Evolution.contract.UnpackLog(event, "NewCollection", log); err != nil {
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

// ParseNewCollection is a log parse operation binding the contract event 0x6eb24fd767a7bcfa417f3fe25a2cb245d2ae52293d3c4a8f8c6450a09795d289.
//
// Solidity: event NewCollection(uint64 collectionId, address indexed owner)
func (_Evolution *EvolutionFilterer) ParseNewCollection(log types.Log) (*EvolutionNewCollection, error) {
	event := new(EvolutionNewCollection)
	if err := _Evolution.contract.UnpackLog(event, "NewCollection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
