package main

import (
	"log"
	"os"
	"strings"

	email "github.com/AgileExecutives/serverbase/modules/email"
	orgmod "github.com/AgileExecutives/serverbase/modules/organizations"
	saas "github.com/AgileExecutives/serverbase/modules/saas"
	settingsModule "github.com/AgileExecutives/serverbase/modules/settings"
	templates "github.com/AgileExecutives/serverbase/modules/templates"
	user "github.com/AgileExecutives/serverbase/modules/user"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/swagger"
	auditmod "github.com/AgileExecutives/shared-modules/audit"
	bookingmod "github.com/AgileExecutives/shared-modules/booking"
	calmod "github.com/AgileExecutives/shared-modules/calendar"
	minimalorg "github.com/AgileExecutives/shared-modules/organization"
	pdf "github.com/AgileExecutives/shared-modules/pdf"
	static "github.com/AgileExecutives/shared-modules/static"
	"github.com/gin-gonic/gin"
	clientmod "github.com/unburdy/unburdy-server-api/modules/client_management"

	models "github.com/AgileExecutives/serverbase/internal/models"
	saasmodels "github.com/AgileExecutives/shared-modules/saas-base/models"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	sbconfig "github.com/AgileExecutives/serverbase/config"
	sbhttp "github.com/AgileExecutives/serverbase/http"
	sbmodule "github.com/AgileExecutives/serverbase/module"
)

func main() {
	// Load .env from server-test directory if present so local dev overrides apply
	loadLocalEnv()

	// create gin engine to satisfy existing core.Module expectations
	gin.SetMode(gin.ReleaseMode)
	// Ensure email verification is disabled for the test harness so
	// registration -> login works without external email flows.
	os.Setenv("FEATURE_EMAIL_VERIFICATION", "false")
	ginEngine := gin.New()
	// Do not auto-redirect requests that differ only by trailing slash — tests
	// expect exact behavior for `/api/v1/static` vs `/api/v1/static/`.
	ginEngine.RedirectTrailingSlash = false

	// create in-memory sqlite DB for modules
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open sqlite db: %v", err)
	}

	// Swagger doc registry – modules register their pre-generated JSON here during Initialize.
	docRegistry := swagger.NewRegistry()

	// build a module context used by core modules
	coreCtx := core.ModuleContext{
		Router:       ginEngine,
		DB:           db,
		EventBus:     core.NewEventBus(),
		Logger:       core.NewLogger(),
		Auth:         &simpleAuthService{db: db},
		TokenService: &simpleTokenService{},
		DocRegistry:  docRegistry,
		Services:     core.NewServiceRegistry(),
	}

	mr := sbmodule.NewRegistry()

	// All routes and swagger docs are registered by the modules themselves.
	modules := []core.Module{
		user.NewUserModule(),
		minimalorg.NewOrganizationModule(), // entity only
		orgmod.NewOrganizationsModule(),    // CRUD routes + swagger docs
		email.NewEmailModule(),
		pdf.NewPDFModule(),
		calmod.NewCoreModule(),
		static.NewStaticModule(),
		saas.NewSaaSModule(), // customers + plans (no newsletter – handled by user module)
		settingsModule.NewSettingsModule(),
		// Modules required by client management
		bookingmod.NewCoreModule(),
		auditmod.NewCoreModule(),
		clientmod.NewCoreModule(),
		// Templates module provides the in-memory template endpoints used by tests
		// and a simple TemplateService. Keep it last so other modules' routes are
		// available when templates are initialized if needed.
		templates.NewTemplatesModule(),
	}

	if err := sbmodule.RegisterCoreModules(mr, modules, db, coreCtx); err != nil {
		log.Fatalf("module registration failed: %v", err)
	}
	cfg := sbconfig.Config{Addr: ":8080"}

	server := sbhttp.New(cfg.Addr)

	// Mount ginEngine handlers under serverbase's http handler by using
	// the standard http Handler from gin and registering it at root.
	server.RegisterRoute("/", ginEngine)

	// Seed test data (users, plans, etc.) used by HURL tests and the harness.
	seedTestData(db)

	// initialize modules against our server adapter
	if err := mr.InitializeAll(server); err != nil {
		log.Fatalf("module init failed: %v", err)
	}

	// Merge swagger docs from all modules and mount the UI endpoint.
	swagger.SetupAndMount(docRegistry, ginEngine, swagger.ServerInfo{
		Title:       "AE SaaS API (test)",
		Description: "Combined API documentation for all registered modules",
		Version:     "1.0.0",
		BasePath:    "/api/v1",
		Schemes:     []string{"http"},
	})

	// Templates are provided by the templates module; create service for other
	// local uses and keep compatibility with older test expectations.

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// keep process alive; server.Start runs in goroutine
	select {}
}

