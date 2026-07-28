package handlers

import (
	"net/http"

	baseAPI "github.com/AgileExecutives/ae-framework/serverbase/api"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/base/models"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/base/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CustomerHandlers struct {
	service *services.CustomerService
}

func NewCustomerHandlers(s *services.CustomerService) *CustomerHandlers {
	return &CustomerHandlers{service: s}
}

func (h *CustomerHandlers) GetCustomers(c *gin.Context) {
	user, err := baseAPI.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", err.Error()))
		return
	}

	page, limit := utils.GetPaginationParams(c)

	customers, err := h.service.GetByTenant(user.TenantID)
	if err != nil {
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
			Total:      len(responses),
			TotalPages: utils.CalculateTotalPages(len(responses), limit),
		},
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customers retrieved successfully", response))
}

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

	customer, err := h.service.GetByID(id)
	if err != nil || customer == nil || customer.TenantID != user.TenantID {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
		return
	}
	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customer retrieved successfully", customer.ToResponse()))
}

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
	customer := &models.Customer{
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

	if err := h.service.Save(customer); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to create customer", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, baseAPI.SuccessResponse("Customer created successfully", customer.ToResponse()))
}

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

	customer, err := h.service.GetByID(id)
	if err != nil || customer == nil || customer.TenantID != user.TenantID {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
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

	if err := h.service.Save(customer); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to update customer", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customer updated successfully", customer.ToResponse()))
}

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

	cust, err := h.service.GetByID(id)
	if err != nil || cust == nil || cust.TenantID != user.TenantID {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Customer not found", "Customer with specified ID does not exist"))
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to delete customer", err.Error()))
		return
	}
	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Customer deleted successfully", nil))
}
