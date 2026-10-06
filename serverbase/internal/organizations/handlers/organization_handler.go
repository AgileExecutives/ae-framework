package handlers

import (
	"io"
	"net/http"
	"strconv"

	baseAPI "github.com/AgileExecutives/ae-framework/serverbase/api"
	"github.com/AgileExecutives/ae-framework/serverbase/internal/models"
	"github.com/AgileExecutives/ae-framework/serverbase/internal/organizations/services"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/formatting"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/utils"
	"github.com/gin-gonic/gin"
)

// OrganizationHandler handles organization-related HTTP requests
type OrganizationHandler struct {
	service *services.OrganizationService
}

const (
	maxLogoSize     = 10 * 1024 * 1024 // 10 MB
	logoFormats     = "image/png,image/jpeg,image/webp,image/svg+xml"
	logoPrintFormats = "image/png,image/jpeg,image/webp,image/svg+xml"
)

var supportedMimeTypes = map[string]bool{
	"image/png":     true,
	"image/jpeg":    true,
	"image/webp":    true,
	"image/svg+xml": true,
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		service: service,
	}
}

// CreateOrganization handles creating a new organization
// @Summary Create a new organization
// @Description Create a new organization with the provided information
// @Tags organizations
// @ID createOrganization
// @Accept json
// @Produce json
// @Param organization body models.CreateOrganizationRequest true "Organization information"
// @Success 201 {object} models.OrganizationAPIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req models.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	tenantID, err := baseAPI.GetTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", "Unable to get tenant ID: "+err.Error()))
		return
	}

	organization, err := h.service.CreateOrganization(req, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Internal server error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, baseAPI.SuccessResponse("Organization created successfully", organization.ToResponse()))
}

// GetOrganization handles retrieving an organization by ID
// @Summary Get an organization by ID
// @Description Retrieve a specific organization by ID
// @Tags organizations
// @ID getOrganizationById
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} models.OrganizationAPIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id} [get]
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Invalid organization ID"))
		return
	}

	tenantID, err := baseAPI.GetTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", "Unable to get tenant ID: "+err.Error()))
		return
	}

	organization, err := h.service.GetOrganizationByID(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("", organization.ToResponse()))
}

// GetAllOrganizations handles retrieving all organizations with pagination
// @Summary Get all organizations
// @Description Retrieve all organizations for the authenticated user with pagination
// @Tags organizations
// @ID getOrganizations
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} models.OrganizationListAPIResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /organizations [get]
func (h *OrganizationHandler) GetAllOrganizations(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	tenantID, err := baseAPI.GetTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", "Unable to get tenant ID: "+err.Error()))
		return
	}

	organizations, total, err := h.service.GetOrganizations(page, limit, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Internal server error", err.Error()))
		return
	}

	responses := make([]models.OrganizationResponse, len(organizations))
	for i, org := range organizations {
		responses[i] = org.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
		"page":    page,
		"limit":   limit,
		"total":   int(total),
	})
}

// UpdateOrganization handles updating an organization
// @Summary Update an organization
// @Description Update an organization's information
// @Tags organizations
// @ID updateOrganization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param organization body models.UpdateOrganizationRequest true "Updated organization information"
// @Success 200 {object} models.OrganizationAPIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Invalid organization ID"))
		return
	}

	var req models.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	tenantID, err := baseAPI.GetTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", "Unable to get tenant ID: "+err.Error()))
		return
	}

	organization, err := h.service.UpdateOrganization(uint(id), tenantID, req)
	if err != nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Organization updated successfully", organization.ToResponse()))
}

// DeleteOrganization handles deleting an organization
// @Summary Delete an organization
// @Description Delete an organization by ID
// @Tags organizations
// @ID deleteOrganization
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} models.OrganizationDeleteResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Invalid organization ID"))
		return
	}

	tenantID, err := baseAPI.GetTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", "Unable to get tenant ID: "+err.Error()))
		return
	}

	if err := h.service.DeleteOrganization(uint(id), tenantID); err != nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessMessageResponse("Organization deleted successfully"))
}

// GetSupportedFormats returns the supported date, time, and amount formats
// @Summary Get supported formats
// @Description Get all supported date, time, and amount formats with examples
// @Tags organizations
// @ID getSupportedFormats
// @Produce json
// @Success 200 {object} map[string]interface{} "Supported formats with examples"
// @Router /organizations/supported-formats [get]
func (h *OrganizationHandler) GetSupportedFormats(c *gin.Context) {
	response := map[string]interface{}{
		"date_formats":   formatting.GetSupportedDateFormats(),
		"time_formats":   formatting.GetSupportedTimeFormats(),
		"amount_formats": formatting.GetSupportedAmountFormats(),
		"locales":        formatting.GetSupportedLocales(),
	}
	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Supported formats retrieved successfully", response))
}

