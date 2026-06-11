package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadFile stores one or more uploaded files and returns their relative paths.
func (h *Handler) UploadFile(c *gin.Context) {
	subfolder := c.PostForm("subfolder")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Файл не передан"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		files = form.File["file"]
	}
	if len(files) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Файл не передан"})
		return
	}

	stored := make([]string, 0, len(files))
	for _, fh := range files {
		if fh == nil || fh.Filename == "" {
			continue
		}
		rel, err := h.Svc.Media.SaveUpload(fh, subfolder)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		stored = append(stored, rel)
	}
	c.JSON(http.StatusCreated, gin.H{"paths": stored})
}

// ServeFile serves an uploaded media file (document scan, photo, QR, logo).
func (h *Handler) ServeFile(c *gin.Context) {
	rel := c.Param("path")
	full := h.Svc.Media.AbsoluteUploadPath(rel)
	// Ensure the resolved path stays within the upload directory.
	if !isWithin(h.Cfg.UploadDir, full) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Файл не найден"})
		return
	}
	c.File(full)
}

func isWithin(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != ".." && !startsWithDotDot(rel)
}

func startsWithDotDot(p string) bool {
	return len(p) >= 2 && p[0] == '.' && p[1] == '.'
}
