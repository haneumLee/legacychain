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

// IndividualVaultVaultConfig is an auto generated low-level Go binding around an user-defined struct.
type IndividualVaultVaultConfig struct {
	Owner                common.Address
	Heirs                []common.Address
	HeirShares           []*big.Int
	HeartbeatInterval    *big.Int
	LastHeartbeat        *big.Int
	UnlockTime           *big.Int
	GracePeriod          *big.Int
	RequiredApprovals    *big.Int
	ApprovalCount        *big.Int
	TotalBalanceAtUnlock *big.Int
	IsLocked             bool
	GracePeriodActive    bool
}

// IndividualVaultMetaData contains all meta data concerning the IndividualVault contract.
var IndividualVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"approveInheritance\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkAndUnlock\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimInheritance\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitHeartbeat\",\"inputs\":[{\"name\":\"_commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"config\",\"inputs\":[],\"outputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"heartbeatInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastHeartbeat\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unlockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gracePeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requiredApprovals\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"approvalCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBalanceAtUnlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isLocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"gracePeriodActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIndividualVault.VaultConfig\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"heirs\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"heirShares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"heartbeatInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastHeartbeat\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unlockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gracePeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requiredApprovals\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"approvalCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBalanceAtUnlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isLocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"gracePeriodActive\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHeir\",\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHeirCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHeirShare\",\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"heirApprovals\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"heirClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_heirs\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"_heartbeatInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_gracePeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_requiredApprovals\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isClaimable\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isHeir\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revealHeartbeat\",\"inputs\":[{\"name\":\"_nonce\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EmergencyPaused\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GracePeriodStarted\",\"inputs\":[{\"name\":\"endTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Heartbeat\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InheritanceApproved\",\"inputs\":[{\"name\":\"heir\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InheritanceClaimed\",\"inputs\":[{\"name\":\"heir\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnlockCancelled\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultUnlocked\",\"inputs\":[{\"name\":\"unlockTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
}

// IndividualVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use IndividualVaultMetaData.ABI instead.
var IndividualVaultABI = IndividualVaultMetaData.ABI

// IndividualVault is an auto generated Go binding around an Ethereum contract.
type IndividualVault struct {
	IndividualVaultCaller     // Read-only binding to the contract
	IndividualVaultTransactor // Write-only binding to the contract
	IndividualVaultFilterer   // Log filterer for contract events
}

// IndividualVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type IndividualVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IndividualVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IndividualVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IndividualVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IndividualVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IndividualVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IndividualVaultSession struct {
	Contract     *IndividualVault  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IndividualVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IndividualVaultCallerSession struct {
	Contract *IndividualVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IndividualVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IndividualVaultTransactorSession struct {
	Contract     *IndividualVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IndividualVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type IndividualVaultRaw struct {
	Contract *IndividualVault // Generic contract binding to access the raw methods on
}

// IndividualVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IndividualVaultCallerRaw struct {
	Contract *IndividualVaultCaller // Generic read-only contract binding to access the raw methods on
}

// IndividualVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IndividualVaultTransactorRaw struct {
	Contract *IndividualVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIndividualVault creates a new instance of IndividualVault, bound to a specific deployed contract.
func NewIndividualVault(address common.Address, backend bind.ContractBackend) (*IndividualVault, error) {
	contract, err := bindIndividualVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IndividualVault{IndividualVaultCaller: IndividualVaultCaller{contract: contract}, IndividualVaultTransactor: IndividualVaultTransactor{contract: contract}, IndividualVaultFilterer: IndividualVaultFilterer{contract: contract}}, nil
}

// NewIndividualVaultCaller creates a new read-only instance of IndividualVault, bound to a specific deployed contract.
func NewIndividualVaultCaller(address common.Address, caller bind.ContractCaller) (*IndividualVaultCaller, error) {
	contract, err := bindIndividualVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultCaller{contract: contract}, nil
}

// NewIndividualVaultTransactor creates a new write-only instance of IndividualVault, bound to a specific deployed contract.
func NewIndividualVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*IndividualVaultTransactor, error) {
	contract, err := bindIndividualVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultTransactor{contract: contract}, nil
}

// NewIndividualVaultFilterer creates a new log filterer instance of IndividualVault, bound to a specific deployed contract.
func NewIndividualVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*IndividualVaultFilterer, error) {
	contract, err := bindIndividualVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultFilterer{contract: contract}, nil
}

