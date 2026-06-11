package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DailyRevenue returns the daily revenue report.
func (h *Handler) DailyRevenue(c *gin.Context) {
	var day *time.Time
	if q := c.Query("date"); q != "" {
		if t, err := time.Parse("2006-01-02", q); err == nil {
			day = &t
		}
	}
	c.JSON(http.StatusOK, h.Svc.DailyRevenue(day))
}

// MonthlyRevenue returns the monthly revenue report.
func (h *Handler) MonthlyRevenue(c *gin.Context) {
	now := time.Now()
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	c.JSON(http.StatusOK, h.Svc.MonthlyRevenue(year, month))
}

// YearlyRevenue returns the yearly revenue report.
func (h *Handler) YearlyRevenue(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	c.JSON(http.StatusOK, h.Svc.YearlyRevenue(year))
}

// ProfitableCars returns the most profitable cars.
func (h *Handler) ProfitableCars(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	c.JSON(http.StatusOK, h.Svc.MostProfitableCars(limit))
}

// ActiveRentalsReport returns all active rentals.
func (h *Handler) ActiveRentalsReport(c *gin.Context) {
	c.JSON(http.StatusOK, h.Svc.ActiveRentals())
}

// DebtorsReport returns the list of debtors.
func (h *Handler) DebtorsReport(c *gin.Context) {
	c.JSON(http.StatusOK, h.Svc.Debtors())
}

// ClientStatisticsReport returns rental statistics per client.
func (h *Handler) ClientStatisticsReport(c *gin.Context) {
	c.JSON(http.StatusOK, h.Svc.ClientStatistics())
}

// ExportReport generates and downloads a styled Excel report.
func (h *Handler) ExportReport(c *gin.Context) {
	reportType := c.Param("report_type")
	params := map[string]string{}
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	path, err := h.Svc.ExportReportExcel(reportType, params)
	if err != nil {
		respondError(c, err)
		return
	}
	c.FileAttachment(path, filepath.Base(path))
}
