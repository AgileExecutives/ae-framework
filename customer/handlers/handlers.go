package handlers
package handlers

import (
    "net/http"

    "github.com/ae/base-server/internal/models"
    customermodels "github.com/ae/base-server/modules/customer/models"
    "github.com/ae/base-server/pkg/core"
    "github.com/ae/base-server/pkg/utils"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type CustomerHandlers struct{ db *gorm.DB; logger core.Logger }
func NewCustomerHandlers(db *gorm.DB, logger core.Logger) *CustomerHandlers { return &CustomerHandlers{db: db, logger: logger} }

func (h *CustomerHandlers) GetCustomers(c *gin.Context) {
    userInterface, exists := c.Get("user")
    if !exists { c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated")); return }
    user := userInterface.(*models.User)
    page, limit := utils.GetPaginationParams(c); offset := utils.GetOffset(page, limit)
    var customers []customermodels.Customer; var total int64
    query := h.db.Model(&customermodels.Customer{}).Where("tenant_id = ?", user.TenantID)
    if activeStr := c.Query("active"); activeStr != "" { if activeStr == "true" { query = query.Where("active = ?", true) } else if activeStr == "false" { query = query.Where("active = ?", false) } }
    if err := query.Count(&total).Error; err != nil { c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to count customers", err.Error())); return }
    if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error; err != nil { c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve customers", err.Error())); return }
    var responses []customermodels.CustomerResponse; for _, customer := range customers { responses = append(responses, customer.ToResponse()) }
    response := models.ListResponse{ Data: responses, Pagination: models.PaginationResponse{ Page: page, Limit: limit, Total: int(total), TotalPages: utils.CalculateTotalPages(int(total), limit) } }
    c.JSON(http.StatusOK, models.SuccessResponse("Customers retrieved successfully", response))
}

func (h *CustomerHandlers) GetCustomer(c *gin.Context) {
    userInterface, exists := c.Get("user")
    if !exists { c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated")); return }
    user := userInterface.(*models.User)
    id, err := utils.ValidateID(c, "id"); if err != nil { c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid customer ID", err.Error())); return }
    var customer customermodels.Customer
    if err := h.db.Where("id = ? AND tenant_id = ?", id, user.TenantID).First(&customer).Error; err != nil {
        if err == gorm.ErrRecordNotFound { c.JSON(http.StatusNotFound, models.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist")); return }
        c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to retrieve customer", err.Error())); return
    }
    c.JSON(http.StatusOK, models.SuccessResponse("Customer retrieved successfully", customer.ToResponse()))
}

// CreateCustomer, UpdateCustomer, DeleteCustomer, PlanHandlers omitted for brevity — copied from original module when needed.
