package media

import (
	"image/color"
	"os"
	"path/filepath"

	qrcode "github.com/skip2/go-qrcode"
)

// brandGreen is the QR foreground colour (#1B5E20) used across the app.
var brandGreen = color.RGBA{R: 0x1B, G: 0x5E, B: 0x20, A: 0xFF}

// GenerateQR creates a QR-code PNG for the given data and returns its path
// relative to the upload directory (e.g. "qr/CL-000001.png").
func (g *Generator) GenerateQR(data, filename string) (string, error) {
	dir := filepath.Join(g.UploadDir, "qr")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	q, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}
	q.ForegroundColor = brandGreen
	q.BackgroundColor = color.White

	if err := q.WriteFile(256, filepath.Join(dir, filename)); err != nil {
		return "", err
	}
	return "qr/" + filename, nil
}
