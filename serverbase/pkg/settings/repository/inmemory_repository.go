package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/AgileExecutives/ae-framework/serverbase/pkg/settings/entities"
)

// InMemorySettingsRepository is a simple in-memory implementation of SettingsRepositoryInterface
// intended for unit tests.
type InMemorySettingsRepository struct {
	mu          sync.RWMutex
	settings    map[uint]map[string]map[string]entities.Setting // tenant -> domain -> key
	definitions map[string]entities.SettingDefinition           // domain:key -> def
}

func NewInMemorySettingsRepository() *InMemorySettingsRepository {
	return &InMemorySettingsRepository{
		settings:    make(map[uint]map[string]map[string]entities.Setting),
		definitions: make(map[string]entities.SettingDefinition),
	}
}

// Helper key for definitions
func defKey(domain, key string) string { return domain + ":" + key }

// Setting definition methods
func (r *InMemorySettingsRepository) GetSettingDefinition(domain, key string) (*entities.SettingDefinition, error) {
	k := defKey(domain, key)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if d, ok := r.definitions[k]; ok {
		return &d, nil
	}
	return nil, nil
}

func (r *InMemorySettingsRepository) CreateSettingDefinition(def *entities.SettingDefinition) error {
	if def == nil {
		return errors.New("nil definition")
	}
	k := defKey(def.Domain, def.Key)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.definitions[k] = *def
	return nil
}

func (r *InMemorySettingsRepository) UpdateSettingDefinition(def *entities.SettingDefinition) error {
	return r.CreateSettingDefinition(def)
}

func (r *InMemorySettingsRepository) GetAllSettingDefinitions() ([]entities.SettingDefinition, error) {
	res := []entities.SettingDefinition{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.definitions {
		res = append(res, v)
	}
	return res, nil
}

func (r *InMemorySettingsRepository) GetSettingDefinitionsByDomain(domain string) ([]entities.SettingDefinition, error) {
	res := []entities.SettingDefinition{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.definitions {
		if v.Domain == domain {
			res = append(res, v)
		}
	}
	return res, nil
}

// Tenant settings
func (r *InMemorySettingsRepository) GetSetting(tenantID uint, domain, key string) (*entities.Setting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.settings[tenantID]; !ok {
		return nil, nil
	}
	if _, ok := r.settings[tenantID][domain]; !ok {
		return nil, nil
	}
	if s, ok := r.settings[tenantID][domain][key]; ok {
		return &s, nil
	}
	return nil, nil
}

func (r *InMemorySettingsRepository) SetSetting(setting *entities.Setting) error {
	if setting == nil {
		return errors.New("nil setting")
	}
	if setting.CreatedAt.IsZero() {
		setting.CreatedAt = time.Now()
	}
	setting.UpdatedAt = time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.settings[setting.TenantID]; !ok {
		r.settings[setting.TenantID] = make(map[string]map[string]entities.Setting)
	}
	if _, ok := r.settings[setting.TenantID][setting.Domain]; !ok {
		r.settings[setting.TenantID][setting.Domain] = make(map[string]entities.Setting)
	}
	r.settings[setting.TenantID][setting.Domain][setting.Key] = *setting
	return nil
}

func (r *InMemorySettingsRepository) GetDomainSettings(tenantID uint, domain string) ([]entities.Setting, error) {
	res := []entities.Setting{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.settings[tenantID]; !ok {
		return res, nil
	}
	if _, ok := r.settings[tenantID][domain]; !ok {
		return res, nil
	}
	for _, s := range r.settings[tenantID][domain] {
		res = append(res, s)
	}
	return res, nil
}

func (r *InMemorySettingsRepository) GetAllSettings(tenantID uint) ([]entities.Setting, error) {
	res := []entities.Setting{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.settings[tenantID]; !ok {
		return res, nil
	}
	for _, domainMap := range r.settings[tenantID] {
		for _, s := range domainMap {
			res = append(res, s)
		}
	}
	return res, nil
}

func (r *InMemorySettingsRepository) DeleteSetting(tenantID uint, domain, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.settings[tenantID]; !ok {
		return errors.New("setting not found")
	}
	if _, ok := r.settings[tenantID][domain]; !ok {
		return errors.New("setting not found")
	}
	if _, ok := r.settings[tenantID][domain][key]; !ok {
		return errors.New("setting not found")
	}
	delete(r.settings[tenantID][domain], key)
	return nil
}

func (r *InMemorySettingsRepository) GetDomains(tenantID uint) ([]string, error) {
	res := []string{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.settings[tenantID]; !ok {
		return res, nil
	}
	for domain := range r.settings[tenantID] {
		res = append(res, domain)
	}
	return res, nil
}

// Lifecycle
func (r *InMemorySettingsRepository) AutoMigrate() error { return nil }
func (r *InMemorySettingsRepository) HealthCheck() error { return nil }
