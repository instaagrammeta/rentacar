package service

import (
	"fmt"
	"time"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/media"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/pricing"
)

// DailyRevenue returns payments and total for a given day (defaults to today).
func (s *Service) DailyRevenue(day *time.Time) map[string]interface{} {
	d := time.Now()
	if day != nil {
		d = *day
	}
	start := startOfDay(d)
	end := start.AddDate(0, 0, 1)

	var payments []models.Payment
	s.DB.Preload("Client").Where("paid_at >= ? AND paid_at < ?", start, end).Order("id desc").Find(&payments)

	total := 0.0
	items := make([]map[string]interface{}, 0, len(payments))
	for i := range payments {
		total += payments[i].Amount
		items = append(items, payments[i].ToMap())
	}
	return map[string]interface{}{
		"date":     start.Format("2006-01-02"),
		"total":    pricing.Round2(total),
		"payments": items,
	}
}

// MonthlyRevenue returns the total revenue for a given year/month.
func (s *Service) MonthlyRevenue(year, month int) map[string]interface{} {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	total := s.revenueBetween(start, end)
	return map[string]interface{}{"year": year, "month": month, "total": pricing.Round2(total)}
}

// YearlyRevenue returns monthly totals and the yearly total.
func (s *Service) YearlyRevenue(year int) map[string]interface{} {
	months := make([]map[string]interface{}, 0, 12)
	total := 0.0
	for m := 1; m <= 12; m++ {
		mr := s.MonthlyRevenue(year, m)
		months = append(months, mr)
		total += mr["total"].(float64)
	}
	return map[string]interface{}{"year": year, "total": pricing.Round2(total), "months": months}
}

// MostProfitableCars returns cars ranked by total payment revenue.
func (s *Service) MostProfitableCars(limit int) []map[string]interface{} {
	if limit <= 0 {
		limit = 10
	}
	type row struct {
		Brand   string
		Model   string
		Plate   string
		Revenue float64
	}
	var rows []row
	s.DB.Model(&models.Car{}).
		Select("cars.brand, cars.model, cars.plate_number as plate, COALESCE(SUM(payments.amount),0) as revenue").
		Joins("JOIN rentals ON rentals.car_id = cars.id").
		Joins("JOIN payments ON payments.rental_id = rentals.id").
		Group("cars.id").
		Order("revenue DESC").
		Limit(limit).
		Scan(&rows)

	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]interface{}{
			"car":     fmt.Sprintf("%s %s (%s)", r.Brand, r.Model, r.Plate),
			"revenue": pricing.Round2(r.Revenue),
		})
	}
	return out
}

// ActiveRentals returns all active rentals.
func (s *Service) ActiveRentals() []map[string]interface{} {
	var rentals []models.Rental
	s.DB.Preload("Client").Preload("Car").Preload("Employee").Preload("VehicleReturn").
		Where("status = ?", models.RentalActive).Order("id desc").Find(&rentals)
	out := make([]map[string]interface{}, 0, len(rentals))
	for i := range rentals {
		out = append(out, rentals[i].ToMap())
	}
	return out
}

// Debtors returns completed rentals with a positive final payment owed.
func (s *Service) Debtors() []map[string]interface{} {
	var rentals []models.Rental
	s.DB.Preload("Client").Preload("VehicleReturn").
		Where("status = ?", models.RentalCompleted).Find(&rentals)
	out := make([]map[string]interface{}, 0)
	for i := range rentals {
		vr := rentals[i].VehicleReturn
		if vr != nil && vr.FinalPayment > 0 {
			name := "—"
			if rentals[i].Client != nil {
				name = rentals[i].Client.FullName()
			}
			out = append(out, map[string]interface{}{
				"client":          name,
				"contract_number": rentals[i].ContractNumber,
				"amount":          pricing.Round2(vr.FinalPayment),
			})
		}
	}
	return out
}

// ClientStatistics returns each client with their rental count.
func (s *Service) ClientStatistics() []map[string]interface{} {
	type row struct {
		FullNameLast  string
		FullNameFirst string
		Code          string
		Rentals       int64
	}
	var rows []row
	s.DB.Model(&models.Client{}).
		Select("clients.last_name as full_name_last, clients.first_name as full_name_first, clients.client_code as code, count(rentals.id) as rentals").
		Joins("LEFT JOIN rentals ON rentals.client_id = clients.id").
		Group("clients.id").
		Order("rentals DESC").
		Scan(&rows)

	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		name := r.FullNameLast + " " + r.FullNameFirst
		out = append(out, map[string]interface{}{
			"client":  name,
			"code":    r.Code,
			"rentals": r.Rentals,
		})
	}
	return out
}

