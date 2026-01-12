package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VaultStatus string

const (
	VaultStatusLocked   VaultStatus = "locked"
	VaultStatusUnlocked VaultStatus = "unlocked"
	VaultStatusClaimed  VaultStatus = "claimed"
)

type Vault struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	VaultID           int64          `gorm:"uniqueIndex;not null" json:"vault_id"`
	ContractAddress   string         `gorm:"type:varchar(42);uniqueIndex;not null" json:"contract_address"`
	OwnerID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	Balance           string         `gorm:"type:numeric(78,0);default:0" json:"balance"`
	Status            VaultStatus    `gorm:"type:varchar(20);not null;default:'locked'" json:"status"`
	HeartbeatInterval int64          `gorm:"not null" json:"heartbeat_interval"`
	GracePeriod       int64          `gorm:"not null" json:"grace_period"`
	RequiredApprovals int            `gorm:"not null" json:"required_approvals"`
	LastHeartbeat     *time.Time     `json:"last_heartbeat,omitempty"`
	UnlockedAt        *time.Time     `json:"unlocked_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Owner      User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Heirs      []Heir      `gorm:"foreignKey:VaultID" json:"heirs,omitempty"`
	Heartbeats []Heartbeat `gorm:"foreignKey:VaultID" json:"heartbeats,omitempty"`
}

// BeforeCreate hook
func (v *Vault) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

func (Vault) TableName() string {
	return "vaults"
}
