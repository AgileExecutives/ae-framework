package entities

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"gorm.io/gorm"
)

type OrganizationEntity struct{}

func NewOrganizationEntity() *OrganizationEntity    { return &OrganizationEntity{} }
func (e *OrganizationEntity) Name() string          { return "organization" }
func (e *OrganizationEntity) TableName() string     { return "organizations" }
func (e *OrganizationEntity) GetModel() interface{} { return &models.Organization{} }
func (e *OrganizationEntity) GetMigrations() []core.Migration {
	return []core.Migration{&OrgMigration001{}}
}

type OrgMigration001 struct{}

func (m *OrgMigration001) Version() string      { return "001_create_organizations" }
func (m *OrgMigration001) Up(db *gorm.DB) error { return db.AutoMigrate(&models.Organization{}) }
func (m *OrgMigration001) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Organization{})
}
