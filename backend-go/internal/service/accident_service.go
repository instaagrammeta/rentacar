package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// GetAccident returns an accident (with relations + photos) or a 404 error.
func (s *Service) GetAccident(id uint) (*models.Accident, error) {
	var a models.Accident
	if err := s.DB.Preload("Car").Preload("Client").Preload("Photos").First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound(fmt.Sprintf("Accident #%d не найден", id))
		}
		return nil, err
	}
	return &a, nil
}

// ListAccidents returns accidents filtered by car.
func (s *Service) ListAccidents(carID uint, page, perPage, offset int) (*ListResult, error) {
	q := s.DB.Model(&models.Accident{})
	if carID > 0 {
		q = q.Where("car_id = ?", carID)
	}
	var total int64
	q.Count(&total)

	var items []models.Accident
	if err := q.Preload("Car").Preload("Client").Preload("Photos").
		Order("id desc").Limit(perPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return &ListResult{Items: out, Total: total, Page: page, PerPage: perPage}, nil
}

// AccidentInput holds create/update fields for an accident.
type AccidentInput struct {
	CarID        *uint     `json:"car_id"`
	ClientID     *uint     `json:"client_id"`
	RentalID     *uint     `json:"rental_id"`
	AccidentDate *string   `json:"accident_date"`
	Description  *string   `json:"description"`
	RepairCost   *float64  `json:"repair_cost"`
	Photos       *[]string `json:"photos"`
}

// CreateAccident records an accident.
func (s *Service) CreateAccident(in AccidentInput, actor *Actor) (*models.Accident, error) {
	if in.CarID == nil {
		return nil, apperr.Validation("Автомобиль обязателен")
	}
	car, err := s.GetCar(*in.CarID)
	if err != nil {
		return nil, err
	}

	accDate := time.Now()
	if in.AccidentDate != nil && *in.AccidentDate != "" {
		if d, err := parseDate(*in.AccidentDate); err == nil && d != nil {
			accDate = *d
		}
	}

	accident := &models.Accident{
		CarID:        car.ID,
		ClientID:     in.ClientID,
		RentalID:     in.RentalID,
		AccidentDate: accDate,
		Description:  in.Description,
		RepairCost:   derefFloat(in.RepairCost),
	}
	if in.Photos != nil {
		for _, p := range *in.Photos {
			accident.Photos = append(accident.Photos, models.AccidentPhoto{FilePath: p})
		}
	}
	if err := s.DB.Create(accident).Error; err != nil {
		return nil, err
	}
	s.record(actor, "create_accident", "accident", uintPtr(accident.ID), "")
	return s.GetAccident(accident.ID)
}

// UpdateAccident updates an accident record.
func (s *Service) UpdateAccident(id uint, in AccidentInput, actor *Actor) (*models.Accident, error) {
	accident, err := s.GetAccident(id)
	if err != nil {
		return nil, err
	}
	if in.Description != nil {
		accident.Description = in.Description
	}
	if in.RepairCost != nil {
		accident.RepairCost = *in.RepairCost
	}
	if in.AccidentDate != nil && *in.AccidentDate != "" {
		if d, err := parseDate(*in.AccidentDate); err == nil && d != nil {
			accident.AccidentDate = *d
		}
	}
	if err := s.DB.Save(accident).Error; err != nil {
		return nil, err
	}
	if in.Photos != nil {
		s.DB.Where("accident_id = ?", accident.ID).Delete(&models.AccidentPhoto{})
		for _, p := range *in.Photos {
			s.DB.Create(&models.AccidentPhoto{AccidentID: accident.ID, FilePath: p})
		}
	}
	return s.GetAccident(accident.ID)
}

// DeleteAccident deletes an accident record.
func (s *Service) DeleteAccident(id uint, actor *Actor) error {
	accident, err := s.GetAccident(id)
	if err != nil {
		return err
	}
	if err := s.DB.Select("Photos").Delete(accident).Error; err != nil {
		return err
	}
	s.record(actor, "delete_accident", "accident", uintPtr(id), "")
	return nil
}
