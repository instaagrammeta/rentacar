package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// SetupRouter builds the Gin engine with all routes and middleware.
func (h *Handler) SetupRouter() *gin.Engine {
	if !h.Cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.MaxMultipartMemory = h.Cfg.MaxUploadBytes

	// CORS
	corsCfg := cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders: []string{"Content-Disposition"},
		MaxAge:        12 * time.Hour,
	}
	if len(h.Cfg.CORSOrigins) == 1 && h.Cfg.CORSOrigins[0] == "*" {
		corsCfg.AllowAllOrigins = true
	} else {
		corsCfg.AllowOrigins = h.Cfg.CORSOrigins
		corsCfg.AllowCredentials = true
	}
	r.Use(cors.New(corsCfg))

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "Rentacar CRM"})
	})

	api := r.Group("/api")

	auth := h.Auth
	admin := models.RoleAdministrator

	// ---- Auth ----
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
		authed := authGroup.Group("", auth.RequireAuth())
		authed.GET("/me", h.Me)
		authed.POST("/logout", h.Logout)
		authed.GET("/roles", h.Roles)
		authed.GET("/users", auth.RequireRole(admin), h.ListUsers)
		authed.POST("/users", auth.RequireRole(admin), h.CreateUser)
		authed.POST("/users/:id/password", auth.RequireRole(admin), h.ChangePassword)
	}

	// ---- Clients ----
	clients := api.Group("/clients", auth.RequireAuth(), auth.RequirePermission("clients"))
	{
		clients.GET("", h.ListClients)
		clients.GET("/:id", h.GetClient)
		clients.GET("/:id/history", h.ClientHistory)
		clients.POST("", h.CreateClient)
		clients.PUT("/:id", h.UpdateClient)
		clients.DELETE("/:id", h.DeleteClient)
	}

	// ---- Cars ----
	cars := api.Group("/cars", auth.RequireAuth(), auth.RequirePermission("cars"))
	{
		cars.GET("", h.ListCars)
		cars.GET("/:id", h.GetCar)
		cars.POST("", h.CreateCar)
		cars.PUT("/:id", h.UpdateCar)
		cars.DELETE("/:id", h.DeleteCar)
	}
	// Registered outside the /cars group: a static segment cannot be a sibling
	// of the "/cars/:id" wildcard in Gin's router tree.
	api.GET("/cars-available", auth.RequireAuth(), auth.RequirePermission("cars"), h.AvailableCars)

	// ---- Reservations ----
	reservations := api.Group("/reservations", auth.RequireAuth(), auth.RequirePermission("reservations"))
	{
		reservations.GET("", h.ListReservations)
		reservations.GET("/:id", h.GetReservation)
		reservations.POST("", h.CreateReservation)
		reservations.POST("/:id/status", h.UpdateReservationStatus)
		reservations.POST("/:id/cancel", h.CancelReservation)
	}

	// ---- Rentals (+ returns) ----
	rentals := api.Group("/rentals", auth.RequireAuth(), auth.RequirePermission("rentals"))
	{
		rentals.GET("", h.ListRentals)
		rentals.GET("/:id", h.GetRental)
		rentals.POST("", h.CreateRental)
		rentals.POST("/:id/cancel", h.CancelRental)
		rentals.GET("/:id/contract", h.DownloadContract)
		rentals.POST("/:id/return/preview", auth.RequirePermission("returns"), h.PreviewReturn)
		rentals.POST("/:id/return", auth.RequirePermission("returns"), h.CreateReturn)
	}

	// ---- Payments ----
	payments := api.Group("/payments", auth.RequireAuth(), auth.RequirePermission("payments"))
	{
		payments.GET("", h.ListPayments)
		payments.GET("/:id", h.GetPayment)
		payments.POST("", h.CreatePayment)
		payments.GET("/:id/receipt", h.DownloadReceipt)
	}

	// ---- Blacklist ----
	blacklist := api.Group("/blacklist", auth.RequireAuth(), auth.RequirePermission("blacklist"))
	{
		blacklist.GET("", h.ListBlacklist)
		blacklist.GET("/reasons", h.BlacklistReasons)
		blacklist.POST("", h.AddToBlacklist)
		blacklist.DELETE("/:id", h.RemoveFromBlacklist)
	}

	// ---- Accidents ----
	accidents := api.Group("/accidents", auth.RequireAuth(), auth.RequirePermission("accidents"))
	{
		accidents.GET("", h.ListAccidents)
		accidents.GET("/:id", h.GetAccident)
		accidents.POST("", h.CreateAccident)
		accidents.PUT("/:id", h.UpdateAccident)
		accidents.DELETE("/:id", h.DeleteAccident)
	}

	// ---- Dashboard ----
	dashboard := api.Group("/dashboard", auth.RequireAuth(), auth.RequirePermission("dashboard"))
	{
		dashboard.GET("/summary", h.DashboardSummary)
		dashboard.GET("/revenue-by-month", h.RevenueByMonth)
		dashboard.GET("/top-cars", h.TopCars)
		dashboard.GET("/rental-statistics", h.RentalStatistics)
	}

	// ---- Reports ----
	reports := api.Group("/reports", auth.RequireAuth(), auth.RequirePermission("reports"))
	{
		reports.GET("/daily-revenue", h.DailyRevenue)
		reports.GET("/monthly-revenue", h.MonthlyRevenue)
		reports.GET("/yearly-revenue", h.YearlyRevenue)
		reports.GET("/profitable-cars", h.ProfitableCars)
		reports.GET("/active-rentals", h.ActiveRentalsReport)
		reports.GET("/debtors", h.DebtorsReport)
		reports.GET("/client-statistics", h.ClientStatisticsReport)
		reports.GET("/export/:report_type", h.ExportReport)
	}

	// ---- Settings ----
	settings := api.Group("/settings", auth.RequireAuth())
	{
		settings.GET("", h.GetSettings)
		settings.PUT("", auth.RequireRole(admin), h.UpdateSettings)
	}

	// ---- Backups ----
	backups := api.Group("/backups", auth.RequireAuth(), auth.RequirePermission("backups"))
	{
		backups.GET("", h.ListBackups)
		backups.POST("", h.CreateBackup)
	}

	// ---- Project import/export ----
	project := api.Group("/project", auth.RequireAuth(), auth.RequireRole(admin))
	{
		project.GET("/export", h.ExportProject)
		project.POST("/import", h.ImportProject)
	}

	// ---- Uploads ----
	uploads := api.Group("/uploads")
	{
		uploads.POST("", auth.RequireAuth(), h.UploadFile)
		uploads.GET("/*path", h.ServeFile)
	}

	// ---- Audit ----
	api.GET("/audit", auth.RequireAuth(), auth.RequireRole(admin), h.ListAuditLogs)

	return r
}
