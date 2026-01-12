package handlers

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/haneumLee/legacychain/backend/internal/service"
	"github.com/haneumLee/legacychain/backend/models"
	"gorm.io/gorm"
)

type HeartbeatHandler struct {
	db         *gorm.DB
	blockchain service.BlockchainService
}

func NewHeartbeatHandler(db *gorm.DB, blockchain service.BlockchainService) *HeartbeatHandler {
	return &HeartbeatHandler{
		db:         db,
		blockchain: blockchain,
	}
}

type CommitHeartbeatRequest struct {
	VaultID string `json:"vault_id" validate:"required,uuid"`
	Nonce   string `json:"nonce" validate:"required"` // Random value from client
}

type CommitHeartbeatResponse struct {
	TxHash     string `json:"tx_hash"`
	CommitHash string `json:"commit_hash"`
	Message    string `json:"message"`
}

type RevealHeartbeatRequest struct {
	VaultID string `json:"vault_id" validate:"required,uuid"`
	Nonce   string `json:"nonce" validate:"required"`
}

type RevealHeartbeatResponse struct {
	TxHash  string `json:"tx_hash"`
	Message string `json:"message"`
}

type HeartbeatStatusResponse struct {
	VaultID        string              `json:"vault_id"`
	LatestCommit   *models.Heartbeat   `json:"latest_commit,omitempty"`
	LastHeartbeat  *big.Int            `json:"last_heartbeat_timestamp,omitempty"`
	OnChainStatus  string              `json:"onchain_status"`
}

// CommitHeartbeat godoc
// @Summary Commit a heartbeat
// @Description Commit a heartbeat hash to the blockchain (phase 1 of commit-reveal)
// @Tags heartbeat
// @Accept json
// @Produce json
// @Param request body CommitHeartbeatRequest true "Commit request"
// @Success 200 {object} CommitHeartbeatResponse
// @Router /heartbeat/commit [post]
// @Security BearerAuth
func (h *HeartbeatHandler) CommitHeartbeat(c fiber.Ctx) error {
	address := c.Locals("address").(string)

	var req CommitHeartbeatRequest
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
	if err := h.db.Where("id = ? AND owner_address = ?", vaultID, address).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found or you don't have permission",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Generate commit hash: keccak256(abi.encodePacked(msg.sender, nonce))
	// In Go, we need to match Solidity's behavior
	nonceBytes, err := hex.DecodeString(req.Nonce)
	if err != nil {
		// If not hex, treat as string
		nonceBytes = []byte(req.Nonce)
	}
	
	senderAddr := common.HexToAddress(address)
	
	// Concatenate: address (20 bytes) + nonce
	data := append(senderAddr.Bytes(), nonceBytes...)
	commitHash := crypto.Keccak256Hash(data)
	
	// Convert to [32]byte for smart contract
	var commitHashArray [32]byte
	copy(commitHashArray[:], commitHash.Bytes())

	// Send transaction to blockchain
	txHash, err := h.blockchain.CommitHeartbeat(c.Context(), common.HexToAddress(vault.ContractAddress), commitHashArray)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to commit heartbeat: %v", err),
		})
	}

	// Save to database
	heartbeat := models.Heartbeat{
		VaultID:      vaultID,
		CommitHash:   commitHash.Hex(),
		CommitTxHash: txHash,
		Nonce:        req.Nonce,
		Status:       models.HeartbeatStatusCommitted,
		CommittedAt:  time.Now(),
	}

	if err := h.db.Create(&heartbeat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save heartbeat record",
		})
	}

	return c.JSON(CommitHeartbeatResponse{
		TxHash:     txHash,
		CommitHash: commitHash.Hex(),
		Message:    "Heartbeat committed successfully. Remember to reveal within the timeout period.",
	})
}