// ExportReportExcel generates a styled Excel report and returns its file path.
func (s *Service) ExportReportExcel(reportType string, params map[string]string) (string, error) {
	settings, _ := s.GetSettings()
	company := "Rentacar CRM"
	if settings != nil {
		company = settings.CompanyName
	}

	switch reportType {
	case "daily_revenue":
		var day *time.Time
		if params["date"] != "" {
			if d, err := parseDate(params["date"]); err == nil {
				day = d
			}
		}
		data := s.DailyRevenue(day)
		payments := data["payments"].([]map[string]interface{})
		rows := make([][]interface{}, 0, len(payments)+1)
		for _, p := range payments {
			rows = append(rows, []interface{}{
				p["receipt_number"], strVal(p["client_name"]), strVal(p["payment_type_label"]),
				strVal(p["method_label"]), p["amount"], p["paid_at"],
			})
		}
		rows = append(rows, []interface{}{"", "", "", "ИТОГО", data["total"], ""})
		return s.Media.ExportReportExcel(media.ExcelReport{
			SheetName: "Выручка за день", Title: "Выручка за " + data["date"].(string), CompanyName: company,
			Headers: []string{"Квитанция", "Клиент", "Тип", "Способ", "Сумма", "Дата"}, Rows: rows,
		})

	case "most_profitable_cars":
		data := s.MostProfitableCars(atoiDefault(params["limit"], 10))
		rows := make([][]interface{}, 0, len(data))
		for _, d := range data {
			rows = append(rows, []interface{}{d["car"], d["revenue"]})
		}
		return s.Media.ExportReportExcel(media.ExcelReport{
			SheetName: "Прибыльные авто", Title: "Самые прибыльные автомобили", CompanyName: company,
			Headers: []string{"Автомобиль", "Выручка"}, Rows: rows,
		})

	case "active_rentals":
		data := s.ActiveRentals()
		rows := make([][]interface{}, 0, len(data))
		for _, r := range data {
			rows = append(rows, []interface{}{
				r["contract_number"], strVal(r["client_name"]), strVal(r["car_name"]),
				r["rental_start"], r["rental_end"], r["total_price"],
			})
		}
		return s.Media.ExportReportExcel(media.ExcelReport{
			SheetName: "Активные аренды", Title: "Активные аренды", CompanyName: company,
			Headers: []string{"Договор", "Клиент", "Автомобиль", "Начало", "Окончание", "Сумма"}, Rows: rows,
		})

	case "debtors":
		data := s.Debtors()
		rows := make([][]interface{}, 0, len(data))
		for _, d := range data {
			rows = append(rows, []interface{}{d["client"], d["contract_number"], d["amount"]})
		}
		return s.Media.ExportReportExcel(media.ExcelReport{
			SheetName: "Должники", Title: "Список должников", CompanyName: company,
			Headers: []string{"Клиент", "Договор", "Сумма долга"}, Rows: rows,
		})

	case "client_statistics":
		data := s.ClientStatistics()
		rows := make([][]interface{}, 0, len(data))
		for _, d := range data {
			rows = append(rows, []interface{}{d["client"], d["code"], d["rentals"]})
		}
		return s.Media.ExportReportExcel(media.ExcelReport{
			SheetName: "Статистика клиентов", Title: "Статистика по клиентам", CompanyName: company,
			Headers: []string{"Клиент", "Код", "Кол-во аренд"}, Rows: rows,
		})

	case "yearly_revenue":
		year := atoiDefault(params["year"], time.Now().Year())
		data := s.YearlyRevenue(year)
		months := data["months"].([]map[string]interface{})
		rows := make([][]interface{}, 0, len(months)+1)
		for _, m := range months {
			rows = append(rows, []interface{}{fmt.Sprintf("%d-%02d", m["year"], m["month"]), m["total"]})
		}
		rows = append(rows, []interface{}{"ИТОГО", data["total"]})
		return s.Media.ExportReportExcel(media.ExcelReport{
			SheetName: "Годовая выручка", Title: fmt.Sprintf("Выручка за %d год", year), CompanyName: company,
			Headers: []string{"Месяц", "Выручка"}, Rows: rows,
		})

	default:
		return "", apperr.Validation("Тип отчёта не поддерживается")
	}
}

func strVal(v interface{}) interface{} {
	if v == nil {
		return ""
	}
	return v
}
