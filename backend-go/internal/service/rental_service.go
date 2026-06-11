package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/pricing"
)

// GetRental returns a rental with relations and its return record.
func (s *Service) GetRental(id uint) (*models.Rental, error) {
	var r models.Rental
	err := s.DB.Preload("Client").Preload("Car").Preload("Employee").Preload("VehicleReturn").First(&r, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound(fmt.Sprintf("Rental #%d не найден", id))
		}
		return nil, err
	}
	return &r, nil
}

// ListRentals returns rentals filtered by status.
func (s *Service) ListRentals(status string, page, perPage, offset int) (*ListResult, error) {
	apply := func(db *gorm.DB) *gorm.DB {
		db = db.Model(&models.Rental{})
		if status != "" {
			db = db.Where("status = ?", status)
		}
		return db
	}
	var total int64
	apply(s.DB).Count(&total)

	var items []models.Rental
	if err := apply(s.DB).Preload("Client").Preload("Car").Preload("Employee").Preload("VehicleReturn").
		Order("id desc").Limit(perPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return &ListResult{Items: out, Total: total, Page: page, PerPage: perPage}, nil
}

func (s *Service) generateContractNumber() string {
	return fmt.Sprintf("DOG-%d-%05d", time.Now().Year(), s.nextSequence(&models.Rental{}))
}

// RentalInput holds fields to create a rental contract.
type RentalInput struct {
	ClientID      uint     `json:"client_id"`
	CarID         uint     `json:"car_id"`
	ReservationID *uint    `json:"reservation_id"`
	EmployeeID    *uint    `json:"employee_id"`
	RentalStart   string   `json:"rental_start"`
	RentalEnd     string   `json:"rental_end"`
	PickupAt      *string  `json:"pickup_at"`
	DueAt         *string  `json:"due_at"`
	Deposit       *float64 `json:"deposit"`
	TotalPrice    *float64 `json:"total_price"`
	StartMileage  *int     `json:"start_mileage"`
	Notes         *string  `json:"notes"`
}

// CreateRental creates a contract, calculates price and generates the PDF.
func (s *Service) CreateRental(in RentalInput, actor *Actor) (*models.Rental, error) {
	client, err := s.GetClient(in.ClientID)
	if err != nil {
		return nil, err
	}
	car, err := s.GetCar(in.CarID)
	if err != nil {
		return nil, err
	}
	if err := s.AssertRentalEligibility(client); err != nil {
		return nil, err
	}
	if car.Status == models.CarRented {
		return nil, apperr.Validation("Автомобиль уже находится в аренде")
	}
	if car.Status == models.CarMaintenance {
		return nil, apperr.Validation("Автомобиль находится на обслуживании")
	}

	start, err := parseDate(in.RentalStart)
	if err != nil || start == nil {
		return nil, apperr.Validation("Дата обязательна")
	}
	end, err := parseDate(in.RentalEnd)
	if err != nil || end == nil {
		return nil, apperr.Validation("Дата обязательна")
	}
	if end.Before(*start) {
		return nil, apperr.Validation("Дата окончания не может быть раньше даты начала")
	}

	breakdown := pricing.Calculate(*start, *end, car.DailyPrice, car.WeeklyPrice, car.MonthlyPrice)
	total := breakdown.Total
	if in.TotalPrice != nil && *in.TotalPrice > 0 {
		total = *in.TotalPrice
	}

	deposit := car.DepositAmount
	if in.Deposit != nil {
		deposit = *in.Deposit
	}
	startMileage := car.Mileage
	if in.StartMileage != nil {
		startMileage = *in.StartMileage
	}

	employeeID := actor.idOrNil()
	if employeeID == nil {
		employeeID = in.EmployeeID
	}

	// Pickup time (when the car was handed over) defaults to now.
	pickup := time.Now()
	if in.PickupAt != nil && *in.PickupAt != "" {
		if t, err := parseDateTime(*in.PickupAt); err == nil {
			pickup = t
		}
	}
	// Due time (precise return deadline) defaults to the end date at 18:00.
	due := time.Date(end.Year(), end.Month(), end.Day(), 18, 0, 0, 0, time.Local)
	if in.DueAt != nil && *in.DueAt != "" {
		if t, err := parseDateTime(*in.DueAt); err == nil {
			due = t
		}
	}

	rental := &models.Rental{
		ContractNumber: s.generateContractNumber(),
		ClientID:       client.ID,
		CarID:          car.ID,
		ReservationID:  in.ReservationID,
		EmployeeID:     employeeID,
		RentalStart:    *start,
		RentalEnd:      *end,
		PickupAt:       &pickup,
		DueAt:          &due,
		Deposit:        deposit,
		DailyPrice:     car.DailyPrice,
		TotalPrice:     total,
		StartMileage:   &startMileage,
		Notes:          in.Notes,
		Status:         models.RentalActive,
	}
	token := uuid.NewString()
	rental.PublicToken = &token
	if err := s.DB.Create(rental).Error; err != nil {
		return nil, err
	}

	s.setCarStatus(car, models.CarRented)
	if rental.ReservationID != nil {
		s.DB.Model(&models.Reservation{}).Where("id = ?", *rental.ReservationID).
			Update("status", models.ReservationConfirmed)
	}

	// Generate the contract PDF and the public QR code (best effort: never
	// fail the contract creation if these auxiliary artefacts fail).
	s.generateContractPDF(rental)
	s.generateRentalQR(rental)

	s.record(actor, "create_rental", "rental", uintPtr(rental.ID), "Договор "+rental.ContractNumber)
	return s.GetRental(rental.ID)
}

// generateContractPDF fills the relations and writes the PDF, ignoring errors.
func (s *Service) generateContractPDF(rental *models.Rental) {
	full, err := s.GetRental(rental.ID)
	if err != nil {
		return
	}
	settings, _ := s.GetSettings()
	path, err := s.Media.GenerateContractPDF(full, settings)
	if err != nil {
		log.Printf("rental: contract PDF generation failed: %v", err)
		return
	}
	s.DB.Model(&models.Rental{}).Where("id = ?", rental.ID).Update("pdf_path", path)
	rental.PDFPath = &path
}

// publicURLForToken builds the absolute public page URL a client opens after
// scanning the QR code, e.g. "https://example.com/r/<token>".
func (s *Service) publicURLForToken(token string) string {
	return s.Cfg.PublicURL + "/r/" + token
}

// generateRentalQR creates the QR-code PNG that points to the public rental
// status page. Best effort: failures are logged and never break the contract.
func (s *Service) generateRentalQR(rental *models.Rental) {
	if rental.PublicToken == nil || *rental.PublicToken == "" {
		token := uuid.NewString()
		rental.PublicToken = &token
		s.DB.Model(&models.Rental{}).Where("id = ?", rental.ID).Update("public_token", token)
	}
	url := s.publicURLForToken(*rental.PublicToken)
	filename := fmt.Sprintf("rental-%d.png", rental.ID)
	path, err := s.Media.GenerateQR(url, filename)
	if err != nil {
		log.Printf("rental: QR generation failed: %v", err)
		return
	}
	s.DB.Model(&models.Rental{}).Where("id = ?", rental.ID).Update("qr_code_path", path)
	rental.QRCodePath = &path
}

// EnsureRentalQR makes sure a rental has a public token and QR code, generating
// them on demand (used for rentals created before this feature existed).
func (s *Service) EnsureRentalQR(id uint) (*models.Rental, error) {
	rental, err := s.GetRental(id)
	if err != nil {
		return nil, err
	}
	if rental.QRCodePath == nil || *rental.QRCodePath == "" ||
		rental.PublicToken == nil || *rental.PublicToken == "" {
		s.generateRentalQR(rental)
	}
	return s.GetRental(id)
}

// PublicRentalURL returns the absolute public link for a rental's QR target.
func (s *Service) PublicRentalURL(rental *models.Rental) string {
	if rental.PublicToken == nil {
		return ""
	}
	return s.publicURLForToken(*rental.PublicToken)
}

// GetRentalByToken looks up a rental by its public token (no auth required).
func (s *Service) GetRentalByToken(token string) (*models.Rental, error) {
	if token == "" {
		return nil, apperr.NotFound("Аренда не найдена")
	}
	var r models.Rental
	err := s.DB.Preload("Car").Preload("Client").Where("public_token = ?", token).First(&r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("Аренда не найдена")
		}
		return nil, err
	}
	return &r, nil
}

