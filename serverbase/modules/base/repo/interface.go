package repo

import (
    "context"

    "github.com/AgileExecutives/serverbase/modules/base/models"
)

// PlanRepo defines the repository responsibilities for plans.
type PlanRepo interface {
    GetByID(id uint) (*models.Plan, error)
    List() ([]models.Plan, error)
    Save(p *models.Plan) error
    Delete(id uint) error
}

// NewsletterRepo defines repository for newsletter subscriptions.
type NewsletterRepo interface {
    Save(ctx context.Context, n *models.Newsletter) error
    FindByEmail(ctx context.Context, email string) (*models.Newsletter, error)
}
