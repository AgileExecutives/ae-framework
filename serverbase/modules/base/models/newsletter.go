package models

import (
	"time"
	"gorm.io/gorm"
)

// Newsletter represents a simple newsletter subscription
type Newsletter struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
	Email         string         `gorm:"not null" json:"email" binding:"required,email"`
	Name          string         `json:"name"`
	TenantID      uint           `json:"tenant_id"`
	Interest      string         `json:"interest"`
	Source        string         `json:"source"`
	LastContact   time.Time      `json:"last_contact"`
}

func (Newsletter) TableName() string { return "newsletters" }
