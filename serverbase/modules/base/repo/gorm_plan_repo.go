package repo

import (
    "gorm.io/gorm"

    "github.com/AgileExecutives/ae-framework/serverbase/modules/base/models"
)

type gormPlanRepo struct{ db *gorm.DB }

func newGormPlanRepo(db *gorm.DB) PlanRepo { return &gormPlanRepo{db: db} }

func (r *gormPlanRepo) GetByID(id uint) (*models.Plan, error) {
    var p models.Plan
    if err := r.db.First(&p, id).Error; err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *gormPlanRepo) List() ([]models.Plan, error) {
    var res []models.Plan
    if err := r.db.Order("created_at DESC").Find(&res).Error; err != nil {
        return nil, err
    }
    return res, nil
}

func (r *gormPlanRepo) Save(p *models.Plan) error {
    return r.db.Save(p).Error
}

func (r *gormPlanRepo) Delete(id uint) error {
    return r.db.Delete(&models.Plan{}, id).Error
}
