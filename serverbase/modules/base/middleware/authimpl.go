package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AgileExecutives/ae-framework/serverbase/internal/models"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/auth"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
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
		logger.Debug("auth: request", "method", c.Request.Method, "path", c.Request.URL.Path, "remote", c.Request.RemoteAddr, "user_agent", c.GetHeader("User-Agent"))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Info("auth: missing Authorization header")
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Authorization required", "Missing authorization header"))
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Info("auth: invalid authorization header format", "header", authHeader)
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid authorization format", "Use Bearer <token> format"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Mask token when logging (show start/end only)
		maskedToken := tokenString
		if len(tokenString) > 16 {
			maskedToken = fmt.Sprintf("%s...%s", tokenString[:8], tokenString[len(tokenString)-4:])
		}
		logger.Debug("auth: received token", "token", maskedToken)

		// Pre-decode header and payload for extra diagnostics (unsafe: does not verify signature)
		parts := strings.Split(tokenString, ".")
		var hdr map[string]interface{}
		var pay map[string]interface{}
		if len(parts) >= 2 {
			if h, err := base64.RawURLEncoding.DecodeString(parts[0]); err == nil {
				_ = json.Unmarshal(h, &hdr)
			}
			if p, err := base64.RawURLEncoding.DecodeString(parts[1]); err == nil {
				_ = json.Unmarshal(p, &pay)
			}
		}

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			// Log the validation error plus decoded header/payload hints to help debug signature/key issues
			if hdr != nil {
				logger.Info("auth: token validation failed", "error", err.Error(), "token_header", hdr)
			} else {
				logger.Info("auth: token validation failed", "error", err.Error())
			}
			if pay != nil {
				logger.Debug("auth: token payload (unverified)", "payload", pay)
			}

			// Return a concise error to client but keep diagnostic details in logs
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid token", "token validation failed"))
			c.Abort()
			return
		}

		logger.Debug("auth: token validated", "token_id", claims.ID, "user_id", claims.UserID, "tenant_id", claims.TenantID)

		// Check if token is blacklisted (compare against current time)
		var blacklistedToken models.TokenBlacklist
		if err := db.Where("token_id = ? AND expires_at > ?", claims.ID, time.Now()).First(&blacklistedToken).Error; err == nil {
			logger.Info("auth: token blacklisted", "token_id", claims.ID)
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Token blacklisted", "Token has been revoked"))
			c.Abort()
			return
		}

		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			logger.Info("auth: user lookup failed", "user_id", claims.UserID, "error", err.Error())
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User associated with token not found"))
			c.Abort()
			return
		}

		if !user.Active {
			logger.Info("auth: user inactive", "user_id", user.ID)
			c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Account disabled", "User account is not active"))
			c.Abort()
			return
		}

		if opts.SingleTenant {
			user.TenantID = singleTenantID
		} else {
			if user.TenantID != claims.TenantID {
				logger.Info("auth: tenant mismatch", "user_tenant", user.TenantID, "token_tenant", claims.TenantID, "user_id", user.ID)
				c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Tenant mismatch", fmt.Sprintf("user tenant %d does not match token tenant %d", user.TenantID, claims.TenantID)))
				c.Abort()
				return
			}
		}

		// Successful authentication: attach user and metadata to context
		logger.Debug("auth: authentication successful", "user_id", user.ID, "tenant_id", user.TenantID)
		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Set("tenantID", user.TenantID)
		c.Set("token", tokenString)
		c.Set("claims", claims)

		c.Next()
	}
}
