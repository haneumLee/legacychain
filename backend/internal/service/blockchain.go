package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/haneumLee/legacychain/backend/config"
	"github.com/haneumLee/legacychain/backend/pkg/bindings"
)

// BlockchainService defines the interface for blockchain operations
type BlockchainService interface {
	// Vault operations
	CreateVault(ctx context.Context, heirs []common.Address, shares []*big.Int, heartbeatInterval, gracePeriod, requiredApprovals *big.Int) (txHash string, err error)
	GetVaultOwner(ctx context.Context, vaultAddress common.Address) (common.Address, error)
	GetVaultConfig(ctx context.Context, vaultAddress common.Address) (*VaultConfig, error)
	
	// Heartbeat operations
	CommitHeartbeat(ctx context.Context, vaultAddr common.Address, commitHash [32]byte) (string, error)
	RevealHeartbeat(ctx context.Context, vaultAddr common.Address, nonce [32]byte) (string, error)
	GetLastHeartbeat(ctx context.Context, vaultAddr common.Address) (*big.Int, error)
	
	// Heir operations
	ApproveInheritance(ctx context.Context, vaultAddr common.Address) (string, error)
	ClaimInheritance(ctx context.Context, vaultAddr common.Address) (string, error)
	GetHeirApprovalStatus(ctx context.Context, vaultAddr common.Address, heirAddr common.Address) (bool, error)
	
	// Event listening
	ListenVaultCreatedEvents(ctx context.Context, handler func(event *bindings.VaultFactoryVaultCreated)) error
	
	// Utility
	GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error)
	Close()
}

// VaultConfig represents the on-chain vault configuration
type VaultConfig struct {
	Owner                 common.Address
	Heirs                 []common.Address
	HeirShares            []*big.Int
	HeartbeatInterval     *big.Int
	LastHeartbeat         *big.Int
	UnlockTime            *big.Int
	GracePeriod           *big.Int
	RequiredApprovals     *big.Int
	ApprovalCount         *big.Int
	TotalBalanceAtUnlock  *big.Int
	IsLocked              bool
	GracePeriodActive     bool
}

// ethBlockchainService is the implementation of BlockchainService
type ethBlockchainService struct {
	client           *ethclient.Client
	wsClient         *ethclient.Client
	chainID          *big.Int
	vaultFactory     *bindings.VaultFactory
	vaultFactoryAddr common.Address
	privateKey       *ecdsa.PrivateKey
	fromAddress      common.Address
}

// NewBlockchainService creates a new BlockchainService instance
func NewBlockchainService(cfg *config.Config) (BlockchainService, error) {
	// 1. Connect HTTP client (for RPC calls)
	client, err := ethclient.Dial(cfg.Blockchain.RpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}

	// 2. Connect WebSocket client (for event listening)
	wsClient, err := ethclient.Dial(cfg.Blockchain.WsURL)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Ethereum WebSocket: %w", err)
	}

	// 3. Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		client.Close()
		wsClient.Close()
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// 4. Validate chain ID
	if chainID.Int64() != cfg.Blockchain.ChainID {
		client.Close()
		wsClient.Close()
		return nil, fmt.Errorf("chain ID mismatch: expected %d, got %d", cfg.Blockchain.ChainID, chainID.Int64())
	}

	// 5. Load VaultFactory contract
	if cfg.Blockchain.VaultFactoryAddress == "" {
		client.Close()
		wsClient.Close()
		return nil, fmt.Errorf("VAULT_FACTORY_ADDRESS not set in config")
	}

	factoryAddr := common.HexToAddress(cfg.Blockchain.VaultFactoryAddress)
	factory, err := bindings.NewVaultFactory(factoryAddr, client)
	if err != nil {
		client.Close()
		wsClient.Close()
		return nil, fmt.Errorf("failed to load VaultFactory contract: %w", err)
	}

	// 6. Load private key
	if cfg.Blockchain.PrivateKey == "" {
		client.Close()
		wsClient.Close()
		return nil, fmt.Errorf("BLOCKCHAIN_PRIVATE_KEY not set in config")
	}

	// Remove 0x prefix if present
	privateKeyHex := strings.TrimPrefix(cfg.Blockchain.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		client.Close()
		wsClient.Close()
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &ethBlockchainService{
		client:           client,
		wsClient:         wsClient,
		chainID:          chainID,
		vaultFactory:     factory,
		vaultFactoryAddr: factoryAddr,
		privateKey:       privateKey,
		fromAddress:      fromAddress,
	}, nil
}

