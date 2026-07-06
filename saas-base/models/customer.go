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
		ID:            c.ID,
		Name:          c.Name,
		Email:         c.Email,
		Phone:         c.Phone,
		Street:        c.Street,
		Zip:           c.Zip,
		City:          c.City,
		Country:       c.Country,
		TaxID:         c.TaxID,
		VAT:           c.VAT,
		PlanID:        c.PlanID,
		TenantID:      c.TenantID,
		Status:        c.Status,
		PaymentMethod: c.PaymentMethod,
		Active:        c.Active,
		CreatedAt:     c.CreatedAt,
	}
}

// CustomerResponse is the API representation of a Customer.
type CustomerResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Street        string    `json:"street"`
	Zip           string    `json:"zip"`
	City          string    `json:"city"`
	Country       string    `json:"country"`
	TaxID         string    `json:"tax_id"`
	VAT           string    `json:"vat"`
	PlanID        uint      `json:"plan_id"`
	TenantID      uint      `json:"tenant_id"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
}

// CustomerRequest is used for swagger docs.
type CustomerRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone"`
	Street        string `json:"street"`
	Zip           string `json:"zip"`
	City          string `json:"city"`
	Country       string `json:"country"`
	TaxID         string `json:"tax_id"`
	VAT           string `json:"vat"`
	PlanID        uint   `json:"plan_id" binding:"required"`
	TenantID      uint   `json:"tenant_id" binding:"required"`
	PaymentMethod string `json:"payment_method"`
}

// CustomerCreateRequest is used for customer creation.
type CustomerCreateRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone"`
	Street        string `json:"street"`
	Zip           string `json:"zip"`
	City          string `json:"city"`
	Country       string `json:"country"`
	TaxID         string `json:"tax_id"`
	VAT           string `json:"vat"`
	PlanID        uint   `json:"plan_id" binding:"required"`
	TenantID      uint   `json:"tenant_id"`
	PaymentMethod string `json:"payment_method"`
}

// CustomerUpdateRequest is used for customer updates.
type CustomerUpdateRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email" binding:"omitempty,email"`
	Phone         string `json:"phone"`
	Street        string `json:"street"`
	Zip           string `json:"zip"`
	City          string `json:"city"`
	Country       string `json:"country"`
	TaxID         string `json:"tax_id"`
	VAT           string `json:"vat"`
	PlanID        *uint  `json:"plan_id"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	Active        *bool  `json:"active"`
}
