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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"}],\"name\":\"EvolvedWithExternalURI\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"_slot\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"}],\"name\":\"MintedWithExternalURI\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"}],\"name\":\"evolveWithExternalURI\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"_slot\",\"type\":\"uint96\"},{\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"}],\"name\":\"mintWithExternalURI\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Evolution *EvolutionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Evolution.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Evolution *EvolutionSession) Owner() (common.Address, error) {
	return _Evolution.Contract.Owner(&_Evolution.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Evolution *EvolutionCallerSession) Owner() (common.Address, error) {
	return _Evolution.Contract.Owner(&_Evolution.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_Evolution *EvolutionCaller) TokenURI(opts *bind.CallOpts, _tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Evolution.contract.Call(opts, &out, "tokenURI", _tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_Evolution *EvolutionSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _Evolution.Contract.TokenURI(&_Evolution.CallOpts, _tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_Evolution *EvolutionCallerSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _Evolution.Contract.TokenURI(&_Evolution.CallOpts, _tokenId)
}

// EvolveWithExternalURI is a paid mutator transaction binding the contract method 0x2fd38f4d.
//
// Solidity: function evolveWithExternalURI(uint256 _tokenId, string _tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactor) EvolveWithExternalURI(opts *bind.TransactOpts, _tokenId *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Evolution.contract.Transact(opts, "evolveWithExternalURI", _tokenId, _tokenURI)
}

// EvolveWithExternalURI is a paid mutator transaction binding the contract method 0x2fd38f4d.
//
// Solidity: function evolveWithExternalURI(uint256 _tokenId, string _tokenURI) returns(uint256)
func (_Evolution *EvolutionSession) EvolveWithExternalURI(_tokenId *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.EvolveWithExternalURI(&_Evolution.TransactOpts, _tokenId, _tokenURI)
}

// EvolveWithExternalURI is a paid mutator transaction binding the contract method 0x2fd38f4d.
//
// Solidity: function evolveWithExternalURI(uint256 _tokenId, string _tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactorSession) EvolveWithExternalURI(_tokenId *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.EvolveWithExternalURI(&_Evolution.TransactOpts, _tokenId, _tokenURI)
}

// MintWithExternalURI is a paid mutator transaction binding the contract method 0xfd024566.
//
// Solidity: function mintWithExternalURI(address _to, uint96 _slot, string _tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactor) MintWithExternalURI(opts *bind.TransactOpts, _to common.Address, _slot *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Evolution.contract.Transact(opts, "mintWithExternalURI", _to, _slot, _tokenURI)
}

// MintWithExternalURI is a paid mutator transaction binding the contract method 0xfd024566.
//
// Solidity: function mintWithExternalURI(address _to, uint96 _slot, string _tokenURI) returns(uint256)
func (_Evolution *EvolutionSession) MintWithExternalURI(_to common.Address, _slot *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.MintWithExternalURI(&_Evolution.TransactOpts, _to, _slot, _tokenURI)
}

// MintWithExternalURI is a paid mutator transaction binding the contract method 0xfd024566.
//
// Solidity: function mintWithExternalURI(address _to, uint96 _slot, string _tokenURI) returns(uint256)
func (_Evolution *EvolutionTransactorSession) MintWithExternalURI(_to common.Address, _slot *big.Int, _tokenURI string) (*types.Transaction, error) {
	return _Evolution.Contract.MintWithExternalURI(&_Evolution.TransactOpts, _to, _slot, _tokenURI)
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
	TokenId  *big.Int
	TokenURI string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterEvolvedWithExternalURI is a free log retrieval operation binding the contract event 0xdde18ad2fe10c12a694de65b920c02b851c382cf63115967ea6f7098902fa1c8.
//
// Solidity: event EvolvedWithExternalURI(uint256 indexed _tokenId, string _tokenURI)
func (_Evolution *EvolutionFilterer) FilterEvolvedWithExternalURI(opts *bind.FilterOpts, _tokenId []*big.Int) (*EvolutionEvolvedWithExternalURIIterator, error) {

	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Evolution.contract.FilterLogs(opts, "EvolvedWithExternalURI", _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &EvolutionEvolvedWithExternalURIIterator{contract: _Evolution.contract, event: "EvolvedWithExternalURI", logs: logs, sub: sub}, nil
}

// WatchEvolvedWithExternalURI is a free log subscription operation binding the contract event 0xdde18ad2fe10c12a694de65b920c02b851c382cf63115967ea6f7098902fa1c8.
//
// Solidity: event EvolvedWithExternalURI(uint256 indexed _tokenId, string _tokenURI)
func (_Evolution *EvolutionFilterer) WatchEvolvedWithExternalURI(opts *bind.WatchOpts, sink chan<- *EvolutionEvolvedWithExternalURI, _tokenId []*big.Int) (event.Subscription, error) {

	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Evolution.contract.WatchLogs(opts, "EvolvedWithExternalURI", _tokenIdRule)
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

// ParseEvolvedWithExternalURI is a log parse operation binding the contract event 0xdde18ad2fe10c12a694de65b920c02b851c382cf63115967ea6f7098902fa1c8.
//
// Solidity: event EvolvedWithExternalURI(uint256 indexed _tokenId, string _tokenURI)
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
	To       common.Address
	Slot     *big.Int
	TokenId  *big.Int
	TokenURI string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMintedWithExternalURI is a free log retrieval operation binding the contract event 0xa7135052b348b0b4e9943bae82d8ef1c5ac225e594ef4271d12f0744cfc98348.
//
// Solidity: event MintedWithExternalURI(address indexed _to, uint96 _slot, uint256 _tokenId, string _tokenURI)
func (_Evolution *EvolutionFilterer) FilterMintedWithExternalURI(opts *bind.FilterOpts, _to []common.Address) (*EvolutionMintedWithExternalURIIterator, error) {

	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Evolution.contract.FilterLogs(opts, "MintedWithExternalURI", _toRule)
	if err != nil {
		return nil, err
	}
	return &EvolutionMintedWithExternalURIIterator{contract: _Evolution.contract, event: "MintedWithExternalURI", logs: logs, sub: sub}, nil
}

// WatchMintedWithExternalURI is a free log subscription operation binding the contract event 0xa7135052b348b0b4e9943bae82d8ef1c5ac225e594ef4271d12f0744cfc98348.
//
// Solidity: event MintedWithExternalURI(address indexed _to, uint96 _slot, uint256 _tokenId, string _tokenURI)
func (_Evolution *EvolutionFilterer) WatchMintedWithExternalURI(opts *bind.WatchOpts, sink chan<- *EvolutionMintedWithExternalURI, _to []common.Address) (event.Subscription, error) {

	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Evolution.contract.WatchLogs(opts, "MintedWithExternalURI", _toRule)
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

// ParseMintedWithExternalURI is a log parse operation binding the contract event 0xa7135052b348b0b4e9943bae82d8ef1c5ac225e594ef4271d12f0744cfc98348.
//
// Solidity: event MintedWithExternalURI(address indexed _to, uint96 _slot, uint256 _tokenId, string _tokenURI)
func (_Evolution *EvolutionFilterer) ParseMintedWithExternalURI(log types.Log) (*EvolutionMintedWithExternalURI, error) {
	event := new(EvolutionMintedWithExternalURI)
	if err := _Evolution.contract.UnpackLog(event, "MintedWithExternalURI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
