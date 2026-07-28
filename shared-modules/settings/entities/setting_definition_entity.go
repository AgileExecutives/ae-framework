package entities

import (
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	settingsEntities "github.com/AgileExecutives/ae-framework/serverbase/pkg/settings/entities"
)

type SettingDefinitionEntity struct{}

func NewSettingDefinitionEntity() core.Entity        { return &SettingDefinitionEntity{} }
func (e *SettingDefinitionEntity) TableName() string { return "setting_definitions" }
func (e *SettingDefinitionEntity) GetModel() interface{} {
	return &settingsEntities.SettingDefinition{}
}
func (e *SettingDefinitionEntity) GetMigrations() []core.Migration { return []core.Migration{} }
