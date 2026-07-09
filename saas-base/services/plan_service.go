package services

import (
	"context"

	"github.com/AgileExecutives/shared-modules/saas-base/models"
	"github.com/AgileExecutives/shared-modules/saas-base/repo"
)

// PlanService provides business logic for plans.
type PlanService struct {
	repo repo.PlanRepo
}

func NewPlanService(r repo.PlanRepo) *PlanService { return &PlanService{repo: r} }

func (s *PlanService) List() ([]models.Plan, error) {
	return s.repo.List(context.Background())
}

func (s *PlanService) GetByID(id uint) (*models.Plan, error) {
	return s.repo.FindByID(context.Background(), id)
}

func (s *PlanService) Save(p *models.Plan) error {
	return s.repo.Save(context.Background(), p)
}

func (s *PlanService) Delete(id uint) error {
	return s.repo.Delete(context.Background(), id)
}
