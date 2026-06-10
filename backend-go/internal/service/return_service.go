package service

import (
	"time"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/pricing"
)

// ReturnCalc holds the financial outcome of a vehicle return.
type ReturnCalc struct {
	ExtraDays    int     `json:"extra_days"`
	LateFee      float64 `json:"late_fee"`
	DamageCost   float64 `json:"damage_cost"`
	Penalties    float64 `json:"penalties"`
	FinalPayment float64 `json:"final_payment"`
}

// ToMap renders the calculation as a JSON map.
func (c ReturnCalc) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"extra_days":    c.ExtraDays,
		"late_fee":      c.LateFee,
		"damage_cost":   c.DamageCost,
		"penalties":     c.Penalties,
		"final_payment": c.FinalPayment,
	}
}

// CalculateReturn computes extra days, late fee and the final payment owed.
// A negative final payment means the company must refund the client.
func (s *Service) CalculateReturn(rental *models.Rental, returnDate time.Time, damageCost, penalties float64) ReturnCalc {
	agreedEnd := rental.RentalEnd
	actualEnd := returnDate

	extraDays := int(actualEnd.Sub(agreedEnd).Hours() / 24)
	if extraDays < 0 {
		extraDays = 0
	}
	lateFee := pricing.Round2(float64(extraDays) * rental.DailyPrice * s.Cfg.LateFeeMultiplier)
	finalPayment := pricing.Round2(lateFee + damageCost + penalties - rental.Deposit)

	return ReturnCalc{
		ExtraDays:    extraDays,
		LateFee:      lateFee,
		DamageCost:   pricing.Round2(damageCost),
		Penalties:    pricing.Round2(penalties),
		FinalPayment: finalPayment,
	}
}

// ReturnInput holds fields for creating a vehicle return.
type ReturnInput struct {
	ReturnDate string   `json:"return_date"`
	Mileage    *int     `json:"mileage"`
	FuelLevel  *int     `json:"fuel_level"`
	Damages    *string  `json:"damages"`
	DamageCost *float64 `json:"damage_cost"`
	Penalties  *float64 `json:"penalties"`
}

// CreateReturn records a vehicle return and completes the rental.
func (s *Service) CreateReturn(rentalID uint, in ReturnInput, actor *Actor) (*models.VehicleReturn, error) {
	rental, err := s.GetRental(rentalID)
	if err != nil {
		return nil, err
	}
	if rental.VehicleReturn != nil {
		return nil, apperr.Conflict("Возврат для этой аренды уже оформлен")
	}
	if rental.Status == models.RentalCancelled {
		return nil, apperr.Validation("Невозможно оформить возврат для отменённой аренды")
	}

	returnDate, _ := parseDateTime(in.ReturnDate)
	damageCost := derefFloat(in.DamageCost)
	penalties := derefFloat(in.Penalties)

	calc := s.CalculateReturn(rental, returnDate, damageCost, penalties)

	mileage := 0
	if rental.Car != nil {
		mileage = rental.Car.Mileage
	}
	if in.Mileage != nil {
		mileage = *in.Mileage
	}
	fuel := 100
	if in.FuelLevel != nil {
		fuel = *in.FuelLevel
	}

	vr := &models.VehicleReturn{
		RentalID:     rental.ID,
		ReturnDate:   returnDate,
		Mileage:      mileage,
		FuelLevel:    fuel,
		Damages:      in.Damages,
		ExtraDays:    calc.ExtraDays,
		LateFee:      calc.LateFee,
		DamageCost:   calc.DamageCost,
		Penalties:    calc.Penalties,
		FinalPayment: calc.FinalPayment,
	}
	if err := s.DB.Create(vr).Error; err != nil {
		return nil, err
	}

	s.DB.Model(&models.Rental{}).Where("id = ?", rental.ID).Update("status", models.RentalCompleted)

	// Update odometer and free the car.
	if rental.Car != nil {
		s.DB.Model(&models.Car{}).Where("id = ?", rental.CarID).
			Updates(map[string]interface{}{"mileage": mileage, "status": models.CarAvailable})
	}

	s.record(actor, "create_return", "rental", uintPtr(rental.ID),
		"Возврат: итог "+formatMoney(calc.FinalPayment))
	return vr, nil
}
