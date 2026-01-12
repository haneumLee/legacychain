package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/haneumLee/legacychain/backend/config"
	"github.com/haneumLee/legacychain/backend/models"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type LoginRequest struct {
	Address   string `json:"address" validate:"required,eth_addr"`
	Signature string `json:"signature" validate:"required"`
	Message   string `json:"message" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
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
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Verify signature (implement signature verification)
	// For now, we'll trust the address

	// Find or create user
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

	// Generate JWT token
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
