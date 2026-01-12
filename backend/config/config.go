package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Blockchain BlockchainConfig
	JWT        JWTConfig
	RateLimit  RateLimitConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type BlockchainConfig struct {
	RpcURL              string
	WsURL               string
	ChainID             int64
	VaultFactoryAddress string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

type RateLimitConfig struct {
	Max    int
	Window time.Duration
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	chainID, _ := strconv.ParseInt(getEnv("CHAIN_ID", "1337"), 10, 64)
	rateLimitMax, _ := strconv.Atoi(getEnv("RATE_LIMIT_MAX", "100"))
	rateLimitWindow, _ := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1m"))
	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "legacychain"),
			Password: getEnv("DB_PASSWORD", "legacychain_password"),
			DBName:   getEnv("DB_NAME", "legacychain"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Blockchain: BlockchainConfig{
			RpcURL:              getEnv("BESU_RPC_URL", "http://localhost:8545"),
			WsURL:               getEnv("BESU_WS_URL", "ws://localhost:8546"),
			ChainID:             chainID,
			VaultFactoryAddress: getEnv("VAULT_FACTORY_ADDRESS", ""),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "change-me-in-production"),
			ExpiresIn: jwtExpiresIn,
		},
		RateLimit: RateLimitConfig{
			Max:    rateLimitMax,
			Window: rateLimitWindow,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
