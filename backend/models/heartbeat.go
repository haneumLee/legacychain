package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Heartbeat struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	VaultID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"vault_id"`
	TxHash    string         `gorm:"type:varchar(66);uniqueIndex;not null" json:"tx_hash"`
	Timestamp time.Time      `gorm:"not null;index" json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

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
