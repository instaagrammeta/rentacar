// Package seed inserts default users, company settings and demo data so the
// application is usable immediately after the first launch.
package seed

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/auth"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// Run seeds the database. It is idempotent: existing rows are left untouched.
func Run(db *gorm.DB, withDemo bool) {
	ensureAdmin(db)
	ensureSettings(db)
	if withDemo {
		ensureDemoUsers(db)
		ensureDemoData(db)
	}
}

func userExists(db *gorm.DB, username string) bool {
	var count int64
	db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func ensureAdmin(db *gorm.DB) {
	if userExists(db, "admin") {
		return
	}
	hash, _ := auth.HashPassword("admin123")
	email := "admin@rentacar.local"
	db.Create(&models.User{
		Username:     "admin",
		FullName:     "Администратор системы",
		Email:        &email,
		Role:         models.RoleAdministrator,
		PasswordHash: hash,
		IsActive:     true,
	})
	log.Println("Создан администратор: admin / admin123")
}

func ensureDemoUsers(db *gorm.DB) {
	demo := []struct {
		Username string
		FullName string
		Role     models.UserRole
	}{
		{"manager", "Менеджер Иванов", models.RoleRentalManager},
		{"cashier", "Кассир Петрова", models.RoleCashier},
		{"operator", "Оператор Сидоров", models.RoleOperator},
	}
	for _, d := range demo {
		if userExists(db, d.Username) {
			continue
		}
		hash, _ := auth.HashPassword(d.Username + "123")
		db.Create(&models.User{
			Username:     d.Username,
			FullName:     d.FullName,
			Role:         d.Role,
			PasswordHash: hash,
			IsActive:     true,
		})
		log.Printf("Создан пользователь: %s / %s123 (%s)", d.Username, d.Username, d.Role.Label())
	}
}

func ensureSettings(db *gorm.DB) {
	var count int64
	db.Model(&models.CompanySettings{}).Count(&count)
	if count > 0 {
		return
	}
	addr := "г. Москва, ул. Примерная, д. 1"
	phone := "+7 (495) 000-00-00"
	email := "info@rentacar.local"
	terms := "Арендатор обязуется бережно использовать транспортное средство и вернуть его в надлежащем состоянии в указанный срок."
	db.Create(&models.CompanySettings{
		CompanyName:   "ООО «Рентакар»",
		Address:       &addr,
		Phone:         &phone,
		Email:         &email,
		Currency:      "RUB",
		ContractTerms: &terms,
	})
	log.Println("Созданы настройки компании")
}

func strPtr(s string) *string { return &s }

func ensureDemoData(db *gorm.DB) {
	var carCount int64
	db.Model(&models.Car{}).Count(&carCount)
	if carCount == 0 {
		cars := []models.Car{
			{Brand: "Toyota", Model: "Camry", Year: 2022, Color: strPtr("Чёрный"), PlateNumber: "А001АА777",
				VIN: strPtr("JT2BF22K1W0123456"), Mileage: 25000, DailyPrice: 3500, WeeklyPrice: 21000,
				MonthlyPrice: 78000, DepositAmount: 15000, Status: models.CarAvailable},
			{Brand: "Kia", Model: "Rio", Year: 2021, Color: strPtr("Белый"), PlateNumber: "В002ВВ777",
				VIN: strPtr("KNADN512AB6123456"), Mileage: 40000, DailyPrice: 2000, WeeklyPrice: 12000,
				MonthlyPrice: 45000, DepositAmount: 8000, Status: models.CarAvailable},
			{Brand: "BMW", Model: "X5", Year: 2023, Color: strPtr("Синий"), PlateNumber: "С003СС777",
				VIN: strPtr("5UXKR0C50J0123456"), Mileage: 12000, DailyPrice: 8000, WeeklyPrice: 49000,
				MonthlyPrice: 185000, DepositAmount: 40000, Status: models.CarAvailable},
		}
		db.Create(&cars)
		log.Printf("Создано автомобилей: %d", len(cars))
	}

	var clientCount int64
	db.Model(&models.Client{}).Count(&clientCount)
	if clientCount == 0 {
		dob30 := time.Now().AddDate(-30, 0, 0)
		dob25 := time.Now().AddDate(-25, 0, 0)
		clients := []models.Client{
			{ClientCode: "CL-000001", FirstName: "Алексей", LastName: "Смирнов", Phone: "+79161112233",
				Email: strPtr("smirnov@example.com"), DateOfBirth: &dob30, DriverLicenseNumber: strPtr("7700123456"),
				DriverExperienceYears: 8, Status: models.ClientActive, IsVIP: true},
			{ClientCode: "CL-000002", FirstName: "Мария", LastName: "Кузнецова", Phone: "+79162223344",
				Email: strPtr("kuznetsova@example.com"), DateOfBirth: &dob25, DriverLicenseNumber: strPtr("7700654321"),
				DriverExperienceYears: 3, Status: models.ClientActive},
		}
		db.Create(&clients)
		log.Printf("Создано клиентов: %d", len(clients))
	}
}
