package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// ListBlacklist returns all blacklist entries.
func (h *Handler) ListBlacklist(c *gin.Context) {
	items, err := h.Svc.ListBlacklist()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// BlacklistReasons lists the available blacklist reasons.
func (h *Handler) BlacklistReasons(c *gin.Context) {
	reasons := make([]gin.H, 0)
	for _, r := range models.AllBlacklistReasons() {
		reasons = append(reasons, gin.H{"value": string(r), "label": r.Label()})
	}
	c.JSON(http.StatusOK, reasons)
}

// AddToBlacklist blacklists a client.
func (h *Handler) AddToBlacklist(c *gin.Context) {
	var in service.BlacklistInput
	_ = c.ShouldBindJSON(&in)
	entry, err := h.Svc.AddToBlacklist(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entry.ToMap())
}

// RemoveFromBlacklist removes a client from the blacklist.
func (h *Handler) RemoveFromBlacklist(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	if err := h.Svc.RemoveFromBlacklist(id, actorFrom(c)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Клиент удалён из чёрного списка"})
}
