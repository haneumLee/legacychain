package crypto

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// VerifySignature verifies an Ethereum signature using EIP-191 (Personal Sign)
// It recovers the public key from the signature and compares the derived address
// with the claimed address.
//
// Parameters:
//   - address: The Ethereum address that supposedly signed the message (0x-prefixed hex)
//   - message: The original message that was signed (plain text)
//   - signature: The signature (0x-prefixed hex, 65 bytes: r, s, v)
//
// Returns:
//   - bool: true if the signature is valid and matches the address
//   - error: any error that occurred during verification
//
// EIP-191 Personal Sign Format:
//
//	"\x19Ethereum Signed Message:\n" + len(message) + message
//
// Example:
//
//	message := "Login to LegacyChain\nNonce: abc123\nTimestamp: 1234567890"
//	signature := "0x1234...abcd"  // 65 bytes from MetaMask
//	valid, err := VerifySignature("0xYourAddress", message, signature)
func VerifySignature(address, message, signature string) (bool, error) {
	// 1. Add Ethereum Personal Sign prefix
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	// 2. Hash the prefixed message with Keccak256
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// 3. Decode the signature (remove 0x prefix if present)
	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	// 4. Validate signature length (must be 65 bytes: r=32, s=32, v=1)
	if len(sigBytes) != 65 {
		return false, fmt.Errorf("invalid signature length: got %d bytes, expected 65", len(sigBytes))
	}

	// 5. Adjust recovery ID (v) from 27/28 to 0/1
	// MetaMask returns v as 27 or 28, but go-ethereum expects 0 or 1
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// 6. Validate recovery ID (must be 0 or 1)
	if sigBytes[64] > 1 {
		return false, fmt.Errorf("invalid recovery id: %d", sigBytes[64])
	}

	// 7. Recover the public key from the signature
	publicKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// 8. Derive the Ethereum address from the public key
	recoveredAddr := crypto.PubkeyToAddress(*publicKey)

	// 9. Compare addresses (case-insensitive)
	claimedAddr := common.HexToAddress(address)
	if !strings.EqualFold(recoveredAddr.Hex(), claimedAddr.Hex()) {
		return false, nil // Signature is valid but address doesn't match
	}

	return true, nil
}

// HashMessage hashes a message using Keccak256 with Ethereum Personal Sign prefix
// This is useful for generating message hashes on the backend that match
// what MetaMask produces when signing.
//
// Parameters:
//   - message: The plain text message to hash
//
// Returns:
//   - string: The 0x-prefixed hex hash
//
// Example:
//
//	hash := HashMessage("Login to LegacyChain")
//	// Returns: "0x..."
func HashMessage(message string) string {
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	return hash.Hex()
}

// RecoverAddress recovers the Ethereum address from a message and signature
// This is a convenience function that combines signature verification
// and address recovery.
//
// Parameters:
//   - message: The original message that was signed
//   - signature: The signature (0x-prefixed hex)
//
// Returns:
//   - string: The recovered Ethereum address (0x-prefixed)
//   - error: any error that occurred during recovery
func RecoverAddress(message, signature string) (string, error) {
	// Add Ethereum Personal Sign prefix
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	// Hash the message
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// Decode signature
	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %w", err)
	}

	if len(sigBytes) != 65 {
		return "", fmt.Errorf("invalid signature length: %d", len(sigBytes))
	}

	// Adjust recovery ID
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Recover public key
	publicKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return "", fmt.Errorf("failed to recover public key: %w", err)
	}

	// Derive address
	recoveredAddr := crypto.PubkeyToAddress(*publicKey)
	return recoveredAddr.Hex(), nil
}
