package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AgileExecutives/shared-modules/booking/entities"
	repo "github.com/AgileExecutives/shared-modules/booking/repo"
	"gorm.io/gorm"
)

type BookingService struct {
	db   *gorm.DB
	repo repo.BookingRepo
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{db: db}
}

// NewBookingServiceWithRepo creates a BookingService backed by a repository
func NewBookingServiceWithRepo(r repo.BookingRepo) *BookingService {
	return &BookingService{repo: r}
}

// NewBookingServiceWithRepoAndDB creates a BookingService backed by a repository and a DB for additional queries
func NewBookingServiceWithRepoAndDB(r repo.BookingRepo, db *gorm.DB) *BookingService {
	return &BookingService{repo: r, db: db}
}

// ClientInfo represents simplified client information returned to handlers
type ClientInfo struct {
	ID              uint       `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	Gender          string     `json:"gender,omitempty"`
	PrimaryLanguage string     `json:"primary_language,omitempty"`
	StreetAddress   string     `json:"street_address,omitempty"`
	Zip             string     `json:"zip,omitempty"`
	City            string     `json:"city,omitempty"`
	Status          string     `json:"status"`
}

// GetClientInfo fetches basic client information by id and tenant. Uses DB when available.
func (s *BookingService) GetClientInfo(clientID, tenantID uint) (*ClientInfo, error) {
	if s.db == nil {
		return nil, errors.New("database not available for client lookup")
	}

	var client ClientInfo
	err := s.db.Table("clients").
		Select("id, first_name, last_name, email, phone, date_of_birth, gender, primary_language, street_address, zip, city, status").
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", clientID, tenantID).
		First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// CreateConfiguration creates a new booking configuration
func (s *BookingService) CreateConfiguration(req entities.CreateBookingTemplateRequest, tenantID uint) (*entities.BookingTemplate, error) {
	config := &entities.BookingTemplate{
		UserID:              req.UserID,
		CalendarID:          req.CalendarID,
		TenantID:            tenantID,
		Name:                req.Name,
		Description:         req.Description,
		SlotDuration:        req.SlotDuration,
		BufferTime:          req.BufferTime,
		MaxSeriesBookings:   req.MaxSeriesBookings,
		AllowedIntervals:    req.AllowedIntervals,
		NumberOfIntervals:   req.NumberOfIntervals,
		WeeklyAvailability:  req.WeeklyAvailability,
		AdvanceBookingDays:  req.AdvanceBookingDays,
		MinNoticeHours:      req.MinNoticeHours,
		Timezone:            req.Timezone,
		MaxBookingsPerDay:   req.MaxBookingsPerDay,
		AllowBackToBack:     req.AllowBackToBack,
		BlockDates:          req.BlockDates,
		AllowedStartMinutes: entities.MinutesArray(req.AllowedStartMinutes),
	}

	if s.repo != nil {
		if err := s.repo.CreateConfiguration(context.Background(), config); err != nil {
			return nil, fmt.Errorf("failed to create booking configuration: %w", err)
		}
		return config, nil
	}

	if err := s.db.Create(config).Error; err != nil {
		return nil, fmt.Errorf("failed to create booking configuration: %w", err)
	}

	return config, nil
}

// GetConfiguration retrieves a booking configuration by ID
func (s *BookingService) GetConfiguration(id uint, tenantID uint) (*entities.BookingTemplate, error) {
	if s.repo != nil {
		cfg, err := s.repo.FindConfiguration(context.Background(), id, tenantID)
		if err != nil {
			return nil, errors.New("booking configuration not found")
		}
		return cfg, nil
	}

	var config entities.BookingTemplate

	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking configuration not found")
		}
		return nil, fmt.Errorf("failed to retrieve booking configuration: %w", err)
	}

	return &config, nil
}

// GetAllConfigurations retrieves all booking configurations for a tenant
func (s *BookingService) GetAllConfigurations(tenantID uint, page, limit int) ([]entities.BookingTemplate, int64, error) {
	if s.repo != nil {
		list, total, err := s.repo.ListConfigurations(context.Background(), tenantID, page, limit)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to retrieve booking configurations: %w", err)
		}
		return list, total, nil
	}

	var configs []entities.BookingTemplate
	var total int64

	query := s.db.Model(&entities.BookingTemplate{}).Where("tenant_id = ?", tenantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count booking configurations: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&configs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve booking configurations: %w", err)
	}

	return configs, total, nil
}

// GetConfigurationsByUser retrieves all booking configurations for a specific user
func (s *BookingService) GetConfigurationsByUser(userID uint, tenantID uint) ([]entities.BookingTemplate, error) {
	if s.repo != nil {
		return s.repo.FindConfigurationsByUser(context.Background(), userID, tenantID)
	}

	var configs []entities.BookingTemplate

	if err := s.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve booking configurations for user: %w", err)
	}

	return configs, nil
}

// GetConfigurationsByCalendar retrieves all booking configurations for a specific calendar
func (s *BookingService) GetConfigurationsByCalendar(calendarID uint, tenantID uint) ([]entities.BookingTemplate, error) {
	if s.repo != nil {
		return s.repo.FindConfigurationsByCalendar(context.Background(), calendarID, tenantID)
	}

	var configs []entities.BookingTemplate

	if err := s.db.Where("calendar_id = ? AND tenant_id = ?", calendarID, tenantID).Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve booking configurations for calendar: %w", err)
	}

	return configs, nil
}

// UpdateConfiguration updates an existing booking configuration
func (s *BookingService) UpdateConfiguration(id uint, tenantID uint, req entities.UpdateBookingTemplateRequest) (*entities.BookingTemplate, error) {
	var config *entities.BookingTemplate
	var err error

	if s.repo != nil {
		config, err = s.repo.FindConfiguration(context.Background(), id, tenantID)
		if err != nil {
			return nil, errors.New("booking configuration not found")
		}
	} else {
		var cfg entities.BookingTemplate
		if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cfg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("booking configuration not found")
			}
			return nil, fmt.Errorf("failed to retrieve booking configuration: %w", err)
		}
		config = &cfg
	}

	// Update fields if provided
	if req.Name != nil {
		config.Name = *req.Name
	}
	if req.Description != nil {
		config.Description = *req.Description
	}
	if req.SlotDuration != nil {
		config.SlotDuration = *req.SlotDuration
	}
	if req.BufferTime != nil {
		config.BufferTime = *req.BufferTime
	}
	if req.MaxSeriesBookings != nil {
		config.MaxSeriesBookings = *req.MaxSeriesBookings
	}
	if req.AllowedIntervals != nil {
		config.AllowedIntervals = req.AllowedIntervals
	}
	if req.NumberOfIntervals != nil {
		config.NumberOfIntervals = *req.NumberOfIntervals
	}
	if req.WeeklyAvailability != nil {
		config.WeeklyAvailability = *req.WeeklyAvailability
	}
	if req.AdvanceBookingDays != nil {
		config.AdvanceBookingDays = *req.AdvanceBookingDays
	}
	if req.MinNoticeHours != nil {
		config.MinNoticeHours = *req.MinNoticeHours
	}
	if req.Timezone != nil {
		config.Timezone = *req.Timezone
	}
	if req.MaxBookingsPerDay != nil {
		config.MaxBookingsPerDay = req.MaxBookingsPerDay
	}
	if req.AllowBackToBack != nil {
		config.AllowBackToBack = req.AllowBackToBack
	}
	if req.BlockDates != nil {
		config.BlockDates = req.BlockDates
	}
	if req.AllowedStartMinutes != nil {
		config.AllowedStartMinutes = entities.MinutesArray(req.AllowedStartMinutes)
	}

	if s.repo != nil {
		if err := s.repo.UpdateConfiguration(context.Background(), config); err != nil {
			return nil, fmt.Errorf("failed to update booking configuration: %w", err)
		}
		return config, nil
	}

	if err := s.db.Save(&config).Error; err != nil {
		return nil, fmt.Errorf("failed to update booking configuration: %w", err)
	}

	return config, nil
}

// DeleteConfiguration soft deletes a booking configuration
func (s *BookingService) DeleteConfiguration(id uint, tenantID uint) error {
	if s.repo != nil {
		if err := s.repo.DeleteConfiguration(context.Background(), id, tenantID); err != nil {
			return fmt.Errorf("failed to delete booking configuration: %w", err)
		}
		return nil
	}

	result := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&entities.BookingTemplate{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete booking configuration: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("booking configuration not found")
	}

	return nil
}
