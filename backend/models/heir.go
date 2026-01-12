package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Heir struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	VaultID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"vault_id"`
	Address     string         `gorm:"type:varchar(42);not null;index" json:"address"`
	ShareBPS    int            `gorm:"not null" json:"share_bps"` // Basis points (0-10000)
	HasApproved bool           `gorm:"default:false" json:"has_approved"`
	HasClaimed  bool           `gorm:"default:false" json:"has_claimed"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Vault Vault `gorm:"foreignKey:VaultID" json:"vault,omitempty"`
}

func (h *Heir) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

func (Heir) TableName() string {
	return "heirs"
}
