package templates

import (
	"net/http"
	"net/http/httptest"
	"testing"

	templateentities "github.com/AgileExecutives/ae-framework/serverbase/modules/templates/entities"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTemplateRoutesExposeSeededDefaultAndContractAlias(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mod := NewTemplatesModule()
	ctx := core.ModuleContext{}

	if err := mod.Initialize(ctx); err != nil {
		t.Fatalf("initialize templates module: %v", err)
	}
	for _, routeProvider := range mod.Routes() {
		routeProvider.RegisterRoutes(router.Group("/api/v1"), ctx)
	}

	assertOK(t, router, "/api/v1/templates/default?template_type=email&channel=EMAIL")
	assertOK(t, router, "/api/v1/templates/contracts/by-key/welcome")
}

func TestTemplateModuleSeedsDatabaseDefaults(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	mod := NewTemplatesModule()
	ctx := core.ModuleContext{DB: db}
	if err := mod.Initialize(ctx); err != nil {
		t.Fatalf("initialize templates module with db: %v", err)
	}

	var templateCount int64
	if err := db.Model(&templateentities.Template{}).Count(&templateCount).Error; err != nil {
		t.Fatalf("count templates: %v", err)
	}
	if templateCount == 0 {
		t.Fatal("expected default templates to be seeded")
	}

	var contractCount int64
	if err := db.Model(&templateentities.TemplateContract{}).Count(&contractCount).Error; err != nil {
		t.Fatalf("count template contracts: %v", err)
	}
	if contractCount == 0 {
		t.Fatal("expected default template contracts to be seeded")
	}

	var welcome templateentities.Template
	if err := db.Where("template_key = ?", "welcome").First(&welcome).Error; err != nil {
		t.Fatalf("load welcome template: %v", err)
	}
	if welcome.Name == "" {
		t.Fatal("welcome template record was not seeded")
	}
}

func assertOK(t *testing.T, handler http.Handler, path string) {
	t.Helper()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET %s returned %d: %s", path, recorder.Code, recorder.Body.String())
	}
}
