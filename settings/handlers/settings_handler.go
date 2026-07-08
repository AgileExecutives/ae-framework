package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AgileExecutives/serverbase/pkg/settings/entities"
	"github.com/AgileExecutives/serverbase/pkg/settings/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DomainSettingsResponse struct {
	Domain   string                 `json:"domain" example:"invoice"`
	Settings map[string]interface{} `json:"settings" swaggertype:"object"`
}
type SettingResponse struct {
	Domain string                 `json:"domain" example:"invoice"`
	Key    string                 `json:"key" example:"invoice_prefix"`
	Data   map[string]interface{} `json:"data" swaggertype:"object"`
}
type UpdateKeysResponse struct {
	Message     string   `json:"message" example:"Settings updated successfully"`
	Domain      string   `json:"domain" example:"invoice"`
	UpdatedKeys []string `json:"updated_keys" example:"invoice_prefix,next_invoice_number"`
}
type SimpleMessageResponse struct {
	Message string `json:"message" example:"Setting updated successfully"`
	Domain  string `json:"domain" example:"invoice"`
	Key     string `json:"key" example:"invoice_prefix"`
}
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid tenant ID"`
}

type SettingsHandler struct {
	repo *repository.SettingsRepository
}

func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{repo: repository.NewSettingsRepository(db)}
}

// Methods follow original implementation
func (h *SettingsHandler) GetDomainSettings(c *gin.Context) {
	tenantIDStr := c.Param("tenant_id")
	domain := c.Param("domain")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tenant ID"})
		return
	}
	settings, err := h.repo.GetDomainSettings(uint(tenantID), domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	result := make(map[string]interface{})
	for _, setting := range settings {
		var data map[string]interface{}
		if err := json.Unmarshal(setting.Data, &data); err == nil {
			for k, v := range data {
				result[k] = v
			}
		}
	}
	c.JSON(http.StatusOK, DomainSettingsResponse{Domain: domain, Settings: result})
}

func (h *SettingsHandler) GetSetting(c *gin.Context) {
	tenantIDStr := c.Param("tenant_id")
	domain := c.Param("domain")
	key := c.Param("key")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tenant ID"})
		return
	}
	setting, err := h.repo.GetSetting(uint(tenantID), domain, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if setting == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Setting not found"})
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(setting.Data, &data); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to parse setting data"})
		return
	}
	c.JSON(http.StatusOK, SettingResponse{Domain: domain, Key: key, Data: data})
}

func (h *SettingsHandler) UpdateSetting(c *gin.Context) {
	tenantIDStr := c.Param("tenant_id")
	domain := c.Param("domain")
	key := c.Param("key")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tenant ID"})
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}
	settingDef, err := h.repo.GetSettingDefinition(domain, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if settingDef == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Setting definition not found"})
		return
	}
	setting, err := h.repo.GetSetting(uint(tenantID), domain, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	dataJSON, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to encode data"})
		return
	}
	if setting == nil {
		setting = &entities.Setting{TenantID: uint(tenantID), Domain: domain, Key: key, Version: settingDef.Version, Data: dataJSON, SettingDefinitionID: settingDef.ID}
	} else {
		setting.Data = dataJSON
		setting.Version = settingDef.Version
	}
	if err := h.repo.SetSetting(setting); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SimpleMessageResponse{Message: "Setting updated successfully", Domain: domain, Key: key})
}

func (h *SettingsHandler) UpdateDomainSettings(c *gin.Context) {
	tenantIDStr := c.Param("tenant_id")
	domain := c.Param("domain")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tenant ID"})
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}
	updatedKeys := make([]string, 0, len(req))
	for key, value := range req {
		settingDef, err := h.repo.GetSettingDefinition(domain, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get setting definition for " + key})
			return
		}
		if settingDef == nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Setting definition not found for " + key})
			return
		}
		setting, err := h.repo.GetSetting(uint(tenantID), domain, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get setting for " + key})
			return
		}
		dataJSON, err := json.Marshal(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to encode data for " + key})
			return
		}
		if setting == nil {
			setting = &entities.Setting{TenantID: uint(tenantID), Domain: domain, Key: key, Version: settingDef.Version, Data: dataJSON, SettingDefinitionID: settingDef.ID}
		} else {
			setting.Data = dataJSON
			setting.Version = settingDef.Version
		}
		if err := h.repo.SetSetting(setting); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save setting for " + key})
			return
		}
		updatedKeys = append(updatedKeys, key)
	}
	c.JSON(http.StatusOK, UpdateKeysResponse{Message: "Settings updated successfully", Domain: domain, UpdatedKeys: updatedKeys})
}
