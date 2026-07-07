package main

import (
	"log"

	email "github.com/AgileExecutives/serverbase/modules/email"
	user "github.com/AgileExecutives/serverbase/modules/user"
	"github.com/AgileExecutives/serverbase/pkg/core"
	calmod "github.com/AgileExecutives/shared-modules/calendar"
	orgmod "github.com/AgileExecutives/shared-modules/organization"
	pdf "github.com/AgileExecutives/shared-modules/pdf"
	static "github.com/AgileExecutives/shared-modules/static"
	"github.com/gin-gonic/gin"

	internalHandlers "github.com/AgileExecutives/serverbase/internal/handlers"
	internalMiddleware "github.com/AgileExecutives/serverbase/internal/middleware"
	saashandlers "github.com/AgileExecutives/shared-modules/saas-base/handlers"

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
	// create gin engine to satisfy existing core.Module expectations
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	// Do not auto-redirect requests that differ only by trailing slash — tests
	// expect exact behavior for `/api/v1/static` vs `/api/v1/static/`.
	ginEngine.RedirectTrailingSlash = false

	// create in-memory sqlite DB for modules
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open sqlite db: %v", err)
	}

	// build a module context used by core modules
	coreCtx := core.ModuleContext{
		Router:       ginEngine,
		DB:           db,
		EventBus:     core.NewEventBus(),
		Logger:       core.NewLogger(),
		Auth:         &simpleAuthService{db: db},
		TokenService: &simpleTokenService{},
	}

	// Ensure saas-base models exist for plans/customers used by tests
	if err := db.AutoMigrate(&saasmodels.Plan{}, &saasmodels.Customer{}); err != nil {
		log.Printf("warning: failed to auto-migrate saas-base models: %v", err)
	}

	mr := sbmodule.NewRegistry()

	// register adapted core modules so they initialize and register their routes
	modules := []core.Module{
		user.NewUserModule(),
		orgmod.NewOrganizationModule(),
		email.NewEmailModule(),
		pdf.NewPDFModule(),
		calmod.NewCoreModule(),
		static.NewStaticModule(),
	}

	for _, m := range modules {
		// run AutoMigrate for module entities so handlers have tables available
		for _, e := range m.Entities() {
			if model := e.GetModel(); model != nil {
				if err := db.AutoMigrate(model); err != nil {
					log.Fatalf("auto-migrate failed for module %s: %v", m.Name(), err)
				}
			}
		}

		mr.RegisterModule(newCoreAdapter(m, coreCtx))
	}

	cfg := sbconfig.Config{Addr: ":8080"}

	server := sbhttp.New(cfg.Addr)

	// Mount ginEngine handlers under serverbase's http handler by using
	// the standard http Handler from gin and registering it at root.
	server.RegisterRoute("/", ginEngine)

	// Seed a test user used by Hurl tests if it doesn't exist
	var existing models.User
	if err := db.Where("email = ?", "testuser@unburdy.de").First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashed, _ := bcrypt.GenerateFromPassword([]byte("newpass123"), bcrypt.DefaultCost)
			u := models.User{
				Username:       "testuser",
				Email:          "testuser@unburdy.de",
				PasswordHash:   string(hashed),
				FirstName:      "Test",
				LastName:       "User",
				TenantID:       1,
				OrganizationID: 1,
				Role:           "admin",
				Active:         true,
				EmailVerified:  true,
			}
			if err := db.Create(&u).Error; err != nil {
				log.Printf("warning: failed to create test user: %v", err)
			} else {
				log.Println("Created test user testuser@unburdy.de")
			}
		}
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

	// initialize modules against our server adapter
	if err := mr.InitializeAll(server); err != nil {
		log.Fatalf("module init failed: %v", err)
	}

	// Register saas-base handlers (customers, plans, newsletter) so tests can hit /customers
	apiGroup := ginEngine.Group("/api/v1")
	protected := apiGroup.Group("")
	protected.Use(internalMiddleware.AuthMiddleware(db))

	customerHandlers := saashandlers.NewCustomerHandlers(db)
	// customer routes
	protected.GET("/customers", customerHandlers.GetCustomers)
	protected.POST("/customers", customerHandlers.CreateCustomer)
	protected.GET("/customers/:id", customerHandlers.GetCustomer)
	protected.PUT("/customers/:id", customerHandlers.UpdateCustomer)
	protected.DELETE("/customers/:id", customerHandlers.DeleteCustomer)

	// plan routes
	planHandler := internalHandlers.NewPlanHandler(db)
	// Public plans endpoints
	apiGroup.GET("/plans", planHandler.GetPlans)
	apiGroup.GET("/plans/:id", planHandler.GetPlan)

	// Admin plan management endpoints
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(internalMiddleware.AuthMiddleware(db), internalMiddleware.RequireAdmin())
	adminGroup.POST("/plans", planHandler.CreatePlan)
	adminGroup.PUT("/plans/:id", planHandler.UpdatePlan)
	adminGroup.DELETE("/plans/:id", planHandler.DeletePlan)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// keep process alive; server.Start runs in goroutine
	select {}
}
