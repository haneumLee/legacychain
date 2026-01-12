// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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
	_ = abi.ConvertType
)

// VaultFactoryMetaData contains all meta data concerning the VaultFactory contract.
var VaultFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createVault\",\"inputs\":[{\"name\":\"_heirs\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"_heartbeatInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_gracePeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_requiredApprovals\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOwnerVaultAt\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOwnerVaultCount\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOwnerVaults\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerVaults\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalVaults\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"vaultImplementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultCreated\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"vaultIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"FailedDeployment\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// VaultFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use VaultFactoryMetaData.ABI instead.
var VaultFactoryABI = VaultFactoryMetaData.ABI

// VaultFactory is an auto generated Go binding around an Ethereum contract.
type VaultFactory struct {
	VaultFactoryCaller     // Read-only binding to the contract
	VaultFactoryTransactor // Write-only binding to the contract
	VaultFactoryFilterer   // Log filterer for contract events
}

// VaultFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type VaultFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VaultFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VaultFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VaultFactorySession struct {
	Contract     *VaultFactory     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VaultFactoryCallerSession struct {
	Contract *VaultFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// VaultFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VaultFactoryTransactorSession struct {
	Contract     *VaultFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// VaultFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type VaultFactoryRaw struct {
	Contract *VaultFactory // Generic contract binding to access the raw methods on
}

// VaultFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VaultFactoryCallerRaw struct {
	Contract *VaultFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// VaultFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VaultFactoryTransactorRaw struct {
	Contract *VaultFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVaultFactory creates a new instance of VaultFactory, bound to a specific deployed contract.
func NewVaultFactory(address common.Address, backend bind.ContractBackend) (*VaultFactory, error) {
	contract, err := bindVaultFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VaultFactory{VaultFactoryCaller: VaultFactoryCaller{contract: contract}, VaultFactoryTransactor: VaultFactoryTransactor{contract: contract}, VaultFactoryFilterer: VaultFactoryFilterer{contract: contract}}, nil
}

// NewVaultFactoryCaller creates a new read-only instance of VaultFactory, bound to a specific deployed contract.
func NewVaultFactoryCaller(address common.Address, caller bind.ContractCaller) (*VaultFactoryCaller, error) {
	contract, err := bindVaultFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VaultFactoryCaller{contract: contract}, nil
}

// NewVaultFactoryTransactor creates a new write-only instance of VaultFactory, bound to a specific deployed contract.
func NewVaultFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*VaultFactoryTransactor, error) {
	contract, err := bindVaultFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VaultFactoryTransactor{contract: contract}, nil
}

// NewVaultFactoryFilterer creates a new log filterer instance of VaultFactory, bound to a specific deployed contract.
func NewVaultFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*VaultFactoryFilterer, error) {
	contract, err := bindVaultFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VaultFactoryFilterer{contract: contract}, nil
}

// bindVaultFactory binds a generic wrapper to an already deployed contract.
func bindVaultFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VaultFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VaultFactory *VaultFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VaultFactory.Contract.VaultFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VaultFactory *VaultFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VaultFactory.Contract.VaultFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VaultFactory *VaultFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VaultFactory.Contract.VaultFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VaultFactory *VaultFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VaultFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VaultFactory *VaultFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VaultFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VaultFactory *VaultFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VaultFactory.Contract.contract.Transact(opts, method, params...)
}

