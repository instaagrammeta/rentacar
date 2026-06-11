package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// ListClients returns a paginated, filtered client list.
func (h *Handler) ListClients(c *gin.Context) {
	page, perPage, offset := pagination(c)
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))
	result, err := h.Svc.ListClients(search, status, page, perPage, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetClient returns a single client.
func (h *Handler) GetClient(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	client, err := h.Svc.GetClient(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, client.ToMap())
}

// ClientHistory returns the aggregated history of a client.
func (h *Handler) ClientHistory(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	history, err := h.Svc.GetClientHistory(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, history)
}

// SearchClient looks up a returning client by phone, code or free text.
func (h *Handler) SearchClient(c *gin.Context) {
	result, err := h.Svc.FindClient(c.Query("phone"), c.Query("code"), c.Query("term"))
	if err != nil {
		respondError(c, err)
		return
	}
	switch v := result.(type) {
	case []models.Client:
		out := make([]map[string]interface{}, 0, len(v))
		for i := range v {
			out = append(out, v[i].ToMap())
		}
		c.JSON(http.StatusOK, out)
	case *models.Client:
		c.JSON(http.StatusOK, v.ToMap())
	default:
		c.JSON(http.StatusOK, gin.H{})
	}
}

// CreateClient creates a new client.
func (h *Handler) CreateClient(c *gin.Context) {
	var in service.ClientInput
	_ = c.ShouldBindJSON(&in)
	client, err := h.Svc.CreateClient(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, client.ToMap())
}

// UpdateClient updates a client.
func (h *Handler) UpdateClient(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	var in service.ClientInput
	_ = c.ShouldBindJSON(&in)
	client, err := h.Svc.UpdateClient(id, in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, client.ToMap())
}

// DeleteClient deletes a client.
func (h *Handler) DeleteClient(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	if err := h.Svc.DeleteClient(id, actorFrom(c)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Клиент удалён"})
}

// SendClientSMS sends a custom SMS message to a client.
func (h *Handler) SendClientSMS(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	var body struct {
		Message string `json:"message"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.Svc.SendClientSMS(id, body.Message, actorFrom(c)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "SMS отправлено"})
}
