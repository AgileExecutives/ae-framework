package handlers
package handlers

import (
    "net/http"
    "os"
    "strconv"
    "time"

    _ "github.com/ae/base-server/modules/user/models" // Import models for swagger
    "github.com/ae/base-server/pkg/core"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type HealthHandlers struct {
    db     *gorm.DB
    logger core.Logger
}

func NewHealthHandlers(db *gorm.DB, logger core.Logger) *HealthHandlers {
    return &HealthHandlers{db: db, logger: logger}
}

func getEnvAsBool(key string, defaultVal bool) bool {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultVal
    }
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultVal
}

func (h *HealthHandlers) HealthCheck(c *gin.Context) {
    response := gin.H{"status": "healthy", "timestamp": time.Now().UTC(), "version": "2.0", "database": "connected"}
    if h.db != nil {
        sqlDB, err := h.db.DB()
        if err != nil || sqlDB.Ping() != nil {
            response["database"] = "disconnected"
            response["status"] = "unhealthy"
            c.JSON(http.StatusServiceUnavailable, response)
            return
        }
    }
    environment := map[string]interface{}{ "mock_email": getEnvAsBool("MOCK_EMAIL", false), "rate_limit_enabled": getEnvAsBool("RATE_LIMIT_ENABLED", true), "email_verification": getEnvAsBool("FEATURE_EMAIL_VERIFICATION", true), "gin_mode": os.Getenv("GIN_MODE"), }
    response["environment"] = environment
    c.JSON(http.StatusOK, response)
}

func (h *HealthHandlers) Ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
