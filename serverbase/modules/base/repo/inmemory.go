package repo

import (
    "errors"
    "sync"

    "github.com/AgileExecutives/serverbase/modules/base/models"
)

type InMemoryPlanRepo struct {
    mu   sync.RWMutex
    byID map[uint]*models.Plan
    next uint
}

func NewInMemoryPlanRepo() *InMemoryPlanRepo {
    return &InMemoryPlanRepo{byID: make(map[uint]*models.Plan), next: 1}
}

func (r *InMemoryPlanRepo) GetByID(id uint) (*models.Plan, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    if p, ok := r.byID[id]; ok {
        return p, nil
    }
    return nil, errors.New("not found")
}

func (r *InMemoryPlanRepo) List() ([]models.Plan, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    res := make([]models.Plan, 0, len(r.byID))
    for _, p := range r.byID {
        res = append(res, *p)
    }
    return res, nil
}

func (r *InMemoryPlanRepo) Save(p *models.Plan) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if p.ID == 0 {
        p.ID = r.next
        r.next++
    }
    copy := *p
    r.byID[p.ID] = &copy
    return nil
}

func (r *InMemoryPlanRepo) Delete(id uint) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.byID, id)
    return nil
}

