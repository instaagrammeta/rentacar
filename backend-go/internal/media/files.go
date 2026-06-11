// Package media handles file uploads, QR codes, PDF documents and Excel
// exports for the backend.
package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Generator produces media artefacts (uploads, QR, PDF, Excel) into the
// configured upload and export directories.
type Generator struct {
	UploadDir string
	ExportDir string
}

// NewGenerator builds a media generator bound to the given directories.
func NewGenerator(uploadDir, exportDir string) *Generator {
	return &Generator{UploadDir: uploadDir, ExportDir: exportDir}
}

var allowedUploadExt = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".webp": true, ".bmp": true, ".pdf": true,
}

// SaveUpload persists an uploaded file and returns its path relative to the
// upload directory. A random prefix avoids collisions while keeping the
// original extension.
func (g *Generator) SaveUpload(fileHeader *multipart.FileHeader, subfolder string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != "" && !allowedUploadExt[ext] {
		return "", fmt.Errorf("недопустимый тип файла: %s", ext)
	}

	subfolder = sanitizeSubfolder(subfolder)
	targetDir := g.UploadDir
	if subfolder != "" {
		targetDir = filepath.Join(g.UploadDir, subfolder)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}

	filename := uuid.NewString() + ext
	dst := filepath.Join(targetDir, filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}

	if subfolder == "" {
		return filename, nil
	}
	return subfolder + "/" + filename, nil
}

// AbsoluteUploadPath resolves a stored relative path to an absolute filesystem
// path inside the upload directory.
func (g *Generator) AbsoluteUploadPath(relative string) string {
	return filepath.Join(g.UploadDir, filepath.Clean("/"+relative))
}

// AbsoluteExportPath resolves a stored relative path inside the export dir.
func (g *Generator) AbsoluteExportPath(relative string) string {
	return filepath.Join(g.ExportDir, filepath.Clean("/"+relative))
}

func sanitizeSubfolder(sub string) string {
	sub = strings.TrimSpace(sub)
	sub = strings.ReplaceAll(sub, "..", "")
	sub = strings.Trim(sub, "/\\")
	return sub
}
