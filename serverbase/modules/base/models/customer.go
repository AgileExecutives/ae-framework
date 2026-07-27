package models

import (
    "time"

    "gorm.io/gorm"
)

// Customer represents a customer account linked to a tenant and plan.
type Customer struct {
    ID            uint           `gorm:"primarykey" json:"id"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
    Name          string         `gorm:"not null" json:"name" binding:"required"`
    Email         string         `gorm:"not null" json:"email" binding:"required,email"`
    Phone         string         `json:"phone"`
    Street        string         `json:"street"`
    Zip           string         `json:"zip"`
    City          string         `json:"city"`
    Country       string         `json:"country"`
    TaxID         string         `json:"tax_id"`
    VAT           string         `json:"vat"`
    PlanID        uint           `gorm:"not null" json:"plan_id" binding:"required"`
    TenantID      uint           `gorm:"not null" json:"tenant_id"`
    Status        string         `gorm:"default:'active'" json:"status"`
    PaymentMethod string         `json:"payment_method"`
    Active        bool           `gorm:"default:true" json:"active"`
}

func (Customer) TableName() string { return "customers" }

func (c *Customer) ToResponse() CustomerResponse {
    return CustomerResponse{
        ID:          int(c.ID),
        Name:        c.Name,
        Email:       c.Email,
        Phone:       c.Phone,
        CompanyName: "",
        CreatedAt:   c.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
    }
}

// API responses and requests exist in shared shapes; reuse types where possible.
