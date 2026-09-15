package services

import (
	"context"
	"fmt"
	"log"

	"github.com/AgileExecutives/ae-framework/serverbase/internal/models"
	"github.com/AgileExecutives/ae-framework/serverbase/modules/tenant/repo"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/AgileExecutives/ae-framework/serverbase/pkg/utils"
)

// TenantService handles tenant operations
type tenantBucketAPI interface {
	CreateTenantBucket(ctx context.Context, tenantID uint) error
	BucketExists(ctx context.Context, tenantID uint) (bool, error)
}

// TenantService handles tenant operations
type TenantService struct {
	repo                repo.TenantRepo
	tenantBucketService tenantBucketAPI
	postCreate          func(tenantID uint) error
	eventBus            core.EventBus
}

// NewTenantService creates a new tenant service using a TenantRepo implementation
func NewTenantService(r repo.TenantRepo, tenantBucketService tenantBucketAPI) *TenantService {
	return &TenantService{repo: r, tenantBucketService: tenantBucketService}
}

// SetEventBus sets an EventBus so the service can publish tenant lifecycle events.
func (s *TenantService) SetEventBus(b core.EventBus) {
	s.eventBus = b
}

// SetPostCreateHook sets a hook that will be invoked after a tenant is created.
// The hook receives the numeric tenant ID and may perform non-critical side effects
// such as contract registration or bucket provisioning. Errors from the hook are
// logged but do not fail tenant creation.
func (s *TenantService) SetPostCreateHook(h func(tenantID uint) error) {
	s.postCreate = h
}

// CreateTenant creates a new tenant and its MinIO bucket
func (s *TenantService) CreateTenant(ctx context.Context, req models.TenantCreateRequest) (*models.Tenant, error) {
	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		base := utils.GenerateSlug(req.Name)
		// fetch existing slugs
		existing := []string{}
		if tenants, err := s.repo.List(ctx); err == nil {
			for _, t := range tenants {
				existing = append(existing, t.Slug)
			}
		}
		slug = utils.EnsureUniqueSlug(base, existing)
	}

	// Create tenant in database
	tenant := models.Tenant{
		CustomerID: req.CustomerID,
		Name:       req.Name,
		Slug:       slug,
	}

	if err := s.repo.Save(ctx, &tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	log.Printf("✅ Created tenant: %s (ID: %d)", tenant.Name, tenant.ID)

	// Create MinIO bucket for the tenant
	if s.tenantBucketService != nil {
		if err := s.tenantBucketService.CreateTenantBucket(ctx, tenant.ID); err != nil {
			log.Printf("❌ Warning: Failed to create MinIO bucket for tenant %d: %v", tenant.ID, err)
			// Don't fail tenant creation if bucket creation fails
		}
	}

	// Invoke post-create hook (e.g., register contracts) if provided. Do not fail tenant creation on hook errors.
	if s.postCreate != nil {
		if err := s.postCreate(tenant.ID); err != nil {
			log.Printf("❌ Warning: Post-create hook failed for tenant %d: %v", tenant.ID, err)
		}
	}

	// Publish tenant.created event if event bus available. Do not fail creation on publish errors.
	if s.eventBus != nil {
		if err := s.eventBus.Publish("tenant.created", tenant.ID); err != nil {
			log.Printf("❌ Warning: failed to publish tenant.created for %d: %v", tenant.ID, err)
		}
	}

	return &tenant, nil
}

// CreateTenantWithoutBucket creates a tenant without creating MinIO bucket (for migration/seeding)
func (s *TenantService) CreateTenantWithoutBucket(ctx context.Context, req models.TenantCreateRequest) (*models.Tenant, error) {
	tenant := models.Tenant{
		CustomerID: req.CustomerID,
		Name:       req.Name,
		Slug:       req.Slug,
	}

	if err := s.repo.Save(ctx, &tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}
	// Invoke post-create hook if set
	if s.postCreate != nil {
		if err := s.postCreate(tenant.ID); err != nil {
			log.Printf("❌ Warning: Post-create hook failed for tenant %d: %v", tenant.ID, err)
		}
	}

	// Publish tenant.created event if event bus available
	if s.eventBus != nil {
		if err := s.eventBus.Publish("tenant.created", tenant.ID); err != nil {
			log.Printf("❌ Warning: failed to publish tenant.created for %d: %v", tenant.ID, err)
		}
	}

	return &tenant, nil
}

// EnsureTenantBucket ensures a tenant has a MinIO bucket created
func (s *TenantService) EnsureTenantBucket(ctx context.Context, tenantID uint) error {
	if s.tenantBucketService == nil {
		return fmt.Errorf("tenant bucket service not available")
	}

	// Check if bucket already exists
	exists, err := s.tenantBucketService.BucketExists(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if exists {
		log.Printf("✅ MinIO bucket already exists for tenant %d", tenantID)
		return nil
	}

	// Create bucket
	return s.tenantBucketService.CreateTenantBucket(ctx, tenantID)
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(ctx context.Context, tenantID uint) (*models.Tenant, error) {
	t, err := s.repo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if t == nil {
		return nil, fmt.Errorf("tenant not found")
	}
	return t, nil
}

// GetAllTenants retrieves all tenants
func (s *TenantService) GetAllTenants(ctx context.Context) ([]models.Tenant, error) {
	return s.repo.List(ctx)
}
