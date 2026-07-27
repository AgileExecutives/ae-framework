package repo

import (
	"context"

	"github.com/AgileExecutives/shared-modules/saas-base/models"
)

// PlanRepo defines repository methods for plans.
type PlanRepo interface {
	List(ctx context.Context) ([]models.Plan, error)
	FindByID(ctx context.Context, id uint) (*models.Plan, error)
	Save(ctx context.Context, p *models.Plan) error
	Delete(ctx context.Context, id uint) error
}

// NewsletterRepo defines repository methods for newsletter subscriptions.
type NewsletterRepo interface {
	List(ctx context.Context) ([]models.Newsletter, error)
	FindByEmail(ctx context.Context, email string) (*models.Newsletter, error)
	FindByID(ctx context.Context, id uint) (*models.Newsletter, error)
	Save(ctx context.Context, n *models.Newsletter) error
	Delete(ctx context.Context, id uint) error
}