// bindIndividualVault binds a generic wrapper to an already deployed contract.
func bindIndividualVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IndividualVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IndividualVault *IndividualVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IndividualVault.Contract.IndividualVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IndividualVault *IndividualVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.Contract.IndividualVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IndividualVault *IndividualVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IndividualVault.Contract.IndividualVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IndividualVault *IndividualVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IndividualVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IndividualVault *IndividualVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IndividualVault *IndividualVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IndividualVault.Contract.contract.Transact(opts, method, params...)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address owner, uint256 heartbeatInterval, uint256 lastHeartbeat, uint256 unlockTime, uint256 gracePeriod, uint256 requiredApprovals, uint256 approvalCount, uint256 totalBalanceAtUnlock, bool isLocked, bool gracePeriodActive)
func (_IndividualVault *IndividualVaultCaller) Config(opts *bind.CallOpts) (struct {
	Owner                common.Address
	HeartbeatInterval    *big.Int
	LastHeartbeat        *big.Int
	UnlockTime           *big.Int
	GracePeriod          *big.Int
	RequiredApprovals    *big.Int
	ApprovalCount        *big.Int
	TotalBalanceAtUnlock *big.Int
	IsLocked             bool
	GracePeriodActive    bool
}, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "config")

	outstruct := new(struct {
		Owner                common.Address
		HeartbeatInterval    *big.Int
		LastHeartbeat        *big.Int
		UnlockTime           *big.Int
		GracePeriod          *big.Int
		RequiredApprovals    *big.Int
		ApprovalCount        *big.Int
		TotalBalanceAtUnlock *big.Int
		IsLocked             bool
		GracePeriodActive    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.HeartbeatInterval = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LastHeartbeat = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UnlockTime = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.GracePeriod = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.RequiredApprovals = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ApprovalCount = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.TotalBalanceAtUnlock = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.IsLocked = *abi.ConvertType(out[8], new(bool)).(*bool)
	outstruct.GracePeriodActive = *abi.ConvertType(out[9], new(bool)).(*bool)

	return *outstruct, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address owner, uint256 heartbeatInterval, uint256 lastHeartbeat, uint256 unlockTime, uint256 gracePeriod, uint256 requiredApprovals, uint256 approvalCount, uint256 totalBalanceAtUnlock, bool isLocked, bool gracePeriodActive)
func (_IndividualVault *IndividualVaultSession) Config() (struct {
	Owner                common.Address
	HeartbeatInterval    *big.Int
	LastHeartbeat        *big.Int
	UnlockTime           *big.Int
	GracePeriod          *big.Int
	RequiredApprovals    *big.Int
	ApprovalCount        *big.Int
	TotalBalanceAtUnlock *big.Int
	IsLocked             bool
	GracePeriodActive    bool
}, error) {
	return _IndividualVault.Contract.Config(&_IndividualVault.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address owner, uint256 heartbeatInterval, uint256 lastHeartbeat, uint256 unlockTime, uint256 gracePeriod, uint256 requiredApprovals, uint256 approvalCount, uint256 totalBalanceAtUnlock, bool isLocked, bool gracePeriodActive)
func (_IndividualVault *IndividualVaultCallerSession) Config() (struct {
	Owner                common.Address
	HeartbeatInterval    *big.Int
	LastHeartbeat        *big.Int
	UnlockTime           *big.Int
	GracePeriod          *big.Int
	RequiredApprovals    *big.Int
	ApprovalCount        *big.Int
	TotalBalanceAtUnlock *big.Int
	IsLocked             bool
	GracePeriodActive    bool
}, error) {
	return _IndividualVault.Contract.Config(&_IndividualVault.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_IndividualVault *IndividualVaultCaller) GetBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "getBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_IndividualVault *IndividualVaultSession) GetBalance() (*big.Int, error) {
	return _IndividualVault.Contract.GetBalance(&_IndividualVault.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_IndividualVault *IndividualVaultCallerSession) GetBalance() (*big.Int, error) {
	return _IndividualVault.Contract.GetBalance(&_IndividualVault.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address[],uint256[],uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_IndividualVault *IndividualVaultCaller) GetConfig(opts *bind.CallOpts) (IndividualVaultVaultConfig, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(IndividualVaultVaultConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IndividualVaultVaultConfig)).(*IndividualVaultVaultConfig)

	return out0, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address[],uint256[],uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_IndividualVault *IndividualVaultSession) GetConfig() (IndividualVaultVaultConfig, error) {
	return _IndividualVault.Contract.GetConfig(&_IndividualVault.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns((address,address[],uint256[],uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_IndividualVault *IndividualVaultCallerSession) GetConfig() (IndividualVaultVaultConfig, error) {
	return _IndividualVault.Contract.GetConfig(&_IndividualVault.CallOpts)
}

// GetHeir is a free data retrieval call binding the contract method 0x975c2308.
//
// Solidity: function getHeir(uint256 _index) view returns(address)
func (_IndividualVault *IndividualVaultCaller) GetHeir(opts *bind.CallOpts, _index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "getHeir", _index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetHeir is a free data retrieval call binding the contract method 0x975c2308.
//
// Solidity: function getHeir(uint256 _index) view returns(address)
func (_IndividualVault *IndividualVaultSession) GetHeir(_index *big.Int) (common.Address, error) {
	return _IndividualVault.Contract.GetHeir(&_IndividualVault.CallOpts, _index)
}

// GetHeir is a free data retrieval call binding the contract method 0x975c2308.
//
// Solidity: function getHeir(uint256 _index) view returns(address)
func (_IndividualVault *IndividualVaultCallerSession) GetHeir(_index *big.Int) (common.Address, error) {
	return _IndividualVault.Contract.GetHeir(&_IndividualVault.CallOpts, _index)
}

// GetHeirCount is a free data retrieval call binding the contract method 0xae69a893.
//
// Solidity: function getHeirCount() view returns(uint256)
func (_IndividualVault *IndividualVaultCaller) GetHeirCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "getHeirCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetHeirCount is a free data retrieval call binding the contract method 0xae69a893.
//
// Solidity: function getHeirCount() view returns(uint256)
func (_IndividualVault *IndividualVaultSession) GetHeirCount() (*big.Int, error) {
	return _IndividualVault.Contract.GetHeirCount(&_IndividualVault.CallOpts)
}

// GetHeirCount is a free data retrieval call binding the contract method 0xae69a893.
//
// Solidity: function getHeirCount() view returns(uint256)
func (_IndividualVault *IndividualVaultCallerSession) GetHeirCount() (*big.Int, error) {
	return _IndividualVault.Contract.GetHeirCount(&_IndividualVault.CallOpts)
}

// GetHeirShare is a free data retrieval call binding the contract method 0x36b9f011.
//
// Solidity: function getHeirShare(uint256 _index) view returns(uint256)
func (_IndividualVault *IndividualVaultCaller) GetHeirShare(opts *bind.CallOpts, _index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "getHeirShare", _index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetHeirShare is a free data retrieval call binding the contract method 0x36b9f011.
//
// Solidity: function getHeirShare(uint256 _index) view returns(uint256)
func (_IndividualVault *IndividualVaultSession) GetHeirShare(_index *big.Int) (*big.Int, error) {
	return _IndividualVault.Contract.GetHeirShare(&_IndividualVault.CallOpts, _index)
}

// GetHeirShare is a free data retrieval call binding the contract method 0x36b9f011.
//
// Solidity: function getHeirShare(uint256 _index) view returns(uint256)
func (_IndividualVault *IndividualVaultCallerSession) GetHeirShare(_index *big.Int) (*big.Int, error) {
	return _IndividualVault.Contract.GetHeirShare(&_IndividualVault.CallOpts, _index)
}

// HeirApprovals is a free data retrieval call binding the contract method 0x6aadb7f0.
//
// Solidity: function heirApprovals(address ) view returns(bool)
func (_IndividualVault *IndividualVaultCaller) HeirApprovals(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "heirApprovals", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HeirApprovals is a free data retrieval call binding the contract method 0x6aadb7f0.
//
// Solidity: function heirApprovals(address ) view returns(bool)
func (_IndividualVault *IndividualVaultSession) HeirApprovals(arg0 common.Address) (bool, error) {
	return _IndividualVault.Contract.HeirApprovals(&_IndividualVault.CallOpts, arg0)
}

// HeirApprovals is a free data retrieval call binding the contract method 0x6aadb7f0.
//
// Solidity: function heirApprovals(address ) view returns(bool)
func (_IndividualVault *IndividualVaultCallerSession) HeirApprovals(arg0 common.Address) (bool, error) {
	return _IndividualVault.Contract.HeirApprovals(&_IndividualVault.CallOpts, arg0)
}

// HeirClaimed is a free data retrieval call binding the contract method 0x9be87bcb.
//
// Solidity: function heirClaimed(address ) view returns(bool)
func (_IndividualVault *IndividualVaultCaller) HeirClaimed(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "heirClaimed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HeirClaimed is a free data retrieval call binding the contract method 0x9be87bcb.
//
// Solidity: function heirClaimed(address ) view returns(bool)
func (_IndividualVault *IndividualVaultSession) HeirClaimed(arg0 common.Address) (bool, error) {
	return _IndividualVault.Contract.HeirClaimed(&_IndividualVault.CallOpts, arg0)
}

// HeirClaimed is a free data retrieval call binding the contract method 0x9be87bcb.
//
// Solidity: function heirClaimed(address ) view returns(bool)
func (_IndividualVault *IndividualVaultCallerSession) HeirClaimed(arg0 common.Address) (bool, error) {
	return _IndividualVault.Contract.HeirClaimed(&_IndividualVault.CallOpts, arg0)
}

// IsClaimable is a free data retrieval call binding the contract method 0x74478bb3.
//
// Solidity: function isClaimable() view returns(bool)
func (_IndividualVault *IndividualVaultCaller) IsClaimable(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "isClaimable")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsClaimable is a free data retrieval call binding the contract method 0x74478bb3.
//
// Solidity: function isClaimable() view returns(bool)
func (_IndividualVault *IndividualVaultSession) IsClaimable() (bool, error) {
	return _IndividualVault.Contract.IsClaimable(&_IndividualVault.CallOpts)
}

// IsClaimable is a free data retrieval call binding the contract method 0x74478bb3.
//
// Solidity: function isClaimable() view returns(bool)
func (_IndividualVault *IndividualVaultCallerSession) IsClaimable() (bool, error) {
	return _IndividualVault.Contract.IsClaimable(&_IndividualVault.CallOpts)
}

// IsHeir is a free data retrieval call binding the contract method 0x568ddfc6.
//
// Solidity: function isHeir(address _address) view returns(bool)
func (_IndividualVault *IndividualVaultCaller) IsHeir(opts *bind.CallOpts, _address common.Address) (bool, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "isHeir", _address)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsHeir is a free data retrieval call binding the contract method 0x568ddfc6.
//
// Solidity: function isHeir(address _address) view returns(bool)
func (_IndividualVault *IndividualVaultSession) IsHeir(_address common.Address) (bool, error) {
	return _IndividualVault.Contract.IsHeir(&_IndividualVault.CallOpts, _address)
}

// IsHeir is a free data retrieval call binding the contract method 0x568ddfc6.
//
// Solidity: function isHeir(address _address) view returns(bool)
func (_IndividualVault *IndividualVaultCallerSession) IsHeir(_address common.Address) (bool, error) {
	return _IndividualVault.Contract.IsHeir(&_IndividualVault.CallOpts, _address)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IndividualVault *IndividualVaultCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IndividualVault.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IndividualVault *IndividualVaultSession) Paused() (bool, error) {
	return _IndividualVault.Contract.Paused(&_IndividualVault.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IndividualVault *IndividualVaultCallerSession) Paused() (bool, error) {
	return _IndividualVault.Contract.Paused(&_IndividualVault.CallOpts)
}

// ApproveInheritance is a paid mutator transaction binding the contract method 0xba58db51.
//
// Solidity: function approveInheritance() returns()
func (_IndividualVault *IndividualVaultTransactor) ApproveInheritance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "approveInheritance")
}

// ApproveInheritance is a paid mutator transaction binding the contract method 0xba58db51.
//
// Solidity: function approveInheritance() returns()
func (_IndividualVault *IndividualVaultSession) ApproveInheritance() (*types.Transaction, error) {
	return _IndividualVault.Contract.ApproveInheritance(&_IndividualVault.TransactOpts)
}

// ApproveInheritance is a paid mutator transaction binding the contract method 0xba58db51.
//
// Solidity: function approveInheritance() returns()
func (_IndividualVault *IndividualVaultTransactorSession) ApproveInheritance() (*types.Transaction, error) {
	return _IndividualVault.Contract.ApproveInheritance(&_IndividualVault.TransactOpts)
}

// CheckAndUnlock is a paid mutator transaction binding the contract method 0xca819c1c.
//
// Solidity: function checkAndUnlock() returns()
func (_IndividualVault *IndividualVaultTransactor) CheckAndUnlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "checkAndUnlock")
}

// CheckAndUnlock is a paid mutator transaction binding the contract method 0xca819c1c.
//
// Solidity: function checkAndUnlock() returns()
func (_IndividualVault *IndividualVaultSession) CheckAndUnlock() (*types.Transaction, error) {
	return _IndividualVault.Contract.CheckAndUnlock(&_IndividualVault.TransactOpts)
}

// CheckAndUnlock is a paid mutator transaction binding the contract method 0xca819c1c.
//
// Solidity: function checkAndUnlock() returns()
func (_IndividualVault *IndividualVaultTransactorSession) CheckAndUnlock() (*types.Transaction, error) {
	return _IndividualVault.Contract.CheckAndUnlock(&_IndividualVault.TransactOpts)
}

// ClaimInheritance is a paid mutator transaction binding the contract method 0x9f3c4416.
//
// Solidity: function claimInheritance() returns()
func (_IndividualVault *IndividualVaultTransactor) ClaimInheritance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "claimInheritance")
}

// ClaimInheritance is a paid mutator transaction binding the contract method 0x9f3c4416.
//
// Solidity: function claimInheritance() returns()
func (_IndividualVault *IndividualVaultSession) ClaimInheritance() (*types.Transaction, error) {
	return _IndividualVault.Contract.ClaimInheritance(&_IndividualVault.TransactOpts)
}

// ClaimInheritance is a paid mutator transaction binding the contract method 0x9f3c4416.
//
// Solidity: function claimInheritance() returns()
func (_IndividualVault *IndividualVaultTransactorSession) ClaimInheritance() (*types.Transaction, error) {
	return _IndividualVault.Contract.ClaimInheritance(&_IndividualVault.TransactOpts)
}

// CommitHeartbeat is a paid mutator transaction binding the contract method 0x14b40a22.
//
// Solidity: function commitHeartbeat(bytes32 _commitment) returns()
func (_IndividualVault *IndividualVaultTransactor) CommitHeartbeat(opts *bind.TransactOpts, _commitment [32]byte) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "commitHeartbeat", _commitment)
}

// CommitHeartbeat is a paid mutator transaction binding the contract method 0x14b40a22.
//
// Solidity: function commitHeartbeat(bytes32 _commitment) returns()
func (_IndividualVault *IndividualVaultSession) CommitHeartbeat(_commitment [32]byte) (*types.Transaction, error) {
	return _IndividualVault.Contract.CommitHeartbeat(&_IndividualVault.TransactOpts, _commitment)
}

// CommitHeartbeat is a paid mutator transaction binding the contract method 0x14b40a22.
//
// Solidity: function commitHeartbeat(bytes32 _commitment) returns()
func (_IndividualVault *IndividualVaultTransactorSession) CommitHeartbeat(_commitment [32]byte) (*types.Transaction, error) {
	return _IndividualVault.Contract.CommitHeartbeat(&_IndividualVault.TransactOpts, _commitment)
}

// Initialize is a paid mutator transaction binding the contract method 0x37f17b7d.
//
// Solidity: function initialize(address _owner, address[] _heirs, uint256[] _shares, uint256 _heartbeatInterval, uint256 _gracePeriod, uint256 _requiredApprovals) returns()
func (_IndividualVault *IndividualVaultTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _heirs []common.Address, _shares []*big.Int, _heartbeatInterval *big.Int, _gracePeriod *big.Int, _requiredApprovals *big.Int) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "initialize", _owner, _heirs, _shares, _heartbeatInterval, _gracePeriod, _requiredApprovals)
}

// Initialize is a paid mutator transaction binding the contract method 0x37f17b7d.
//
// Solidity: function initialize(address _owner, address[] _heirs, uint256[] _shares, uint256 _heartbeatInterval, uint256 _gracePeriod, uint256 _requiredApprovals) returns()
func (_IndividualVault *IndividualVaultSession) Initialize(_owner common.Address, _heirs []common.Address, _shares []*big.Int, _heartbeatInterval *big.Int, _gracePeriod *big.Int, _requiredApprovals *big.Int) (*types.Transaction, error) {
	return _IndividualVault.Contract.Initialize(&_IndividualVault.TransactOpts, _owner, _heirs, _shares, _heartbeatInterval, _gracePeriod, _requiredApprovals)
}

// Initialize is a paid mutator transaction binding the contract method 0x37f17b7d.
//
// Solidity: function initialize(address _owner, address[] _heirs, uint256[] _shares, uint256 _heartbeatInterval, uint256 _gracePeriod, uint256 _requiredApprovals) returns()
func (_IndividualVault *IndividualVaultTransactorSession) Initialize(_owner common.Address, _heirs []common.Address, _shares []*big.Int, _heartbeatInterval *big.Int, _gracePeriod *big.Int, _requiredApprovals *big.Int) (*types.Transaction, error) {
	return _IndividualVault.Contract.Initialize(&_IndividualVault.TransactOpts, _owner, _heirs, _shares, _heartbeatInterval, _gracePeriod, _requiredApprovals)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IndividualVault *IndividualVaultTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IndividualVault *IndividualVaultSession) Pause() (*types.Transaction, error) {
	return _IndividualVault.Contract.Pause(&_IndividualVault.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IndividualVault *IndividualVaultTransactorSession) Pause() (*types.Transaction, error) {
	return _IndividualVault.Contract.Pause(&_IndividualVault.TransactOpts)
}

// RevealHeartbeat is a paid mutator transaction binding the contract method 0xa127268d.
//
// Solidity: function revealHeartbeat(bytes32 _nonce) returns()
func (_IndividualVault *IndividualVaultTransactor) RevealHeartbeat(opts *bind.TransactOpts, _nonce [32]byte) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "revealHeartbeat", _nonce)
}

