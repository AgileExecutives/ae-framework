package main

import (
	"log"
	"os"
	"strings"

	sbconfig "github.com/AgileExecutives/ae-framework/serverbase/config"
	sbhttp "github.com/AgileExecutives/ae-framework/serverbase/http"
	sbmodule "github.com/AgileExecutives/ae-framework/serverbase/module"
	email "github.com/AgileExecutives/ae-framework/serverbase/modules/email"
	orgmod "github.com/AgileExecutives/ae-framework/serverbase/modules/organizations"
	saas "github.com/AgileExecutives/ae-framework/serverbase/modules/saas"
	settingsModule "github.com/AgileExecutives/ae-framework/serverbase/modules/settings"
	templates "github.com/AgileExecutives/ae-framework/serverbase/modules/templates"
	user "github.com/AgileExecutives/ae-framework/serverbase/modules/user"
	pkgconfig "github.com/AgileExecutives/ae-framework/serverbase/pkg/config"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/database"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/swagger"
	minimalorg "github.com/AgileExecutives/ae-framework/shared-modules/organization"
	pdf "github.com/AgileExecutives/ae-framework/shared-modules/pdf"
	static "github.com/AgileExecutives/ae-framework/shared-modules/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load .env from server-test directory if present so local dev overrides apply
	loadLocalEnv()

	// create gin engine to satisfy existing core.Module expectations
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	// Do not auto-redirect requests that differ only by trailing slash — tests
	// expect exact behavior for `/api/v1/static` vs `/api/v1/static/`.
	ginEngine.RedirectTrailingSlash = false

	// Create database connection. Postgres is the default; set USE_IN_MEMORY_DB=true
	// to run the server-test harness against SQLite.
	dbConfig := pkgconfig.Load().Database
	normalizeDatabaseEnvCompatibility(&dbConfig)
	db, err := database.ConnectWithAutoCreate(dbConfig)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
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
		static.NewStaticModule(),
		saas.NewSaaSModule(), // customers + plans (no newsletter – handled by user module)
		settingsModule.NewSettingsModule(),
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

	// Seed test data used by HURL tests and the harness when the DB is empty.
	if err := RunIfEmpty(db); err != nil {
		log.Fatalf("server-test seed failed: %v", err)
	}

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
	// Try common local env file locations so the harness works when run
	// from the repo root (`go run server-test/*.go`) and when run from
	// the `server-test` directory (`go run main.go`). Prefer
	// `server-test/.env` when running from repo root, otherwise fall back
	// to `.env` in the current directory.
	candidates := []string{"server-test/.env", ".env"}
	var data []byte
	var err error
	var loaded string
	for _, path := range candidates {
		data, err = os.ReadFile(path)
		if err == nil {
			loaded = path
			break
		}
	}
	if err != nil {
		return
	}
	log.Printf("Loaded local env from %s", loaded)
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

func normalizeDatabaseEnvCompatibility(dbConfig *database.Config) {
	if sslMode := strings.TrimSpace(os.Getenv("DB_SSL_MODE")); sslMode != "" && strings.TrimSpace(os.Getenv("DB_SSLMODE")) == "" {
		dbConfig.SSLMode = sslMode
	}
}
