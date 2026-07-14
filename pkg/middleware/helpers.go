package middleware

import (
	"fmt"
	"net/http"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/gin-gonic/gin"
)

// RequireRole middleware checks if user has required role
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "error": "User not found", "detail": "User not authenticated"})
			c.Abort()
			return
		}

		user := userInterface.(*models.User)

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range requiredRoles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, map[string]interface{}{"success": false, "error": "Insufficient permissions", "detail": "User does not have required role"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware checks if user is admin
func RequireAdmin() gin.HandlerFunc { return RequireRole("admin", "super-admin") }

// TenantIsolation middleware ensures data access is limited to user's organization
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false, "error": "User not found", "detail": "User not authenticated"})
			c.Abort()
			return
		}

		user := userInterface.(*models.User)
		c.Set("tenant_id", user.TenantID)

		c.Next()
	}
}

// GetUser retrieves the authenticated user from context
func GetUser(c *gin.Context) (*models.User, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		return nil, fmt.Errorf("invalid user type in context")
	}

	return user, nil
}

// GetUserID retrieves the authenticated user ID from context
func GetUserID(c *gin.Context) (uint, error) {
	user, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetTenantID retrieves the tenant ID from context
func GetTenantID(c *gin.Context) (uint, error) {
	user, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return user.TenantID, nil
}
