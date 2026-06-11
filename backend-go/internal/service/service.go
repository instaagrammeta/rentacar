// Package service contains the business logic of the backend, mirroring the
// Flask service layer. A single Service struct shares the database, cache,
// configuration and media generator across all domain operations.
package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/auth"
	"github.com/instaagrammeta/rentacar/backend-go/internal/config"
	"github.com/instaagrammeta/rentacar/backend-go/internal/media"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/sms"
)

// Service bundles dependencies for the business logic.
type Service struct {
	DB     *gorm.DB
	Cfg    *config.Config
	Redis  *redis.Client
	Tokens *auth.Manager
	Media  *media.Generator
	SMS    *sms.Sender
}

// New builds a Service.
func New(db *gorm.DB, cfg *config.Config, rds *redis.Client, tokens *auth.Manager, gen *media.Generator) *Service {
	return &Service{
		DB:     db,
		Cfg:    cfg,
		Redis:  rds,
		Tokens: tokens,
		Media:  gen,
		SMS:    sms.New(cfg.SMSEnabled, cfg.SMSURL, cfg.SMSLogin, cfg.SMSSender, cfg.SMSSecret),
	}
}

// Actor identifies the user performing an action, for the audit trail.
type Actor struct {
	ID       *uint
	Username string
	IP       string
}

// record writes an audit log entry. Failures never break the main operation.
func (s *Service) record(actor *Actor, action, entity string, entityID *uint, details string) {
	entry := &models.AuditLog{Action: action}
	if entity != "" {
		entry.Entity = &entity
	}
	entry.EntityID = entityID
	if details != "" {
		entry.Details = &details
	}
	if actor != nil {
		entry.UserID = actor.ID
		if actor.Username != "" {
			u := actor.Username
			entry.Username = &u
		}
		if actor.IP != "" {
			ip := actor.IP
			entry.IPAddress = &ip
		}
	}
	if err := s.DB.Create(entry).Error; err != nil {
		log.Printf("audit: failed to record %q: %v", action, err)
	}
}

// GetSettings returns the singleton company settings row, creating it if absent.
func (s *Service) GetSettings() (*models.CompanySettings, error) {
	var settings models.CompanySettings
	err := s.DB.Order("id asc").First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		settings = models.CompanySettings{CompanyName: "Rentacar CRM", Currency: "TJS"}
		if err := s.DB.Create(&settings).Error; err != nil {
			return nil, err
		}
		return &settings, nil
	}
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// ---------------------------------------------------------------- date helpers

func uintPtr(v uint) *uint { return &v }

// parseDate parses a "YYYY-MM-DD" (optionally longer ISO) string into a date.
func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	if len(value) > 10 {
		value = value[:10]
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// parseDateTime parses an ISO datetime, falling back to date-only.
func parseDateTime(value string) (time.Time, error) {
	if value == "" {
		return time.Now(), nil
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04:05", value); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04", value); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", value[:min(10, len(value))]); err == nil {
		return t, nil
	}
	return time.Now(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// formatMoney renders a currency amount with two decimals (e.g. "1500.00").
func formatMoney(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

// atoiDefault parses an integer string, returning a fallback on failure.
func atoiDefault(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return fallback
}
