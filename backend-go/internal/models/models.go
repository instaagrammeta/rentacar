// Package models defines the GORM domain models mirroring the original
// SQLAlchemy schema, together with ToMap helpers that reproduce the exact JSON
// contract used by the API (including computed labels and related names).
package models

import "time"

// dateLayout is the ISO date used by the original backend (date.isoformat()).
const dateLayout = "2006-01-02"

// isoDateTime mirrors Python's datetime.isoformat() closely enough for the API.
func isoDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04:05Z07:00")
}

func isoDate(t *time.Time) interface{} {
	if t == nil || t.IsZero() {
		return nil
	}
	return t.Format(dateLayout)
}

// Base contains the surrogate key and timestamps shared by every entity.
type Base struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b Base) baseMap(m map[string]interface{}) {
	m["id"] = b.ID
	m["created_at"] = isoDateTime(b.CreatedAt)
	m["updated_at"] = isoDateTime(b.UpdatedAt)
}

// ----------------------------------------------------------------------- User

// User is an employee account that can log into the CRM.
type User struct {
	Base
	Username     string   `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email        *string  `gorm:"size:120;uniqueIndex" json:"email"`
	FullName     string   `gorm:"size:160;not null;default:''" json:"full_name"`
	PasswordHash string   `gorm:"size:255;not null" json:"-"`
	Role         UserRole `gorm:"size:32;not null;default:operator" json:"role"`
	IsActive     bool     `gorm:"not null;default:true" json:"is_active"`
}

func (User) TableName() string { return "users" }

// ToMap serialises a user, excluding the password hash.
func (u *User) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	u.baseMap(m)
	m["username"] = u.Username
	m["email"] = u.Email
	m["full_name"] = u.FullName
	m["role"] = string(u.Role)
	m["role_label"] = u.Role.Label()
	m["is_active"] = u.IsActive
	return m
}

// --------------------------------------------------------------------- Client

// Client is a rental customer with documents and verification state.
type Client struct {
	Base
	ClientCode             string       `gorm:"size:32;uniqueIndex;not null" json:"client_code"`
	FirstName              string       `gorm:"size:80;not null" json:"first_name"`
	LastName               string       `gorm:"size:80;not null" json:"last_name"`
	Phone                  string       `gorm:"size:32;index;not null" json:"phone"`
	Email                  *string      `gorm:"size:120;index" json:"email"`
	DateOfBirth            *time.Time   `gorm:"type:date" json:"date_of_birth"`
	PassportNumber         *string      `gorm:"size:64;index" json:"passport_number"`
	DriverLicenseNumber    *string      `gorm:"size:64;index" json:"driver_license_number"`
	DriverLicenseIssueDate *time.Time   `gorm:"type:date" json:"driver_license_issue_date"`
	DriverExperienceYears  int          `gorm:"not null;default:0" json:"driver_experience_years"`
	PassportScan           *string      `gorm:"size:255" json:"passport_scan"`
	DriverLicenseScan      *string      `gorm:"size:255" json:"driver_license_scan"`
	QRCodePath             *string      `gorm:"size:255" json:"qr_code_path"`
	IsVIP                  bool         `gorm:"not null;default:false" json:"is_vip"`
	Notes                  *string      `gorm:"type:text" json:"notes"`
	Status                 ClientStatus `gorm:"size:32;not null;default:pending_verification" json:"status"`
}

func (Client) TableName() string { return "clients" }

// FullName mirrors the Python property: "<last> <first>".
func (c *Client) FullName() string {
	name := c.LastName + " " + c.FirstName
	if name == " " {
		return ""
	}
	return name
}

func (c *Client) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	c.baseMap(m)
	m["client_code"] = c.ClientCode
	m["first_name"] = c.FirstName
	m["last_name"] = c.LastName
	m["phone"] = c.Phone
	m["email"] = c.Email
	m["date_of_birth"] = isoDate(c.DateOfBirth)
	m["passport_number"] = c.PassportNumber
	m["driver_license_number"] = c.DriverLicenseNumber
	m["driver_license_issue_date"] = isoDate(c.DriverLicenseIssueDate)
	m["driver_experience_years"] = c.DriverExperienceYears
	m["passport_scan"] = c.PassportScan
	m["driver_license_scan"] = c.DriverLicenseScan
	m["qr_code_path"] = c.QRCodePath
	m["is_vip"] = c.IsVIP
	m["notes"] = c.Notes
	m["status"] = string(c.Status)
	m["full_name"] = c.FullName()
	m["status_label"] = c.Status.Label()
	return m
}

// ------------------------------------------------------------------------ Car

// Car is a vehicle available for rental.
type Car struct {
	Base
	Brand         string     `gorm:"size:80;not null;index" json:"brand"`
	Model         string     `gorm:"size:80;not null;index" json:"model"`
	Year          int        `gorm:"not null" json:"year"`
	Color         *string    `gorm:"size:40" json:"color"`
	VIN           *string    `gorm:"size:32;uniqueIndex" json:"vin"`
	PlateNumber   string     `gorm:"size:20;uniqueIndex;not null" json:"plate_number"`
	Mileage       int        `gorm:"not null;default:0" json:"mileage"`
	DailyPrice    float64    `gorm:"not null;default:0" json:"daily_price"`
	WeeklyPrice   float64    `gorm:"not null;default:0" json:"weekly_price"`
	MonthlyPrice  float64    `gorm:"not null;default:0" json:"monthly_price"`
	DepositAmount float64    `gorm:"not null;default:0" json:"deposit_amount"`
	Status        CarStatus  `gorm:"size:32;not null;default:available" json:"status"`
	Photos        []CarPhoto `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

func (Car) TableName() string { return "cars" }

// DisplayName mirrors the Python property: "<brand> <model> (<plate>)".
func (c *Car) DisplayName() string {
	return c.Brand + " " + c.Model + " (" + c.PlateNumber + ")"
}

func (c *Car) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	c.baseMap(m)
	m["brand"] = c.Brand
	m["model"] = c.Model
	m["year"] = c.Year
	m["color"] = c.Color
	m["vin"] = c.VIN
	m["plate_number"] = c.PlateNumber
	m["mileage"] = c.Mileage
	m["daily_price"] = c.DailyPrice
	m["weekly_price"] = c.WeeklyPrice
	m["monthly_price"] = c.MonthlyPrice
	m["deposit_amount"] = c.DepositAmount
	m["status"] = string(c.Status)
	m["display_name"] = c.DisplayName()
	m["status_label"] = c.Status.Label()
	photos := make([]string, 0, len(c.Photos))
	for _, p := range c.Photos {
		photos = append(photos, p.FilePath)
	}
	m["photos"] = photos
	return m
}