// CreateVault creates a new vault on the blockchain
func (s *ethBlockchainService) CreateVault(
	ctx context.Context,
	heirs []common.Address,
	shares []*big.Int,
	heartbeatInterval, gracePeriod, requiredApprovals *big.Int,
) (string, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	// Estimate gas (optional but recommended)
	auth.GasLimit = 5000000 // 5M gas limit

	tx, err := s.vaultFactory.CreateVault(auth, heirs, shares, heartbeatInterval, gracePeriod, requiredApprovals)
	if err != nil {
		return "", fmt.Errorf("failed to create vault: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// GetVaultOwner returns the owner of a vault
func (s *ethBlockchainService) GetVaultOwner(ctx context.Context, vaultAddress common.Address) (common.Address, error) {
	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to load vault contract: %w", err)
	}

	config, err := vault.Config(&bind.CallOpts{Context: ctx})
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get vault config: %w", err)
	}

	return config.Owner, nil
}

// GetVaultConfig returns the full vault configuration
func (s *ethBlockchainService) GetVaultConfig(ctx context.Context, vaultAddress common.Address) (*VaultConfig, error) {
	vault, err := bindings.NewIndividualVault(vaultAddress, s.client)
	if err != nil {
		return nil, fmt.Errorf("failed to load vault contract: %w", err)
	}

	config, err := vault.Config(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get vault config: %w", err)
	}

	return &VaultConfig{
		Owner:                config.Owner,
		Heirs:                nil, // Will be populated separately if needed
		HeirShares:           nil, // Will be populated separately if needed
		HeartbeatInterval:    config.HeartbeatInterval,
		LastHeartbeat:        config.LastHeartbeat,
		UnlockTime:           config.UnlockTime,
		GracePeriod:          config.GracePeriod,
		RequiredApprovals:    config.RequiredApprovals,
		ApprovalCount:        config.ApprovalCount,
		TotalBalanceAtUnlock: config.TotalBalanceAtUnlock,
		IsLocked:             config.IsLocked,
		GracePeriodActive:    config.GracePeriodActive,
	}, nil
}

// CommitHeartbeat commits a heartbeat hash
func (s *ethBlockchainService) CommitHeartbeat(ctx context.Context, vaultAddr common.Address, commitHash [32]byte) (string, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	vault, err := bindings.NewIndividualVault(vaultAddr, s.client)
	if err != nil {
		return "", fmt.Errorf("failed to load vault contract: %w", err)
	}

	auth.GasLimit = 200000

	tx, err := vault.CommitHeartbeat(auth, commitHash)
	if err != nil {
		return "", fmt.Errorf("failed to commit heartbeat: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// RevealHeartbeat reveals a committed heartbeat
func (s *ethBlockchainService) RevealHeartbeat(ctx context.Context, vaultAddr common.Address, nonce [32]byte) (string, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	vault, err := bindings.NewIndividualVault(vaultAddr, s.client)
	if err != nil {
		return "", fmt.Errorf("failed to load vault contract: %w", err)
	}

	auth.GasLimit = 200000

	tx, err := vault.RevealHeartbeat(auth, nonce)
	if err != nil {
		return "", fmt.Errorf("failed to reveal heartbeat: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// GetLastHeartbeat returns the timestamp of the last heartbeat
func (s *ethBlockchainService) GetLastHeartbeat(ctx context.Context, vaultAddr common.Address) (*big.Int, error) {
	vault, err := bindings.NewIndividualVault(vaultAddr, s.client)
	if err != nil {
		return nil, fmt.Errorf("failed to load vault contract: %w", err)
	}

	config, err := vault.Config(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get vault config: %w", err)
	}

	return config.LastHeartbeat, nil
}

// ApproveInheritance approves inheritance as an heir
func (s *ethBlockchainService) ApproveInheritance(ctx context.Context, vaultAddr common.Address) (string, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	vault, err := bindings.NewIndividualVault(vaultAddr, s.client)
	if err != nil {
		return "", fmt.Errorf("failed to load vault contract: %w", err)
	}

	auth.GasLimit = 150000

	tx, err := vault.ApproveInheritance(auth)
	if err != nil {
		return "", fmt.Errorf("failed to approve inheritance: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// ClaimInheritance claims inheritance as an heir
func (s *ethBlockchainService) ClaimInheritance(ctx context.Context, vaultAddr common.Address) (string, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	vault, err := bindings.NewIndividualVault(vaultAddr, s.client)
	if err != nil {
		return "", fmt.Errorf("failed to load vault contract: %w", err)
	}

	auth.GasLimit = 300000

	tx, err := vault.ClaimInheritance(auth)
	if err != nil {
		return "", fmt.Errorf("failed to claim inheritance: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// GetHeirApprovalStatus returns whether an heir has approved
func (s *ethBlockchainService) GetHeirApprovalStatus(ctx context.Context, vaultAddr common.Address, heirAddr common.Address) (bool, error) {
	vault, err := bindings.NewIndividualVault(vaultAddr, s.client)
	if err != nil {
		return false, fmt.Errorf("failed to load vault contract: %w", err)
	}

	approved, err := vault.HeirApprovals(&bind.CallOpts{Context: ctx}, heirAddr)
	if err != nil {
		return false, fmt.Errorf("failed to get heir approval status: %w", err)
	}

	return approved, nil
}

// ListenVaultCreatedEvents listens for VaultCreated events
func (s *ethBlockchainService) ListenVaultCreatedEvents(ctx context.Context, handler func(event *bindings.VaultFactoryVaultCreated)) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{s.vaultFactoryAddr},
	}

	logs := make(chan types.Log)
	sub, err := s.wsClient.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()

		for {
			select {
			case err := <-sub.Err():
				fmt.Printf("Event subscription error: %v\n", err)
				return
			case vLog := <-logs:
				event, err := s.vaultFactory.ParseVaultCreated(vLog)
				if err != nil {
					fmt.Printf("Failed to parse VaultCreated event: %v\n", err)
					continue
				}

				handler(event)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// GetTransactionReceipt gets the receipt of a transaction
func (s *ethBlockchainService) GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	receipt, err := s.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	return receipt, nil
}

// Close closes the blockchain service connections
func (s *ethBlockchainService) Close() {
	s.client.Close()
	s.wsClient.Close()
}

// getTransactor creates a new transactor with current nonce and gas price
func (s *ethBlockchainService) getTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := s.client.PendingNonceAt(ctx, s.fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // No ETH transfer
	auth.GasPrice = gasPrice
	auth.Context = ctx

	return auth, nil
}
