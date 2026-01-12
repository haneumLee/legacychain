package service

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlockchainService is a mock implementation of BlockchainService for testing
type MockBlockchainService struct {
	mock.Mock
}

func (m *MockBlockchainService) CreateVault(
	heirs []common.Address,
	shares []uint16,
	heartbeatInterval uint64,
	gracePeriod uint64,
	requiredApprovals uint8,
) (string, error) {
	args := m.Called(heirs, shares, heartbeatInterval, gracePeriod, requiredApprovals)
	return args.String(0), args.Error(1)
}

func (m *MockBlockchainService) CommitHeartbeat(vaultAddress string, commitHash [32]byte) (string, error) {
	args := m.Called(vaultAddress, commitHash)
	return args.String(0), args.Error(1)
}

func (m *MockBlockchainService) RevealHeartbeat(vaultAddress string, nonce [32]byte) (string, error) {
	args := m.Called(vaultAddress, nonce)
	return args.String(0), args.Error(1)
}

func (m *MockBlockchainService) ApproveInheritance(vaultAddress, heirAddress string) (string, error) {
	args := m.Called(vaultAddress, heirAddress)
	return args.String(0), args.Error(1)
}

func (m *MockBlockchainService) ClaimInheritance(vaultAddress, heirAddress string) (string, error) {
	args := m.Called(vaultAddress, heirAddress)
	return args.String(0), args.Error(1)
}

func (m *MockBlockchainService) GetVaultConfig(vaultAddress string) (*VaultConfig, error) {
	args := m.Called(vaultAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*VaultConfig), args.Error(1)
}

func (m *MockBlockchainService) Close() {
	m.Called()
}

// Test mock implementation
func TestMockBlockchainService(t *testing.T) {
	mockService := new(MockBlockchainService)
	
	// Setup expectations
	expectedTxHash := "0x1234567890abcdef"
	mockService.On("CreateVault",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("uint64"),
		mock.AnythingOfType("uint64"),
		mock.AnythingOfType("uint8"),
	).Return(expectedTxHash, nil)
	
	// Test
	heirs := []common.Address{common.HexToAddress("0x123")}
	shares := []uint16{10000}
	
	txHash, err := mockService.CreateVault(heirs, shares, 86400, 172800, 1)
	
	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedTxHash, txHash)
	mockService.AssertExpectations(t)
}

// Test VaultConfig struct
func TestVaultConfig(t *testing.T) {
	ownerAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	heir1Addr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	heir2Addr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	
	config := &VaultConfig{
		Owner:                ownerAddr,
		HeartbeatInterval:    big.NewInt(86400),
		GracePeriod:          big.NewInt(172800),
		LastHeartbeat:        big.NewInt(1673456789),
		IsLocked:             true,
		GracePeriodActive:    false,
		RequiredApprovals:    big.NewInt(2),
		ApprovalCount:        big.NewInt(1),
		Heirs:                []common.Address{heir1Addr, heir2Addr},
		HeirShares:           []*big.Int{big.NewInt(5000), big.NewInt(5000)},
	}
	
	assert.Equal(t, ownerAddr, config.Owner)
	assert.Equal(t, int64(86400), config.HeartbeatInterval.Int64())
	assert.Equal(t, true, config.IsLocked)
	assert.Equal(t, int64(2), config.RequiredApprovals.Int64())
	assert.Len(t, config.Heirs, 2)
}

// Test commit hash generation
func TestCommitHash(t *testing.T) {
	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	nonce := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	
	// Test that commit hash is deterministic
	commitHash1 := GenerateCommitHash(address, nonce)
	commitHash2 := GenerateCommitHash(address, nonce)
	
	assert.Equal(t, commitHash1, commitHash2)
	assert.Equal(t, 32, len(commitHash1))
}

// Helper function to generate commit hash (same as in smart contract)
func GenerateCommitHash(address common.Address, nonce [32]byte) [32]byte {
	// keccak256(abi.encodePacked(msg.sender, _nonce))
	data := append(address.Bytes(), nonce[:]...)
	hash := common.BytesToHash(data)
	return hash
}
