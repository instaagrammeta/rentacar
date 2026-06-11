package service

import "github.com/instaagrammeta/rentacar/backend-go/internal/models"

// SettingsInput holds updatable company settings fields.
type SettingsInput struct {
	CompanyName   *string `json:"company_name"`
	Address       *string `json:"address"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	LogoPath      *string `json:"logo_path"`
	Currency      *string `json:"currency"`
	ContractTerms *string `json:"contract_terms"`
}

// UpdateSettings updates the singleton company settings row.
func (s *Service) UpdateSettings(in SettingsInput, actor *Actor) (*models.CompanySettings, error) {
	settings, err := s.GetSettings()
	if err != nil {
		return nil, err
	}
	if in.CompanyName != nil {
		settings.CompanyName = *in.CompanyName
	}
	if in.Address != nil {
		settings.Address = in.Address
	}
	if in.Phone != nil {
		settings.Phone = in.Phone
	}
	if in.Email != nil {
		settings.Email = in.Email
	}
	if in.LogoPath != nil {
		settings.LogoPath = in.LogoPath
	}
	if in.Currency != nil {
		settings.Currency = *in.Currency
	}
	if in.ContractTerms != nil {
		settings.ContractTerms = in.ContractTerms
	}
	if err := s.DB.Save(settings).Error; err != nil {
		return nil, err
	}
	s.record(actor, "update_settings", "settings", uintPtr(settings.ID), "")
	return settings, nil
}
