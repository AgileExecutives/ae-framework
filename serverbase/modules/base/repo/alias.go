package repo

import (
    "gorm.io/gorm"
)

// Provide lightweight local constructors to avoid importing shared-modules during migration.
// Use a real Gorm-backed repo when a DB is provided; otherwise fall back to in-memory.
func NewGormPlanRepo(db *gorm.DB) PlanRepo {
    if db != nil {
        return newGormPlanRepo(db)
    }
    return NewInMemoryPlanRepo()
}

func NewGormNewsletterRepo(db *gorm.DB) NewsletterRepo { return nil }
