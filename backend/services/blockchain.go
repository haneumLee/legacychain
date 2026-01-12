package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/haneumLee/legacychain/backend/config"
	"github.com/haneumLee/legacychain/backend/pkg/bindings"
)

// BlockchainService handles all blockchain interactions
type BlockchainService struct {
	client         *ethclient.Client
	factoryAddress common.Address
	factory        *bindings.VaultFactory
	chainID        *big.Int
	cfg            *config.Config
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService(cfg *config.Config) (*BlockchainService, error) {
	// Connect to Ethereum node
	client, err := ethclient.Dial(cfg.Blockchain.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Parse factory address
	factoryAddress := common.HexToAddress(cfg.Blockchain.VaultFactoryAddress)

	// Create factory contract instance
	factory, err := bindings.NewVaultFactory(factoryAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create factory instance: %w", err)
	}

	return &BlockchainService{
		client:         client,
		factoryAddress: factoryAddress,
		factory:        factory,
		chainID:        chainID,
		cfg:            cfg,
	}, nil
}

// CreateVault creates a new vault on the blockchain
func (s *BlockchainService) CreateVault(
	privateKey *ecdsa.PrivateKey,
	heirs []common.Address,
	heirShares []*big.Int,
	heartbeatInterval *big.Int,
	gracePeriod *big.Int,
	requiredApprovals *big.Int,
) (vaultAddress common.Address, txHash common.Hash, err error) {
	ctx := context.Background()

	// Get auth transactor
	auth, err := s.getTransactor(privateKey)
	if err != nil {
		return common.Address{}, common.Hash{}, err
	}

	// Call createVault
	tx, err := s.factory.CreateVault(auth, heirs, heirShares, heartbeatInterval, gracePeriod, requiredApprovals)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("failed to create vault: %w", err)
	}

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	// Check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Address{}, common.Hash{}, fmt.Errorf("transaction failed")
	}

	// Parse VaultCreated event to get vault address
	vaultCreatedEvent, err := s.parseVaultCreatedEvent(receipt.Logs)
	if err != nil {
		return common.Address{}, common.Hash{}, err
	}

	return vaultCreatedEvent.Vault, tx.Hash(), nil
}

// GetVaultImplementation gets the vault implementation address
func (s *BlockchainService) GetVaultImplementation() (common.Address, error) {
	opts := &bind.CallOpts{Context: context.Background()}
	return s.factory.VaultImplementation(opts)
}

// GetOwnerVaults gets all vaults owned by an address
func (s *BlockchainService) GetOwnerVaults(owner common.Address) ([]common.Address, error) {
	opts := &bind.CallOpts{Context: context.Background()}
	return s.factory.GetOwnerVaults(opts, owner)
}

// CommitHeartbeat commits a heartbeat to a vault
func (s *BlockchainService) CommitHeartbeat(
	privateKey *ecdsa.PrivateKey,
	vaultAddress common.Address,
	commitment [32]byte,
) (txHash common.Hash, err error) {
	ctx := context.Background()

	// Create vault instance
	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create vault instance: %w", err)
	}

	// Get auth transactor
	auth, err := s.getTransactor(privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	// Call commitHeartbeat
	tx, err := vault.CommitHeartbeat(auth, commitment)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to commit heartbeat: %w", err)
	}

	// Wait for transaction
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Hash{}, fmt.Errorf("transaction failed")
	}

	return tx.Hash(), nil
}

// RevealHeartbeat reveals a heartbeat
func (s *BlockchainService) RevealHeartbeat(
	privateKey *ecdsa.PrivateKey,
	vaultAddress common.Address,
	nonce [32]byte,
) (txHash common.Hash, err error) {
	ctx := context.Background()

	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create vault instance: %w", err)
	}

	auth, err := s.getTransactor(privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := vault.RevealHeartbeat(auth, nonce)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to reveal heartbeat: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Hash{}, fmt.Errorf("transaction failed")
	}

	return tx.Hash(), nil
}

// GetVaultConfig gets the configuration of a vault
func (s *BlockchainService) GetVaultConfig(vaultAddress common.Address) (*bindings.IndividualVaultVaultConfig, error) {
	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault instance: %w", err)
	}

	opts := &bind.CallOpts{Context: context.Background()}
	config, err := vault.GetConfig(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault config: %w", err)
	}

	return &config, nil
}

