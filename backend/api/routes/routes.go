package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haneumLee/legacychain/backend/api/handlers"
	"github.com/haneumLee/legacychain/backend/api/middleware"
	"github.com/haneumLee/legacychain/backend/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, redisClient *redis.Client, cfg *config.Config) {
	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "legacychain-backend",
		})
	})

	// API v1 group
	api := app.Group("/api/v1")

	// Apply rate limiter to all API routes
	api.Use(middleware.RateLimiter(cfg, redisClient))

	// Auth routes (no JWT required)
	authHandler := handlers.NewAuthHandler(db, cfg)
	auth := api.Group("/auth")
	{
		auth.Post("/login", authHandler.Login)
		auth.Get("/me", middleware.JWTAuth(cfg), authHandler.GetMe)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg))

	// Vault routes
	vaultHandler := handlers.NewVaultHandler(db)
	vaults := protected.Group("/vaults")
	{
		vaults.Post("", vaultHandler.CreateVault)
		vaults.Get("", vaultHandler.ListVaults)
		vaults.Get("/:id", vaultHandler.GetVault)
	}
}