// CarPhoto is a single photo belonging to a car.
type CarPhoto struct {
	Base
	CarID    uint   `gorm:"not null;index" json:"car_id"`
	FilePath string `gorm:"size:255;not null" json:"file_path"`
}

func (CarPhoto) TableName() string { return "car_photos" }

// ---------------------------------------------------------------- Reservation

// Reservation reserves a car for a client for a date range.
type Reservation struct {
	Base
	ClientID  uint              `gorm:"not null;index" json:"client_id"`
	CarID     uint              `gorm:"not null;index" json:"car_id"`
	StartDate time.Time         `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time         `gorm:"type:date;not null" json:"end_date"`
	Deposit   float64           `gorm:"not null;default:0" json:"deposit"`
	Notes     *string           `gorm:"type:text" json:"notes"`
	Status    ReservationStatus `gorm:"size:32;not null;default:reserved" json:"status"`

	Client *Client `gorm:"foreignKey:ClientID" json:"-"`
	Car    *Car    `gorm:"foreignKey:CarID" json:"-"`
}

func (Reservation) TableName() string { return "reservations" }

func (r *Reservation) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	r.baseMap(m)
	m["client_id"] = r.ClientID
	m["car_id"] = r.CarID
	m["start_date"] = r.StartDate.Format(dateLayout)
	m["end_date"] = r.EndDate.Format(dateLayout)
	m["deposit"] = r.Deposit
	m["notes"] = r.Notes
	m["status"] = string(r.Status)
	m["status_label"] = r.Status.Label()
	if r.Client != nil {
		m["client_name"] = r.Client.FullName()
	}
	if r.Car != nil {
		m["car_name"] = r.Car.DisplayName()
	}
	return m
}

// --------------------------------------------------------------------- Rental

// Rental is a signed rental contract.
type Rental struct {
	Base
	ContractNumber string       `gorm:"size:32;uniqueIndex;not null" json:"contract_number"`
	ClientID       uint         `gorm:"not null;index" json:"client_id"`
	CarID          uint         `gorm:"not null;index" json:"car_id"`
	ReservationID  *uint        `json:"reservation_id"`
	EmployeeID     *uint        `json:"employee_id"`
	RentalStart    time.Time    `gorm:"type:date;not null" json:"rental_start"`
	RentalEnd      time.Time    `gorm:"type:date;not null" json:"rental_end"`
	Deposit        float64      `gorm:"not null;default:0" json:"deposit"`
	DailyPrice     float64      `gorm:"not null;default:0" json:"daily_price"`
	TotalPrice     float64      `gorm:"not null;default:0" json:"total_price"`
	StartMileage   *int         `json:"start_mileage"`
	PDFPath        *string      `gorm:"size:255" json:"pdf_path"`
	Notes          *string      `gorm:"type:text" json:"notes"`
	Status         RentalStatus `gorm:"size:32;not null;default:active" json:"status"`

	Client        *Client        `gorm:"foreignKey:ClientID" json:"-"`
	Car           *Car           `gorm:"foreignKey:CarID" json:"-"`
	Employee      *User          `gorm:"foreignKey:EmployeeID" json:"-"`
	VehicleReturn *VehicleReturn `gorm:"foreignKey:RentalID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Rental) TableName() string { return "rentals" }

func (r *Rental) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	r.baseMap(m)
	m["contract_number"] = r.ContractNumber
	m["client_id"] = r.ClientID
	m["car_id"] = r.CarID
	m["reservation_id"] = r.ReservationID
	m["employee_id"] = r.EmployeeID
	m["rental_start"] = r.RentalStart.Format(dateLayout)
	m["rental_end"] = r.RentalEnd.Format(dateLayout)
	m["deposit"] = r.Deposit
	m["daily_price"] = r.DailyPrice
	m["total_price"] = r.TotalPrice
	m["start_mileage"] = r.StartMileage
	m["pdf_path"] = r.PDFPath
	m["notes"] = r.Notes
	m["status"] = string(r.Status)
	m["status_label"] = r.Status.Label()
	if r.Client != nil {
		m["client_name"] = r.Client.FullName()
	}
	if r.Car != nil {
		m["car_name"] = r.Car.DisplayName()
	}
	if r.Employee != nil {
		m["employee_name"] = r.Employee.FullName
	}
	m["has_return"] = r.VehicleReturn != nil
	return m
}