// loadLocalEnv reads server-test/.env if it exists and sets environment variables.
// This is a lightweight loader used for the test harness; it ignores empty lines and comments.
func loadLocalEnv() {
	path := "server-test/.env"
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// split at first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove optional surrounding quotes
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		os.Setenv(key, val)
	}
}

// seedTestData ensures the test harness has the minimal required data present
// so HURL tests and local development work without external dependencies.
func seedTestData(db *gorm.DB) {
	// Seed a test user used by Hurl tests if it doesn't exist
	seedUsername := "testuser"
	seedEmail := "testuser@unburdy.de"
	seedPassword := "newpass123"
	seedRole := "admin"
	seedTenantID := uint(1)
	seedOrgID := uint(1)

	var existing models.User
	if err := db.Where("email = ?", seedEmail).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashed, _ := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
			u := models.User{
				Username:       seedUsername,
				Email:          seedEmail,
				PasswordHash:   string(hashed),
				FirstName:      "Test",
				LastName:       "User",
				TenantID:       seedTenantID,
				OrganizationID: seedOrgID,
				Role:           seedRole,
				Active:         true,
				EmailVerified:  true,
			}
			if err := db.Create(&u).Error; err != nil {
				log.Printf("warning: failed to create test user: %v", err)
			} else {
				log.Printf("Created test user %s", seedEmail)
				// Explicitly show the seed values used for creating the test user
				log.Println("--- Test Server Seed User ---")
				log.Printf("Username: %s", seedUsername)
				log.Printf("Email: %s", seedEmail)
				log.Printf("Password: %s", seedPassword)
				log.Printf("Role: %s", seedRole)
				log.Printf("TenantID: %d", seedTenantID)
				log.Printf("OrganizationID: %d", seedOrgID)
				log.Println("-----------------------------")
			}
		}
	} else {
		// If a user already exists in the in-memory DB, show what was found
		log.Println("--- Existing User Found in Test DB ---")
		log.Printf("Username: %s", existing.Username)
		log.Printf("Email: %s", existing.Email)
		log.Printf("Role: %s", existing.Role)
		log.Printf("TenantID: %d", existing.TenantID)
		log.Printf("OrganizationID: %d", existing.OrganizationID)
		log.Println("(Password not available for existing user)")
		log.Println("--------------------------------------")
	}

	// Seed default plans if none exist
	var planCount int64
	db.Model(&saasmodels.Plan{}).Count(&planCount)
	if planCount == 0 {
		plans := []saasmodels.Plan{
			{Name: "Free", Slug: "free", Description: "Free tier", Price: 0.0, Currency: "EUR", InvoicePeriod: "monthly", Active: true},
			{Name: "Pro", Slug: "pro", Description: "Pro tier", Price: 29.0, Currency: "EUR", InvoicePeriod: "monthly", Active: true},
		}
		for _, p := range plans {
			if err := db.Create(&p).Error; err != nil {
				log.Printf("warning: failed to create default plan %s: %v", p.Name, err)
			}
		}
		log.Println("Seeded default plans")
	}
}
