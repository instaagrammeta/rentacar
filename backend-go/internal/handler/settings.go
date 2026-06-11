package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// GetSettings returns the company settings (any authenticated user).
func (h *Handler) GetSettings(c *gin.Context) {
	settings, err := h.Svc.GetSettings()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings.ToMap())
}

// UpdateSettings updates the company settings (admin only).
func (h *Handler) UpdateSettings(c *gin.Context) {
	var in service.SettingsInput
	_ = c.ShouldBindJSON(&in)
	settings, err := h.Svc.UpdateSettings(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings.ToMap())
}
