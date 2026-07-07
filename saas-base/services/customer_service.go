package services

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/saas-base/models"
	"gorm.io/gorm"
)

// CustomerService provides customer business logic.
type CustomerService struct {
	db     *gorm.DB
	logger core.Logger
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(db *gorm.DB, logger core.Logger) *CustomerService {
	return &CustomerService{db: db, logger: logger}
}

// GetByTenant returns all customers for a given tenant.
func (s *CustomerService) GetByTenant(tenantID uint) ([]models.Customer, error) {
	var customers []models.Customer
	if err := s.db.Where("tenant_id = ?", tenantID).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
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
