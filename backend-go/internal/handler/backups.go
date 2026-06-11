package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListBackups returns the available backups.
func (h *Handler) ListBackups(c *gin.Context) {
	c.JSON(http.StatusOK, h.Svc.ListBackups())
}

// CreateBackup triggers a manual backup.
func (h *Handler) CreateBackup(c *gin.Context) {
	path, err := h.Svc.CreateBackup()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Резервная копия создана", "path": path})
}