// GetVaultBalance gets the balance of a vault
func (s *BlockchainService) GetVaultBalance(vaultAddress common.Address) (*big.Int, error) {
	ctx := context.Background()
	balance, err := s.client.BalanceAt(ctx, vaultAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault balance: %w", err)
	}
	return balance, nil
}

// IsHeir checks if an address is an heir of a vault
func (s *BlockchainService) IsHeir(vaultAddress, heirAddress common.Address) (bool, error) {
	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return false, fmt.Errorf("failed to create vault instance: %w", err)
	}

	opts := &bind.CallOpts{Context: context.Background()}
	isHeir, err := vault.IsHeir(opts, heirAddress)
	if err != nil {
		return false, fmt.Errorf("failed to check if heir: %w", err)
	}

	return isHeir, nil
}

// ApproveInheritance approves inheritance for a vault
func (s *BlockchainService) ApproveInheritance(
	privateKey *ecdsa.PrivateKey,
	vaultAddress common.Address,
) (txHash common.Hash, err error) {
	ctx := context.Background()

	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create vault instance: %w", err)
	}

	auth, err := s.getTransactor(privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := vault.ApproveInheritance(auth)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to approve inheritance: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Hash{}, fmt.Errorf("transaction failed")
	}

	return tx.Hash(), nil
}

// ClaimInheritance claims inheritance from a vault
func (s *BlockchainService) ClaimInheritance(
	privateKey *ecdsa.PrivateKey,
	vaultAddress common.Address,
) (txHash common.Hash, amount *big.Int, err error) {
	ctx := context.Background()

	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to create vault instance: %w", err)
	}

	auth, err := s.getTransactor(privateKey)
	if err != nil {
		return common.Hash{}, nil, err
	}

	tx, err := vault.ClaimInheritance(auth)
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to claim inheritance: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Hash{}, nil, fmt.Errorf("transaction failed")
	}

	// Parse InheritanceClaimed event to get claimed amount
	claimedAmount, err := s.parseInheritanceClaimedEvent(receipt.Logs)
	if err != nil {
		return common.Hash{}, nil, err
	}

	return tx.Hash(), claimedAmount, nil
}

// Helper functions

// getTransactor creates a transactor from a private key
func (s *BlockchainService) getTransactor(privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, s.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas price and limit (optional, can be estimated)
	ctx := context.Background()
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(3000000) // 3M gas limit

	return auth, nil
}

// parseVaultCreatedEvent parses VaultCreated event from logs
func (s *BlockchainService) parseVaultCreatedEvent(logs []*types.Log) (*bindings.VaultFactoryVaultCreated, error) {
	factoryABI, err := abi.JSON(strings.NewReader(bindings.VaultFactoryMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse factory ABI: %w", err)
	}

	for _, vLog := range logs {
		if vLog.Address == s.factoryAddress {
			event := new(bindings.VaultFactoryVaultCreated)
			err := factoryABI.UnpackIntoInterface(event, "VaultCreated", vLog.Data)
			if err != nil {
				continue
			}

			// Copy indexed topics (owner, vault address)
			if len(vLog.Topics) >= 3 {
				event.Owner = common.HexToAddress(vLog.Topics[1].Hex())
				event.Vault = common.HexToAddress(vLog.Topics[2].Hex())
			}

			return event, nil
		}
	}

	return nil, fmt.Errorf("VaultCreated event not found")
}

// parseInheritanceClaimedEvent parses InheritanceClaimed event
func (s *BlockchainService) parseInheritanceClaimedEvent(logs []*types.Log) (*big.Int, error) {
	vaultABI, err := abi.JSON(strings.NewReader(bindings.IndividualVaultMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse vault ABI: %w", err)
	}

	for _, vLog := range logs {
		event := new(bindings.IndividualVaultInheritanceClaimed)
		err := vaultABI.UnpackIntoInterface(event, "InheritanceClaimed", vLog.Data)
		if err != nil {
			continue
		}

		return event.Amount, nil
	}

	return nil, fmt.Errorf("InheritanceClaimed event not found")
}

// ParsePrivateKey parses a hex-encoded private key
func ParsePrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	// Remove 0x prefix if present
	hexKey = strings.TrimPrefix(hexKey, "0x")

	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// GetAddressFromPrivateKey derives Ethereum address from private key
func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
