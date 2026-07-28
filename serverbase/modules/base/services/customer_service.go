package services

import (
	"context"

	"github.com/AgileExecutives/ae-framework/serverbase/modules/base/models"
	custrepo "github.com/AgileExecutives/ae-framework/serverbase/modules/customers/repo"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"gorm.io/gorm"
)

// CustomerService provides customer business logic for serverbase during migration.
type CustomerService struct {
	repo   custrepo.CustomerRepo
	logger core.Logger
	db     *gorm.DB
}

func NewCustomerServiceWithDB(r custrepo.CustomerRepo, db *gorm.DB, logger core.Logger) *CustomerService {
	return &CustomerService{repo: r, logger: logger, db: db}
}

func (s *CustomerService) GetByTenant(tenantID uint) ([]models.Customer, error) {
	return s.repo.FindByTenant(context.Background(), tenantID)
}

func (s *CustomerService) GetByID(id uint) (*models.Customer, error) {
	return s.repo.FindByID(context.Background(), id)
}

func (s *CustomerService) FindByEmail(email string) (*models.Customer, error) {
	return s.repo.FindByEmail(context.Background(), email)
}

func (s *CustomerService) Save(c *models.Customer) error {
	return s.repo.Save(context.Background(), c)
}

func (s *CustomerService) Delete(id uint) error {
	if s.db != nil {
		return s.db.Delete(&models.Customer{}, id).Error
	}
	c, err := s.repo.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	if c == nil {
		return gorm.ErrRecordNotFound
	}
	c.Active = false
	return s.repo.Save(context.Background(), c)
}
