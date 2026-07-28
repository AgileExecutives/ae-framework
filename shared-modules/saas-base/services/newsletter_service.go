package services

import (
	"context"

	"github.com/AgileExecutives/ae-framwork/shared-modules/saas-base/models"
	"github.com/AgileExecutives/ae-framwork/shared-modules/saas-base/repo"
)

// NewsletterService provides business logic for newsletter subscriptions.
type NewsletterService struct {
	repo repo.NewsletterRepo
}

func NewNewsletterService(r repo.NewsletterRepo) *NewsletterService {
	return &NewsletterService{repo: r}
}

func (s *NewsletterService) List() ([]models.Newsletter, error) {
	return s.repo.List(context.Background())
}

func (s *NewsletterService) FindByEmail(email string) (*models.Newsletter, error) {
	return s.repo.FindByEmail(context.Background(), email)
}

func (s *NewsletterService) GetByID(id uint) (*models.Newsletter, error) {
	return s.repo.FindByID(context.Background(), id)
}

func (s *NewsletterService) Save(n *models.Newsletter) error {
	return s.repo.Save(context.Background(), n)
}

func (s *NewsletterService) Delete(id uint) error {
	return s.repo.Delete(context.Background(), id)
}
