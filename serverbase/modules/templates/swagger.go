package templates

import (
	"github.com/gin-gonic/gin"
)

// The following handler stubs are documentation-only and are not used at
// runtime. They exist so swag can discover the API paths and models for the
// templates module when generating swagger docs.

// ListTemplates godoc
// @Summary List templates
// @ID listTemplates
// @Tags templates
// @Produce json
// @Param template_type query string false "Template type"
// @Param channel query string false "Channel"
// @Success 200 {object} ListTemplatesResponse
// @Router /templates [get]
func listTemplatesDoc(c *gin.Context) {}

// CreateTemplate godoc
// @Summary Create a template
// @ID createTemplate
// @Tags templates
// @Accept json
// @Produce json
// @Param payload body CreateTemplateRequest true "Create template payload"
// @Success 201 {object} TemplateWithContentResponse
// @Router /templates [post]
func createTemplateDoc(c *gin.Context) {}

// GetTemplate godoc
// @Summary Get template by id
// @ID getTemplateById
// @Tags templates
// @Produce json
// @Param id path int true "Template ID"
// @Success 200 {object} TemplateWithContentResponse
// @Router /templates/{id} [get]
func getTemplateDoc(c *gin.Context) {}

// UpdateTemplate godoc
// @Summary Update template
// @ID updateTemplate
// @Tags templates
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Param payload body UpdateTemplateRequest true "Update payload"
// @Success 200 {object} TemplateWithContentResponse
// @Router /templates/{id} [put]
func updateTemplateDoc(c *gin.Context) {}

// DeleteTemplate godoc
// @Summary Delete template
// @ID deleteTemplate
// @Tags templates
// @Produce json
// @Param id path int true "Template ID"
// @Success 204 {object} SuccessResponse
// @Router /templates/{id} [delete]
func deleteTemplateDoc(c *gin.Context) {}

// RenderTemplateByID godoc
// @Summary Render template by id
// @ID renderTemplateById
// @Tags templates
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Param payload body RenderTemplateRequest true "Render payload"
// @Success 200 {object} TemplateWithContentResponse
// @Router /templates/{id}/render [post]
func renderTemplateByIDDoc(c *gin.Context) {}

// RenderTemplateByKey godoc
// @Summary Render template by key
// @ID renderTemplateByKey
// @Tags templates
// @Accept json
// @Produce json
// @Param payload body RenderTemplateRequest true "Render payload"
// @Success 200 {object} TemplateWithContentResponse
// @Router /templates/render [post]
func renderTemplateByKeyDoc(c *gin.Context) {}

// DuplicateTemplate godoc
// @Summary Duplicate a template
// @ID duplicateTemplate
// @Tags templates
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Param payload body DuplicateTemplateRequest true "Duplicate overrides"
// @Success 201 {object} TemplateWithContentResponse
// @Router /templates/{id}/duplicate [post]
func duplicateTemplateDoc(c *gin.Context) {}

// DefaultTemplate godoc
// @Summary Get default template
// @ID getDefaultTemplate
// @Tags templates
// @Produce json
// @Param template_type query string true "Template type"
// @Param channel query string true "Channel"
// @Success 200 {object} TemplateWithContentResponse
// @Router /templates/default [get]
func defaultTemplateDoc(c *gin.Context) {}

// ListContracts godoc
// @Summary List template contracts
// @ID listTemplateContracts
// @Tags templates
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /templates/contracts [get]
func listContractsDoc(c *gin.Context) {}

// GetContract godoc
// @Summary Get contract by key
// @ID getTemplateContract
// @Tags templates
// @Produce json
// @Param key path string true "Contract key"
// @Success 200 {object} map[string]interface{}
// @Router /templates/contracts/{key} [get]
func getContractDoc(c *gin.Context) {}

// ValidateContract godoc
// @Summary Validate data against contract
// @ID validateTemplateContract
// @Tags templates
// @Accept json
// @Produce json
// @Param key path string true "Contract key"
// @Param payload body map[string]interface{} true "Data to validate"
// @Success 200 {object} map[string]interface{}
// @Router /templates/contracts/{key}/validate [post]
func validateContractDoc(c *gin.Context) {}

// ensure local types are referenced so they appear in generated docs
var _ = ListTemplatesResponse{}
