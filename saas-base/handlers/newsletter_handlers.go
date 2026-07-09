package handlers

import (
	"net/http"
	"time"

	baseAPI "github.com/AgileExecutives/serverbase/api"
	"github.com/AgileExecutives/serverbase/pkg/utils"
	"github.com/AgileExecutives/shared-modules/saas-base/models"
	"github.com/AgileExecutives/shared-modules/saas-base/services"
	"github.com/gin-gonic/gin"
)

// NewsletterHandlers provides newsletter subscription management handlers.
type NewsletterHandlers struct {
	service *services.NewsletterService
}

// NewNewsletterHandlers creates new newsletter handlers.
func NewNewsletterHandlers(s *services.NewsletterService) *NewsletterHandlers {
	return &NewsletterHandlers{service: s}
}

// GetSubscribers lists all newsletter subscribers with pagination.
// @Summary Get all newsletter subscribers
// @ID getNewsletterSubscribers
// @Description Get a paginated list of newsletter subscribers (admin only)
// @Tags newsletter
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} baseAPI.APIResponse{data=baseAPI.ListResponse}
// @Failure 401 {object} baseAPI.ErrorResponse
// @Failure 500 {object} baseAPI.ErrorResponse
// @Router /newsletter [get]
func (h *NewsletterHandlers) GetSubscribers(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	subs, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to retrieve subscribers", err.Error()))
		return
	}

	var responses []models.NewsletterResponse
	for _, s := range subs {
		responses = append(responses, s.ToResponse())
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

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Subscribers retrieved successfully", response))
}

// Subscribe adds a new newsletter subscription.
// @Summary Subscribe to newsletter
// @ID subscribeNewsletter
// @Description Subscribe an email address to the newsletter
// @Tags newsletter
// @Accept json
// @Produce json
// @Param subscription body models.NewsletterSubscribeRequest true "Subscription data"
// @Success 201 {object} baseAPI.APIResponse{data=models.NewsletterResponse}
// @Failure 400 {object} baseAPI.ErrorResponse
// @Router /newsletter/subscribe [post]
func (h *NewsletterHandlers) Subscribe(c *gin.Context) {
	var req models.NewsletterSubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	if req.Interest == "" {
		req.Interest = "general"
	}

	// Check if already subscribed
	existing, err := h.service.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to check subscription", err.Error()))
		return
	}
	if existing != nil {
		existing.Name = req.Name
		existing.Interest = req.Interest
		existing.Source = req.Source
		existing.LastContact = time.Now()
		if err := h.service.Save(existing); err != nil {
			c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to update subscription", err.Error()))
			return
		}
		c.JSON(http.StatusOK, baseAPI.SuccessResponse("Subscription updated successfully", existing.ToResponse()))
		return
	}

	subscriber := &models.Newsletter{
		Name:        req.Name,
		Email:       req.Email,
		Interest:    req.Interest,
		Source:      req.Source,
		LastContact: time.Now(),
	}

	if err := h.service.Save(subscriber); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to create subscription", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, baseAPI.SuccessResponse("Subscribed successfully", subscriber.ToResponse()))
}

// Unsubscribe removes a newsletter subscription.
// @Summary Unsubscribe from newsletter
// @ID unsubscribeNewsletter
// @Description Remove an email address from the newsletter
// @Tags newsletter
// @Accept json
// @Produce json
// @Param request body models.NewsletterUnsubscribeRequest true "Unsubscribe data"
// @Success 200 {object} baseAPI.APIResponse
// @Failure 400 {object} baseAPI.ErrorResponse
// @Failure 404 {object} baseAPI.ErrorResponse
// @Router /newsletter/unsubscribe [post]
func (h *NewsletterHandlers) Unsubscribe(c *gin.Context) {
	var req models.NewsletterUnsubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	sub, err := h.service.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to find subscriber", err.Error()))
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Subscriber not found", "Email is not subscribed"))
		return
	}

	if err := h.service.Delete(sub.ID); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to unsubscribe", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Unsubscribed successfully", nil))
}

// DeleteSubscriber hard-deletes a subscriber by ID (admin only).
// @Summary Delete subscriber by ID
// @ID deleteNewsletterSubscriber
// @Description Permanently delete a subscriber record (admin only)
// @Tags newsletter
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscriber ID"
// @Success 200 {object} baseAPI.APIResponse
// @Failure 400 {object} baseAPI.ErrorResponse
// @Failure 404 {object} baseAPI.ErrorResponse
// @Router /newsletter/{id} [delete]
func (h *NewsletterHandlers) DeleteSubscriber(c *gin.Context) {
	id, err := utils.ValidateID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid subscriber ID", err.Error()))
		return
	}

	sub, err := h.service.GetByID(id)
	if err != nil || sub == nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Subscriber not found", "Subscriber with specified ID does not exist"))
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Failed to delete subscriber", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Subscriber deleted successfully", nil))
}
