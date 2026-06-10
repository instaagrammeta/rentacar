package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// ListBlacklist returns all blacklist entries (with client) newest first.
func (s *Service) ListBlacklist() ([]map[string]interface{}, error) {
	var items []models.BlacklistEntry
	if err := s.DB.Preload("Client").Order("id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return out, nil
}

// BlacklistInput holds fields to blacklist a client.
type BlacklistInput struct {
	ClientID uint    `json:"client_id"`
	Reason   string  `json:"reason"`
	Comment  *string `json:"comment"`
}

// AddToBlacklist blacklists a client.
func (s *Service) AddToBlacklist(in BlacklistInput, actor *Actor) (*models.BlacklistEntry, error) {
	client, err := s.GetClient(in.ClientID)
	if err != nil {
		return nil, err
	}

	var existing int64
	s.DB.Model(&models.BlacklistEntry{}).Where("client_id = ?", client.ID).Count(&existing)
	if existing > 0 {
		return nil, apperr.Conflict("Клиент уже в чёрном списке")
	}

	reason := models.BlacklistReason(in.Reason)
	if !reason.Valid() {
		return nil, apperr.Validation("Укажите корректную причину")
	}

	entry := &models.BlacklistEntry{
		ClientID:  client.ID,
		Reason:    reason,
		Comment:   in.Comment,
		CreatedBy: actor.idOrNil(),
	}
	if err := s.DB.Create(entry).Error; err != nil {
		return nil, err
	}

	s.DB.Model(&models.Client{}).Where("id = ?", client.ID).Update("status", models.ClientBlacklisted)
	entry.Client = client

	s.record(actor, "blacklist_add", "client", uintPtr(client.ID), "Причина: "+reason.Label())
	return entry, nil
}

// RemoveFromBlacklist removes a client from the blacklist.
func (s *Service) RemoveFromBlacklist(clientID uint, actor *Actor) error {
	var entry models.BlacklistEntry
	err := s.DB.Where("client_id = ?", clientID).First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.Conflict("Клиент не находится в чёрном списке")
	}
	if err != nil {
		return err
	}
	if err := s.DB.Delete(&entry).Error; err != nil {
		return err
	}
	s.DB.Model(&models.Client{}).Where("id = ?", clientID).Update("status", models.ClientActive)
	s.record(actor, "blacklist_remove", "client", uintPtr(clientID), "")
	return nil
}
