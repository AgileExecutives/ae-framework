package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/ae-framwork/shared-modules/saas-base/models"
)

type inMemoryStore struct {
	mu          sync.RWMutex
	plans       map[uint]models.Plan
	newsletters map[uint]models.Newsletter
	nextID      uint
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{plans: make(map[uint]models.Plan), newsletters: make(map[uint]models.Newsletter), nextID: 1}
}

type InMemoryPlanRepo struct{ store *inMemoryStore }
type InMemoryNewsletterRepo struct{ store *inMemoryStore }

func NewInMemoryPlanRepo() PlanRepo { return &InMemoryPlanRepo{store: newInMemoryStore()} }
func NewInMemoryNewsletterRepo() NewsletterRepo {
	return &InMemoryNewsletterRepo{store: newInMemoryStore()}
}

func (r *InMemoryPlanRepo) List(ctx context.Context) ([]models.Plan, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	var out []models.Plan
	for _, p := range r.store.plans {
		out = append(out, p)
	}
	return out, nil
}

func (r *InMemoryPlanRepo) FindByID(ctx context.Context, id uint) (*models.Plan, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	p, ok := r.store.plans[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &p, nil
}

func (r *InMemoryPlanRepo) Save(ctx context.Context, p *models.Plan) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if p.ID == 0 {
		p.ID = r.store.nextID
		r.store.nextID++
	}
	r.store.plans[p.ID] = *p
	return nil
}

func (r *InMemoryPlanRepo) Delete(ctx context.Context, id uint) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	delete(r.store.plans, id)
	return nil
}

func (r *InMemoryNewsletterRepo) List(ctx context.Context) ([]models.Newsletter, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	var out []models.Newsletter
	for _, n := range r.store.newsletters {
		out = append(out, n)
	}
	return out, nil
}

func (r *InMemoryNewsletterRepo) FindByEmail(ctx context.Context, email string) (*models.Newsletter, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, n := range r.store.newsletters {
		if n.Email == email {
			return &n, nil
		}
	}
	return nil, nil
}

func (r *InMemoryNewsletterRepo) FindByID(ctx context.Context, id uint) (*models.Newsletter, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	n, ok := r.store.newsletters[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &n, nil
}

func (r *InMemoryNewsletterRepo) Save(ctx context.Context, n *models.Newsletter) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if n.ID == 0 {
		n.ID = r.store.nextID
		r.store.nextID++
	}
	r.store.newsletters[n.ID] = *n
	return nil
}

func (r *InMemoryNewsletterRepo) Delete(ctx context.Context, id uint) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	delete(r.store.newsletters, id)
	return nil
}

var _ PlanRepo = (*InMemoryPlanRepo)(nil)
var _ NewsletterRepo = (*InMemoryNewsletterRepo)(nil)
