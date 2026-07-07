package handlers
package handlers

import (
    "net/http"

    "github.com/ae/base-server/internal/models"
    "github.com/ae/base-server/pkg/auth"
    "github.com/gin-gonic/gin"
)

func (h *AuthHandlers) CheckResetToken(c *gin.Context) {
    token := c.Param("token")
    email, err := auth.ValidateResetToken(token)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid or expired reset token", err.Error()))
        return
    }
    c.JSON(http.StatusOK, models.SuccessResponse("Token is valid", gin.H{"valid": true, "email": email}))
}
