package crypto

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// Test VerifySignature with valid and invalid signatures
func TestVerifySignature(t *testing.T) {
	// Generate a test private key and address
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	message := "Login to LegacyChain\nNonce: test-nonce-123\nTimestamp: 1673456789"
	
	// Create valid signature using proper EIP-191 format
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	validSignature, err := crypto.Sign(hash.Bytes(), privateKey)
	assert.NoError(t, err)
	
	validSigHex := hexutil.Encode(validSignature)

	tests := []struct {
		name      string
		address   string
		message   string
		signature string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "Valid signature",
			address:   address,
			message:   message,
			signature: validSigHex,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Wrong address",
			address:   "0x0000000000000000000000000000000000000000",
			message:   message,
			signature: validSigHex,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "Wrong message",
			address:   address,
			message:   "Different message",
			signature: validSigHex,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "Invalid signature format (too short)",
			address:   address,
			message:   message,
			signature: "0x1234",
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "Invalid signature format (not hex)",
			address:   address,
			message:   message,
			signature: "not-a-hex-string",
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := VerifySignature(tt.address, tt.message, tt.signature)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValid, valid)
			}
		})
	}
}

// Test HashMessage function
func TestHashMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Simple message",
			message: "Hello, World!",
		},
		{
			name:    "Login message with nonce",
			message: "Login to LegacyChain\nNonce: abc123\nTimestamp: 1234567890",
		},
		{
			name:    "Empty message",
			message: "",
		},
		{
			name:    "Unicode message",
			message: "안녕하세요 🎉",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashMessage(tt.message)
			
			// Verify hash is hex string with 0x prefix
			assert.True(t, len(hash) > 0)
			assert.Equal(t, "0x", hash[:2])
			assert.Equal(t, 66, len(hash)) // "0x" + 64 hex chars = 66
			
			// Verify hash is deterministic (same input = same output)
			hash2 := HashMessage(tt.message)
			assert.Equal(t, hash, hash2)
		})
	}
}

// Test RecoverAddress function
func TestRecoverAddress(t *testing.T) {
	// Generate test key pair
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	message := "Test message for address recovery"
	
	// Sign the message using proper EIP-191 format
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	assert.NoError(t, err)
	
	signatureHex := hexutil.Encode(signature)

	tests := []struct {
		name      string
		message   string
		signature string
		wantAddr  string
		wantErr   bool
	}{
		{
			name:      "Valid signature recovery",
			message:   message,
			signature: signatureHex,
			wantAddr:  expectedAddress.Hex(),
			wantErr:   false,
		},
		{
			name:      "Invalid signature format",
			message:   message,
			signature: "0xinvalid",
			wantAddr:  "",
			wantErr:   true,
		},
		{
			name:      "Wrong message (different hash)",
			message:   "Different message",
			signature: signatureHex,
			wantAddr:  "", // Will recover different address
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := RecoverAddress(tt.message, tt.signature)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.wantAddr != "" {
					assert.Equal(t, tt.wantAddr, addr)
				}
			}
		})
	}
}

// Benchmark VerifySignature
func BenchmarkVerifySignature(b *testing.B) {
	privateKey, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	message := "Login to LegacyChain\nNonce: test-nonce\nTimestamp: 1673456789"
	
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	signature, _ := crypto.Sign(hash.Bytes(), privateKey)
	signatureHex := hexutil.Encode(signature)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifySignature(address, message, signatureHex)
	}
}

// Benchmark HashMessage
func BenchmarkHashMessage(b *testing.B) {
	message := "Login to LegacyChain\nNonce: test-nonce\nTimestamp: 1673456789"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashMessage(message)
	}
}
