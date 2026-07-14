package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/pkg/auth"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const singleTenantID = uint(1)

// Options for the auth handler
type Options struct {
	SingleTenant bool
}

// BuildAuthHandler constructs the auth middleware handler used by the base
// module middleware provider and other callers.
func BuildAuthHandler(db *gorm.DB, logger core.Logger, opts Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Authorization required", "Missing authorization header"))
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid authorization format", "Use Bearer <token> format"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid token", err.Error()))
			c.Abort()
			return
		}

		// Check if token is blacklisted (compare against current time)
		var blacklistedToken models.TokenBlacklist
		if err := db.Where("token_id = ? AND expires_at > ?", claims.ID, time.Now()).First(&blacklistedToken).Error; err == nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Token blacklisted", "Token has been revoked"))
			c.Abort()
			return
		}

		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User associated with token not found"))
			c.Abort()
			return
		}

		if !user.Active {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Account disabled", "User account is not active"))
			c.Abort()
			return
		}

		if opts.SingleTenant {
			user.TenantID = singleTenantID
		} else {
			if user.TenantID != claims.TenantID {
				c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Tenant mismatch", fmt.Sprintf("user tenant %d does not match token tenant %d", user.TenantID, claims.TenantID)))
				c.Abort()
				return
			}
		}

		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Set("tenantID", user.TenantID)
		c.Set("token", tokenString)
		c.Set("claims", claims)

		c.Next()
	}
}