// VehicleReturn records the state of a returned car and final charges.
type VehicleReturn struct {
	Base
	RentalID     uint      `gorm:"uniqueIndex;not null" json:"rental_id"`
	ReturnDate   time.Time `gorm:"not null" json:"return_date"`
	Mileage      int       `gorm:"not null;default:0" json:"mileage"`
	FuelLevel    int       `gorm:"not null;default:100" json:"fuel_level"`
	Damages      *string   `gorm:"type:text" json:"damages"`
	ExtraDays    int       `gorm:"not null;default:0" json:"extra_days"`
	LateFee      float64   `gorm:"not null;default:0" json:"late_fee"`
	DamageCost   float64   `gorm:"not null;default:0" json:"damage_cost"`
	Penalties    float64   `gorm:"not null;default:0" json:"penalties"`
	FinalPayment float64   `gorm:"not null;default:0" json:"final_payment"`
}

func (VehicleReturn) TableName() string { return "vehicle_returns" }

func (v *VehicleReturn) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	v.baseMap(m)
	m["rental_id"] = v.RentalID
	m["return_date"] = isoDateTime(v.ReturnDate)
	m["mileage"] = v.Mileage
	m["fuel_level"] = v.FuelLevel
	m["damages"] = v.Damages
	m["extra_days"] = v.ExtraDays
	m["late_fee"] = v.LateFee
	m["damage_cost"] = v.DamageCost
	m["penalties"] = v.Penalties
	m["final_payment"] = v.FinalPayment
	return m
}

