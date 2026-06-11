package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// nextSequence returns max(id)+1 for the given model, mirroring the Flask code
// that predicts the next id to build human-readable codes.
func (s *Service) nextSequence(model interface{}) int {
	var maxID int64
	s.DB.Model(model).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
	return int(maxID) + 1
}

func (s *Service) generateClientCode() string {
	return fmt.Sprintf("CL-%06d", s.nextSequence(&models.Client{}))
}

func calculateAge(born *time.Time, on time.Time) *int {
	if born == nil {
		return nil
	}
	age := on.Year() - born.Year()
	if on.Month() < born.Month() || (on.Month() == born.Month() && on.Day() < born.Day()) {
		age--
	}
	return &age
}

// GetClient returns a client by id or a 404 error.
func (s *Service) GetClient(id uint) (*models.Client, error) {
	var client models.Client
	if err := s.DB.First(&client, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound(fmt.Sprintf("Client #%d не найден", id))
		}
		return nil, err
	}
	return &client, nil
}

// AssertRentalEligibility enforces age, experience and blacklist rules.
func (s *Service) AssertRentalEligibility(client *models.Client) error {
	if client.Status == models.ClientBlacklisted {
		return apperr.BusinessRule("Клиент находится в чёрном списке и не может арендовать автомобиль")
	}
	var blCount int64
	s.DB.Model(&models.BlacklistEntry{}).Where("client_id = ?", client.ID).Count(&blCount)
	if blCount > 0 {
		return apperr.BusinessRule("Клиент находится в чёрном списке и не может арендовать автомобиль")
	}

	age := calculateAge(client.DateOfBirth, time.Now())
	if age == nil || *age < s.Cfg.MinClientAge {
		current := "не указан"
		if age != nil {
			current = fmt.Sprintf("%d", *age)
		}
		return apperr.BusinessRule(fmt.Sprintf("Минимальный возраст клиента — %d лет (текущий: %s)", s.Cfg.MinClientAge, current))
	}
	if client.DriverExperienceYears < s.Cfg.MinDrivingExperienceYrs {
		return apperr.BusinessRule(fmt.Sprintf("Минимальный стаж вождения — %d год(а) (текущий: %d)", s.Cfg.MinDrivingExperienceYrs, client.DriverExperienceYears))
	}
	return nil
}

