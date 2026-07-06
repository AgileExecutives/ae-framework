package handlers
package handlers

import (
	"net/http"

	baseAPI "github.com/ae/base-server/api"
	"github.com/ae/base-server/pkg/utils"
	"github.com/ae/shared-modules/saas-base/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CustomerHandlers provides customer management handlers.
type CustomerHandlers struct {
	db *gorm.DB
}

// NewCustomerHandlers creates new customer handlers.
func NewCustomerHandlers(db *gorm.DB) *CustomerHandlers {
	return &CustomerHandlers{db: db}
}

// GetCustomers retrieves all customers with pagination.
// @Summary Get all customers
// @ID getCustomers
// @Description Get a paginated list of all customers for the authenticated tenant
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param active query bool false "Filter by active status"
// @Success 200 {object} baseAPI.APIResponse{data=baseAPI.ListResponse}
// @Failure 401 {object} baseAPI.ErrorResponse
// @Failure 500 {object} baseAPI.ErrorResponse
// @Router /customers [get]
func (h *CustomerHandlers) GetCustomers(c *gin.Context) {
	user, err := baseAPI.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", err.Error()))
		return
	}

	page, limit := utils.GetPaginationParams(c)
	offset := utils.GetOffset(page, limit)

	var customers []models.Customer
	var total int64

	query := h.db.Model(&models.Customer{}).Where("tenant_id = ?", user.TenantID)

	if activeStr := c.Query("active"); activeStr != "" {
		if activeStr == "true" {
			query = query.Where("active = ?", true)
		} else if activeStr == "false" {
			query = query.Where("active = ?", false)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to count customers", err.Error()))
		return
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to retrieve customers", err.Error()))
		return
	}

	var responses []models.CustomerResponse
	for _, customer := range customers {
		responses = append(responses, customer.ToResponse())
	}

	response := baseAPI.ListResponse{
		Data: responses,
		Pagination: baseAPI.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: utils.CalculateTotalPages(int(total), limit),
		},
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customers retrieved successfully", response))
}

// GetCustomer retrieves a specific customer by ID.
// @Summary Get customer by ID
// @ID getCustomerById
// @Description Get a specific customer by its ID
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} baseAPI.APIResponse{data=models.CustomerResponse}
// @Failure 400 {object} baseAPI.ErrorResponse
// @Failure 401 {object} baseAPI.ErrorResponse
// @Failure 404 {object} baseAPI.ErrorResponse
// @Router /customers/{id} [get]
func (h *CustomerHandlers) GetCustomer(c *gin.Context) {
	user, err := baseAPI.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", err.Error()))
		return
	}

	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid customer ID", err.Error()))
		return
	}

	var customer models.Customer
	if err := h.db.Where("id = ? AND tenant_id = ?", id, user.TenantID).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to retrieve customer", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customer retrieved successfully", customer.ToResponse()))
}

// CreateCustomer creates a new customer.
// @Summary Create a new customer
// @ID createCustomer
// @Description Create a new customer within the authenticated tenant
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param customer body models.CustomerRequest true "Customer data"
// @Success 201 {object} baseAPI.APIResponse{data=models.CustomerResponse}
// @Failure 400 {object} baseAPI.ErrorResponse
// @Failure 401 {object} baseAPI.ErrorResponse
// @Router /customers [post]
func (h *CustomerHandlers) CreateCustomer(c *gin.Context) {
	user, err := baseAPI.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", err.Error()))
		return
	}

	var req models.CustomerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	req.TenantID = user.TenantID

	// Verify the plan exists
	var plan models.Plan
	if err := h.db.First(&plan, req.PlanID).Error; err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Plan not found", "Invalid plan ID"))
		return
	}

	customer := models.Customer{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Street:        req.Street,
		Zip:           req.Zip,
		City:          req.City,
		Country:       req.Country,
		TaxID:         req.TaxID,
		VAT:           req.VAT,
		PlanID:        req.PlanID,
		TenantID:      req.TenantID,
		Status:        "active",
		PaymentMethod: req.PaymentMethod,
		Active:        true,
	}

	if err := h.db.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to create customer", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, baseAPI.SuccessResponse("Customer created successfully", customer.ToResponse()))
}

// UpdateCustomer updates an existing customer.
// @Summary Update a customer
// @ID updateCustomer
// @Description Update an existing customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param customer body models.CustomerRequest true "Updated customer data"
// @Success 200 {object} baseAPI.APIResponse{data=models.CustomerResponse}
// @Failure 400 {object} baseAPI.ErrorResponse
// @Failure 401 {object} baseAPI.ErrorResponse
// @Failure 404 {object} baseAPI.ErrorResponse
// @Router /customers/{id} [put]
func (h *CustomerHandlers) UpdateCustomer(c *gin.Context) {
	user, err := baseAPI.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", err.Error()))
		return
	}

	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid customer ID", err.Error()))
		return
	}

	var req models.CustomerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	var customer models.Customer
	if err := h.db.Where("id = ? AND tenant_id = ?", id, user.TenantID).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to retrieve customer", err.Error()))
		return
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Street != "" {
		customer.Street = req.Street
	}
	if req.Zip != "" {
		customer.Zip = req.Zip
	}
	if req.City != "" {
		customer.City = req.City
	}
	if req.Country != "" {
		customer.Country = req.Country
	}
	if req.TaxID != "" {
		customer.TaxID = req.TaxID
	}
	if req.VAT != "" {
		customer.VAT = req.VAT
	}
	if req.PlanID != nil {
		var plan models.Plan
		if err := h.db.First(&plan, *req.PlanID).Error; err != nil {
			c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Plan not found", "Invalid plan ID"))
			return
		}
		customer.PlanID = *req.PlanID
	}
	if req.Status != "" {
		customer.Status = req.Status
	}
	if req.PaymentMethod != "" {
		customer.PaymentMethod = req.PaymentMethod
	}
	if req.Active != nil {
		customer.Active = *req.Active
	}

	if err := h.db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to update customer", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customer updated successfully", customer.ToResponse()))
}

// DeleteCustomer soft-deletes a customer.
// @Summary Delete a customer
// @ID deleteCustomer
// @Description Soft delete a customer by ID
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} baseAPI.APIResponse
// @Failure 400 {object} baseAPI.ErrorResponse
// @Failure 401 {object} baseAPI.ErrorResponse
// @Failure 404 {object} baseAPI.ErrorResponse
// @Router /customers/{id} [delete]
func (h *CustomerHandlers) DeleteCustomer(c *gin.Context) {
	user, err := baseAPI.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", err.Error()))
		return
	}

	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid customer ID", err.Error()))
		return
	}

	var customer models.Customer
	if err := h.db.Where("id = ? AND tenant_id = ?", id, user.TenantID).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to retrieve customer", err.Error()))
		return
	}

	if err := h.db.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to delete customer", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customer deleted successfully", nil))
}
