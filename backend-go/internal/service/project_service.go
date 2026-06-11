package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

const projectFormatVersion = "1.0"

// exportOrder lists tables in dependency (FK) order for export/import.
var exportOrder = []string{
	"users", "company_settings", "clients", "cars", "car_photos",
	"reservations", "rentals", "vehicle_returns", "payments",
	"accidents", "accident_photos", "blacklist_entries",
}

// dumpData returns every table serialised to JSON-friendly column maps.
func (s *Service) dumpData() map[string][]map[string]interface{} {
	data := map[string][]map[string]interface{}{}

	var users []models.User
	s.DB.Find(&users)
	for i := range users {
		m := users[i].ToMap()
		m["password_hash"] = users[i].PasswordHash // include for full restore
		data["users"] = append(data["users"], m)
	}

	var settings []models.CompanySettings
	s.DB.Find(&settings)
	for i := range settings {
		data["company_settings"] = append(data["company_settings"], settings[i].ToMap())
	}

	var clients []models.Client
	s.DB.Find(&clients)
	for i := range clients {
		data["clients"] = append(data["clients"], clients[i].ToMap())
	}

	var cars []models.Car
	s.DB.Find(&cars)
	for i := range cars {
		data["cars"] = append(data["cars"], cars[i].ToMap())
	}

	var carPhotos []models.CarPhoto
	s.DB.Find(&carPhotos)
	for i := range carPhotos {
		data["car_photos"] = append(data["car_photos"], map[string]interface{}{
			"id": carPhotos[i].ID, "car_id": carPhotos[i].CarID, "file_path": carPhotos[i].FilePath,
		})
	}

	var reservations []models.Reservation
	s.DB.Find(&reservations)
	for i := range reservations {
		data["reservations"] = append(data["reservations"], reservations[i].ToMap())
	}

	var rentals []models.Rental
	s.DB.Find(&rentals)
	for i := range rentals {
		data["rentals"] = append(data["rentals"], rentals[i].ToMap())
	}

	var returns []models.VehicleReturn
	s.DB.Find(&returns)
	for i := range returns {
		data["vehicle_returns"] = append(data["vehicle_returns"], returns[i].ToMap())
	}

	var payments []models.Payment
	s.DB.Find(&payments)
	for i := range payments {
		data["payments"] = append(data["payments"], payments[i].ToMap())
	}

	var accidents []models.Accident
	s.DB.Find(&accidents)
	for i := range accidents {
		data["accidents"] = append(data["accidents"], accidents[i].ToMap())
	}

	var accidentPhotos []models.AccidentPhoto
	s.DB.Find(&accidentPhotos)
	for i := range accidentPhotos {
		data["accident_photos"] = append(data["accident_photos"], map[string]interface{}{
			"id": accidentPhotos[i].ID, "accident_id": accidentPhotos[i].AccidentID, "file_path": accidentPhotos[i].FilePath,
		})
	}

	var blacklist []models.BlacklistEntry
	s.DB.Find(&blacklist)
	for i := range blacklist {
		data["blacklist_entries"] = append(data["blacklist_entries"], blacklist[i].ToMap())
	}

	return data
}

// collectUploadPaths gathers all relative upload paths referenced by the data.
func collectUploadPaths(data map[string][]map[string]interface{}) map[string]bool {
	paths := map[string]bool{}
	keys := []string{"passport_scan", "passport_front", "passport_back", "driver_license_scan", "driver_photo", "qr_code_path", "file_path", "logo_path"}
	for _, rows := range data {
		for _, row := range rows {
			for _, k := range keys {
				if v, ok := row[k].(string); ok && v != "" {
					paths[v] = true
				}
			}
		}
	}
	return paths
}

// ExportProject writes the whole database to a .rentacar ZIP and returns its path.
func (s *Service) ExportProject(actor *Actor) (string, error) {
	data := s.dumpData()

	tables := map[string]int{}
	for name, rows := range data {
		tables[name] = len(rows)
	}
	manifest := map[string]interface{}{
		"format":      "rentacar",
		"version":     projectFormatVersion,
		"exported_at": time.Now().Format(time.RFC3339),
		"tables":      tables,
	}

	dir := filepath.Join(s.Cfg.ExportDir, "projects")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	target := filepath.Join(dir, fmt.Sprintf("project_%s.rentacar", time.Now().Format("20060102_150405")))

	out, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer out.Close()

	archive := zip.NewWriter(out)
	defer archive.Close()

	if err := writeJSONEntry(archive, "manifest.json", manifest); err != nil {
		return "", err
	}
	if err := writeJSONEntry(archive, "data.json", data); err != nil {
		return "", err
	}
	for rel := range collectUploadPaths(data) {
		src := filepath.Join(s.Cfg.UploadDir, filepath.FromSlash(rel))
		if f, err := os.Open(src); err == nil {
			w, _ := archive.Create("uploads/" + rel)
			_, _ = io.Copy(w, f)
			f.Close()
		}
	}

	s.record(actor, "export_project", "project", nil, target)
	return target, nil
}

