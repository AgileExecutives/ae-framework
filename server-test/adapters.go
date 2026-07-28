package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	models "github.com/AgileExecutives/ae-framework/serverbase/internal/models"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/auth"
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
