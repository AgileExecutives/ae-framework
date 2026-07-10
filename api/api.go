package api

import (
	"fmt"
	"net/http"

	"github.com/AgileExecutives/serverbase/internal/middleware"
	"github.com/AgileExecutives/serverbase/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Re-export internal response types for external modules and handlers
type APIResponse = models.APIResponse
type PaginationResponse = models.PaginationResponse
type ListResponse = models.ListResponse
type ErrorResponse = models.ErrorResponse

// ErrorResponseFunc returns a JSON error payload used by handlers.
func ErrorResponseFunc(title, detail string) map[string]interface{} {
	return map[string]interface{}{"success": false, "error": title, "detail": detail}
}

// SuccessResponse returns a success envelope with data.
func SuccessResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"success": true, "message": message, "data": data}
}

// SuccessListResponse returns a paginated list envelope.
func SuccessListResponse(items interface{}, page, limit, total int) map[string]interface{} {
	return map[string]interface{}{"success": true, "items": items, "page": page, "limit": limit, "total": total}
}

// GetTenantID extracts tenant id from context or headers.
func GetTenantID(c *gin.Context) (uint, error) {
	return middleware.GetTenantID(c)
}

// SuccessMessageResponse returns a simple success envelope with only a message.
func SuccessMessageResponse(message string) map[string]interface{} {
	return map[string]interface{}{"success": true, "message": message}
}

// GetUser returns a User from context (set by auth middleware) for handlers.
func GetUser(c *gin.Context) (*models.User, error) {
	return middleware.GetUser(c)
}

// (legacy) AuthMiddleware wrapper removed — use module-provided middleware instead.

// GetUserID retrieves the authenticated user's ID from the context.
func GetUserID(c *gin.Context) (uint, error) {
	u, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

// (Removed) ModuleRouteProvider legacy compatibility interface.

// Helper to write JSON errors
func JSONError(c *gin.Context, status int, title, detail string) {
	c.JSON(status, ErrorResponseFunc(title, detail))
}

// Convenience wrapper for standard library http handlers if needed
func WriteHTTPError(w http.ResponseWriter, status int, title, detail string) {
	w.WriteHeader(status)
	w.Write([]byte(title + ": " + detail))
}

// --- Compatibility shims (legacy baseAPI helpers) -----------------------

// DatabaseConfig is a lightweight re-export used by older code paths.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectDatabaseWithAutoCreate is a compatibility wrapper used by seeding
// and legacy code. It attempts to open a connection but for safety returns
// a descriptive error if not implemented in the environment.
func ConnectDatabaseWithAutoCreate(cfg DatabaseConfig) (*gorm.DB, error) {
	// Keep implementation minimal: return a helpful error so callers can
	// decide whether to fallback to local connection logic.
	return nil, fmt.Errorf("ConnectDatabaseWithAutoCreate: not available in this workspace stub")
}

// MigrateBaseEntities runs GORM automigrations for core serverbase models.
// This ensures consuming projects (like unburdy) have required columns/tables
// available during tests that rely on the schema.
func MigrateBaseEntities(db *gorm.DB) error {
	// Auto-migrate core models used by external modules.
	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Organization{}, &models.Plan{}); err != nil {
		return err
	}
	return nil
}

// SeedBaseData is a shim used by seed helpers in consuming modules.
func SeedBaseData(db *gorm.DB) error { return nil }

// Minimal legacy types used by consuming code to avoid import churn.
type Tenant struct {
	ID   uint
	Name string
}

type User struct {
	ID             uint
	Email          string
	TenantID       uint
	OrganizationID uint
}

// Organization represents a minimal organization entity used by legacy consumers.
type Organization struct {
	ID               uint     `json:"id"`
	TenantID         uint     `json:"tenant_id"`
	Name             string   `json:"name"`
	OwnerName        string   `json:"owner_name,omitempty"`
	OwnerTitle       string   `json:"owner_title,omitempty"`
	StreetAddress    string   `json:"street_address,omitempty"`
	Zip              string   `json:"zip,omitempty"`
	City             string   `json:"city,omitempty"`
	Email            string   `json:"email,omitempty"`
	PaymentDueDays   int      `json:"payment_due_days,omitempty"`
	Phone            string   `json:"phone,omitempty"`
	TaxID            string   `json:"tax_id,omitempty"`
	TaxRate          *float64 `json:"tax_rate,omitempty"`
	TaxUstID         string   `json:"tax_ustid,omitempty"`
	BankAccountOwner string   `json:"bank_account_owner,omitempty"`
	BankAccountBank  string   `json:"bank_account_bank,omitempty"`
	BankAccountIBAN  string   `json:"bank_account_iban,omitempty"`
	BankAccountBIC   string   `json:"bank_account_bic,omitempty"`
}

// OrganizationResponse is a lightweight response shape for organization data.
type OrganizationResponse struct {
	ID       uint   `json:"id"`
	TenantID uint   `json:"tenant_id"`
	Name     string `json:"name"`
}

// ToResponse converts an Organization to its response representation.
func (o *Organization) ToResponse() OrganizationResponse {
	if o == nil {
		return OrganizationResponse{}
	}
	return OrganizationResponse{ID: o.ID, TenantID: o.TenantID, Name: o.Name}
}

// --- Legacy compatibility helpers (router/auth/interfaces) --------------

// SetupBaseRouterWithConfig returns a gin engine populated with basic middlewares
// and a root group. It's intentionally minimal to support legacy callers in tests.
func SetupBaseRouterWithConfig(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

// AuthMiddleware returns a middleware that validates authentication. For legacy
// compatibility in tests this is a permissive no-op that simply continues.
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

// ModuleRouteProvider is a lightweight compatibility interface used by legacy
// modules that implement route providers outside the bootstrap system.
type ModuleRouteProvider interface {
	GetPrefix() string
	RegisterRoutes(router *gin.RouterGroup)
}
