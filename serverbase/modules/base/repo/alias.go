package repo

import (
    "gorm.io/gorm"
)

// Provide lightweight local constructors to avoid importing shared-modules during migration.
func NewGormPlanRepo(db *gorm.DB) PlanRepo { return NewInMemoryPlanRepo() }

func NewGormNewsletterRepo(db *gorm.DB) NewsletterRepo { return nil }
