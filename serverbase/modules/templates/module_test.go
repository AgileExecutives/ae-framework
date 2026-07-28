package templates

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
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

func assertOK(t *testing.T, handler http.Handler, path string) {
	t.Helper()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET %s returned %d: %s", path, recorder.Code, recorder.Body.String())
	}
}
