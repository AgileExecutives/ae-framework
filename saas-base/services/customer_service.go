package services

import (
	"context"

	"github.com/AgileExecutives/serverbase/modules/customers/repo"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/saas-base/models"
)

// CustomerService provides customer business logic.
type CustomerService struct {
	repo   repo.CustomerRepo
	logger core.Logger
}

// NewCustomerService creates a new CustomerService accepting a CustomerRepo.
func NewCustomerService(r repo.CustomerRepo, logger core.Logger) *CustomerService {
	return &CustomerService{repo: r, logger: logger}
}

// GetByTenant returns all customers for a given tenant.
func (s *CustomerService) GetByTenant(tenantID uint) ([]models.Customer, error) {
	return s.repo.FindByTenant(context.Background(), tenantID)
}

// CustomerServiceProvider implements core.ServiceProvider.
type CustomerServiceProvider struct {
	service *CustomerService
}

// NewCustomerServiceProvider creates a new CustomerServiceProvider.
func NewCustomerServiceProvider(service *CustomerService) core.ServiceProvider {
	return &CustomerServiceProvider{service: service}
}

func (p *CustomerServiceProvider) ServiceName() string {
	return "saas-base-customer"
}

func (p *CustomerServiceProvider) ServiceInterface() interface{} {
	return (*CustomerService)(nil)
}

func (p *CustomerServiceProvider) Factory(ctx core.ModuleContext) (interface{}, error) {
	return p.service, nil
}
