package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// ListReservations returns reservations filtered by status.
func (h *Handler) ListReservations(c *gin.Context) {
	page, perPage, offset := pagination(c)
	status := strings.TrimSpace(c.Query("status"))
	result, err := h.Svc.ListReservations(status, page, perPage, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetReservation returns a single reservation.
func (h *Handler) GetReservation(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	r, err := h.Svc.GetReservation(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, r.ToMap())
}

// CreateReservation books a car for a client.
func (h *Handler) CreateReservation(c *gin.Context) {
	var in service.ReservationInput
	_ = c.ShouldBindJSON(&in)
	r, err := h.Svc.CreateReservation(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, r.ToMap())
}

// UpdateReservationStatus changes a reservation's status.
func (h *Handler) UpdateReservationStatus(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	_ = c.ShouldBindJSON(&body)
	r, err := h.Svc.UpdateReservationStatus(id, body.Status, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, r.ToMap())
}

// CancelReservation cancels a reservation.
func (h *Handler) CancelReservation(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	r, err := h.Svc.CancelReservation(id, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, r.ToMap())
}