// UploadLogo handles uploading an organization logo
// @Summary Upload organization logo
// @Description Upload a new logo for the organization
// @Tags organizations
// @ID uploadOrganizationLogo
// @Accept mpfd
// @Produce json
// @Param id path int true "Organization ID"
// @Param logo formData file true "Logo image (PNG, JPEG, WebP, SVG, max 10MB)"
// @Success 200 {object} models.OrganizationAPIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 415 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id}/logo [post]
func (h *OrganizationHandler) UploadLogo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Invalid organization ID"))
		return
	}

	tenantID, err := baseAPI.GetTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", "Unable to get tenant ID: "+err.Error()))
		return
	}

	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Logo file is required"))
		return
	}

	if file.Size > maxLogoSize {
		c.JSON(http.StatusRequestEntityTooLarge, baseAPI.ErrorResponseFunc("File too large", "Logo must be smaller than 10MB"))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Internal server error", "Unable to read logo file"))
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Internal server error", "Unable to read logo file"))
		return
	}

	contentType := file.Header.Get("Content-Type")
	if !supportedMimeTypes[contentType] {
		c.JSON(http.StatusUnsupportedMediaType, baseAPI.ErrorResponseFunc("Unsupported media type", "Only PNG, JPEG, WebP, and SVG formats are supported"))
		return
	}

	organization, err := h.service.UpdateOrganizationLogo(uint(id), tenantID, fileBytes, contentType)
	if err != nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Logo uploaded successfully", organization.ToResponse()))
}

// UploadLogoPrint handles uploading an organization print logo
// @Summary Upload organization print logo
// @Description Upload a new print logo for the organization
// @Tags organizations
// @ID uploadOrganizationLogoPrint
// @Accept mpfd
// @Produce json
// @Param id path int true "Organization ID"
// @Param logo_print formData file true "Print logo image (PNG, JPEG, WebP, SVG, max 10MB)"
// @Success 200 {object} models.OrganizationAPIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 415 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id}/printlogo [post]
func (h *OrganizationHandler) UploadLogoPrint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Invalid organization ID"))
		return
	}

	tenantID, err := baseAPI.GetTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, baseAPI.ErrorResponseFunc("Unauthorized", "Unable to get tenant ID: "+err.Error()))
		return
	}

	file, err := c.FormFile("logo_print")
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Print logo file is required"))
		return
	}

	if file.Size > maxLogoSize {
		c.JSON(http.StatusRequestEntityTooLarge, baseAPI.ErrorResponseFunc("File too large", "Print logo must be smaller than 10MB"))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Internal server error", "Unable to read print logo file"))
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, baseAPI.ErrorResponseFunc("Internal server error", "Unable to read print logo file"))
		return
	}

	contentType := file.Header.Get("Content-Type")
	if !supportedMimeTypes[contentType] {
		c.JSON(http.StatusUnsupportedMediaType, baseAPI.ErrorResponseFunc("Unsupported media type", "Only PNG, JPEG, WebP, and SVG formats are supported"))
		return
	}

	organization, err := h.service.UpdateOrganizationLogoPrint(uint(id), tenantID, fileBytes, contentType)
	if err != nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, baseAPI.SuccessResponse("Print logo uploaded successfully", organization.ToResponse()))
}

// GetLogo handles retrieving an organization logo without authentication
// @Summary Get organization logo
// @Description Get the logo image for an organization (public, no auth required)
// @Tags organizations
// @ID getOrganizationLogo
// @Produce image/png,image/jpeg,image/webp,image/svg+xml
// @Param organization_id query int true "Organization ID"
// @Success 200 {file} binary "Logo image"
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organization/logo [get]
func (h *OrganizationHandler) GetLogo(c *gin.Context) {
	organizationID := c.Query("organization_id")
	if organizationID == "" {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Organization ID is required"))
		return
	}

	id, err := strconv.ParseUint(organizationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Invalid organization ID"))
		return
	}

	organization, err := h.service.GetOrganizationByIDPublic(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", "Organization logo not found"))
		return
	}

	if len(organization.Logo) == 0 {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", "Logo not found"))
		return
	}

	if organization.LogoContentType == "" {
		organization.LogoContentType = "image/png"
	}

	c.Header("Content-Type", organization.LogoContentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, organization.LogoContentType, organization.Logo)
}

// GetLogoPrint handles retrieving an organization print logo without authentication
// @Summary Get organization print logo
// @Description Get the print logo image for an organization (public, no auth required)
// @Tags organizations
// @ID getOrganizationLogoPrint
// @Produce image/png,image/jpeg,image/webp,image/svg+xml
// @Param organization_id query int true "Organization ID"
// @Success 200 {file} binary "Print logo image"
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organization/printlogo [get]
func (h *OrganizationHandler) GetLogoPrint(c *gin.Context) {
	organizationID := c.Query("organization_id")
	if organizationID == "" {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Organization ID is required"))
		return
	}

	id, err := strconv.ParseUint(organizationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, baseAPI.ErrorResponseFunc("Invalid request", "Invalid organization ID"))
		return
	}

	organization, err := h.service.GetOrganizationByIDPublic(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", "Organization print logo not found"))
		return
	}

	if len(organization.LogoPrint) == 0 {
		c.JSON(http.StatusNotFound, baseAPI.ErrorResponseFunc("Not found", "Print logo not found"))
		return
	}

	if organization.LogoPrintContentType == "" {
		organization.LogoPrintContentType = "image/png"
	}

	c.Header("Content-Type", organization.LogoPrintContentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, organization.LogoPrintContentType, organization.LogoPrint)
}