// ClientListResult is the paginated client list payload.
type ClientListResult struct {
	Items   []map[string]interface{} `json:"items"`
	Total   int64                    `json:"total"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
}

// ListClients returns clients filtered by search term and status.
func (s *Service) ListClients(search, status string, page, perPage, offset int) (*ClientListResult, error) {
	apply := func(db *gorm.DB) *gorm.DB {
		db = db.Model(&models.Client{})
		if search != "" {
			like := "%" + search + "%"
			db = db.Where(
				s.DB.Where("first_name ILIKE ?", like).
					Or("last_name ILIKE ?", like).
					Or("phone ILIKE ?", like).
					Or("client_code ILIKE ?", like).
					Or("email ILIKE ?", like),
			)
		} else if status != "" {
			db = db.Where("status = ?", status)
		}
		return db
	}

	var total int64
	apply(s.DB).Count(&total)

	var items []models.Client
	if err := apply(s.DB).Order("last_name asc").Limit(perPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}

	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return &ClientListResult{Items: out, Total: total, Page: page, PerPage: perPage}, nil
}

// ClientInput holds create/update fields for a client.
type ClientInput struct {
	FirstName              *string `json:"first_name"`
	LastName               *string `json:"last_name"`
	Phone                  *string `json:"phone"`
	Email                  *string `json:"email"`
	DateOfBirth            *string `json:"date_of_birth"`
	PassportNumber         *string `json:"passport_number"`
	DriverLicenseNumber    *string `json:"driver_license_number"`
	DriverLicenseIssueDate *string `json:"driver_license_issue_date"`
	DriverExperienceYears  *int    `json:"driver_experience_years"`
	PassportScan           *string `json:"passport_scan"`
	PassportFront          *string `json:"passport_front"`
	PassportBack           *string `json:"passport_back"`
	DriverLicenseScan      *string `json:"driver_license_scan"`
	DriverPhoto            *string `json:"driver_photo"`
	IsVIP                  *bool   `json:"is_vip"`
	Notes                  *string `json:"notes"`
	Status                 *string `json:"status"`
}

// CreateClient creates a client with a unique code and QR identifier.
func (s *Service) CreateClient(in ClientInput, actor *Actor) (*models.Client, error) {
	if in.FirstName == nil || *in.FirstName == "" || in.LastName == nil || *in.LastName == "" {
		return nil, apperr.Validation("Имя и фамилия обязательны")
	}
	if in.Phone == nil || *in.Phone == "" {
		return nil, apperr.Validation("Телефон обязателен")
	}

	var existing int64
	s.DB.Model(&models.Client{}).Where("phone = ?", *in.Phone).Count(&existing)
	if existing > 0 {
		return nil, apperr.Conflict("Клиент с таким телефоном уже существует")
	}

	dob, _ := parseDate(deref(in.DateOfBirth))
	licIssue, _ := parseDate(deref(in.DriverLicenseIssueDate))

	status := models.ClientPendingVerification
	if in.Status != nil && *in.Status != "" {
		status = models.ClientStatus(*in.Status)
	}

	client := &models.Client{
		ClientCode:             s.generateClientCode(),
		FirstName:              strings.TrimSpace(*in.FirstName),
		LastName:               strings.TrimSpace(*in.LastName),
		Phone:                  strings.TrimSpace(*in.Phone),
		Email:                  optStr(in.Email),
		DateOfBirth:            dob,
		PassportNumber:         optStr(in.PassportNumber),
		DriverLicenseNumber:    optStr(in.DriverLicenseNumber),
		DriverLicenseIssueDate: licIssue,
		DriverExperienceYears:  derefInt(in.DriverExperienceYears),
		PassportScan:           optStr(in.PassportScan),
		PassportFront:          optStr(in.PassportFront),
		PassportBack:           optStr(in.PassportBack),
		DriverLicenseScan:      optStr(in.DriverLicenseScan),
		DriverPhoto:            optStr(in.DriverPhoto),
		Notes:                  optStr(in.Notes),
		Status:                 status,
	}
	if in.IsVIP != nil {
		client.IsVIP = *in.IsVIP
	}

	if err := s.DB.Create(client).Error; err != nil {
		return nil, err
	}

	// Generate a QR code now that the code is known.
	if path, err := s.Media.GenerateQR(client.ClientCode, client.ClientCode+".png"); err == nil {
		client.QRCodePath = &path
		s.DB.Save(client)
	} else {
		log.Printf("client: QR generation failed: %v", err)
	}

	s.record(actor, "create_client", "client", uintPtr(client.ID), "Создан клиент "+client.FullName())
	return client, nil
}

// UpdateClient updates an existing client.
func (s *Service) UpdateClient(id uint, in ClientInput, actor *Actor) (*models.Client, error) {
	client, err := s.GetClient(id)
	if err != nil {
		return nil, err
	}

	if in.FirstName != nil {
		client.FirstName = *in.FirstName
	}
	if in.LastName != nil {
		client.LastName = *in.LastName
	}
	if in.Phone != nil {
		client.Phone = *in.Phone
	}
	if in.Email != nil {
		client.Email = optStr(in.Email)
	}
	if in.PassportNumber != nil {
		client.PassportNumber = optStr(in.PassportNumber)
	}
	if in.DriverLicenseNumber != nil {
		client.DriverLicenseNumber = optStr(in.DriverLicenseNumber)
	}
	if in.PassportScan != nil {
		client.PassportScan = optStr(in.PassportScan)
	}
	if in.PassportFront != nil {
		client.PassportFront = optStr(in.PassportFront)
	}
	if in.PassportBack != nil {
		client.PassportBack = optStr(in.PassportBack)
	}
	if in.DriverPhoto != nil {
		client.DriverPhoto = optStr(in.DriverPhoto)
	}
	if in.DriverLicenseScan != nil {
		client.DriverLicenseScan = optStr(in.DriverLicenseScan)
	}
	if in.Notes != nil {
		client.Notes = optStr(in.Notes)
	}
	if in.IsVIP != nil {
		client.IsVIP = *in.IsVIP
	}
	if in.DriverExperienceYears != nil {
		client.DriverExperienceYears = *in.DriverExperienceYears
	}
	if in.DateOfBirth != nil {
		client.DateOfBirth, _ = parseDate(*in.DateOfBirth)
	}
	if in.DriverLicenseIssueDate != nil {
		client.DriverLicenseIssueDate, _ = parseDate(*in.DriverLicenseIssueDate)
	}
	if in.Status != nil && *in.Status != "" {
		client.Status = models.ClientStatus(*in.Status)
	}

	if err := s.DB.Save(client).Error; err != nil {
		return nil, err
	}
	s.record(actor, "update_client", "client", uintPtr(client.ID), "")
	return client, nil
}

// DeleteClient deletes a client without rental history.
func (s *Service) DeleteClient(id uint, actor *Actor) error {
	client, err := s.GetClient(id)
	if err != nil {
		return err
	}
	var rentalCount int64
	s.DB.Model(&models.Rental{}).Where("client_id = ?", id).Count(&rentalCount)
	if rentalCount > 0 {
		return apperr.Conflict("Невозможно удалить клиента с историей аренды")
	}
	if err := s.DB.Delete(client).Error; err != nil {
		return err
	}
	s.record(actor, "delete_client", "client", uintPtr(id), "")
	return nil
}

// FindClient looks up a returning client by code, phone, or free-text term.
// It returns either a single client (*models.Client) or a slice ([]).
func (s *Service) FindClient(phone, code, term string) (interface{}, error) {
	if code != "" {
		var client models.Client
		if err := s.DB.Where("client_code = ?", code).First(&client).Error; err == nil {
			return &client, nil
		}
	}
	if phone != "" {
		var client models.Client
		if err := s.DB.Where("phone = ?", phone).First(&client).Error; err == nil {
			return &client, nil
		}
	}
	if term != "" {
		like := "%" + term + "%"
		var items []models.Client
		s.DB.Where("first_name ILIKE ?", like).
			Or("last_name ILIKE ?", like).
			Or("phone ILIKE ?", like).
			Or("client_code ILIKE ?", like).
			Or("email ILIKE ?", like).
			Find(&items)
		return items, nil
	}
	return nil, apperr.NotFound("Клиент не найден")
}

// GetClientHistory aggregates a client's rentals, payments, penalties, accidents.
func (s *Service) GetClientHistory(id uint) (map[string]interface{}, error) {
	client, err := s.GetClient(id)
	if err != nil {
		return nil, err
	}

	var rentals []models.Rental
	s.DB.Preload("Client").Preload("Car").Preload("Employee").Preload("VehicleReturn").
		Where("client_id = ?", id).Order("id desc").Find(&rentals)

	var payments []models.Payment
	s.DB.Preload("Client").Where("client_id = ?", id).Order("id desc").Find(&payments)

	var accidents []models.Accident
	s.DB.Preload("Car").Preload("Client").Preload("Photos").Where("client_id = ?", id).Order("id desc").Find(&accidents)

	rentalMaps := make([]map[string]interface{}, 0, len(rentals))
	penalties := make([]map[string]interface{}, 0)
	for i := range rentals {
		rentalMaps = append(rentalMaps, rentals[i].ToMap())
		vr := rentals[i].VehicleReturn
		if vr != nil && (vr.Penalties != 0 || vr.LateFee != 0) {
			penalties = append(penalties, map[string]interface{}{
				"rental_id":       rentals[i].ID,
				"contract_number": rentals[i].ContractNumber,
				"amount":          vr.Penalties + vr.LateFee,
			})
		}
	}
	paymentMaps := make([]map[string]interface{}, 0, len(payments))
	for i := range payments {
		paymentMaps = append(paymentMaps, payments[i].ToMap())
	}
	accidentMaps := make([]map[string]interface{}, 0, len(accidents))
	for i := range accidents {
		accidentMaps = append(accidentMaps, accidents[i].ToMap())
	}

	return map[string]interface{}{
		"client":    client.ToMap(),
		"rentals":   rentalMaps,
		"payments":  paymentMaps,
		"penalties": penalties,
		"accidents": accidentMaps,
	}, nil
}

// ----------------------------------------------------------------- small utils

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// optStr returns nil for empty strings, otherwise a pointer to the value.
func optStr(p *string) *string {
	if p == nil {
		return nil
	}
	v := strings.TrimSpace(*p)
	if v == "" {
		return nil
	}
	return &v
}
