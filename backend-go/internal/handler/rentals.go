package handler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// parseReturnDate parses an ISO datetime, defaulting to now.
func parseReturnDate(value string) time.Time {
	if value == "" {
		return time.Now()
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02T15:04:05", value); err == nil {
		return t
	}
	if len(value) >= 10 {
		if t, err := time.Parse("2006-01-02", value[:10]); err == nil {
			return t
		}
	}
	return time.Now()
}

// ListRentals returns rentals filtered by status.
func (h *Handler) ListRentals(c *gin.Context) {
	page, perPage, offset := pagination(c)
	status := strings.TrimSpace(c.Query("status"))
	result, err := h.Svc.ListRentals(status, page, perPage, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetRental returns a single rental including its return record.
func (h *Handler) GetRental(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	rental, err := h.Svc.GetRental(id)
	if err != nil {
		respondError(c, err)
		return
	}
	data := rental.ToMap()
	if rental.VehicleReturn != nil {
		data["vehicle_return"] = rental.VehicleReturn.ToMap()
	}
	c.JSON(http.StatusOK, data)
}

// CreateRental creates a rental contract.
func (h *Handler) CreateRental(c *gin.Context) {
	var in service.RentalInput
	_ = c.ShouldBindJSON(&in)
	rental, err := h.Svc.CreateRental(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, rental.ToMap())
}

// CancelRental cancels a rental.
func (h *Handler) CancelRental(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	rental, err := h.Svc.CancelRental(id, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, rental.ToMap())
}

// DownloadContract serves the rental contract PDF.
func (h *Handler) DownloadContract(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	rental, err := h.Svc.GetRental(id)
	if err != nil {
		respondError(c, err)
		return
	}
	if rental.PDFPath == nil || *rental.PDFPath == "" {
		rental, err = h.Svc.RegeneratePDF(id)
		if err != nil {
			respondError(c, err)
			return
		}
	}
	if rental.PDFPath == nil || *rental.PDFPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Договор не найден"})
		return
	}
	full := h.Svc.Media.AbsoluteExportPath(*rental.PDFPath)
	c.FileAttachment(full, rental.ContractNumber+".pdf")
}

// frontendBaseURL derives the public frontend base URL (scheme://host[:port])
// from the incoming request, so generated QR links point to whatever address
// the staff member is actually using. Falls back to the configured PublicURL.
func frontendBaseURL(c *gin.Context, fallback string) string {
	if o := strings.TrimSpace(c.GetHeader("Origin")); o != "" {
		return strings.TrimRight(o, "/")
	}
	if ref := strings.TrimSpace(c.GetHeader("Referer")); ref != "" {
		if u, err := url.Parse(ref); err == nil && u.Scheme != "" && u.Host != "" {
			return u.Scheme + "://" + u.Host
		}
	}
	return fallback
}

// RentalQR returns (and regenerates) the public QR code for a rental. The QR
// encodes a link to the public status page the client opens to see how much
// time is left. The link host is taken from the request origin so it always
// points to a reachable address.
func (h *Handler) RentalQR(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	base := frontendBaseURL(c, h.Cfg.PublicURL)
	rental, err := h.Svc.RegenerateRentalQR(id, base)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"qr_code_path": rental.QRCodePath,
		"public_token": rental.PublicToken,
		"public_url":   h.Svc.PublicRentalURLWithBase(rental, base),
	})
}

// PublicRental returns the public, unauthenticated view of a rental looked up
// by its opaque token — used by the page a client reaches via the QR code.
func (h *Handler) PublicRental(c *gin.Context) {
	token := c.Param("token")
	rental, err := h.Svc.GetRentalByToken(token)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, h.Svc.PublicRentalView(rental))
}

// PreviewReturn calculates return charges without persisting.
func (h *Handler) PreviewReturn(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	rental, err := h.Svc.GetRental(id)
	if err != nil {
		respondError(c, err)
		return
	}
	var body struct {
		ReturnDate string  `json:"return_date"`
		DamageCost float64 `json:"damage_cost"`
		Penalties  float64 `json:"penalties"`
	}
	_ = c.ShouldBindJSON(&body)
	rd := parseReturnDate(body.ReturnDate)
	calc := h.Svc.CalculateReturn(rental, rd, body.DamageCost, body.Penalties)
	c.JSON(http.StatusOK, calc.ToMap())
}

// CreateReturn records a vehicle return.
func (h *Handler) CreateReturn(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	var in service.ReturnInput
	_ = c.ShouldBindJSON(&in)
	vr, err := h.Svc.CreateReturn(id, in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, vr.ToMap())
}
