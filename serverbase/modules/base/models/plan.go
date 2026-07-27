package models

import (
    "time"

    "gorm.io/gorm"
)

// Plan represents a subscription plan in the system.
type Plan struct {
    ID            uint           `gorm:"primarykey" json:"id"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
    Name          string         `gorm:"not null" json:"name" binding:"required"`
    Slug          string         `gorm:"not null;uniqueIndex" json:"slug" binding:"required"`
    Description   string         `json:"description"`
    Price         float64        `gorm:"not null" json:"price" binding:"required"`
    Currency      string         `gorm:"not null;default:'EUR'" json:"currency"`
    InvoicePeriod string         `gorm:"not null;default:'monthly'" json:"invoice_period"`
    MaxUsers      int            `gorm:"default:10" json:"max_users"`
    MaxClients    int            `gorm:"default:100" json:"max_clients"`
    Features      string         `gorm:"type:text" json:"features"`
    Active        bool           `gorm:"default:true" json:"active"`
}

func (Plan) TableName() string { return "plans" }

func (p *Plan) ToResponse() PlanResponse {
    return PlanResponse{
        ID:            p.ID,
        Name:          p.Name,
        Slug:          p.Slug,
        Description:   p.Description,
        Price:         p.Price,
        Currency:      p.Currency,
        InvoicePeriod: p.InvoicePeriod,
        MaxUsers:      p.MaxUsers,
        MaxClients:    p.MaxClients,
        Features:      p.Features,
        Active:        p.Active,
        CreatedAt:     p.CreatedAt,
    }
}

// PlanResponse is the API representation of a Plan.
type PlanResponse struct {
    ID            uint      `json:"id"`
    Name          string    `json:"name"`
    Slug          string    `json:"slug"`
    Description   string    `json:"description"`
    Price         float64   `json:"price"`
    Currency      string    `json:"currency"`
    InvoicePeriod string    `json:"invoice_period"`
    MaxUsers      int       `json:"max_users"`
    MaxClients    int       `json:"max_clients"`
    Features      string    `json:"features"`
    Active        bool      `json:"active"`
    CreatedAt     time.Time `json:"created_at"`
}

// PlanRequest is used for create/update operations.
type PlanRequest struct {
    Name          string  `json:"name" binding:"required"`
    Slug          string  `json:"slug" binding:"required"`
    Description   string  `json:"description"`
    Price         float64 `json:"price" binding:"required"`
    Currency      string  `json:"currency"`
    InvoicePeriod string  `json:"invoice_period"`
    MaxUsers      int     `json:"max_users"`
    MaxClients    int     `json:"max_clients"`
    Features      string  `json:"features"`
    Active        *bool   `json:"active"`
}
