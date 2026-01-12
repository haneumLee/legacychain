package crypto

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	// NonceExpiration is the duration a nonce is valid (5 minutes)
	NonceExpiration = 5 * time.Minute
	
	// NonceKeyPrefix is the Redis key prefix for nonces
	NonceKeyPrefix = "nonce:"
	
	// SignatureMaxAge is the maximum age of a signature timestamp (5 minutes)
	SignatureMaxAge = 5 * time.Minute
)

// NonceManager handles nonce generation, storage, and validation
// for preventing replay attacks in signature-based authentication.
type NonceManager struct {
	redis *redis.Client
}

// NewNonceManager creates a new NonceManager
func NewNonceManager(redisClient *redis.Client) *NonceManager {
	return &NonceManager{
		redis: redisClient,
	}
}

// GenerateNonce generates a new UUID-based nonce and stores it in Redis
// with an expiration time.
//
// Returns:
//   - nonce: A unique identifier (UUID v4)
//   - timestamp: Unix timestamp when the nonce was created
//   - error: any error during Redis operation
func (nm *NonceManager) GenerateNonce(ctx context.Context) (nonce string, timestamp int64, err error) {
	nonce = uuid.New().String()
	timestamp = time.Now().Unix()

	// Store nonce in Redis with expiration
	key := NonceKeyPrefix + nonce
	err = nm.redis.Set(ctx, key, timestamp, NonceExpiration).Err()
	if err != nil {
		return "", 0, fmt.Errorf("failed to store nonce in Redis: %w", err)
	}

	return nonce, timestamp, nil
}

// ValidateNonce checks if a nonce exists in Redis and has not been used.
// If valid, it deletes the nonce to prevent reuse (one-time use).
//
// Parameters:
//   - ctx: Context for Redis operations
//   - nonce: The nonce to validate
//
// Returns:
//   - bool: true if the nonce is valid and has not been used
//   - error: any error during Redis operation
func (nm *NonceManager) ValidateNonce(ctx context.Context, nonce string) (bool, error) {
	key := NonceKeyPrefix + nonce

	// Check if nonce exists
	exists, err := nm.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check nonce in Redis: %w", err)
	}

	if exists == 0 {
		return false, nil // Nonce doesn't exist or already used
	}

	// Delete nonce immediately (one-time use)
	err = nm.redis.Del(ctx, key).Err()
	if err != nil {
		return false, fmt.Errorf("failed to delete nonce from Redis: %w", err)
	}

	return true, nil
}

// ValidateTimestamp checks if a timestamp is within the acceptable range
// (not too old and not in the future).
//
// Parameters:
//   - timestamp: Unix timestamp to validate
//
// Returns:
//   - bool: true if the timestamp is within acceptable range
//   - error: error describing why the timestamp is invalid
func ValidateTimestamp(timestamp int64) (bool, error) {
	now := time.Now().Unix()
	age := now - timestamp

	// Check if timestamp is too old
	if age > int64(SignatureMaxAge.Seconds()) {
		return false, fmt.Errorf("signature expired: %d seconds old (max: %d)", age, int64(SignatureMaxAge.Seconds()))
	}

	// Check if timestamp is in the future (with 1 minute tolerance for clock skew)
	if age < -60 {
		return false, fmt.Errorf("signature timestamp is in the future: %d seconds ahead", -age)
	}

	return true, nil
}

// FormatLoginMessage creates a standardized login message for signing
//
// Parameters:
//   - nonce: Unique nonce
//   - timestamp: Unix timestamp
//
// Returns:
//   - string: Formatted message for signing
//
// Example output:
//
//	Login to LegacyChain
//	Nonce: 550e8400-e29b-41d4-a716-446655440000
//	Timestamp: 1673456789
func FormatLoginMessage(nonce string, timestamp int64) string {
	return fmt.Sprintf("Login to LegacyChain\nNonce: %s\nTimestamp: %d", nonce, timestamp)
}
