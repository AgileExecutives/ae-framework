package services

import (
	"context"

	basemodels "github.com/AgileExecutives/serverbase/modules/base/models"
	"github.com/AgileExecutives/serverbase/modules/customers/repo"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/AgileExecutives/shared-modules/saas-base/models"
	"gorm.io/gorm"
)

// CustomerService provides customer business logic.
type CustomerService struct {
	repo   repo.CustomerRepo
	logger core.Logger
	db     *gorm.DB
}

// NewCustomerService creates a new CustomerService accepting a CustomerRepo.
func NewCustomerService(r repo.CustomerRepo, logger core.Logger) *CustomerService {
	return &CustomerService{repo: r, logger: logger}
}

// NewCustomerServiceWithDB creates a CustomerService that also has access to the raw DB
// for operations (like soft delete) that are not exposed by the repo interface.
func NewCustomerServiceWithDB(r repo.CustomerRepo, db *gorm.DB, logger core.Logger) *CustomerService {
	return &CustomerService{repo: r, logger: logger, db: db}
}

// GetByTenant returns all customers for a given tenant.
func (s *CustomerService) GetByTenant(tenantID uint) ([]models.Customer, error) {
	baseList, err := s.repo.FindByTenant(context.Background(), tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Customer, 0, len(baseList))
	for _, b := range baseList {
		out = append(out, convertBaseCustomerToShared(&b))
	}
	return out, nil
}

// GetByID returns a single customer by ID.
func (s *CustomerService) GetByID(id uint) (*models.Customer, error) {
	b, err := s.repo.FindByID(context.Background(), id)
	if err != nil || b == nil {
		return nil, err
	}
	shared := convertBaseCustomerToShared(b)
	return &shared, nil
}

// FindByEmail finds a customer by email.
func (s *CustomerService) FindByEmail(email string) (*models.Customer, error) {
	b, err := s.repo.FindByEmail(context.Background(), email)
	if err != nil || b == nil {
		return nil, err
	}
	shared := convertBaseCustomerToShared(b)
	return &shared, nil
}

// Save persists a customer (create or update).
func (s *CustomerService) Save(c *models.Customer) error {
	// convert to base model expected by repo
	b := convertSharedCustomerToBase(c)
	return s.repo.Save(context.Background(), &b)
}

// Delete removes a customer by ID.
func (s *CustomerService) Delete(id uint) error {
	if s.db != nil {
		return s.db.Delete(&models.Customer{}, id).Error
	}
	// If no DB available, attempt to mark inactive as fallback
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

// convertBaseCustomerToShared converts a serverbase base Customer to the shared saas Customer
func convertBaseCustomerToShared(b *basemodels.Customer) models.Customer {
	return models.Customer{
		ID:            b.ID,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		DeletedAt:     b.DeletedAt,
		Name:          b.Name,
		Email:         b.Email,
		Phone:         b.Phone,
		Street:        b.Street,
		Zip:           b.Zip,
		City:          b.City,
		Country:       b.Country,
		TaxID:         b.TaxID,
		VAT:           b.VAT,
		PlanID:        b.PlanID,
		TenantID:      b.TenantID,
		Status:        b.Status,
		PaymentMethod: b.PaymentMethod,
		Active:        b.Active,
	}
}

// convertSharedCustomerToBase converts a shared saas Customer to the serverbase base Customer
func convertSharedCustomerToBase(s *models.Customer) basemodels.Customer {
	return basemodels.Customer{
		ID:            s.ID,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
		DeletedAt:     s.DeletedAt,
		Name:          s.Name,
		Email:         s.Email,
		Phone:         s.Phone,
		Street:        s.Street,
		Zip:           s.Zip,
		City:          s.City,
		Country:       s.Country,
		TaxID:         s.TaxID,
		VAT:           s.VAT,
		PlanID:        s.PlanID,
		TenantID:      s.TenantID,
		Status:        s.Status,
		PaymentMethod: s.PaymentMethod,
		Active:        s.Active,
	}
}
