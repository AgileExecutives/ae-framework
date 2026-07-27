package handlers

import (
    "net/http"

    baseAPI "github.com/AgileExecutives/serverbase/api"
    "github.com/AgileExecutives/serverbase/pkg/utils"
    "github.com/AgileExecutives/serverbase/modules/base/models"
    "github.com/AgileExecutives/serverbase/modules/base/services"
    "github.com/gin-gonic/gin"
)

type PlanHandlers struct {
    service *services.PlanService
}

func NewPlanHandlers(s *services.PlanService) *PlanHandlers { return &PlanHandlers{service: s} }

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

func (h *PlanHandlers) CreatePlan(c *gin.Context) {
    var req models.PlanRequest
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

func (h *PlanHandlers) UpdatePlan(c *gin.Context) {
    id, err := utils.ValidateID(c, "id")
    if err != nil {
        c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid plan ID", err.Error()))
        return
    }

    var req models.PlanRequest
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
    if req.Price != 0 {
        plan.Price = req.Price
    }
    if req.Currency != "" {
        plan.Currency = req.Currency
    }
    if req.InvoicePeriod != "" {
        plan.InvoicePeriod = req.InvoicePeriod
    }
    if req.MaxUsers != 0 {
        plan.MaxUsers = req.MaxUsers
    }
    if req.MaxClients != 0 {
        plan.MaxClients = req.MaxClients
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
