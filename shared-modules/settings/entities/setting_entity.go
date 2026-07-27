package entities

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	settingsEntities "github.com/AgileExecutives/serverbase/pkg/settings/entities"
)

type SettingEntity struct{}

func NewSettingEntity() core.Entity                      { return &SettingEntity{} }
func (e *SettingEntity) TableName() string               { return "settings" }
func (e *SettingEntity) GetModel() interface{}           { return &settingsEntities.Setting{} }
func (e *SettingEntity) GetMigrations() []core.Migration { return []core.Migration{} }