func writeJSONEntry(archive *zip.Writer, name string, value interface{}) error {
	w, err := archive.Create(name)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}

// ImportProject restores a .rentacar archive, replacing all current data.
func (s *Service) ImportProject(filePath string, actor *Actor) (map[string]interface{}, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, apperr.Validation("Некорректный формат файла .rentacar")
	}
	defer r.Close()

	var data map[string][]map[string]interface{}
	found := false
	for _, f := range r.File {
		if f.Name == "data.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			dec := json.NewDecoder(rc)
			if err := dec.Decode(&data); err != nil {
				rc.Close()
				return nil, apperr.Validation("Не удалось прочитать data.json")
			}
			rc.Close()
			found = true
		}
	}
	if !found {
		return nil, apperr.Validation("В файле отсутствует data.json")
	}

	// Restore uploaded files.
	for _, f := range r.File {
		if len(f.Name) > len("uploads/") && f.Name[:len("uploads/")] == "uploads/" && f.Name[len(f.Name)-1] != '/' {
			rel := f.Name[len("uploads/"):]
			target := filepath.Join(s.Cfg.UploadDir, filepath.FromSlash(rel))
			_ = os.MkdirAll(filepath.Dir(target), 0o755)
			if rc, err := f.Open(); err == nil {
				if out, err := os.Create(target); err == nil {
					_, _ = io.Copy(out, rc)
					out.Close()
				}
				rc.Close()
			}
		}
	}

	// Clear existing data in reverse FK order.
	for i := len(exportOrder) - 1; i >= 0; i-- {
		s.DB.Exec("DELETE FROM " + exportOrder[i])
	}

	counts := map[string]interface{}{}
	for _, table := range exportOrder {
		rows := data[table]
		for _, row := range rows {
			s.insertRow(table, row)
		}
		counts[table] = len(rows)
		s.resetSequence(table)
	}

	s.record(actor, "import_project", "project", nil, fmt.Sprintf("Импортировано: %v", counts))
	return map[string]interface{}{"imported": counts}, nil
}

