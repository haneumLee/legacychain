package handlers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/haneumLee/legacychain/backend/internal/service"
	"github.com/haneumLee/legacychain/backend/models"
	"gorm.io/gorm"
)

type HeirHandler struct {
	db         *gorm.DB
	blockchain service.BlockchainService
}

func NewHeirHandler(db *gorm.DB, blockchain service.BlockchainService) *HeirHandler {
	return &HeirHandler{
		db:         db,
		blockchain: blockchain,
	}
}

type ApproveHeirRequest struct {
	VaultID string `json:"vault_id" validate:"required,uuid"`
}

type ApproveHeirResponse struct {
	TxHash  string `json:"tx_hash"`
	Message string `json:"message"`
}

type ClaimInheritanceRequest struct {
	VaultID string `json:"vault_id" validate:"required,uuid"`
}

type ClaimInheritanceResponse struct {
	TxHash  string `json:"tx_hash"`
	Message string `json:"message"`
}

type HeirApprovalStatus struct {
	VaultID        string   `json:"vault_id"`
	HeirAddress    string   `json:"heir_address"`
	ApprovalCount  string   `json:"approval_count"`
	RequiredCount  int      `json:"required_count"`
	CanClaim       bool     `json:"can_claim"`
}

// ApproveHeir godoc
// @Summary Approve inheritance
// @Description As an heir, approve the inheritance claim (multi-sig approval)
// @Tags heir
// @Accept json
// @Produce json
// @Param request body ApproveHeirRequest true "Approve request"
// @Success 200 {object} ApproveHeirResponse
// @Router /heir/approve [post]
// @Security BearerAuth
func (h *HeirHandler) ApproveHeir(c fiber.Ctx) error {
	address := c.Locals("address").(string)

	var req ApproveHeirRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse vault ID
	vaultID, err := uuid.Parse(req.VaultID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vault ID format",
		})
	}

	// Find vault
	var vault models.Vault
	if err := h.db.Where("id = ?", vaultID).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Check if the caller is an heir of this vault
	var heir models.Heir
	if err := h.db.Where("vault_id = ? AND address = ?", vaultID, address).First(&heir).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You are not an heir of this vault",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify heir status",
		})
	}

	// Send approval transaction
	txHash, err := h.blockchain.ApproveInheritance(
		c.Context(),
		common.HexToAddress(vault.ContractAddress),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to approve inheritance: %v", err),
		})
	}

	return c.JSON(ApproveHeirResponse{
		TxHash:  txHash,
		Message: "Inheritance approved successfully",
	})
}

// ClaimInheritance godoc
// @Summary Claim inheritance as an heir
// @Description Claim the inheritance from a vault (requires sufficient approvals from other heirs)
// @Tags heir
// @Accept json
// @Produce json
// @Param request body ClaimInheritanceRequest true "Claim request"
// @Success 200 {object} ClaimInheritanceResponse
// @Router /heir/claim [post]
// @Security BearerAuth
func (h *HeirHandler) ClaimInheritance(c fiber.Ctx) error {
	address := c.Locals("address").(string)

	var req ClaimInheritanceRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse vault ID
	vaultID, err := uuid.Parse(req.VaultID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vault ID format",
		})
	}

	// Find vault
	var vault models.Vault
	if err := h.db.Where("id = ?", vaultID).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Check if the caller is an heir
	var heir models.Heir
	if err := h.db.Where("vault_id = ? AND address = ?", vaultID, address).First(&heir).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You are not an heir of this vault",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify heir status",
		})
	}

	// Check approval status from blockchain (using vault config approval count)
	vaultConfig, err := h.blockchain.GetVaultConfig(c.Context(), common.HexToAddress(vault.ContractAddress))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get vault config: %v", err),
		})
	}

	// Get total number of heirs
	var heirCount int64
	if err := h.db.Model(&models.Heir{}).Where("vault_id = ?", vaultID).Count(&heirCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count heirs",
		})
	}

	// Check if caller has already approved
	hasApproved, err := h.blockchain.GetHeirApprovalStatus(c.Context(), common.HexToAddress(vault.ContractAddress), common.HexToAddress(address))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to check approval status: %v", err),
		})
	}

	approvalCount := vaultConfig.ApprovalCount

	// Calculate required approvals (majority: > 50%)
	requiredApprovals := (heirCount / 2) + 1

	if approvalCount.Int64() < requiredApprovals {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Insufficient approvals: %d/%d (need %d). You have%s approved.", approvalCount.Int64(), heirCount, requiredApprovals, map[bool]string{true: "", false: " not"}[hasApproved]),
		})
	}

	// Send claim transaction
	txHash, err := h.blockchain.ClaimInheritance(
		c.Context(),
		common.HexToAddress(vault.ContractAddress),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to claim inheritance: %v", err),
		})
	}

	// Update vault status
	vault.Status = "claimed"
	if err := h.db.Save(&vault).Error; err != nil {
		// Log error but don't fail the response since blockchain transaction succeeded
		fmt.Printf("Warning: Failed to update vault status: %v\n", err)
	}

	return c.JSON(ClaimInheritanceResponse{
		TxHash:  txHash,
		Message: "Inheritance claimed successfully",
	})
}