// RevealHeartbeat godoc
// @Summary Reveal a committed heartbeat
// @Description Reveal a previously committed heartbeat (phase 2 of commit-reveal)
// @Tags heartbeat
// @Accept json
// @Produce json
// @Param request body RevealHeartbeatRequest true "Reveal request"
// @Success 200 {object} RevealHeartbeatResponse
// @Router /heartbeat/reveal [post]
// @Security BearerAuth
func (h *HeartbeatHandler) RevealHeartbeat(c fiber.Ctx) error {
	address := c.Locals("address").(string)

	var req RevealHeartbeatRequest
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
	if err := h.db.Where("id = ? AND owner_address = ?", vaultID, address).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found or you don't have permission",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Find the latest committed heartbeat
	var heartbeat models.Heartbeat
	if err := h.db.Where("vault_id = ? AND nonce = ? AND status = ?", vaultID, req.Nonce, models.HeartbeatStatusCommitted).
		Order("committed_at DESC").
		First(&heartbeat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No matching committed heartbeat found with this nonce",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query heartbeat",
		})
	}

	// Convert nonce to [32]byte
	var nonceArray [32]byte
	nonceBytes, err := hex.DecodeString(req.Nonce)
	if err != nil {
		// If not hex, convert string to bytes
		nonceBytes = []byte(req.Nonce)
	}
	
	// Pad or truncate to 32 bytes
	if len(nonceBytes) > 32 {
		copy(nonceArray[:], nonceBytes[:32])
	} else {
		copy(nonceArray[:], nonceBytes)
	}

	// Send reveal transaction
	txHash, err := h.blockchain.RevealHeartbeat(c.Context(), common.HexToAddress(vault.ContractAddress), nonceArray)
	if err != nil {
		// Mark as failed
		heartbeat.Status = models.HeartbeatStatusFailed
		h.db.Save(&heartbeat)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to reveal heartbeat: %v", err),
		})
	}

	// Update heartbeat record
	now := time.Now()
	heartbeat.RevealTxHash = txHash
	heartbeat.Status = models.HeartbeatStatusRevealed
	heartbeat.RevealedAt = &now

	if err := h.db.Save(&heartbeat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update heartbeat record",
		})
	}

	return c.JSON(RevealHeartbeatResponse{
		TxHash:  txHash,
		Message: "Heartbeat revealed successfully",
	})
}

// GetHeartbeatStatus godoc
// @Summary Get heartbeat status
// @Description Get the current heartbeat status for a vault
// @Tags heartbeat
// @Produce json
// @Param vault_id path string true "Vault ID (UUID)"
// @Success 200 {object} HeartbeatStatusResponse
// @Router /heartbeat/status/{vault_id} [get]
// @Security BearerAuth
func (h *HeartbeatHandler) GetHeartbeatStatus(c fiber.Ctx) error {
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
	if err := h.db.Where("id = ? AND owner_address = ?", vaultID, address).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found or you don't have permission",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Get latest heartbeat from database
	var latestHeartbeat models.Heartbeat
	err = h.db.Where("vault_id = ?", vaultID).
		Order("committed_at DESC").
		First(&latestHeartbeat).Error
	
	var latestHeartbeatPtr *models.Heartbeat
	if err != gorm.ErrRecordNotFound {
		latestHeartbeatPtr = &latestHeartbeat
	}

	// Get on-chain last heartbeat timestamp
	lastHeartbeatTime, err := h.blockchain.GetLastHeartbeat(c.Context(), common.HexToAddress(vault.ContractAddress))
	var onchainStatus string
	if err != nil {
		onchainStatus = "error: " + err.Error()
	} else {
		if lastHeartbeatTime.Cmp(big.NewInt(0)) == 0 {
			onchainStatus = "No heartbeat recorded on-chain yet"
		} else {
			onchainStatus = fmt.Sprintf("Last heartbeat: %s", time.Unix(lastHeartbeatTime.Int64(), 0).Format(time.RFC3339))
		}
	}

	return c.JSON(HeartbeatStatusResponse{
		VaultID:       vaultIDStr,
		LatestCommit:  latestHeartbeatPtr,
		LastHeartbeat: lastHeartbeatTime,
		OnChainStatus: onchainStatus,
	})
}

// ListHeartbeats godoc
// @Summary List heartbeats for a vault
// @Description Get all heartbeat records for a vault
// @Tags heartbeat
// @Produce json
// @Param vault_id path string true "Vault ID (UUID)"
// @Success 200 {array} models.Heartbeat
// @Router /heartbeat/list/{vault_id} [get]
// @Security BearerAuth
func (h *HeartbeatHandler) ListHeartbeats(c fiber.Ctx) error {
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
	if err := h.db.Where("id = ? AND owner_address = ?", vaultID, address).First(&vault).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Vault not found or you don't have permission",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query vault",
		})
	}

	// Get all heartbeats
	var heartbeats []models.Heartbeat
	if err := h.db.Where("vault_id = ?", vaultID).
		Order("committed_at DESC").
		Find(&heartbeats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query heartbeats",
		})
	}

	return c.JSON(heartbeats)
}