// RevealHeartbeat is a paid mutator transaction binding the contract method 0xa127268d.
//
// Solidity: function revealHeartbeat(bytes32 _nonce) returns()
func (_IndividualVault *IndividualVaultSession) RevealHeartbeat(_nonce [32]byte) (*types.Transaction, error) {
	return _IndividualVault.Contract.RevealHeartbeat(&_IndividualVault.TransactOpts, _nonce)
}

// RevealHeartbeat is a paid mutator transaction binding the contract method 0xa127268d.
//
// Solidity: function revealHeartbeat(bytes32 _nonce) returns()
func (_IndividualVault *IndividualVaultTransactorSession) RevealHeartbeat(_nonce [32]byte) (*types.Transaction, error) {
	return _IndividualVault.Contract.RevealHeartbeat(&_IndividualVault.TransactOpts, _nonce)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IndividualVault *IndividualVaultTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IndividualVault *IndividualVaultSession) Unpause() (*types.Transaction, error) {
	return _IndividualVault.Contract.Unpause(&_IndividualVault.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IndividualVault *IndividualVaultTransactorSession) Unpause() (*types.Transaction, error) {
	return _IndividualVault.Contract.Unpause(&_IndividualVault.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_IndividualVault *IndividualVaultTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _IndividualVault.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_IndividualVault *IndividualVaultSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _IndividualVault.Contract.Withdraw(&_IndividualVault.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_IndividualVault *IndividualVaultTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _IndividualVault.Contract.Withdraw(&_IndividualVault.TransactOpts, _amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IndividualVault *IndividualVaultTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IndividualVault.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IndividualVault *IndividualVaultSession) Receive() (*types.Transaction, error) {
	return _IndividualVault.Contract.Receive(&_IndividualVault.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IndividualVault *IndividualVaultTransactorSession) Receive() (*types.Transaction, error) {
	return _IndividualVault.Contract.Receive(&_IndividualVault.TransactOpts)
}

// IndividualVaultDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the IndividualVault contract.
type IndividualVaultDepositedIterator struct {
	Event *IndividualVaultDeposited // Event containing the contract specifics and raw log

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
func (it *IndividualVaultDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultDeposited)
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
		it.Event = new(IndividualVaultDeposited)
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
func (it *IndividualVaultDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultDeposited represents a Deposited event raised by the IndividualVault contract.
type IndividualVaultDeposited struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed from, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) FilterDeposited(opts *bind.FilterOpts, from []common.Address) (*IndividualVaultDepositedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "Deposited", fromRule)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultDepositedIterator{contract: _IndividualVault.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed from, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *IndividualVaultDeposited, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "Deposited", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultDeposited)
				if err := _IndividualVault.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed from, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) ParseDeposited(log types.Log) (*IndividualVaultDeposited, error) {
	event := new(IndividualVaultDeposited)
	if err := _IndividualVault.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultEmergencyPausedIterator is returned from FilterEmergencyPaused and is used to iterate over the raw logs and unpacked data for EmergencyPaused events raised by the IndividualVault contract.
type IndividualVaultEmergencyPausedIterator struct {
	Event *IndividualVaultEmergencyPaused // Event containing the contract specifics and raw log

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
func (it *IndividualVaultEmergencyPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultEmergencyPaused)
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
		it.Event = new(IndividualVaultEmergencyPaused)
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
func (it *IndividualVaultEmergencyPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultEmergencyPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultEmergencyPaused represents a EmergencyPaused event raised by the IndividualVault contract.
type IndividualVaultEmergencyPaused struct {
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEmergencyPaused is a free log retrieval operation binding the contract event 0x11d8c430ae8dff60717610a0356ce3ac35187fc9569e6adfe01d281cc1c167d5.
//
// Solidity: event EmergencyPaused(uint256 timestamp)
func (_IndividualVault *IndividualVaultFilterer) FilterEmergencyPaused(opts *bind.FilterOpts) (*IndividualVaultEmergencyPausedIterator, error) {

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "EmergencyPaused")
	if err != nil {
		return nil, err
	}
	return &IndividualVaultEmergencyPausedIterator{contract: _IndividualVault.contract, event: "EmergencyPaused", logs: logs, sub: sub}, nil
}

// WatchEmergencyPaused is a free log subscription operation binding the contract event 0x11d8c430ae8dff60717610a0356ce3ac35187fc9569e6adfe01d281cc1c167d5.
//
// Solidity: event EmergencyPaused(uint256 timestamp)
func (_IndividualVault *IndividualVaultFilterer) WatchEmergencyPaused(opts *bind.WatchOpts, sink chan<- *IndividualVaultEmergencyPaused) (event.Subscription, error) {

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "EmergencyPaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultEmergencyPaused)
				if err := _IndividualVault.contract.UnpackLog(event, "EmergencyPaused", log); err != nil {
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

// ParseEmergencyPaused is a log parse operation binding the contract event 0x11d8c430ae8dff60717610a0356ce3ac35187fc9569e6adfe01d281cc1c167d5.
//
// Solidity: event EmergencyPaused(uint256 timestamp)
func (_IndividualVault *IndividualVaultFilterer) ParseEmergencyPaused(log types.Log) (*IndividualVaultEmergencyPaused, error) {
	event := new(IndividualVaultEmergencyPaused)
	if err := _IndividualVault.contract.UnpackLog(event, "EmergencyPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultGracePeriodStartedIterator is returned from FilterGracePeriodStarted and is used to iterate over the raw logs and unpacked data for GracePeriodStarted events raised by the IndividualVault contract.
type IndividualVaultGracePeriodStartedIterator struct {
	Event *IndividualVaultGracePeriodStarted // Event containing the contract specifics and raw log

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
func (it *IndividualVaultGracePeriodStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultGracePeriodStarted)
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
		it.Event = new(IndividualVaultGracePeriodStarted)
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
func (it *IndividualVaultGracePeriodStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultGracePeriodStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultGracePeriodStarted represents a GracePeriodStarted event raised by the IndividualVault contract.
type IndividualVaultGracePeriodStarted struct {
	EndTime *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGracePeriodStarted is a free log retrieval operation binding the contract event 0x8f8e13a8081be398c5032b04ccda337e916cf2d86c3708db9a734e6db8eb5490.
//
// Solidity: event GracePeriodStarted(uint256 endTime)
func (_IndividualVault *IndividualVaultFilterer) FilterGracePeriodStarted(opts *bind.FilterOpts) (*IndividualVaultGracePeriodStartedIterator, error) {

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "GracePeriodStarted")
	if err != nil {
		return nil, err
	}
	return &IndividualVaultGracePeriodStartedIterator{contract: _IndividualVault.contract, event: "GracePeriodStarted", logs: logs, sub: sub}, nil
}

// WatchGracePeriodStarted is a free log subscription operation binding the contract event 0x8f8e13a8081be398c5032b04ccda337e916cf2d86c3708db9a734e6db8eb5490.
//
// Solidity: event GracePeriodStarted(uint256 endTime)
func (_IndividualVault *IndividualVaultFilterer) WatchGracePeriodStarted(opts *bind.WatchOpts, sink chan<- *IndividualVaultGracePeriodStarted) (event.Subscription, error) {

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "GracePeriodStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultGracePeriodStarted)
				if err := _IndividualVault.contract.UnpackLog(event, "GracePeriodStarted", log); err != nil {
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

// ParseGracePeriodStarted is a log parse operation binding the contract event 0x8f8e13a8081be398c5032b04ccda337e916cf2d86c3708db9a734e6db8eb5490.
//
// Solidity: event GracePeriodStarted(uint256 endTime)
func (_IndividualVault *IndividualVaultFilterer) ParseGracePeriodStarted(log types.Log) (*IndividualVaultGracePeriodStarted, error) {
	event := new(IndividualVaultGracePeriodStarted)
	if err := _IndividualVault.contract.UnpackLog(event, "GracePeriodStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultHeartbeatIterator is returned from FilterHeartbeat and is used to iterate over the raw logs and unpacked data for Heartbeat events raised by the IndividualVault contract.
type IndividualVaultHeartbeatIterator struct {
	Event *IndividualVaultHeartbeat // Event containing the contract specifics and raw log

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
func (it *IndividualVaultHeartbeatIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultHeartbeat)
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
		it.Event = new(IndividualVaultHeartbeat)
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
func (it *IndividualVaultHeartbeatIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultHeartbeatIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultHeartbeat represents a Heartbeat event raised by the IndividualVault contract.
type IndividualVaultHeartbeat struct {
	Timestamp  *big.Int
	Commitment [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterHeartbeat is a free log retrieval operation binding the contract event 0xda505ce43338f29af55a60fece8fb6fe412d9fbdcd5b29b18e387b314dd8ba25.
//
// Solidity: event Heartbeat(uint256 timestamp, bytes32 commitment)
func (_IndividualVault *IndividualVaultFilterer) FilterHeartbeat(opts *bind.FilterOpts) (*IndividualVaultHeartbeatIterator, error) {

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "Heartbeat")
	if err != nil {
		return nil, err
	}
	return &IndividualVaultHeartbeatIterator{contract: _IndividualVault.contract, event: "Heartbeat", logs: logs, sub: sub}, nil
}

// WatchHeartbeat is a free log subscription operation binding the contract event 0xda505ce43338f29af55a60fece8fb6fe412d9fbdcd5b29b18e387b314dd8ba25.
//
// Solidity: event Heartbeat(uint256 timestamp, bytes32 commitment)
func (_IndividualVault *IndividualVaultFilterer) WatchHeartbeat(opts *bind.WatchOpts, sink chan<- *IndividualVaultHeartbeat) (event.Subscription, error) {

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "Heartbeat")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultHeartbeat)
				if err := _IndividualVault.contract.UnpackLog(event, "Heartbeat", log); err != nil {
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

// ParseHeartbeat is a log parse operation binding the contract event 0xda505ce43338f29af55a60fece8fb6fe412d9fbdcd5b29b18e387b314dd8ba25.
//
// Solidity: event Heartbeat(uint256 timestamp, bytes32 commitment)
func (_IndividualVault *IndividualVaultFilterer) ParseHeartbeat(log types.Log) (*IndividualVaultHeartbeat, error) {
	event := new(IndividualVaultHeartbeat)
	if err := _IndividualVault.contract.UnpackLog(event, "Heartbeat", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultInheritanceApprovedIterator is returned from FilterInheritanceApproved and is used to iterate over the raw logs and unpacked data for InheritanceApproved events raised by the IndividualVault contract.
type IndividualVaultInheritanceApprovedIterator struct {
	Event *IndividualVaultInheritanceApproved // Event containing the contract specifics and raw log

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
func (it *IndividualVaultInheritanceApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultInheritanceApproved)
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
		it.Event = new(IndividualVaultInheritanceApproved)
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
func (it *IndividualVaultInheritanceApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultInheritanceApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultInheritanceApproved represents a InheritanceApproved event raised by the IndividualVault contract.
type IndividualVaultInheritanceApproved struct {
	Heir common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterInheritanceApproved is a free log retrieval operation binding the contract event 0xb7369063899dd4b20f32f32aeb10e936a86a5c7221e4bbfb38e7aa697ada7d82.
//
// Solidity: event InheritanceApproved(address indexed heir)
func (_IndividualVault *IndividualVaultFilterer) FilterInheritanceApproved(opts *bind.FilterOpts, heir []common.Address) (*IndividualVaultInheritanceApprovedIterator, error) {

	var heirRule []interface{}
	for _, heirItem := range heir {
		heirRule = append(heirRule, heirItem)
	}

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "InheritanceApproved", heirRule)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultInheritanceApprovedIterator{contract: _IndividualVault.contract, event: "InheritanceApproved", logs: logs, sub: sub}, nil
}

// WatchInheritanceApproved is a free log subscription operation binding the contract event 0xb7369063899dd4b20f32f32aeb10e936a86a5c7221e4bbfb38e7aa697ada7d82.
//
// Solidity: event InheritanceApproved(address indexed heir)
func (_IndividualVault *IndividualVaultFilterer) WatchInheritanceApproved(opts *bind.WatchOpts, sink chan<- *IndividualVaultInheritanceApproved, heir []common.Address) (event.Subscription, error) {

	var heirRule []interface{}
	for _, heirItem := range heir {
		heirRule = append(heirRule, heirItem)
	}

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "InheritanceApproved", heirRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultInheritanceApproved)
				if err := _IndividualVault.contract.UnpackLog(event, "InheritanceApproved", log); err != nil {
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

// ParseInheritanceApproved is a log parse operation binding the contract event 0xb7369063899dd4b20f32f32aeb10e936a86a5c7221e4bbfb38e7aa697ada7d82.
//
// Solidity: event InheritanceApproved(address indexed heir)
func (_IndividualVault *IndividualVaultFilterer) ParseInheritanceApproved(log types.Log) (*IndividualVaultInheritanceApproved, error) {
	event := new(IndividualVaultInheritanceApproved)
	if err := _IndividualVault.contract.UnpackLog(event, "InheritanceApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultInheritanceClaimedIterator is returned from FilterInheritanceClaimed and is used to iterate over the raw logs and unpacked data for InheritanceClaimed events raised by the IndividualVault contract.
type IndividualVaultInheritanceClaimedIterator struct {
	Event *IndividualVaultInheritanceClaimed // Event containing the contract specifics and raw log

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
func (it *IndividualVaultInheritanceClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultInheritanceClaimed)
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
		it.Event = new(IndividualVaultInheritanceClaimed)
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
func (it *IndividualVaultInheritanceClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultInheritanceClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultInheritanceClaimed represents a InheritanceClaimed event raised by the IndividualVault contract.
type IndividualVaultInheritanceClaimed struct {
	Heir   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterInheritanceClaimed is a free log retrieval operation binding the contract event 0x92a193f55349470137c2eab4bd04427f216036ec3ff6ba54aed58e17337af6a2.
//
// Solidity: event InheritanceClaimed(address indexed heir, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) FilterInheritanceClaimed(opts *bind.FilterOpts, heir []common.Address) (*IndividualVaultInheritanceClaimedIterator, error) {

	var heirRule []interface{}
	for _, heirItem := range heir {
		heirRule = append(heirRule, heirItem)
	}

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "InheritanceClaimed", heirRule)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultInheritanceClaimedIterator{contract: _IndividualVault.contract, event: "InheritanceClaimed", logs: logs, sub: sub}, nil
}

// WatchInheritanceClaimed is a free log subscription operation binding the contract event 0x92a193f55349470137c2eab4bd04427f216036ec3ff6ba54aed58e17337af6a2.
//
// Solidity: event InheritanceClaimed(address indexed heir, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) WatchInheritanceClaimed(opts *bind.WatchOpts, sink chan<- *IndividualVaultInheritanceClaimed, heir []common.Address) (event.Subscription, error) {

	var heirRule []interface{}
	for _, heirItem := range heir {
		heirRule = append(heirRule, heirItem)
	}

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "InheritanceClaimed", heirRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultInheritanceClaimed)
				if err := _IndividualVault.contract.UnpackLog(event, "InheritanceClaimed", log); err != nil {
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

// ParseInheritanceClaimed is a log parse operation binding the contract event 0x92a193f55349470137c2eab4bd04427f216036ec3ff6ba54aed58e17337af6a2.
//
// Solidity: event InheritanceClaimed(address indexed heir, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) ParseInheritanceClaimed(log types.Log) (*IndividualVaultInheritanceClaimed, error) {
	event := new(IndividualVaultInheritanceClaimed)
	if err := _IndividualVault.contract.UnpackLog(event, "InheritanceClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IndividualVault contract.
type IndividualVaultInitializedIterator struct {
	Event *IndividualVaultInitialized // Event containing the contract specifics and raw log

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
func (it *IndividualVaultInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultInitialized)
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
		it.Event = new(IndividualVaultInitialized)
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
func (it *IndividualVaultInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultInitialized represents a Initialized event raised by the IndividualVault contract.
type IndividualVaultInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IndividualVault *IndividualVaultFilterer) FilterInitialized(opts *bind.FilterOpts) (*IndividualVaultInitializedIterator, error) {

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IndividualVaultInitializedIterator{contract: _IndividualVault.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IndividualVault *IndividualVaultFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IndividualVaultInitialized) (event.Subscription, error) {

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultInitialized)
				if err := _IndividualVault.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IndividualVault *IndividualVaultFilterer) ParseInitialized(log types.Log) (*IndividualVaultInitialized, error) {
	event := new(IndividualVaultInitialized)
	if err := _IndividualVault.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the IndividualVault contract.
type IndividualVaultPausedIterator struct {
	Event *IndividualVaultPaused // Event containing the contract specifics and raw log

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
func (it *IndividualVaultPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultPaused)
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
		it.Event = new(IndividualVaultPaused)
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
func (it *IndividualVaultPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultPaused represents a Paused event raised by the IndividualVault contract.
type IndividualVaultPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IndividualVault *IndividualVaultFilterer) FilterPaused(opts *bind.FilterOpts) (*IndividualVaultPausedIterator, error) {

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IndividualVaultPausedIterator{contract: _IndividualVault.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IndividualVault *IndividualVaultFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IndividualVaultPaused) (event.Subscription, error) {

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultPaused)
				if err := _IndividualVault.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IndividualVault *IndividualVaultFilterer) ParsePaused(log types.Log) (*IndividualVaultPaused, error) {
	event := new(IndividualVaultPaused)
	if err := _IndividualVault.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultUnlockCancelledIterator is returned from FilterUnlockCancelled and is used to iterate over the raw logs and unpacked data for UnlockCancelled events raised by the IndividualVault contract.
type IndividualVaultUnlockCancelledIterator struct {
	Event *IndividualVaultUnlockCancelled // Event containing the contract specifics and raw log

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
func (it *IndividualVaultUnlockCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultUnlockCancelled)
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
		it.Event = new(IndividualVaultUnlockCancelled)
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
func (it *IndividualVaultUnlockCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultUnlockCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultUnlockCancelled represents a UnlockCancelled event raised by the IndividualVault contract.
type IndividualVaultUnlockCancelled struct {
	Owner     common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnlockCancelled is a free log retrieval operation binding the contract event 0xd1bbb16d648df2dfd0b1b3d88a5bc3f262e0323d3af8fd678115b4f1e7b4b88a.
//
// Solidity: event UnlockCancelled(address indexed owner, uint256 timestamp)
func (_IndividualVault *IndividualVaultFilterer) FilterUnlockCancelled(opts *bind.FilterOpts, owner []common.Address) (*IndividualVaultUnlockCancelledIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "UnlockCancelled", ownerRule)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultUnlockCancelledIterator{contract: _IndividualVault.contract, event: "UnlockCancelled", logs: logs, sub: sub}, nil
}

// WatchUnlockCancelled is a free log subscription operation binding the contract event 0xd1bbb16d648df2dfd0b1b3d88a5bc3f262e0323d3af8fd678115b4f1e7b4b88a.
//
// Solidity: event UnlockCancelled(address indexed owner, uint256 timestamp)
func (_IndividualVault *IndividualVaultFilterer) WatchUnlockCancelled(opts *bind.WatchOpts, sink chan<- *IndividualVaultUnlockCancelled, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "UnlockCancelled", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultUnlockCancelled)
				if err := _IndividualVault.contract.UnpackLog(event, "UnlockCancelled", log); err != nil {
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

// ParseUnlockCancelled is a log parse operation binding the contract event 0xd1bbb16d648df2dfd0b1b3d88a5bc3f262e0323d3af8fd678115b4f1e7b4b88a.
//
// Solidity: event UnlockCancelled(address indexed owner, uint256 timestamp)
func (_IndividualVault *IndividualVaultFilterer) ParseUnlockCancelled(log types.Log) (*IndividualVaultUnlockCancelled, error) {
	event := new(IndividualVaultUnlockCancelled)
	if err := _IndividualVault.contract.UnpackLog(event, "UnlockCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the IndividualVault contract.
type IndividualVaultUnpausedIterator struct {
	Event *IndividualVaultUnpaused // Event containing the contract specifics and raw log

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
func (it *IndividualVaultUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultUnpaused)
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
		it.Event = new(IndividualVaultUnpaused)
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
func (it *IndividualVaultUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultUnpaused represents a Unpaused event raised by the IndividualVault contract.
type IndividualVaultUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IndividualVault *IndividualVaultFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IndividualVaultUnpausedIterator, error) {

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IndividualVaultUnpausedIterator{contract: _IndividualVault.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IndividualVault *IndividualVaultFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IndividualVaultUnpaused) (event.Subscription, error) {

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultUnpaused)
				if err := _IndividualVault.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IndividualVault *IndividualVaultFilterer) ParseUnpaused(log types.Log) (*IndividualVaultUnpaused, error) {
	event := new(IndividualVaultUnpaused)
	if err := _IndividualVault.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultVaultUnlockedIterator is returned from FilterVaultUnlocked and is used to iterate over the raw logs and unpacked data for VaultUnlocked events raised by the IndividualVault contract.
type IndividualVaultVaultUnlockedIterator struct {
	Event *IndividualVaultVaultUnlocked // Event containing the contract specifics and raw log

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
func (it *IndividualVaultVaultUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultVaultUnlocked)
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
		it.Event = new(IndividualVaultVaultUnlocked)
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
func (it *IndividualVaultVaultUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultVaultUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultVaultUnlocked represents a VaultUnlocked event raised by the IndividualVault contract.
type IndividualVaultVaultUnlocked struct {
	UnlockTime *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVaultUnlocked is a free log retrieval operation binding the contract event 0x3a119f7ba8b2d0b50b2052497850fdad1464c8f48288498812ad2673e06f6c83.
//
// Solidity: event VaultUnlocked(uint256 unlockTime)
func (_IndividualVault *IndividualVaultFilterer) FilterVaultUnlocked(opts *bind.FilterOpts) (*IndividualVaultVaultUnlockedIterator, error) {

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "VaultUnlocked")
	if err != nil {
		return nil, err
	}
	return &IndividualVaultVaultUnlockedIterator{contract: _IndividualVault.contract, event: "VaultUnlocked", logs: logs, sub: sub}, nil
}

// WatchVaultUnlocked is a free log subscription operation binding the contract event 0x3a119f7ba8b2d0b50b2052497850fdad1464c8f48288498812ad2673e06f6c83.
//
// Solidity: event VaultUnlocked(uint256 unlockTime)
func (_IndividualVault *IndividualVaultFilterer) WatchVaultUnlocked(opts *bind.WatchOpts, sink chan<- *IndividualVaultVaultUnlocked) (event.Subscription, error) {

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "VaultUnlocked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultVaultUnlocked)
				if err := _IndividualVault.contract.UnpackLog(event, "VaultUnlocked", log); err != nil {
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

// ParseVaultUnlocked is a log parse operation binding the contract event 0x3a119f7ba8b2d0b50b2052497850fdad1464c8f48288498812ad2673e06f6c83.
//
// Solidity: event VaultUnlocked(uint256 unlockTime)
func (_IndividualVault *IndividualVaultFilterer) ParseVaultUnlocked(log types.Log) (*IndividualVaultVaultUnlocked, error) {
	event := new(IndividualVaultVaultUnlocked)
	if err := _IndividualVault.contract.UnpackLog(event, "VaultUnlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IndividualVaultWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the IndividualVault contract.
type IndividualVaultWithdrawnIterator struct {
	Event *IndividualVaultWithdrawn // Event containing the contract specifics and raw log

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
func (it *IndividualVaultWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IndividualVaultWithdrawn)
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
		it.Event = new(IndividualVaultWithdrawn)
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
func (it *IndividualVaultWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IndividualVaultWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IndividualVaultWithdrawn represents a Withdrawn event raised by the IndividualVault contract.
type IndividualVaultWithdrawn struct {
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed owner, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) FilterWithdrawn(opts *bind.FilterOpts, owner []common.Address) (*IndividualVaultWithdrawnIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _IndividualVault.contract.FilterLogs(opts, "Withdrawn", ownerRule)
	if err != nil {
		return nil, err
	}
	return &IndividualVaultWithdrawnIterator{contract: _IndividualVault.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed owner, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *IndividualVaultWithdrawn, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _IndividualVault.contract.WatchLogs(opts, "Withdrawn", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IndividualVaultWithdrawn)
				if err := _IndividualVault.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed owner, uint256 amount)
func (_IndividualVault *IndividualVaultFilterer) ParseWithdrawn(log types.Log) (*IndividualVaultWithdrawn, error) {
	event := new(IndividualVaultWithdrawn)
	if err := _IndividualVault.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
