package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// ListAccidents returns accidents filtered by car.
func (h *Handler) ListAccidents(c *gin.Context) {
	page, perPage, offset := pagination(c)
	result, err := h.Svc.ListAccidents(queryUint(c, "car_id"), page, perPage, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAccident returns a single accident.
func (h *Handler) GetAccident(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	accident, err := h.Svc.GetAccident(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, accident.ToMap())
}

// CreateAccident records an accident.
func (h *Handler) CreateAccident(c *gin.Context) {
	var in service.AccidentInput
	_ = c.ShouldBindJSON(&in)
	accident, err := h.Svc.CreateAccident(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, accident.ToMap())
}

// UpdateAccident updates an accident.
func (h *Handler) UpdateAccident(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	var in service.AccidentInput
	_ = c.ShouldBindJSON(&in)
	accident, err := h.Svc.UpdateAccident(id, in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, accident.ToMap())
}

// DeleteAccident deletes an accident.
func (h *Handler) DeleteAccident(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	if err := h.Svc.DeleteAccident(id, actorFrom(c)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Запись о ДТП удалена"})
}
