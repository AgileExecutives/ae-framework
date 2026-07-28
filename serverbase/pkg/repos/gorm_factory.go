package repos

import (
	"github.com/AgileExecutives/ae-framework/serverbase/modules/customers/repo"
	emailrepo "github.com/AgileExecutives/ae-framework/serverbase/modules/email/repo"
	tenantrepo "github.com/AgileExecutives/ae-framework/serverbase/modules/tenant/repo"
	userrepo "github.com/AgileExecutives/ae-framework/serverbase/modules/user/repo"
	"gorm.io/gorm"
)

// GormRepoFactory centralizes creation of GORM-backed repo implementations.
type GormRepoFactory struct {
	db *gorm.DB
}

// NewGormRepoFactory returns a factory bound to the provided DB.
func NewGormRepoFactory(db *gorm.DB) *GormRepoFactory { return &GormRepoFactory{db: db} }

func (f *GormRepoFactory) TenantRepo() tenantrepo.TenantRepo {
	return tenantrepo.NewGormTenantRepo(f.db)
}
func (f *GormRepoFactory) CustomerRepo() repo.CustomerRepo { return repo.NewGormCustomerRepo(f.db) }
func (f *GormRepoFactory) UserRepo() userrepo.UserRepo     { return userrepo.NewGormUserRepo(f.db) }
func (f *GormRepoFactory) ContactRepo() userrepo.ContactRepo {
	return userrepo.NewGormContactRepo(f.db)
}
func (f *GormRepoFactory) NewsletterRepo() userrepo.NewsletterRepo {
	return userrepo.NewGormUserRepo(f.db)
}
func (f *GormRepoFactory) TokenBlacklistRepo() userrepo.TokenBlacklistRepo {
	return userrepo.NewGormUserRepo(f.db)
}
func (f *GormRepoFactory) EmailRepo() emailrepo.EmailRepo { return emailrepo.NewGormEmailRepo(f.db) }
