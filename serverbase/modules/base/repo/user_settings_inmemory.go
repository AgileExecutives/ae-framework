package repo

import (
	"errors"
	"sync"

	"github.com/AgileExecutives/ae-framework/serverbase/internal/models"
)

type InMemoryUserSettingsRepo struct {
	mu     sync.RWMutex
	byUser map[uint]*models.UserSettings
	next   uint
}

func NewInMemoryUserSettingsRepo() *InMemoryUserSettingsRepo {
	return &InMemoryUserSettingsRepo{byUser: make(map[uint]*models.UserSettings), next: 1}
}

func (r *InMemoryUserSettingsRepo) FindByUserID(userID uint) (models.UserSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.byUser[userID]; ok {
		// return a copy
		copy := *s
		return copy, nil
	}
	return models.UserSettings{}, errors.New("not found")
}

func (r *InMemoryUserSettingsRepo) Create(settings *models.UserSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byUser[settings.UserID]; exists {
		return errors.New("already exists")
	}
	if settings.ID == 0 {
		settings.ID = r.next
		r.next++
	}
	copy := *settings
	r.byUser[settings.UserID] = &copy
	return nil
}

func (r *InMemoryUserSettingsRepo) Save(settings *models.UserSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if settings.ID == 0 {
		settings.ID = r.next
		r.next++
	}
	copy := *settings
	r.byUser[settings.UserID] = &copy
	return nil
}
