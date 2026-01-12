package crypto

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis creates an in-memory Redis server for testing
func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	
	return client, mr
}

// TestGenerateNonce tests nonce generation and storage
func TestGenerateNonce(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	
	nm := NewNonceManager(client)
	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Generate valid nonce",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, timestamp, err := nm.GenerateNonce(ctx)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify nonce is not empty (UUID format)
				assert.NotEmpty(t, nonce)
				assert.Len(t, nonce, 36) // UUID v4 length with hyphens
				
				// Verify timestamp is recent
				now := time.Now().Unix()
				assert.InDelta(t, now, timestamp, 2) // Within 2 seconds
				
				// Verify nonce is stored in Redis
				key := NonceKeyPrefix + nonce
				exists, err := client.Exists(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), exists)
				
				// Verify TTL is set correctly
				ttl, err := client.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.Greater(t, ttl, 4*time.Minute) // At least 4 minutes left
				assert.LessOrEqual(t, ttl, NonceExpiration)
			}
		})
	}
}

// TestGenerateNonce_Uniqueness tests that generated nonces are unique
func TestGenerateNonce_Uniqueness(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	
	nm := NewNonceManager(client)
	ctx := context.Background()
	
	nonces := make(map[string]bool)
	iterations := 100
	
	for i := 0; i < iterations; i++ {
		nonce, _, err := nm.GenerateNonce(ctx)
		assert.NoError(t, err)
		
		// Verify nonce is unique
		assert.False(t, nonces[nonce], "Duplicate nonce generated: %s", nonce)
		nonces[nonce] = true
	}
	
	assert.Len(t, nonces, iterations)
}

// TestValidateNonce tests nonce validation and one-time use
func TestValidateNonce(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	
	nm := NewNonceManager(client)
	ctx := context.Background()
	
	// Generate a valid nonce
	validNonce, _, err := nm.GenerateNonce(ctx)
	require.NoError(t, err)

	tests := []struct {
		name      string
		nonce     string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "Valid nonce (first use)",
			nonce:     validNonce,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Invalid nonce (non-existent)",
			nonce:     "non-existent-nonce",
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := nm.ValidateNonce(ctx, tt.nonce)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValid, valid)
			}
		})
	}
	
	// Test one-time use (replay attack prevention)
	t.Run("Nonce reuse should fail (replay attack)", func(t *testing.T) {
		valid, err := nm.ValidateNonce(ctx, validNonce)
		assert.NoError(t, err)
		assert.False(t, valid, "Nonce should not be valid after first use")
	})
}

// TestValidateNonce_Expiration tests that expired nonces are automatically removed
func TestValidateNonce_Expiration(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	
	nm := NewNonceManager(client)
	ctx := context.Background()
	
	// Generate a nonce
	nonce, _, err := nm.GenerateNonce(ctx)
	require.NoError(t, err)
	
	// Fast-forward time in miniredis
	mr.FastForward(NonceExpiration + time.Second)
	
	// Validate nonce (should be expired)
	valid, err := nm.ValidateNonce(ctx, nonce)
	assert.NoError(t, err)
	assert.False(t, valid, "Expired nonce should not be valid")
}

// TestValidateTimestamp tests timestamp validation with clock skew tolerance
func TestValidateTimestamp(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name      string
		timestamp int64
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "Current timestamp (valid)",
			timestamp: now,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Recent timestamp (1 minute ago)",
			timestamp: now - 60,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Old timestamp (10 minutes ago - expired)",
			timestamp: now - 600,
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "Future timestamp within tolerance (30 seconds ahead)",
			timestamp: now + 30,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Future timestamp beyond tolerance (2 minutes ahead)",
			timestamp: now + 120,
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "Very old timestamp (1 hour ago)",
			timestamp: now - 3600,
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := ValidateTimestamp(tt.timestamp)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, valid)
			} else {
				assert.NoError(t, err)
				assert.True(t, valid)
			}
		})
	}
}

// TestFormatLoginMessage tests login message formatting
func TestFormatLoginMessage(t *testing.T) {
	tests := []struct {
		name      string
		nonce     string
		timestamp int64
		wantMsg   string
	}{
		{
			name:      "Standard login message",
			nonce:     "test-nonce-123",
			timestamp: 1673456789,
			wantMsg:   "Login to LegacyChain\nNonce: test-nonce-123\nTimestamp: 1673456789",
		},
		{
			name:      "UUID nonce",
			nonce:     "550e8400-e29b-41d4-a716-446655440000",
			timestamp: 1234567890,
			wantMsg:   "Login to LegacyChain\nNonce: 550e8400-e29b-41d4-a716-446655440000\nTimestamp: 1234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := FormatLoginMessage(tt.nonce, tt.timestamp)
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}

// Benchmark nonce generation
func BenchmarkGenerateNonce(b *testing.B) {
	mr, err := miniredis.Run()
	if err != nil {
		b.Fatal(err)
	}
	defer mr.Close()
	
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	nm := NewNonceManager(client)
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nm.GenerateNonce(ctx)
	}
}

// Benchmark nonce validation
func BenchmarkValidateNonce(b *testing.B) {
	mr, err := miniredis.Run()
	if err != nil {
		b.Fatal(err)
	}
	defer mr.Close()
	
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	nm := NewNonceManager(client)
	ctx := context.Background()
	
	// Pre-generate nonces
	nonces := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		nonce, _, _ := nm.GenerateNonce(ctx)
		nonces[i] = nonce
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nm.ValidateNonce(ctx, nonces[i])
	}
}
