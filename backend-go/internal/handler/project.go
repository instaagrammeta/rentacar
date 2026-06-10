package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExportProject exports the whole database to a downloadable .rentacar file.
func (h *Handler) ExportProject(c *gin.Context) {
	path, err := h.Svc.ExportProject(actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.FileAttachment(path, filepath.Base(path))
}

// ImportProject restores a .rentacar file, replacing all data.
func (h *Handler) ImportProject(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Файл проекта не передан"})
		return
	}
	if !strings.HasSuffix(fileHeader.Filename, ".rentacar") {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Ожидается файл с расширением .rentacar"})
		return
	}

	tmp, err := os.CreateTemp("", "import_*.rentacar")
	if err != nil {
		respondError(c, err)
		return
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	if err := c.SaveUploadedFile(fileHeader, tmpPath); err != nil {
		respondError(c, err)
		return
	}

	result, err := h.Svc.ImportProject(tmpPath, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	response := gin.H{"message": "Проект импортирован"}
	for k, v := range result {
		response[k] = v
	}
	c.JSON(http.StatusOK, response)
}
