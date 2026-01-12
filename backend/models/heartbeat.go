package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HeartbeatStatus string

const (
	HeartbeatStatusCommitted HeartbeatStatus = "committed"
	HeartbeatStatusRevealed  HeartbeatStatus = "revealed"
	HeartbeatStatusFailed    HeartbeatStatus = "failed"
)

type Heartbeat struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	VaultID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"vault_id"`
	CommitHash  string          `gorm:"type:varchar(66)" json:"commit_hash"`  // Hash of the commit
	CommitTxHash string         `gorm:"type:varchar(66)" json:"commit_tx_hash"` // Commit transaction hash
	RevealTxHash string         `gorm:"type:varchar(66)" json:"reveal_tx_hash"` // Reveal transaction hash
	Nonce       string          `gorm:"type:varchar(100)" json:"nonce"`       // Nonce used for commit
	Status      HeartbeatStatus `gorm:"type:varchar(20);not null;default:'committed'" json:"status"`
	CommittedAt time.Time       `gorm:"index" json:"committed_at"`
	RevealedAt  *time.Time      `json:"revealed_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relationships
	Vault Vault `gorm:"foreignKey:VaultID" json:"vault,omitempty"`
}

func (h *Heartbeat) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

func (Heartbeat) TableName() string {
	return "heartbeats"
}
