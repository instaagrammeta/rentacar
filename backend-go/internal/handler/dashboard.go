package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DashboardSummary returns the headline metrics.
func (h *Handler) DashboardSummary(c *gin.Context) {
	c.JSON(http.StatusOK, h.Svc.DashboardSummary())
}

// RevenueByMonth returns revenue grouped by month.
func (h *Handler) RevenueByMonth(c *gin.Context) {
	months, _ := strconv.Atoi(c.DefaultQuery("months", "12"))
	c.JSON(http.StatusOK, h.Svc.RevenueByMonth(months))
}

// TopCars returns the most frequently rented cars.
func (h *Handler) TopCars(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	c.JSON(http.StatusOK, h.Svc.TopRentedCars(limit))
}

// RentalStatistics returns rental counts by status.
func (h *Handler) RentalStatistics(c *gin.Context) {
	c.JSON(http.StatusOK, h.Svc.RentalStatistics())
}