// -------------------------------------------------------------------- Payment

// Payment is a money transaction linked to a client and optionally a rental.
type Payment struct {
	Base
	ReceiptNumber string        `gorm:"size:32;uniqueIndex;not null" json:"receipt_number"`
	ClientID      uint          `gorm:"not null;index" json:"client_id"`
	RentalID      *uint         `json:"rental_id"`
	CashierID     *uint         `json:"cashier_id"`
	Amount        float64       `gorm:"not null" json:"amount"`
	Method        PaymentMethod `gorm:"size:32;not null" json:"method"`
	PaymentType   PaymentType   `gorm:"size:32;not null;default:rental" json:"payment_type"`
	PaidAt        time.Time     `gorm:"not null" json:"paid_at"`
	ReceiptPath   *string       `gorm:"size:255" json:"receipt_path"`
	Notes         *string       `gorm:"type:text" json:"notes"`

	Client  *Client `gorm:"foreignKey:ClientID" json:"-"`
	Rental  *Rental `gorm:"foreignKey:RentalID" json:"-"`
	Cashier *User   `gorm:"foreignKey:CashierID" json:"-"`
}

func (Payment) TableName() string { return "payments" }

func (p *Payment) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	p.baseMap(m)
	m["receipt_number"] = p.ReceiptNumber
	m["client_id"] = p.ClientID
	m["rental_id"] = p.RentalID
	m["cashier_id"] = p.CashierID
	m["amount"] = p.Amount
	m["method"] = string(p.Method)
	m["method_label"] = p.Method.Label()
	m["payment_type"] = string(p.PaymentType)
	m["payment_type_label"] = p.PaymentType.Label()
	m["paid_at"] = isoDateTime(p.PaidAt)
	m["receipt_path"] = p.ReceiptPath
	m["notes"] = p.Notes
	if p.Client != nil {
		m["client_name"] = p.Client.FullName()
	}
	return m
}

// ------------------------------------------------------------------ Blacklist

// BlacklistEntry marks a client as blacklisted.
type BlacklistEntry struct {
	Base
	ClientID  uint            `gorm:"uniqueIndex;not null" json:"client_id"`
	Reason    BlacklistReason `gorm:"size:32;not null" json:"reason"`
	Comment   *string         `gorm:"type:text" json:"comment"`
	CreatedBy *uint           `json:"created_by"`

	Client *Client `gorm:"foreignKey:ClientID" json:"-"`
}

func (BlacklistEntry) TableName() string { return "blacklist_entries" }

func (b *BlacklistEntry) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	b.baseMap(m)
	m["client_id"] = b.ClientID
	m["reason"] = string(b.Reason)
	m["reason_label"] = b.Reason.Label()
	m["comment"] = b.Comment
	m["created_by"] = b.CreatedBy
	if b.Client != nil {
		m["client_name"] = b.Client.FullName()
	}
	return m
}

// ------------------------------------------------------------------- Accident

