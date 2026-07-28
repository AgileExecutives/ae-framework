package handlers

import (
	"net/http"

	baseAPI "github.com/AgileExecutives/ae-framework/serverbase/api"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/utils"
	"github.com/AgileExecutives/ae-framwork/shared-modules/saas-base/models"
	"github.com/AgileExecutives/ae-framwork/shared-modules/saas-base/services"
	"github.com/gin-gonic/gin"
)

// PlanHandlers provides subscription plan management handlers.
type PlanHandlers struct {
	service *services.PlanService
}

// NewPlanHandlers creates new plan handlers.
func NewPlanHandlers(s *services.PlanService) *PlanHandlers { return &PlanHandlers{service: s} }

// GetPlans retrieves all available plans with pagination.
// @Summary Get all plans
// @ID getPlans
// @Description Get a paginated list of all subscription plans
// @Tags plans
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param active query bool false "Filter by active status"
// @Success 200 {object} handlers.APIResponse{data=handlers.ListResponse}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /plans [get]
func (h *PlanHandlers) GetPlans(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	_ = utils.GetOffset(page, limit)

	plans, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to retrieve plans", err.Error()))
		return
	}

	var responses []models.PlanResponse
	for _, plan := range plans {
		responses = append(responses, plan.ToResponse())
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

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Plans retrieved successfully", response))
}

// GetPlan retrieves a specific plan by ID.
// @Summary Get plan by ID
// @ID getPlanById
// @Description Get a specific subscription plan by its ID
// @Tags plans
// @Produce json
// @Param id path int true "Plan ID"
// @Success 200 {object} handlers.APIResponse{data=models.PlanResponse}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /plans/{id} [get]
func (h *PlanHandlers) GetPlan(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid plan ID", err.Error()))
		return
	}

	plan, err := h.service.GetByID(id)
	if err != nil || plan == nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Plan not found", "Plan with specified ID does not exist"))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Plan retrieved successfully", plan.ToResponse()))
}

// CreatePlan creates a new subscription plan.
// @Summary Create a new plan
// @ID createPlan
// @Description Create a new subscription plan (admin only)
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan body models.PlanRequest true "Plan data"
// @Success 201 {object} handlers.APIResponse{data=models.PlanResponse}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 409 {object} handlers.ErrorResponse
// @Router /plans [post]
func (h *PlanHandlers) CreatePlan(c *gin.Context) {
	var req models.PlanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	if req.Currency == "" {
		req.Currency = "EUR"
	}
	if req.InvoicePeriod == "" {
		req.InvoicePeriod = "monthly"
	}
	if req.MaxUsers == 0 {
		req.MaxUsers = 10
	}
	if req.MaxClients == 0 {
		req.MaxClients = 100
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	plan := models.Plan{
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		Price:         req.Price,
		Currency:      req.Currency,
		InvoicePeriod: req.InvoicePeriod,
		MaxUsers:      req.MaxUsers,
		MaxClients:    req.MaxClients,
		Features:      req.Features,
		Active:        active,
	}

	if err := h.service.Save(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to create plan", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, baseAPI.SuccessResponse("Plan created successfully", plan.ToResponse()))
}

// UpdatePlan updates an existing subscription plan.
// @Summary Update a plan
// @ID updatePlan
// @Description Update an existing subscription plan (admin only)
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Param plan body models.PlanRequest true "Updated plan data"
// @Success 200 {object} handlers.APIResponse{data=models.PlanResponse}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /plans/{id} [put]
func (h *PlanHandlers) UpdatePlan(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid plan ID", err.Error()))
		return
	}

	var req models.PlanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	plan, err := h.service.GetByID(id)
	if err != nil || plan == nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Plan not found", "Plan with specified ID does not exist"))
		return
	}

	if req.Name != "" {
		plan.Name = req.Name
	}
	if req.Description != "" {
		plan.Description = req.Description
	}
	if req.Price != nil {
		plan.Price = *req.Price
	}
	if req.Currency != "" {
		plan.Currency = req.Currency
	}
	if req.InvoicePeriod != "" {
		plan.InvoicePeriod = req.InvoicePeriod
	}
	if req.MaxUsers != nil {
		plan.MaxUsers = *req.MaxUsers
	}
	if req.MaxClients != nil {
		plan.MaxClients = *req.MaxClients
	}
	if req.Features != "" {
		plan.Features = req.Features
	}
	if req.Active != nil {
		plan.Active = *req.Active
	}

	if err := h.service.Save(plan); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to update plan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Plan updated successfully", plan.ToResponse()))
}

// DeletePlan soft-deletes a subscription plan.
// @Summary Delete a plan
// @ID deletePlan
// @Description Delete a subscription plan (admin only)
// @Tags plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} handlers.APIResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /plans/{id} [delete]
func (h *PlanHandlers) DeletePlan(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid plan ID", err.Error()))
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to delete plan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Plan deleted successfully", nil))
}