// GetOwnerVaultAt is a free data retrieval call binding the contract method 0x2e8b66ce.
//
// Solidity: function getOwnerVaultAt(address _owner, uint256 _index) view returns(address)
func (_VaultFactory *VaultFactoryCaller) GetOwnerVaultAt(opts *bind.CallOpts, _owner common.Address, _index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VaultFactory.contract.Call(opts, &out, "getOwnerVaultAt", _owner, _index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwnerVaultAt is a free data retrieval call binding the contract method 0x2e8b66ce.
//
// Solidity: function getOwnerVaultAt(address _owner, uint256 _index) view returns(address)
func (_VaultFactory *VaultFactorySession) GetOwnerVaultAt(_owner common.Address, _index *big.Int) (common.Address, error) {
	return _VaultFactory.Contract.GetOwnerVaultAt(&_VaultFactory.CallOpts, _owner, _index)
}

// GetOwnerVaultAt is a free data retrieval call binding the contract method 0x2e8b66ce.
//
// Solidity: function getOwnerVaultAt(address _owner, uint256 _index) view returns(address)
func (_VaultFactory *VaultFactoryCallerSession) GetOwnerVaultAt(_owner common.Address, _index *big.Int) (common.Address, error) {
	return _VaultFactory.Contract.GetOwnerVaultAt(&_VaultFactory.CallOpts, _owner, _index)
}

// GetOwnerVaultCount is a free data retrieval call binding the contract method 0xaaa76043.
//
// Solidity: function getOwnerVaultCount(address _owner) view returns(uint256)
func (_VaultFactory *VaultFactoryCaller) GetOwnerVaultCount(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VaultFactory.contract.Call(opts, &out, "getOwnerVaultCount", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOwnerVaultCount is a free data retrieval call binding the contract method 0xaaa76043.
//
// Solidity: function getOwnerVaultCount(address _owner) view returns(uint256)
func (_VaultFactory *VaultFactorySession) GetOwnerVaultCount(_owner common.Address) (*big.Int, error) {
	return _VaultFactory.Contract.GetOwnerVaultCount(&_VaultFactory.CallOpts, _owner)
}

// GetOwnerVaultCount is a free data retrieval call binding the contract method 0xaaa76043.
//
// Solidity: function getOwnerVaultCount(address _owner) view returns(uint256)
func (_VaultFactory *VaultFactoryCallerSession) GetOwnerVaultCount(_owner common.Address) (*big.Int, error) {
	return _VaultFactory.Contract.GetOwnerVaultCount(&_VaultFactory.CallOpts, _owner)
}

// GetOwnerVaults is a free data retrieval call binding the contract method 0x3a0f555d.
//
// Solidity: function getOwnerVaults(address _owner) view returns(address[])
func (_VaultFactory *VaultFactoryCaller) GetOwnerVaults(opts *bind.CallOpts, _owner common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _VaultFactory.contract.Call(opts, &out, "getOwnerVaults", _owner)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwnerVaults is a free data retrieval call binding the contract method 0x3a0f555d.
//
// Solidity: function getOwnerVaults(address _owner) view returns(address[])
func (_VaultFactory *VaultFactorySession) GetOwnerVaults(_owner common.Address) ([]common.Address, error) {
	return _VaultFactory.Contract.GetOwnerVaults(&_VaultFactory.CallOpts, _owner)
}

// GetOwnerVaults is a free data retrieval call binding the contract method 0x3a0f555d.
//
// Solidity: function getOwnerVaults(address _owner) view returns(address[])
func (_VaultFactory *VaultFactoryCallerSession) GetOwnerVaults(_owner common.Address) ([]common.Address, error) {
	return _VaultFactory.Contract.GetOwnerVaults(&_VaultFactory.CallOpts, _owner)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VaultFactory *VaultFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VaultFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VaultFactory *VaultFactorySession) Owner() (common.Address, error) {
	return _VaultFactory.Contract.Owner(&_VaultFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_VaultFactory *VaultFactoryCallerSession) Owner() (common.Address, error) {
	return _VaultFactory.Contract.Owner(&_VaultFactory.CallOpts)
}

// OwnerVaults is a free data retrieval call binding the contract method 0x0c90d7b5.
//
// Solidity: function ownerVaults(address , uint256 ) view returns(address)
func (_VaultFactory *VaultFactoryCaller) OwnerVaults(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VaultFactory.contract.Call(opts, &out, "ownerVaults", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerVaults is a free data retrieval call binding the contract method 0x0c90d7b5.
//
// Solidity: function ownerVaults(address , uint256 ) view returns(address)
func (_VaultFactory *VaultFactorySession) OwnerVaults(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _VaultFactory.Contract.OwnerVaults(&_VaultFactory.CallOpts, arg0, arg1)
}

// OwnerVaults is a free data retrieval call binding the contract method 0x0c90d7b5.
//
// Solidity: function ownerVaults(address , uint256 ) view returns(address)
func (_VaultFactory *VaultFactoryCallerSession) OwnerVaults(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _VaultFactory.Contract.OwnerVaults(&_VaultFactory.CallOpts, arg0, arg1)
}

// TotalVaults is a free data retrieval call binding the contract method 0x8d654023.
//
// Solidity: function totalVaults() view returns(uint256)
func (_VaultFactory *VaultFactoryCaller) TotalVaults(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VaultFactory.contract.Call(opts, &out, "totalVaults")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVaults is a free data retrieval call binding the contract method 0x8d654023.
//
// Solidity: function totalVaults() view returns(uint256)
func (_VaultFactory *VaultFactorySession) TotalVaults() (*big.Int, error) {
	return _VaultFactory.Contract.TotalVaults(&_VaultFactory.CallOpts)
}

// TotalVaults is a free data retrieval call binding the contract method 0x8d654023.
//
// Solidity: function totalVaults() view returns(uint256)
func (_VaultFactory *VaultFactoryCallerSession) TotalVaults() (*big.Int, error) {
	return _VaultFactory.Contract.TotalVaults(&_VaultFactory.CallOpts)
}

// VaultImplementation is a free data retrieval call binding the contract method 0xbba48a90.
//
// Solidity: function vaultImplementation() view returns(address)
func (_VaultFactory *VaultFactoryCaller) VaultImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VaultFactory.contract.Call(opts, &out, "vaultImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VaultImplementation is a free data retrieval call binding the contract method 0xbba48a90.
//
// Solidity: function vaultImplementation() view returns(address)
func (_VaultFactory *VaultFactorySession) VaultImplementation() (common.Address, error) {
	return _VaultFactory.Contract.VaultImplementation(&_VaultFactory.CallOpts)
}

// VaultImplementation is a free data retrieval call binding the contract method 0xbba48a90.
//
// Solidity: function vaultImplementation() view returns(address)
func (_VaultFactory *VaultFactoryCallerSession) VaultImplementation() (common.Address, error) {
	return _VaultFactory.Contract.VaultImplementation(&_VaultFactory.CallOpts)
}

// CreateVault is a paid mutator transaction binding the contract method 0x8b927014.
//
// Solidity: function createVault(address[] _heirs, uint256[] _shares, uint256 _heartbeatInterval, uint256 _gracePeriod, uint256 _requiredApprovals) returns(address vaultAddress)
func (_VaultFactory *VaultFactoryTransactor) CreateVault(opts *bind.TransactOpts, _heirs []common.Address, _shares []*big.Int, _heartbeatInterval *big.Int, _gracePeriod *big.Int, _requiredApprovals *big.Int) (*types.Transaction, error) {
	return _VaultFactory.contract.Transact(opts, "createVault", _heirs, _shares, _heartbeatInterval, _gracePeriod, _requiredApprovals)
}

// CreateVault is a paid mutator transaction binding the contract method 0x8b927014.
//
// Solidity: function createVault(address[] _heirs, uint256[] _shares, uint256 _heartbeatInterval, uint256 _gracePeriod, uint256 _requiredApprovals) returns(address vaultAddress)
func (_VaultFactory *VaultFactorySession) CreateVault(_heirs []common.Address, _shares []*big.Int, _heartbeatInterval *big.Int, _gracePeriod *big.Int, _requiredApprovals *big.Int) (*types.Transaction, error) {
	return _VaultFactory.Contract.CreateVault(&_VaultFactory.TransactOpts, _heirs, _shares, _heartbeatInterval, _gracePeriod, _requiredApprovals)
}

// CreateVault is a paid mutator transaction binding the contract method 0x8b927014.
//
// Solidity: function createVault(address[] _heirs, uint256[] _shares, uint256 _heartbeatInterval, uint256 _gracePeriod, uint256 _requiredApprovals) returns(address vaultAddress)
func (_VaultFactory *VaultFactoryTransactorSession) CreateVault(_heirs []common.Address, _shares []*big.Int, _heartbeatInterval *big.Int, _gracePeriod *big.Int, _requiredApprovals *big.Int) (*types.Transaction, error) {
	return _VaultFactory.Contract.CreateVault(&_VaultFactory.TransactOpts, _heirs, _shares, _heartbeatInterval, _gracePeriod, _requiredApprovals)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_VaultFactory *VaultFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VaultFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_VaultFactory *VaultFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _VaultFactory.Contract.RenounceOwnership(&_VaultFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_VaultFactory *VaultFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _VaultFactory.Contract.RenounceOwnership(&_VaultFactory.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VaultFactory *VaultFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VaultFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VaultFactory *VaultFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VaultFactory.Contract.TransferOwnership(&_VaultFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_VaultFactory *VaultFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VaultFactory.Contract.TransferOwnership(&_VaultFactory.TransactOpts, newOwner)
}

// VaultFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the VaultFactory contract.
type VaultFactoryOwnershipTransferredIterator struct {
	Event *VaultFactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VaultFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultFactoryOwnershipTransferred)
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
		it.Event = new(VaultFactoryOwnershipTransferred)
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
func (it *VaultFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the VaultFactory contract.
type VaultFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VaultFactory *VaultFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VaultFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VaultFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VaultFactoryOwnershipTransferredIterator{contract: _VaultFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VaultFactory *VaultFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VaultFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _VaultFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultFactoryOwnershipTransferred)
				if err := _VaultFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_VaultFactory *VaultFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*VaultFactoryOwnershipTransferred, error) {
	event := new(VaultFactoryOwnershipTransferred)
	if err := _VaultFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultFactoryVaultCreatedIterator is returned from FilterVaultCreated and is used to iterate over the raw logs and unpacked data for VaultCreated events raised by the VaultFactory contract.
type VaultFactoryVaultCreatedIterator struct {
	Event *VaultFactoryVaultCreated // Event containing the contract specifics and raw log

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
func (it *VaultFactoryVaultCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultFactoryVaultCreated)
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
		it.Event = new(VaultFactoryVaultCreated)
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
func (it *VaultFactoryVaultCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultFactoryVaultCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultFactoryVaultCreated represents a VaultCreated event raised by the VaultFactory contract.
type VaultFactoryVaultCreated struct {
	VaultAddress common.Address
	Owner        common.Address
	VaultIndex   *big.Int
	Timestamp    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterVaultCreated is a free log retrieval operation binding the contract event 0xa28dca8eb29fc00aa90ca34a94ceb2d4eef6c73d46ee64d3b3834ff1fb6d1668.
//
// Solidity: event VaultCreated(address indexed vaultAddress, address indexed owner, uint256 vaultIndex, uint256 timestamp)
func (_VaultFactory *VaultFactoryFilterer) FilterVaultCreated(opts *bind.FilterOpts, vaultAddress []common.Address, owner []common.Address) (*VaultFactoryVaultCreatedIterator, error) {

	var vaultAddressRule []interface{}
	for _, vaultAddressItem := range vaultAddress {
		vaultAddressRule = append(vaultAddressRule, vaultAddressItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _VaultFactory.contract.FilterLogs(opts, "VaultCreated", vaultAddressRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &VaultFactoryVaultCreatedIterator{contract: _VaultFactory.contract, event: "VaultCreated", logs: logs, sub: sub}, nil
}

// WatchVaultCreated is a free log subscription operation binding the contract event 0xa28dca8eb29fc00aa90ca34a94ceb2d4eef6c73d46ee64d3b3834ff1fb6d1668.
//
// Solidity: event VaultCreated(address indexed vaultAddress, address indexed owner, uint256 vaultIndex, uint256 timestamp)
func (_VaultFactory *VaultFactoryFilterer) WatchVaultCreated(opts *bind.WatchOpts, sink chan<- *VaultFactoryVaultCreated, vaultAddress []common.Address, owner []common.Address) (event.Subscription, error) {

	var vaultAddressRule []interface{}
	for _, vaultAddressItem := range vaultAddress {
		vaultAddressRule = append(vaultAddressRule, vaultAddressItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _VaultFactory.contract.WatchLogs(opts, "VaultCreated", vaultAddressRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultFactoryVaultCreated)
				if err := _VaultFactory.contract.UnpackLog(event, "VaultCreated", log); err != nil {
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

// ParseVaultCreated is a log parse operation binding the contract event 0xa28dca8eb29fc00aa90ca34a94ceb2d4eef6c73d46ee64d3b3834ff1fb6d1668.
//
// Solidity: event VaultCreated(address indexed vaultAddress, address indexed owner, uint256 vaultIndex, uint256 timestamp)
func (_VaultFactory *VaultFactoryFilterer) ParseVaultCreated(log types.Log) (*VaultFactoryVaultCreated, error) {
	event := new(VaultFactoryVaultCreated)
	if err := _VaultFactory.contract.UnpackLog(event, "VaultCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
