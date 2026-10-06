package seed

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	internalmodels "github.com/AgileExecutives/ae-framework/serverbase/internal/models"
	templateentities "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/entities"
	pkgmodels "github.com/AgileExecutives/ae-framework/serverbase/pkg/models"
	saasmodels "github.com/AgileExecutives/ae-framework/shared-modules/saas-base/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	defaultAdminEmail    = "testuser@unburdy.de"
	defaultAdminPassword = "newpass123"
	defaultTenantName    = "Server Test Tenant"
	defaultTenantSlug    = "server-test"
	defaultOrgName       = "Server Test Organization"
)

// RunIfEmpty bootstraps the server-test database only when no users exist.
func RunIfEmpty(db *gorm.DB) error {
	var userCount int64
	if err := db.Model(&internalmodels.User{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to check users table: %w", err)
	}
	if userCount > 0 {
		log.Printf("server-test seed skipped: database already has %d user(s)", userCount)
		return nil
	}

	if err := db.AutoMigrate(&templateentities.Template{}, &templateentities.TemplateContract{}); err != nil {
		return fmt.Errorf("failed to migrate server-test template tables: %w", err)
	}

	plan, err := seedPlans(db)
	if err != nil {
		return err
	}
	customer, err := seedCustomer(db, plan.ID)
	if err != nil {
		return err
	}
	tenant, err := seedTenant(db, customer.ID)
	if err != nil {
		return err
	}
	organization, err := seedOrganization(db, tenant.ID)
	if err != nil {
		return err
	}
	user, password, err := seedAdminUser(db, tenant.ID, organization.ID)
	if err != nil {
		return err
	}
	if err := seedEmailTemplates(db, tenant.ID, organization.ID); err != nil {
		return err
	}

	log.Println("--- Server Test Seed ---")
	log.Printf("User: %s", user.Email)
	log.Printf("Password: %s", maskPassword(password))
	log.Printf("TenantID: %d", tenant.ID)
	log.Printf("OrganizationID: %d", organization.ID)
	log.Println("------------------------")
	return nil
}

func seedPlans(db *gorm.DB) (saasmodels.Plan, error) {
	plans := []saasmodels.Plan{
		{Name: "Free", Slug: "free", Description: "Free tier", Price: 0, Currency: "EUR", InvoicePeriod: "monthly", MaxUsers: 3, MaxClients: 10, Features: `{"tier":"free"}`, Active: true},
		{Name: "Pro", Slug: "pro", Description: "Pro tier", Price: 29, Currency: "EUR", InvoicePeriod: "monthly", MaxUsers: 20, MaxClients: 500, Features: `{"tier":"pro"}`, Active: true},
	}

	for _, plan := range plans {
		var planCount int64
		if err := db.Model(&saasmodels.Plan{}).Where("slug = ?", plan.Slug).Count(&planCount).Error; err != nil {
			return saasmodels.Plan{}, fmt.Errorf("failed to look up plan %s: %w", plan.Slug, err)
		}
		if planCount == 0 {
			if err := db.Create(&plan).Error; err != nil {
				return saasmodels.Plan{}, fmt.Errorf("failed to create plan %s: %w", plan.Slug, err)
			}
		}
	}

	var freePlan saasmodels.Plan
	if err := db.Where("slug = ?", "free").First(&freePlan).Error; err != nil {
		return saasmodels.Plan{}, fmt.Errorf("failed to load seeded free plan: %w", err)
	}
	return freePlan, nil
}

func seedCustomer(db *gorm.DB, planID uint) (saasmodels.Customer, error) {
	customer := saasmodels.Customer{
		Name:     "Server Test Customer",
		Email:    "server-test@example.com",
		PlanID:   planID,
		TenantID: 1,
		Status:   "active",
		Active:   true,
	}
	if err := db.Create(&customer).Error; err != nil {
		return saasmodels.Customer{}, fmt.Errorf("failed to create server-test customer: %w", err)
	}
	return customer, nil
}

func seedTenant(db *gorm.DB, customerID uint) (internalmodels.Tenant, error) {
	tenant := internalmodels.Tenant{
		CustomerID: customerID,
		Name:       defaultTenantName,
		Slug:       defaultTenantSlug,
	}
	if err := db.Create(&tenant).Error; err != nil {
		return internalmodels.Tenant{}, fmt.Errorf("failed to create server-test tenant: %w", err)
	}
	return tenant, nil
}

func seedOrganization(db *gorm.DB, tenantID uint) (pkgmodels.Organization, error) {
	organization := pkgmodels.Organization{
		TenantID: tenantID,
		Name:     defaultOrgName,
		Email:    "server-test@example.com",
	}
	if err := db.Create(&organization).Error; err != nil {
		return pkgmodels.Organization{}, fmt.Errorf("failed to create server-test organization: %w", err)
	}
	return organization, nil
}

func seedAdminUser(db *gorm.DB, tenantID, organizationID uint) (internalmodels.User, string, error) {
	adminUser := strings.TrimSpace(os.Getenv("ADMIN_USER"))
	password := os.Getenv("ADMIN_PASSWORD")
	if adminUser == "" {
		adminUser = defaultAdminEmail
	}
	if password == "" {
		password = defaultAdminPassword
	}

	username, email := adminIdentity(adminUser)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return internalmodels.User{}, "", fmt.Errorf("failed to hash server-test admin password: %w", err)
	}
	now := time.Now()
	user := internalmodels.User{
		Username:        username,
		Email:           email,
		PasswordHash:    string(hashedPassword),
		FirstName:       "Test",
		LastName:        "User",
		TenantID:        tenantID,
		OrganizationID:  organizationID,
		Role:            "admin",
		Active:          true,
		EmailVerified:   true,
		EmailVerifiedAt: &now,
	}
	if err := db.Create(&user).Error; err != nil {
		return internalmodels.User{}, "", fmt.Errorf("failed to create server-test admin user: %w", err)
	}
	return user, password, nil
}

