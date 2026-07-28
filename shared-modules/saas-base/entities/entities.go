package entities

import (
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framework/shared-modules/saas-base/models"
)

// PlanEntity implements core.Entity for the Plan model.
type PlanEntity struct{}

func NewPlanEntity() core.Entity {
	return &PlanEntity{}
}

func (e *PlanEntity) TableName() string {
	return "plans"
}

func (e *PlanEntity) GetModel() interface{} {
	return &models.Plan{}
}

func (e *PlanEntity) GetMigrations() []core.Migration {
	return []core.Migration{}
}

// CustomerEntity implements core.Entity for the Customer model.
type CustomerEntity struct{}

func NewCustomerEntity() core.Entity {
	return &CustomerEntity{}
}

func (e *CustomerEntity) TableName() string {
	return "customers"
}

func (e *CustomerEntity) GetModel() interface{} {
	return &models.Customer{}
}

func (e *CustomerEntity) GetMigrations() []core.Migration {
	return []core.Migration{}
}

// NewsletterEntity implements core.Entity for the Newsletter model.
type NewsletterEntity struct{}

func NewNewsletterEntity() core.Entity {
	return &NewsletterEntity{}
}

func (e *NewsletterEntity) TableName() string {
	return "newsletters"
}

func (e *NewsletterEntity) GetModel() interface{} {
	return &models.Newsletter{}
}

func (e *NewsletterEntity) GetMigrations() []core.Migration {
	return []core.Migration{}
}
