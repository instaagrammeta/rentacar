package service

import (
	"errors"
	"fmt"
	"log"
	"time"

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

	rental := &models.Rental{
		ContractNumber: s.generateContractNumber(),
		ClientID:       client.ID,
		CarID:          car.ID,
		ReservationID:  in.ReservationID,
		EmployeeID:     employeeID,
		RentalStart:    *start,
		RentalEnd:      *end,
		Deposit:        deposit,
		DailyPrice:     car.DailyPrice,
		TotalPrice:     total,
		StartMileage:   &startMileage,
		Notes:          in.Notes,
		Status:         models.RentalActive,
	}
	if err := s.DB.Create(rental).Error; err != nil {
		return nil, err
	}

	s.setCarStatus(car, models.CarRented)
	if rental.ReservationID != nil {
		s.DB.Model(&models.Reservation{}).Where("id = ?", *rental.ReservationID).
			Update("status", models.ReservationConfirmed)
	}

	// Generate the contract PDF (best effort: never fail the contract creation).
	s.generateContractPDF(rental)

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
