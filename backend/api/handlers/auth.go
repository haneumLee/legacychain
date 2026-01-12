package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/haneumLee/legacychain/backend/config"
	"github.com/haneumLee/legacychain/backend/models"
	"github.com/haneumLee/legacychain/backend/pkg/crypto"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db           *gorm.DB
	cfg          *config.Config
	nonceManager *crypto.NonceManager
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		db:           db,
		cfg:          cfg,
		nonceManager: crypto.NewNonceManager(redisClient),
	}
}

type NonceResponse struct {
	Nonce     string `json:"nonce"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type LoginRequest struct {
	Address   string `json:"address" validate:"required,eth_addr"`
	Signature string `json:"signature" validate:"required"`
	Message   string `json:"message" validate:"required"`
	Nonce     string `json:"nonce" validate:"required"`
	Timestamp int64  `json:"timestamp" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// GetNonce godoc
// @Summary Get nonce for signature
// @Description Generate a nonce for signing the login message
// @Tags auth
// @Produce json
// @Success 200 {object} NonceResponse
// @Router /auth/nonce [get]
func (h *AuthHandler) GetNonce(c fiber.Ctx) error {
	ctx := context.Background()

	// Generate nonce
	nonce, timestamp, err := h.nonceManager.GenerateNonce(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate nonce",
		})
	}

	// Create login message
	message := crypto.FormatLoginMessage(nonce, timestamp)

	return c.JSON(NonceResponse{
		Nonce:     nonce,
		Message:   message,
		Timestamp: timestamp,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user with Ethereum address and signature
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	ctx := context.Background()

	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 1. Validate nonce (check if it exists and hasn't been used)
	valid, err := h.nonceManager.ValidateNonce(ctx, req.Nonce)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate nonce",
		})
	}
	if !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired nonce",
		})
	}

	// 2. Validate timestamp (not too old, not in the future)
	valid, err = crypto.ValidateTimestamp(req.Timestamp)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 3. Reconstruct the message that should have been signed
	expectedMessage := crypto.FormatLoginMessage(req.Nonce, req.Timestamp)
	if strings.TrimSpace(req.Message) != strings.TrimSpace(expectedMessage) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Message mismatch",
		})
	}

	// 4. Verify signature using EIP-191 (Ethereum Personal Sign)
	valid, err = crypto.VerifySignature(req.Address, req.Message, req.Signature)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Signature verification failed: %v", err),
		})
	}
	if !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid signature: address mismatch",
		})
	}

	// 5. Find or create user
	var user models.User
	result := h.db.Where("address = ?", req.Address).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		user = models.User{
			Address: req.Address,
		}
		if err := h.db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
	}

	// 6. Generate JWT token
	claims := jwt.MapClaims{
		"address": user.Address,
		"exp":     time.Now().Add(h.cfg.JWT.ExpiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWT.Secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(LoginResponse{
		Token: tokenString,
		User:  &user,
	})
}

// GetMe godoc
// @Summary Get current user
// @Description Get current authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/me [get]
// @Security BearerAuth
func (h *AuthHandler) GetMe(c fiber.Ctx) error {
	address := c.Locals("address").(string)

	var user models.User
	if err := h.db.Where("address = ?", address).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}
