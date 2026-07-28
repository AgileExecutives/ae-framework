// Package database provides serverbase startup seeding.
package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AgileExecutives/ae-framework/serverbase/internal/models"
	"github.com/AgileExecutives/ae-framework/serverbase/internal/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	initialAdminRole           = "admin"
	initialAdminTenantID       = uint(1)
	initialAdminOrganizationID = uint(1)
)

// Seed creates only the initial admin user configured through the environment.
func Seed(db *gorm.DB, tenantService *services.TenantService) error {
	return SeedWithEventBus(db, tenantService, nil)
}

// SeedWithEventBus creates only the initial admin user configured through the environment.
func SeedWithEventBus(db *gorm.DB, _ *services.TenantService, _ core.EventBus) error {
	return SeedInitialAdminUser(db)
}

// SeedInitialAdminUser creates or refreshes the initial admin user from ADMIN_USER and ADMIN_PASSWORD.
func SeedInitialAdminUser(db *gorm.DB) error {
	adminUser := strings.TrimSpace(os.Getenv("ADMIN_USER"))
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUser == "" && adminPassword == "" {
		log.Println("Skipping initial admin seed: ADMIN_USER and ADMIN_PASSWORD are not set")
		return nil
	}
	if adminUser == "" || adminPassword == "" {
		return fmt.Errorf("ADMIN_USER and ADMIN_PASSWORD must both be set to seed the initial admin user")
	}

	username, email := adminIdentity(adminUser)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash ADMIN_PASSWORD: %w", err)
	}

	var user models.User
	err = db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to look up initial admin user: %w", err)
	}

	now := time.Now()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = models.User{
			Username:        username,
			Email:           email,
			PasswordHash:    string(hashedPassword),
			FirstName:       "Admin",
			LastName:        "User",
			Role:            initialAdminRole,
			TenantID:        initialAdminTenantID,
			OrganizationID:  initialAdminOrganizationID,
			Active:          true,
			EmailVerified:   true,
			EmailVerifiedAt: &now,
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create initial admin user %s: %w", username, err)
		}
		log.Printf("Created initial admin user: %s", username)
		return nil
	}

	user.PasswordHash = string(hashedPassword)
	user.Role = initialAdminRole
	user.Active = true
	user.EmailVerified = true
	if user.EmailVerifiedAt == nil {
		user.EmailVerifiedAt = &now
	}
	if user.FirstName == "" {
		user.FirstName = "Admin"
	}
	if user.LastName == "" {
		user.LastName = "User"
	}
	if user.TenantID == 0 {
		user.TenantID = initialAdminTenantID
	}
	if user.OrganizationID == 0 {
		user.OrganizationID = initialAdminOrganizationID
	}

	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update initial admin user %s: %w", username, err)
	}
	log.Printf("Updated initial admin user: %s", username)
	return nil
}

func adminIdentity(adminUser string) (string, string) {
	if strings.Contains(adminUser, "@") {
		username := strings.TrimSpace(strings.Split(adminUser, "@")[0])
		if username == "" {
			username = "admin"
		}
		return username, adminUser
	}

	return adminUser, adminUser + "@localhost"
}
