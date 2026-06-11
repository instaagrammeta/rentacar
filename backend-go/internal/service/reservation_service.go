package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// GetReservation returns a reservation (with client + car) or a 404 error.
func (s *Service) GetReservation(id uint) (*models.Reservation, error) {
	var r models.Reservation
	if err := s.DB.Preload("Client").Preload("Car").First(&r, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound(fmt.Sprintf("Reservation #%d не найден", id))
		}
		return nil, err
	}
	return &r, nil
}

// ListReservations returns reservations filtered by status.
func (s *Service) ListReservations(status string, page, perPage, offset int) (*ListResult, error) {
	apply := func(db *gorm.DB) *gorm.DB {
		db = db.Model(&models.Reservation{})
		if status != "" {
			db = db.Where("status = ?", status)
		}
		return db
	}
	var total int64
	apply(s.DB).Count(&total)

	var items []models.Reservation
	if err := apply(s.DB).Preload("Client").Preload("Car").Order("id desc").Limit(perPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return &ListResult{Items: out, Total: total, Page: page, PerPage: perPage}, nil
}

// ReservationInput holds fields to create a reservation.
type ReservationInput struct {
	ClientID  uint     `json:"client_id"`
	CarID     uint     `json:"car_id"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
	Deposit   *float64 `json:"deposit"`
	Notes     *string  `json:"notes"`
}

// CreateReservation books a car for a client for a date range.
func (s *Service) CreateReservation(in ReservationInput, actor *Actor) (*models.Reservation, error) {
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

	start, err := parseDate(in.StartDate)
	if err != nil || start == nil {
		return nil, apperr.Validation("Дата обязательна")
	}
	end, err := parseDate(in.EndDate)
	if err != nil || end == nil {
		return nil, apperr.Validation("Дата обязательна")
	}
	if end.Before(*start) {
		return nil, apperr.Validation("Дата окончания не может быть раньше даты начала")
	}
	if car.Status == models.CarRented || car.Status == models.CarMaintenance {
		return nil, apperr.Validation("Автомобиль недоступен для бронирования")
	}

	deposit := car.DepositAmount
	if in.Deposit != nil {
		deposit = *in.Deposit
	}

	reservation := &models.Reservation{
		ClientID:  client.ID,
		CarID:     car.ID,
		StartDate: *start,
		EndDate:   *end,
		Deposit:   deposit,
		Notes:     in.Notes,
		Status:    models.ReservationReserved,
	}
	if err := s.DB.Create(reservation).Error; err != nil {
		return nil, err
	}
	s.setCarStatus(car, models.CarReserved)

	reservation.Client = client
	reservation.Car = car
	s.record(actor, "create_reservation", "reservation", uintPtr(reservation.ID), "")
	return reservation, nil
}

// UpdateReservationStatus changes a reservation's status and syncs the car.
func (s *Service) UpdateReservationStatus(id uint, status string, actor *Actor) (*models.Reservation, error) {
	reservation, err := s.GetReservation(id)
	if err != nil {
		return nil, err
	}
	newStatus := models.ReservationStatus(status)
	if !newStatus.Valid() {
		return nil, apperr.Validation("Некорректный статус брони")
	}
	reservation.Status = newStatus
	if err := s.DB.Model(reservation).Update("status", newStatus).Error; err != nil {
		return nil, err
	}

	var car models.Car
	if err := s.DB.First(&car, reservation.CarID).Error; err == nil {
		if newStatus == models.ReservationCancelled && car.Status == models.CarReserved {
			s.setCarStatus(&car, models.CarAvailable)
		} else if newStatus == models.ReservationReserved || newStatus == models.ReservationConfirmed {
			s.setCarStatus(&car, models.CarReserved)
		}
	}

	s.record(actor, "reservation_status", "reservation", uintPtr(reservation.ID), "Статус: "+newStatus.Label())
	return reservation, nil
}

// CancelReservation cancels a reservation.
func (s *Service) CancelReservation(id uint, actor *Actor) (*models.Reservation, error) {
	return s.UpdateReservationStatus(id, string(models.ReservationCancelled), actor)
}
