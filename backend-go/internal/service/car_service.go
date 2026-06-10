package service

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// GetCar returns a car (with photos) by id or a 404 error.
func (s *Service) GetCar(id uint) (*models.Car, error) {
	var car models.Car
	if err := s.DB.Preload("Photos").First(&car, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound(fmt.Sprintf("Car #%d не найден", id))
		}
		return nil, err
	}
	return &car, nil
}

// ListResult is a generic paginated list payload.
type ListResult struct {
	Items   []map[string]interface{} `json:"items"`
	Total   int64                    `json:"total"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
}

// ListCars returns cars filtered by status and search term.
func (s *Service) ListCars(status, search string, page, perPage, offset int) (*ListResult, error) {
	q := s.DB.Model(&models.Car{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		like := "%" + search + "%"
		q = q.Where(
			s.DB.Where("brand ILIKE ?", like).Or("model ILIKE ?", like).Or("plate_number ILIKE ?", like),
		)
	}

	var total int64
	q.Count(&total)

	var items []models.Car
	if err := q.Preload("Photos").Order("id desc").Limit(perPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return &ListResult{Items: out, Total: total, Page: page, PerPage: perPage}, nil
}

// AvailableCars returns cars whose status is "available".
func (s *Service) AvailableCars() ([]map[string]interface{}, error) {
	var items []models.Car
	if err := s.DB.Preload("Photos").Where("status = ?", models.CarAvailable).Limit(500).Order("id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return out, nil
}

// CarInput holds create/update fields for a car.
type CarInput struct {
	Brand         *string   `json:"brand"`
	Model         *string   `json:"model"`
	Year          *int      `json:"year"`
	Color         *string   `json:"color"`
	VIN           *string   `json:"vin"`
	PlateNumber   *string   `json:"plate_number"`
	Mileage       *int      `json:"mileage"`
	DailyPrice    *float64  `json:"daily_price"`
	WeeklyPrice   *float64  `json:"weekly_price"`
	MonthlyPrice  *float64  `json:"monthly_price"`
	DepositAmount *float64  `json:"deposit_amount"`
	Status        *string   `json:"status"`
	Photos        *[]string `json:"photos"`
}

// CreateCar creates a vehicle.
func (s *Service) CreateCar(in CarInput, actor *Actor) (*models.Car, error) {
	if in.Brand == nil || *in.Brand == "" || in.Model == nil || *in.Model == "" {
		return nil, apperr.Validation("Марка и модель обязательны")
	}
	if in.PlateNumber == nil || *in.PlateNumber == "" {
		return nil, apperr.Validation("Гос. номер обязателен")
	}

	var existing int64
	s.DB.Model(&models.Car{}).Where("plate_number = ?", *in.PlateNumber).Count(&existing)
	if existing > 0 {
		return nil, apperr.Conflict("Автомобиль с таким гос. номером уже существует")
	}

	status := models.CarAvailable
	if in.Status != nil && *in.Status != "" {
		status = models.CarStatus(*in.Status)
	}

	car := &models.Car{
		Brand:         strings.TrimSpace(*in.Brand),
		Model:         strings.TrimSpace(*in.Model),
		Year:          derefInt(in.Year),
		Color:         optStr(in.Color),
		VIN:           optStr(in.VIN),
		PlateNumber:   strings.TrimSpace(*in.PlateNumber),
		Mileage:       derefInt(in.Mileage),
		DailyPrice:    derefFloat(in.DailyPrice),
		WeeklyPrice:   derefFloat(in.WeeklyPrice),
		MonthlyPrice:  derefFloat(in.MonthlyPrice),
		DepositAmount: derefFloat(in.DepositAmount),
		Status:        status,
	}
	if in.Photos != nil {
		for _, p := range *in.Photos {
			car.Photos = append(car.Photos, models.CarPhoto{FilePath: p})
		}
	}
	if err := s.DB.Create(car).Error; err != nil {
		return nil, err
	}
	s.record(actor, "create_car", "car", uintPtr(car.ID), "Добавлен автомобиль "+car.DisplayName())
	return car, nil
}

// UpdateCar updates a vehicle.
func (s *Service) UpdateCar(id uint, in CarInput, actor *Actor) (*models.Car, error) {
	car, err := s.GetCar(id)
	if err != nil {
		return nil, err
	}

	if in.Brand != nil {
		car.Brand = *in.Brand
	}
	if in.Model != nil {
		car.Model = *in.Model
	}
	if in.Color != nil {
		car.Color = optStr(in.Color)
	}
	if in.VIN != nil {
		car.VIN = optStr(in.VIN)
	}
	if in.PlateNumber != nil {
		car.PlateNumber = *in.PlateNumber
	}
	if in.Year != nil {
		car.Year = *in.Year
	}
	if in.Mileage != nil {
		car.Mileage = *in.Mileage
	}
	if in.DailyPrice != nil {
		car.DailyPrice = *in.DailyPrice
	}
	if in.WeeklyPrice != nil {
		car.WeeklyPrice = *in.WeeklyPrice
	}
	if in.MonthlyPrice != nil {
		car.MonthlyPrice = *in.MonthlyPrice
	}
	if in.DepositAmount != nil {
		car.DepositAmount = *in.DepositAmount
	}
	if in.Status != nil && *in.Status != "" {
		car.Status = models.CarStatus(*in.Status)
	}

	if err := s.DB.Save(car).Error; err != nil {
		return nil, err
	}

	// Replace photos if provided.
	if in.Photos != nil {
		s.DB.Where("car_id = ?", car.ID).Delete(&models.CarPhoto{})
		for _, p := range *in.Photos {
			s.DB.Create(&models.CarPhoto{CarID: car.ID, FilePath: p})
		}
		s.DB.Preload("Photos").First(car, car.ID)
	}

	s.record(actor, "update_car", "car", uintPtr(car.ID), "")
	return car, nil
}

// setCarStatus updates a car's availability status.
func (s *Service) setCarStatus(car *models.Car, status models.CarStatus) {
	car.Status = status
	s.DB.Model(car).Update("status", status)
}

// DeleteCar deletes a car without rental history.
func (s *Service) DeleteCar(id uint, actor *Actor) error {
	car, err := s.GetCar(id)
	if err != nil {
		return err
	}
	var rentalCount int64
	s.DB.Model(&models.Rental{}).Where("car_id = ?", id).Count(&rentalCount)
	if rentalCount > 0 {
		return apperr.Conflict("Невозможно удалить автомобиль с историей аренды")
	}
	if err := s.DB.Select("Photos").Delete(car).Error; err != nil {
		return err
	}
	s.record(actor, "delete_car", "car", uintPtr(id), "")
	return nil
}

func derefFloat(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}
