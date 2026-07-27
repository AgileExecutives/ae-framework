package repo

import (
	"context"

	"github.com/AgileExecutives/shared-modules/saas-base/models"
	"gorm.io/gorm"
)

type gormPlanRepo struct{ db *gorm.DB }
type gormNewsletterRepo struct{ db *gorm.DB }

func NewGormPlanRepo(db *gorm.DB) PlanRepo             { return &gormPlanRepo{db: db} }
func NewGormNewsletterRepo(db *gorm.DB) NewsletterRepo { return &gormNewsletterRepo{db: db} }

func (r *gormPlanRepo) List(ctx context.Context) ([]models.Plan, error) {
	var res []models.Plan
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *gormPlanRepo) FindByID(ctx context.Context, id uint) (*models.Plan, error) {
	var p models.Plan
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *gormPlanRepo) Save(ctx context.Context, p *models.Plan) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *gormPlanRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Plan{}, id).Error
}

func (r *gormNewsletterRepo) List(ctx context.Context) ([]models.Newsletter, error) {
	var res []models.Newsletter
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *gormNewsletterRepo) FindByEmail(ctx context.Context, email string) (*models.Newsletter, error) {
	var n models.Newsletter
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&n).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func (r *gormNewsletterRepo) FindByID(ctx context.Context, id uint) (*models.Newsletter, error) {
	var n models.Newsletter
	if err := r.db.WithContext(ctx).First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *gormNewsletterRepo) Save(ctx context.Context, n *models.Newsletter) error {
	return r.db.WithContext(ctx).Save(n).Error
}

func (r *gormNewsletterRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Newsletter{}, id).Error
}

var _ PlanRepo = (*gormPlanRepo)(nil)
var _ NewsletterRepo = (*gormNewsletterRepo)(nil)
