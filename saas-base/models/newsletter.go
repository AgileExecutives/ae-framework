package models

import (
	"time"

	"gorm.io/gorm"
)

// Newsletter represents a newsletter subscription.
type Newsletter struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Email       string         `json:"email" gorm:"not null;uniqueIndex"`
	Interest    string         `json:"interest" gorm:"default:'general'"`
	Source      string         `json:"source" gorm:"not null"`
	LastContact time.Time      `json:"last_contact" gorm:"autoUpdateTime"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}

func (Newsletter) TableName() string { return "newsletters" }

// NewsletterResponse is the API representation of a Newsletter subscription.
type NewsletterResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Interest    string    `json:"interest"`
	Source      string    `json:"source"`
	LastContact time.Time `json:"last_contact"`
	CreatedAt   time.Time `json:"created_at"`
}

func (n *Newsletter) ToResponse() NewsletterResponse {
	return NewsletterResponse{
		ID:          n.ID,
		Name:        n.Name,
		Email:       n.Email,
		Interest:    n.Interest,
		Source:      n.Source,
		LastContact: n.LastContact,
		CreatedAt:   n.CreatedAt,
	}
}

// NewsletterSubscribeRequest is used to subscribe to the newsletter.
type NewsletterSubscribeRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Interest string `json:"interest"`
	Source   string `json:"source" binding:"required"`
}

// NewsletterUnsubscribeRequest is used to unsubscribe from the newsletter.
type NewsletterUnsubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}