// insertRow constructs a model from a column map and inserts it (preserving id).
func (s *Service) insertRow(table string, m map[string]interface{}) {
	switch table {
	case "users":
		s.DB.Create(&models.User{
			Base: base(m), Username: asStr(m, "username"), Email: asStrPtr(m, "email"),
			FullName: asStr(m, "full_name"), PasswordHash: asStr(m, "password_hash"),
			Role: models.UserRole(asStr(m, "role")), IsActive: asBool(m, "is_active"),
		})
	case "company_settings":
		s.DB.Create(&models.CompanySettings{
			Base: base(m), CompanyName: asStr(m, "company_name"), Address: asStrPtr(m, "address"),
			Phone: asStrPtr(m, "phone"), Email: asStrPtr(m, "email"), LogoPath: asStrPtr(m, "logo_path"),
			Currency: asStr(m, "currency"), ContractTerms: asStrPtr(m, "contract_terms"),
		})
	case "clients":
		s.DB.Create(&models.Client{
			Base: base(m), ClientCode: asStr(m, "client_code"), FirstName: asStr(m, "first_name"),
			LastName: asStr(m, "last_name"), Phone: asStr(m, "phone"), Email: asStrPtr(m, "email"),
			DateOfBirth: asDatePtr(m, "date_of_birth"), PassportNumber: asStrPtr(m, "passport_number"),
			DriverLicenseNumber: asStrPtr(m, "driver_license_number"), DriverLicenseIssueDate: asDatePtr(m, "driver_license_issue_date"),
			DriverExperienceYears: asInt(m, "driver_experience_years"), PassportScan: asStrPtr(m, "passport_scan"),
			PassportFront: asStrPtr(m, "passport_front"), PassportBack: asStrPtr(m, "passport_back"),
			DriverLicenseScan: asStrPtr(m, "driver_license_scan"), DriverPhoto: asStrPtr(m, "driver_photo"), QRCodePath: asStrPtr(m, "qr_code_path"),
			IsVIP: asBool(m, "is_vip"), Notes: asStrPtr(m, "notes"), Status: models.ClientStatus(asStr(m, "status")),
		})
	case "cars":
		s.DB.Create(&models.Car{
			Base: base(m), Brand: asStr(m, "brand"), Model: asStr(m, "model"), Year: asInt(m, "year"),
			Color: asStrPtr(m, "color"), VIN: asStrPtr(m, "vin"), PlateNumber: asStr(m, "plate_number"),
			Mileage: asInt(m, "mileage"), DailyPrice: asFloat(m, "daily_price"), WeeklyPrice: asFloat(m, "weekly_price"),
			MonthlyPrice: asFloat(m, "monthly_price"), DepositAmount: asFloat(m, "deposit_amount"),
			Status: models.CarStatus(asStr(m, "status")),
		})
	case "car_photos":
		s.DB.Create(&models.CarPhoto{Base: base(m), CarID: asUint(m, "car_id"), FilePath: asStr(m, "file_path")})
	case "reservations":
		s.DB.Create(&models.Reservation{
			Base: base(m), ClientID: asUint(m, "client_id"), CarID: asUint(m, "car_id"),
			StartDate: asDateVal(m, "start_date"), EndDate: asDateVal(m, "end_date"),
			Deposit: asFloat(m, "deposit"), Notes: asStrPtr(m, "notes"), Status: models.ReservationStatus(asStr(m, "status")),
		})
	case "rentals":
		s.DB.Create(&models.Rental{
			Base: base(m), ContractNumber: asStr(m, "contract_number"), ClientID: asUint(m, "client_id"),
			CarID: asUint(m, "car_id"), ReservationID: asUintPtr(m, "reservation_id"), EmployeeID: asUintPtr(m, "employee_id"),
			RentalStart: asDateVal(m, "rental_start"), RentalEnd: asDateVal(m, "rental_end"),
			PickupAt: asDateTimePtr(m, "pickup_at"), DueAt: asDateTimePtr(m, "due_at"),
			Reminder1hSent: asBool(m, "reminder_1h_sent"), Reminder30mSent: asBool(m, "reminder_30m_sent"),
			Deposit: asFloat(m, "deposit"), DailyPrice: asFloat(m, "daily_price"), TotalPrice: asFloat(m, "total_price"),
			StartMileage: asIntPtr(m, "start_mileage"), PDFPath: asStrPtr(m, "pdf_path"), Notes: asStrPtr(m, "notes"),
			Status: models.RentalStatus(asStr(m, "status")),
		})
	case "vehicle_returns":
		s.DB.Create(&models.VehicleReturn{
			Base: base(m), RentalID: asUint(m, "rental_id"), ReturnDate: asDateTime(m, "return_date"),
			Mileage: asInt(m, "mileage"), FuelLevel: asInt(m, "fuel_level"), Damages: asStrPtr(m, "damages"),
			ExtraDays: asInt(m, "extra_days"), LateFee: asFloat(m, "late_fee"), DamageCost: asFloat(m, "damage_cost"),
			Penalties: asFloat(m, "penalties"), FinalPayment: asFloat(m, "final_payment"),
		})
	case "payments":
		s.DB.Create(&models.Payment{
			Base: base(m), ReceiptNumber: asStr(m, "receipt_number"), ClientID: asUint(m, "client_id"),
			RentalID: asUintPtr(m, "rental_id"), CashierID: asUintPtr(m, "cashier_id"), Amount: asFloat(m, "amount"),
			Method: models.PaymentMethod(asStr(m, "method")), PaymentType: models.PaymentType(asStr(m, "payment_type")),
			PaidAt: asDateTime(m, "paid_at"), ReceiptPath: asStrPtr(m, "receipt_path"), Notes: asStrPtr(m, "notes"),
		})
	case "accidents":
		s.DB.Create(&models.Accident{
			Base: base(m), CarID: asUint(m, "car_id"), ClientID: asUintPtr(m, "client_id"), RentalID: asUintPtr(m, "rental_id"),
			AccidentDate: asDateVal(m, "accident_date"), Description: asStrPtr(m, "description"), RepairCost: asFloat(m, "repair_cost"),
		})
	case "accident_photos":
		s.DB.Create(&models.AccidentPhoto{Base: base(m), AccidentID: asUint(m, "accident_id"), FilePath: asStr(m, "file_path")})
	case "blacklist_entries":
		s.DB.Create(&models.BlacklistEntry{
			Base: base(m), ClientID: asUint(m, "client_id"), Reason: models.BlacklistReason(asStr(m, "reason")),
			Comment: asStrPtr(m, "comment"), CreatedBy: asUintPtr(m, "created_by"),
		})
	}
}

// resetSequence aligns the Postgres identity sequence with the max id present.
func (s *Service) resetSequence(table string) {
	s.DB.Exec(fmt.Sprintf(
		"SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE((SELECT MAX(id) FROM %s), 1))",
		table, table,
	))
}
