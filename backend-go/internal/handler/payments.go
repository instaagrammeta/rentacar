package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// ListPayments returns payments filtered by client and rental.
func (h *Handler) ListPayments(c *gin.Context) {
	page, perPage, offset := pagination(c)
	result, err := h.Svc.ListPayments(queryUint(c, "client_id"), queryUint(c, "rental_id"), page, perPage, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetPayment returns a single payment.
func (h *Handler) GetPayment(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	payment, err := h.Svc.GetPayment(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, payment.ToMap())
}

// CreatePayment records a payment.
func (h *Handler) CreatePayment(c *gin.Context) {
	var in service.PaymentInput
	_ = c.ShouldBindJSON(&in)
	payment, err := h.Svc.CreatePayment(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, payment.ToMap())
}

// DownloadReceipt serves the payment receipt PDF.
func (h *Handler) DownloadReceipt(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	payment, err := h.Svc.GetPayment(id)
	if err != nil {
		respondError(c, err)
		return
	}
	if payment.ReceiptPath == nil || *payment.ReceiptPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Квитанция не найдена"})
		return
	}
	full := h.Svc.Media.AbsoluteExportPath(*payment.ReceiptPath)
	c.FileAttachment(full, payment.ReceiptNumber+".pdf")
}