// GetApprovalStatus godoc
// @Summary Get heir approval status
// @Description Get the current approval status for an heir to claim inheritance
// @Tags heir
// @Produce json
// @Param vault_id path string true "Vault ID (UUID)"
// @Success 200 {object} HeirApprovalStatus
// @Router /heir/status/{vault_id} [get]
// @Security BearerAuth
func (h *HeirHandler) GetApprovalStatus(c fiber.Ctx) error {
	address := c.Locals("address").(string)
	vaultIDStr := c.Params("vault_id")

	// Parse vault ID
	vaultID, err := uuid.Parse(vaultIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vault ID format",
		})
	}

	// Find vault
	var vault models.Vault
	if err := h.db.Where("id = ?", vaultID).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Check if the caller is an heir
	var heir models.Heir
	if err := h.db.Where("vault_id = ? AND address = ?", vaultID, address).First(&heir).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You are not an heir of this vault",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify heir status",
		})
	}

	// Get vault config from blockchain
	vaultConfig, err := h.blockchain.GetVaultConfig(c.Context(), common.HexToAddress(vault.ContractAddress))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get vault config: %v", err),
		})
	}

	// Check if caller has approved
	hasApproved, err := h.blockchain.GetHeirApprovalStatus(c.Context(), common.HexToAddress(vault.ContractAddress), common.HexToAddress(address))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get approval status: %v", err),
		})
	}

	approvalCount := vaultConfig.ApprovalCount

	// Get total number of heirs
	var heirCount int64
	if err := h.db.Model(&models.Heir{}).Where("vault_id = ?", vaultID).Count(&heirCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count heirs",
		})
	}

	// Calculate required approvals
	requiredApprovals := int((heirCount / 2) + 1)
	canClaim := approvalCount.Int64() >= int64(requiredApprovals)

	return c.JSON(HeirApprovalStatus{
		VaultID:       vaultIDStr,
		HeirAddress:   address,
		ApprovalCount: fmt.Sprintf("%d (You have%s approved)", approvalCount.Int64(), map[bool]string{true: "", false: " not"}[hasApproved]),
		RequiredCount: requiredApprovals,
		CanClaim:      canClaim,
	})
}

// ListHeirs godoc
// @Summary List all heirs for a vault
// @Description Get all heirs and their shares for a specific vault
// @Tags heir
// @Produce json
// @Param vault_id path string true "Vault ID (UUID)"
// @Success 200 {array} models.Heir
// @Router /heir/list/{vault_id} [get]
// @Security BearerAuth
func (h *HeirHandler) ListHeirs(c fiber.Ctx) error {
	vaultIDStr := c.Params("vault_id")

	// Parse vault ID
	vaultID, err := uuid.Parse(vaultIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vault ID format",
		})
	}

	// Find vault
	var vault models.Vault
	if err := h.db.Where("id = ?", vaultID).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Get all heirs
	var heirs []models.Heir
	if err := h.db.Where("vault_id = ?", vaultID).Find(&heirs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query heirs",
		})
	}

	return c.JSON(heirs)
}
