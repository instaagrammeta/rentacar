package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListAuditLogs returns recent audit log entries (admin only).
func (h *Handler) ListAuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "200"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	logs, err := h.Svc.ListAuditLogs(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, logs)
}
