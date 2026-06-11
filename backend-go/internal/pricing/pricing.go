// Package pricing implements the rental price calculation. It chooses the most
// favourable tariff for the customer: monthly blocks first, then whole weeks,
// then remaining days, mirroring the original Flask implementation.
package pricing

import (
	"math"
	"time"
)

// Breakdown describes how a total rental price was composed.
type Breakdown struct {
	Days      int     `json:"days"`
	Months    int     `json:"months"`
	Weeks     int     `json:"weeks"`
	ExtraDays int     `json:"extra_days"`
	Total     float64 `json:"total"`
}

// ToMap returns the breakdown with the total rounded to 2 decimals.
func (b Breakdown) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"days":       b.Days,
		"months":     b.Months,
		"weeks":      b.Weeks,
		"extra_days": b.ExtraDays,
		"total":      Round2(b.Total),
	}
}

// RentalDays returns the number of billable days (at least one).
func RentalDays(start, end time.Time) int {
	delta := int(end.Sub(start).Hours() / 24)
	if delta < 1 {
		return 1
	}
	return delta
}

// Calculate computes the total rental price for a date range.
func Calculate(start, end time.Time, dailyPrice, weeklyPrice, monthlyPrice float64) Breakdown {
	days := RentalDays(start, end)
	remaining := days
	total := 0.0

	months := 0
	if monthlyPrice > 0 {
		months = remaining / 30
		total += float64(months) * monthlyPrice
		remaining -= months * 30
	}

	weeks := 0
	if weeklyPrice > 0 {
		weeks = remaining / 7
		total += float64(weeks) * weeklyPrice
		remaining -= weeks * 7
	}

	total += float64(remaining) * dailyPrice
	return Breakdown{Days: days, Months: months, Weeks: weeks, ExtraDays: remaining, Total: total}
}

// Round2 rounds a float to two decimal places (banker-free, matches Python round
// closely enough for currency display).
func Round2(v float64) float64 {
	return math.Round(v*100) / 100
}
