package services

import (
	"github.com/AgileExecutives/ae-framework/serverbase/modules/base/models"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/base/repo"
)

type PlanService struct {
	repo repo.PlanRepo
}

func NewPlanService(r repo.PlanRepo) *PlanService { return &PlanService{repo: r} }

func (s *PlanService) List() ([]models.Plan, error) { return s.repo.List() }

func (s *PlanService) GetByID(id uint) (*models.Plan, error) { return s.repo.GetByID(id) }

func (s *PlanService) Save(p *models.Plan) error { return s.repo.Save(p) }

func (s *PlanService) Delete(id uint) error { return s.repo.Delete(id) }