// PublicRentalView builds the minimal, privacy-safe payload shown on the public
// status page, including how much time is left until the return deadline.
func (s *Service) PublicRentalView(r *models.Rental) map[string]interface{} {
	now := time.Now()

	var due *time.Time
	if r.DueAt != nil {
		due = r.DueAt
	}

	remainingSeconds := int64(0)
	isOverdue := false
	if due != nil {
		diff := due.Sub(now)
		remainingSeconds = int64(diff.Seconds())
		if remainingSeconds < 0 {
			isOverdue = true
		}
	}
	// Once the rental is completed or cancelled, the countdown is irrelevant.
	if r.Status != models.RentalActive {
		remainingSeconds = 0
		isOverdue = false
	}

	carName := ""
	if r.Car != nil {
		carName = r.Car.DisplayName()
	}
	clientName := ""
	if r.Client != nil {
		clientName = r.Client.FirstName
	}

	out := map[string]interface{}{
		"contract_number":   r.ContractNumber,
		"car_name":          carName,
		"client_name":       clientName,
		"status":            string(r.Status),
		"status_label":      r.Status.Label(),
		"rental_start":      r.RentalStart.Format("2006-01-02"),
		"rental_end":        r.RentalEnd.Format("2006-01-02"),
		"pickup_at":         nil,
		"due_at":            nil,
		"server_time":       now.UTC().Format("2006-01-02T15:04:05Z07:00"),
		"remaining_seconds": remainingSeconds,
		"is_overdue":        isOverdue,
	}
	if r.PickupAt != nil {
		out["pickup_at"] = r.PickupAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	if due != nil {
		out["due_at"] = due.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return out
}

// CancelRental cancels an active rental and frees the car.
func (s *Service) CancelRental(id uint, actor *Actor) (*models.Rental, error) {
	rental, err := s.GetRental(id)
	if err != nil {
		return nil, err
	}
	if rental.Status == models.RentalCompleted {
		return nil, apperr.Validation("Невозможно отменить завершённую аренду")
	}
	rental.Status = models.RentalCancelled
	s.DB.Model(rental).Update("status", models.RentalCancelled)

	var car models.Car
	if err := s.DB.First(&car, rental.CarID).Error; err == nil {
		s.setCarStatus(&car, models.CarAvailable)
	}
	s.record(actor, "cancel_rental", "rental", uintPtr(rental.ID), "")
	return rental, nil
}

// RegeneratePDF regenerates the contract PDF for a rental.
func (s *Service) RegeneratePDF(id uint) (*models.Rental, error) {
	rental, err := s.GetRental(id)
	if err != nil {
		return nil, err
	}
	s.generateContractPDF(rental)
	return s.GetRental(id)
}

// idOrNil returns the actor id as a pointer (nil when unauthenticated).
func (a *Actor) idOrNil() *uint {
	if a == nil {
		return nil
	}
	return a.ID
}
