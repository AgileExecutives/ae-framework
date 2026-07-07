package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	models "github.com/AgileExecutives/serverbase/internal/models"
	sbmodule "github.com/AgileExecutives/serverbase/module"
	"github.com/AgileExecutives/serverbase/pkg/auth"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// simpleAuthService is a minimal AuthService that validates JWTs and loads
// the corresponding user from the provided DB so handlers can read `user` from
// Gin context as expected by the modules.
type simpleAuthService struct {
	db *gorm.DB
}

func (s *simpleAuthService) ValidateToken(token string) (interface{}, error) {
	return auth.ValidateJWT(token)
}
func (s *simpleAuthService) GenerateToken(user interface{}) (string, error) { return "", nil }
func (s *simpleAuthService) GetCurrentUser(c *gin.Context) (interface{}, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, nil
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, nil
	}
	claims, err := auth.ValidateJWT(parts[1])
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *simpleAuthService) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Authorization required", "Missing authorization header"))
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid authorization format", "Use Bearer <token> format"))
			c.Abort()
			return
		}
		claims, err := auth.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid token", err.Error()))
			c.Abort()
			return
		}
		var user models.User
		if err := s.db.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User associated with token not found"))
			c.Abort()
			return
		}
		// Debug log to verify middleware sets context keys
		log.Printf("simpleAuthService: authenticated user id=%d tenant=%d", user.ID, user.TenantID)
		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Set("tenantID", user.TenantID)
		c.Set("token", parts[1])
		c.Set("claims", claims)
		c.Next()
	}
}

func (s *simpleAuthService) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userI, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
			c.Abort()
			return
		}
		user := userI.(*models.User)
		for _, r := range roles {
			if user.Role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(403, gin.H{"success": false, "message": "Forbidden"})
	}
}

// simpleTokenService is a minimal TokenService stub.
type simpleTokenService struct{}

func (s *simpleTokenService) GenerateToken(claims interface{}) (string, error)           { return "", nil }
func (s *simpleTokenService) ValidateToken(tokenString string, claims interface{}) error { return nil }
func (s *simpleTokenService) ParseTokenID(tokenString string) (string, error)            { return "", nil }
func (s *simpleTokenService) GetTokenExpiration(tokenString string) (time.Time, error) {
	return time.Time{}, nil
}

// coreModuleAdapter adapts a core.Module to the serverbase Module interface.
type coreModuleAdapter struct {
	mod core.Module
	ctx core.ModuleContext
}

func newCoreAdapter(m core.Module, ctx core.ModuleContext) sbmodule.Module {
	return &coreModuleAdapter{mod: m, ctx: ctx}
}

func (a *coreModuleAdapter) Name() string { return a.mod.Name() }

// Register will initialize the underlying core module and register its routes
// directly on the provided ModuleContext.Router (gin.Engine) under /api/v1.
func (a *coreModuleAdapter) Register(reg sbmodule.Registry) error {
	// Initialize the module with the prepared context
	if err := a.mod.Initialize(a.ctx); err != nil {
		return err
	}

	// Register each route provider on the gin router
	apiV1 := a.ctx.Router.Group("/api/v1")
	for _, rp := range a.mod.Routes() {
		group := apiV1.Group(rp.GetPrefix())
		for _, m := range rp.GetMiddleware() {
			group.Use(m)
		}
		rp.RegisterRoutes(group, a.ctx)
	}

	// Register event handlers / services would go here if needed
	return nil
}

func (a *coreModuleAdapter) Start(ctx context.Context) error { return a.mod.Start(ctx) }
func (a *coreModuleAdapter) Stop(ctx context.Context) error  { return a.mod.Stop(ctx) }