// Accident is a damage event involving a car and optionally a client.
type Accident struct {
	Base
	CarID        uint            `gorm:"not null;index" json:"car_id"`
	ClientID     *uint           `json:"client_id"`
	RentalID     *uint           `json:"rental_id"`
	AccidentDate time.Time       `gorm:"type:date;not null" json:"accident_date"`
	Description  *string         `gorm:"type:text" json:"description"`
	RepairCost   float64         `gorm:"not null;default:0" json:"repair_cost"`
	Photos       []AccidentPhoto `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	Car    *Car    `gorm:"foreignKey:CarID" json:"-"`
	Client *Client `gorm:"foreignKey:ClientID" json:"-"`
}

func (Accident) TableName() string { return "accidents" }

func (a *Accident) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	a.baseMap(m)
	m["car_id"] = a.CarID
	m["client_id"] = a.ClientID
	m["rental_id"] = a.RentalID
	m["accident_date"] = a.AccidentDate.Format(dateLayout)
	m["description"] = a.Description
	m["repair_cost"] = a.RepairCost
	photos := make([]string, 0, len(a.Photos))
	for _, p := range a.Photos {
		photos = append(photos, p.FilePath)
	}
	m["photos"] = photos
	if a.Car != nil {
		m["car_name"] = a.Car.DisplayName()
	}
	if a.Client != nil {
		m["client_name"] = a.Client.FullName()
	}
	return m
}

// AccidentPhoto is a photo attached to an accident record.
type AccidentPhoto struct {
	Base
	AccidentID uint   `gorm:"not null;index" json:"accident_id"`
	FilePath   string `gorm:"size:255;not null" json:"file_path"`
}

func (AccidentPhoto) TableName() string { return "accident_photos" }

// ------------------------------------------------------------------ AuditLog

// AuditLog records logins, actions and data changes performed by users.
type AuditLog struct {
	Base
	UserID    *uint   `json:"user_id"`
	Username  *string `gorm:"size:64" json:"username"`
	Action    string  `gorm:"size:64;not null;index" json:"action"`
	Entity    *string `gorm:"size:64;index" json:"entity"`
	EntityID  *uint   `json:"entity_id"`
	IPAddress *string `gorm:"size:64" json:"ip_address"`
	Details   *string `gorm:"type:text" json:"details"`
}

func (AuditLog) TableName() string { return "audit_logs" }

func (a *AuditLog) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	a.baseMap(m)
	m["user_id"] = a.UserID
	m["username"] = a.Username
	m["action"] = a.Action
	m["entity"] = a.Entity
	m["entity_id"] = a.EntityID
	m["ip_address"] = a.IPAddress
	m["details"] = a.Details
	return m
}

// ------------------------------------------------------------------ Settings

// CompanySettings is a single-row table storing company information.
type CompanySettings struct {
	Base
	CompanyName   string  `gorm:"size:160;not null;default:'Rentacar CRM'" json:"company_name"`
	Address       *string `gorm:"type:text" json:"address"`
	Phone         *string `gorm:"size:64" json:"phone"`
	Email         *string `gorm:"size:120" json:"email"`
	LogoPath      *string `gorm:"size:255" json:"logo_path"`
	Currency      string  `gorm:"size:8;not null;default:'RUB'" json:"currency"`
	ContractTerms *string `gorm:"type:text" json:"contract_terms"`
}

func (CompanySettings) TableName() string { return "company_settings" }

func (s *CompanySettings) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	s.baseMap(m)
	m["company_name"] = s.CompanyName
	m["address"] = s.Address
	m["phone"] = s.Phone
	m["email"] = s.Email
	m["logo_path"] = s.LogoPath
	m["currency"] = s.Currency
	m["contract_terms"] = s.ContractTerms
	return m
}

// AllModels returns every model for AutoMigrate, in dependency order.
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&CompanySettings{},
		&Client{},
		&Car{},
		&CarPhoto{},
		&Reservation{},
		&Rental{},
		&VehicleReturn{},
		&Payment{},
		&Accident{},
		&AccidentPhoto{},
		&BlacklistEntry{},
		&AuditLog{},
	}
}
