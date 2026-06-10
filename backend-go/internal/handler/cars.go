package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// ListCars returns a paginated, filtered car list.
func (h *Handler) ListCars(c *gin.Context) {
	page, perPage, offset := pagination(c)
	status := strings.TrimSpace(c.Query("status"))
	search := strings.TrimSpace(c.Query("search"))
	result, err := h.Svc.ListCars(status, search, page, perPage, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// AvailableCars returns all available cars.
func (h *Handler) AvailableCars(c *gin.Context) {
	cars, err := h.Svc.AvailableCars()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, cars)
}

// GetCar returns a single car.
func (h *Handler) GetCar(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	car, err := h.Svc.GetCar(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, car.ToMap())
}

// CreateCar creates a vehicle.
func (h *Handler) CreateCar(c *gin.Context) {
	var in service.CarInput
	_ = c.ShouldBindJSON(&in)
	car, err := h.Svc.CreateCar(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, car.ToMap())
}

// UpdateCar updates a vehicle.
func (h *Handler) UpdateCar(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	var in service.CarInput
	_ = c.ShouldBindJSON(&in)
	car, err := h.Svc.UpdateCar(id, in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, car.ToMap())
}

// DeleteCar deletes a vehicle.
func (h *Handler) DeleteCar(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	if err := h.Svc.DeleteCar(id, actorFrom(c)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Автомобиль удалён"})
}