func seedEmailTemplates(db *gorm.DB, tenantID, organizationID uint) error {
	organizationIDPtr := organizationID
	templates := []templateentities.Template{
		{
			TenantID:       tenantID,
			OrganizationID: &organizationIDPtr,
			Module:         "user",
			TemplateKey:    "welcome",
			Channel:        templateentities.ChannelEmail,
			Name:           "Welcome Email",
			Description:    "Default welcome email for server-test",
			StorageKey:     "server-test/templates/welcome.html",
			Version:        1,
			IsActive:       true,
			IsDefault:      true,
			Variables:      datatypes.JSON([]byte(`["FirstName","LastName","OrganizationName","ActivationLink"]`)),
			SampleData:     datatypes.JSON([]byte(`{"FirstName":"Test","LastName":"User","OrganizationName":"Server Test Organization","ActivationLink":"http://localhost:5173/activate"}`)),
			TemplateType:   "email",
		},
		{
			TenantID:       tenantID,
			OrganizationID: &organizationIDPtr,
			Module:         "user",
			TemplateKey:    "password_reset",
			Channel:        templateentities.ChannelEmail,
			Name:           "Password Reset Email",
			Description:    "Default password reset email for server-test",
			StorageKey:     "server-test/templates/password_reset.html",
			Version:        1,
			IsActive:       true,
			IsDefault:      true,
			Variables:      datatypes.JSON([]byte(`["FirstName","LastName","ResetLink","ExpirationTime"]`)),
			SampleData:     datatypes.JSON([]byte(`{"FirstName":"Test","LastName":"User","ResetLink":"http://localhost:5173/reset","ExpirationTime":"24 hours"}`)),
			TemplateType:   "email",
		},
		{
			TenantID:       tenantID,
			OrganizationID: &organizationIDPtr,
			Module:         "booking",
			TemplateKey:    "booking_confirmation",
			Channel:        templateentities.ChannelEmail,
			Name:           "Booking Confirmation Email",
			Description:    "Default booking confirmation email for server-test",
			StorageKey:     "server-test/templates/booking_confirmation.html",
			Version:        1,
			IsActive:       true,
			IsDefault:      true,
			Variables:      datatypes.JSON([]byte(`["CustomerFirstName","CustomerLastName","ServiceName","BookingReference","TotalAmount","Currency"]`)),
			SampleData:     datatypes.JSON([]byte(`{"CustomerFirstName":"Test","CustomerLastName":"Customer","ServiceName":"Test Service","BookingReference":"BK-2024-001","TotalAmount":100.00,"Currency":"EUR"}`)),
			TemplateType:   "email",
		},
		{
			TenantID:       tenantID,
			OrganizationID: &organizationIDPtr,
			Module:         "invoice",
			TemplateKey:    "invoice",
			Channel:        templateentities.ChannelPDF,
			Name:           "Invoice PDF",
			Description:    "Default invoice PDF template for server-test",
			StorageKey:     "server-test/templates/invoice.html",
			Version:        1,
			IsActive:       true,
			IsDefault:      true,
			Variables:      datatypes.JSON([]byte(`["Customer","InvoiceData","OrganizationData"]`)),
			SampleData:     datatypes.JSON([]byte(`{"Customer":{"Name":"Test Customer"},"InvoiceData":{"InvoiceNumber":"INV-2024-001","Total":1000.00},"OrganizationData":{"Name":"Test Organization"}}`)),
			TemplateType:   "pdf",
		},
	}
	for _, template := range templates {
		if err := db.Create(&template).Error; err != nil {
			return fmt.Errorf("failed to create template %s: %w", template.TemplateKey, err)
		}
	}

	contracts := []templateentities.TemplateContract{
		{
			TenantID:          tenantID,
			Module:            "user",
			TemplateKey:       "welcome",
			VariableSchema:    datatypes.JSON([]byte(`{"type":"object","properties":{"FirstName":{"type":"string"},"LastName":{"type":"string"},"OrganizationName":{"type":"string"},"ActivationLink":{"type":"string"}}}`)),
			DefaultSampleData: datatypes.JSON([]byte(`{"FirstName":"Test","LastName":"User","OrganizationName":"Server Test Organization","ActivationLink":"http://localhost:5173/activate"}`)),
		},
		{
			TenantID:          tenantID,
			Module:            "user",
			TemplateKey:       "password_reset",
			VariableSchema:    datatypes.JSON([]byte(`{"type":"object","properties":{"FirstName":{"type":"string"},"LastName":{"type":"string"},"ResetLink":{"type":"string"},"ExpirationTime":{"type":"string"}}}`)),
			DefaultSampleData: datatypes.JSON([]byte(`{"FirstName":"Test","LastName":"User","ResetLink":"http://localhost:5173/reset","ExpirationTime":"24 hours"}`)),
		},
		{
			TenantID:          tenantID,
			Module:            "booking",
			TemplateKey:       "booking_confirmation",
			VariableSchema:    datatypes.JSON([]byte(`{"type":"object","properties":{"CustomerFirstName":{"type":"string"},"CustomerLastName":{"type":"string"},"ServiceName":{"type":"string"},"BookingReference":{"type":"string"},"TotalAmount":{"type":"number"},"Currency":{"type":"string"}}}`)),
			DefaultSampleData: datatypes.JSON([]byte(`{"CustomerFirstName":"Test","CustomerLastName":"Customer","ServiceName":"Test Service","BookingReference":"BK-2024-001","TotalAmount":100.00,"Currency":"EUR"}`)),
		},
		{
			TenantID:          tenantID,
			Module:            "invoice",
			TemplateKey:       "invoice",
			VariableSchema:    datatypes.JSON([]byte(`{"type":"object","properties":{"Customer":{"type":"object"},"InvoiceData":{"type":"object"},"OrganizationData":{"type":"object"}}}`)),
			DefaultSampleData: datatypes.JSON([]byte(`{"Customer":{"Name":"Test Customer"},"InvoiceData":{"InvoiceNumber":"INV-2024-001","Total":1000.00},"OrganizationData":{"Name":"Test Organization"}}`)),
		},
	}
	for _, contract := range contracts {
		if err := db.Create(&contract).Error; err != nil {
			return fmt.Errorf("failed to create template contract %s: %w", contract.TemplateKey, err)
		}
	}
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

func maskPassword(password string) string {
	if password == "" {
		return "<empty>"
	}
	return "********"
}
