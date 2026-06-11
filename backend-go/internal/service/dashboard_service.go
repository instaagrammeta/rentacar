package service

import (
	"fmt"
	"time"

	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/pricing"
)

// revenueBetween sums payment amounts in the half-open interval [start, end).
func (s *Service) revenueBetween(start, end time.Time) float64 {
	var total float64
	s.DB.Model(&models.Payment{}).
		Where("paid_at >= ? AND paid_at < ?", start, end).
		Select("COALESCE(SUM(amount), 0)").Scan(&total)
	return total
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// DashboardSummary returns the headline dashboard metrics.
func (s *Service) DashboardSummary() map[string]interface{} {
	counts := map[string]int64{}
	for _, st := range models.AllCarStatuses() {
		counts[string(st)] = 0
	}
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	s.DB.Model(&models.Car{}).Select("status, count(id) as count").Group("status").Scan(&rows)
	for _, r := range rows {
		counts[r.Status] = r.Count
	}

	now := time.Now()
	dayStart := startOfDay(now)
	monthStart := startOfMonth(now)
	nextMonth := monthStart.AddDate(0, 1, 0)

	todayRevenue := s.revenueBetween(dayStart, dayStart.AddDate(0, 0, 1))
	monthlyRevenue := s.revenueBetween(monthStart, nextMonth)

	var activeRentals int64
	s.DB.Model(&models.Rental{}).Where("status = ?", models.RentalActive).Count(&activeRentals)

	var upcomingReturns int64
	limit := startOfDay(now).AddDate(0, 0, 4) // rental_end <= today+3 (dates)
	s.DB.Model(&models.Rental{}).
		Where("status = ? AND rental_end < ?", models.RentalActive, limit).
		Count(&upcomingReturns)

	return map[string]interface{}{
		"cars_available":   counts[string(models.CarAvailable)],
		"cars_rented":      counts[string(models.CarRented)],
		"cars_reserved":    counts[string(models.CarReserved)],
		"cars_maintenance": counts[string(models.CarMaintenance)],
		"today_revenue":    pricing.Round2(todayRevenue),
		"monthly_revenue":  pricing.Round2(monthlyRevenue),
		"active_rentals":   activeRentals,
		"upcoming_returns": upcomingReturns,
	}
}

// RevenueByMonth returns revenue totals for the last `months` calendar months.
func (s *Service) RevenueByMonth(months int) []map[string]interface{} {
	if months <= 0 {
		months = 12
	}
	now := startOfMonth(time.Now())
	results := make([]map[string]interface{}, 0, months)
	for i := months - 1; i >= 0; i-- {
		monthStart := now.AddDate(0, -i, 0)
		nextMonth := monthStart.AddDate(0, 1, 0)
		amount := s.revenueBetween(monthStart, nextMonth)
		results = append(results, map[string]interface{}{
			"month":   monthStart.Format("2006-01"),
			"revenue": pricing.Round2(amount),
		})
	}
	return results
}

// TopRentedCars returns the most frequently rented cars.
func (s *Service) TopRentedCars(limit int) []map[string]interface{} {
	if limit <= 0 {
		limit = 5
	}
	type row struct {
		ID      uint
		Brand   string
		Model   string
		Plate   string
		Rentals int64
	}
	var rows []row
	s.DB.Model(&models.Car{}).
		Select("cars.id, cars.brand, cars.model, cars.plate_number as plate, count(rentals.id) as rentals").
		Joins("JOIN rentals ON rentals.car_id = cars.id").
		Group("cars.id").
		Order("rentals DESC").
		Limit(limit).
		Scan(&rows)

	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]interface{}{
			"car":     fmt.Sprintf("%s %s (%s)", r.Brand, r.Model, r.Plate),
			"rentals": r.Rentals,
		})
	}
	return out
}

// RentalStatistics returns rental counts grouped by status.
func (s *Service) RentalStatistics() map[string]interface{} {
	stats := map[string]interface{}{}
	for _, st := range models.AllRentalStatuses() {
		stats[string(st)] = int64(0)
	}
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	s.DB.Model(&models.Rental{}).Select("status, count(id) as count").Group("status").Scan(&rows)
	for _, r := range rows {
		stats[r.Status] = r.Count
	}
	return stats
}
