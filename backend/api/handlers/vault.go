package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/haneumLee/legacychain/backend/models"
	"gorm.io/gorm"
)

type VaultHandler struct {
	db *gorm.DB
}

func NewVaultHandler(db *gorm.DB) *VaultHandler {
	return &VaultHandler{db: db}
}

type CreateVaultRequest struct {
	VaultID           int64    `json:"vault_id" validate:"required"`
	ContractAddress   string   `json:"contract_address" validate:"required,eth_addr"`
	HeartbeatInterval int64    `json:"heartbeat_interval" validate:"required,min=1"`
	GracePeriod       int64    `json:"grace_period" validate:"required,min=1"`
	RequiredApprovals int      `json:"required_approvals" validate:"required,min=1"`
	HeirAddresses     []string `json:"heir_addresses" validate:"required,min=1"`
	HeirShares        []int    `json:"heir_shares" validate:"required,min=1"`
}

// CreateVault godoc
// @Summary Create new vault
// @Description Create a new vault after on-chain deployment
// @Tags vaults
// @Accept json
// @Produce json
// @Param request body CreateVaultRequest true "Vault data"
// @Success 201 {object} models.Vault
// @Router /vaults [post]
// @Security BearerAuth
func (h *VaultHandler) CreateVault(c fiber.Ctx) error {
	address := c.Locals("address").(string)

	var req CreateVaultRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate heir arrays length
	if len(req.HeirAddresses) != len(req.HeirShares) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Heir addresses and shares must have same length",
		})
	}

	// Find user
	var user models.User
	if err := h.db.Where("address = ?", address).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Create vault
	vault := models.Vault{
		VaultID:           req.VaultID,
		ContractAddress:   req.ContractAddress,
		OwnerID:           user.ID,
		HeartbeatInterval: req.HeartbeatInterval,
		GracePeriod:       req.GracePeriod,
		RequiredApprovals: req.RequiredApprovals,
		Status:            models.VaultStatusLocked,
	}

	// Start transaction
	tx := h.db.Begin()

	if err := tx.Create(&vault).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create vault",
		})
	}

	// Create heirs
	for i, heirAddr := range req.HeirAddresses {
		heir := models.Heir{
			VaultID:  vault.ID,
			Address:  heirAddr,
			ShareBPS: req.HeirShares[i],
		}
		if err := tx.Create(&heir).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create heirs",
			})
		}
	}

	tx.Commit()

	// Load relationships
	h.db.Preload("Owner").Preload("Heirs").First(&vault, vault.ID)

	return c.Status(fiber.StatusCreated).JSON(vault)
}

// GetVault godoc
// @Summary Get vault by ID
// @Description Get vault details with heirs and heartbeats
// @Tags vaults
// @Produce json
// @Param id path string true "Vault UUID"
// @Success 200 {object} models.Vault
// @Router /vaults/{id} [get]
// @Security BearerAuth
func (h *VaultHandler) GetVault(c fiber.Ctx) error {
	vaultID := c.Params("id")

	uid, err := uuid.Parse(vaultID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vault ID",
		})
	}

	var vault models.Vault
	if err := h.db.Preload("Owner").Preload("Heirs").Preload("Heartbeats").
		First(&vault, uid).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Vault not found",
		})
	}

	return c.JSON(vault)
}

// ListVaults godoc
// @Summary List user's vaults
// @Description Get all vaults owned by the authenticated user
// @Tags vaults
// @Produce json
// @Success 200 {array} models.Vault
// @Router /vaults [get]
// @Security BearerAuth
func (h *VaultHandler) ListVaults(c fiber.Ctx) error {
	address := c.Locals("address").(string)

	var user models.User
	if err := h.db.Where("address = ?", address).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var vaults []models.Vault
	if err := h.db.Where("owner_id = ?", user.ID).
		Preload("Heirs").
		Find(&vaults).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch vaults",
		})
	}

	return c.JSON(vaults)
}
