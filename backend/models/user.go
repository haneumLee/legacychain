package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Address   string         `gorm:"type:varchar(42);uniqueIndex;not null" json:"address"`
	Email     string         `gorm:"type:varchar(255);index" json:"email,omitempty"`
	Nickname  string         `gorm:"type:varchar(100)" json:"nickname,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Vaults []Vault `gorm:"foreignKey:OwnerID" json:"vaults,omitempty"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}
