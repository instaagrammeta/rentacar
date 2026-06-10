package media

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

const (
	fontFamily = "App"
)

// newPDF builds a PDF document with the embedded Cyrillic font registered.
func newPDF(title string) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(title, true)
	pdf.AddUTF8FontFromBytes(fontFamily, "", fontRegular)
	pdf.AddUTF8FontFromBytes(fontFamily, "B", fontBold)
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	return pdf
}

func (g *Generator) companyHeader(pdf *fpdf.Fpdf, settings *models.CompanySettings) {
	name := "Rentacar CRM"
	if settings != nil && settings.CompanyName != "" {
		name = settings.CompanyName
	}
	pdf.SetFont(fontFamily, "B", 18)
	pdf.SetTextColor(27, 94, 32)
	pdf.CellFormat(0, 10, name, "", 1, "L", false, 0, "")

	pdf.SetFont(fontFamily, "", 10)
	pdf.SetTextColor(40, 40, 40)
	if settings != nil {
		if settings.Address != nil && *settings.Address != "" {
			pdf.CellFormat(0, 6, *settings.Address, "", 1, "L", false, 0, "")
		}
		if settings.Phone != nil && *settings.Phone != "" {
			pdf.CellFormat(0, 6, "Телефон: "+*settings.Phone, "", 1, "L", false, 0, "")
		}
	}
	pdf.Ln(6)
}

// kvTable renders a two-column label/value table.
func kvTable(pdf *fpdf.Fpdf, rows [][2]string) {
	const labelW, valueW, h = 60.0, 115.0, 8.0
	for _, row := range rows {
		pdf.SetFont(fontFamily, "B", 10)
		pdf.SetFillColor(232, 245, 233)
		pdf.SetTextColor(27, 94, 32)
		pdf.CellFormat(labelW, h, row[0], "1", 0, "L", true, 0, "")

		pdf.SetFont(fontFamily, "", 10)
		pdf.SetTextColor(20, 20, 20)
		pdf.CellFormat(valueW, h, row[1], "1", 1, "L", false, 0, "")
	}
}

func strOr(p *string, fallback string) string {
	if p != nil && *p != "" {
		return *p
	}
	return fallback
}

// GenerateContractPDF renders a rental contract and returns its relative path.
func (g *Generator) GenerateContractPDF(rental *models.Rental, settings *models.CompanySettings) (string, error) {
	dir := filepath.Join(g.ExportDir, "contracts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	relative := "contracts/" + rental.ContractNumber + ".pdf"
	fullPath := filepath.Join(g.ExportDir, filepath.FromSlash(relative))

	pdf := newPDF("Договор " + rental.ContractNumber)
	g.companyHeader(pdf, settings)

	pdf.SetFont(fontFamily, "B", 14)
	pdf.SetTextColor(27, 94, 32)
	pdf.CellFormat(0, 9, "Договор аренды № "+rental.ContractNumber, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	clientName, clientPhone, license := "—", "—", "—"
	if rental.Client != nil {
		clientName = rental.Client.FullName()
		clientPhone = rental.Client.Phone
		license = strOr(rental.Client.DriverLicenseNumber, "—")
	}
	carName, vin := "—", "—"
	if rental.Car != nil {
		carName = rental.Car.DisplayName()
		vin = strOr(rental.Car.VIN, "—")
	}
	employee := "—"
	if rental.Employee != nil {
		employee = rental.Employee.FullName
	}

	kvTable(pdf, [][2]string{
		{"Клиент", clientName},
		{"Телефон", clientPhone},
		{"Вод. удостоверение", license},
		{"Автомобиль", carName},
		{"VIN", vin},
		{"Начало аренды", rental.RentalStart.Format("2006-01-02")},
		{"Окончание аренды", rental.RentalEnd.Format("2006-01-02")},
		{"Цена за сутки", fmt.Sprintf("%.2f", rental.DailyPrice)},
		{"Депозит", fmt.Sprintf("%.2f", rental.Deposit)},
		{"Итоговая стоимость", fmt.Sprintf("%.2f", rental.TotalPrice)},
		{"Ответственный сотрудник", employee},
	})
	pdf.Ln(8)

	if settings != nil && settings.ContractTerms != nil && *settings.ContractTerms != "" {
		pdf.SetFont(fontFamily, "B", 12)
		pdf.SetTextColor(27, 94, 32)
		pdf.CellFormat(0, 8, "Условия договора", "", 1, "L", false, 0, "")
		pdf.SetFont(fontFamily, "", 10)
		pdf.SetTextColor(20, 20, 20)
		pdf.MultiCell(0, 6, *settings.ContractTerms, "", "L", false)
		pdf.Ln(6)
	}

	pdf.Ln(12)
	kvTable(pdf, [][2]string{
		{"Подпись клиента", "______________________"},
		{"Подпись сотрудника", "______________________"},
	})

	if err := pdf.OutputFileAndClose(fullPath); err != nil {
		return "", err
	}
	return relative, nil
}

// GenerateReceiptPDF renders a payment receipt and returns its relative path.
func (g *Generator) GenerateReceiptPDF(payment *models.Payment, settings *models.CompanySettings) (string, error) {
	dir := filepath.Join(g.ExportDir, "receipts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	relative := "receipts/" + payment.ReceiptNumber + ".pdf"
	fullPath := filepath.Join(g.ExportDir, filepath.FromSlash(relative))

	pdf := newPDF("Квитанция " + payment.ReceiptNumber)
	g.companyHeader(pdf, settings)

	pdf.SetFont(fontFamily, "B", 14)
	pdf.SetTextColor(27, 94, 32)
	pdf.CellFormat(0, 9, "Квитанция № "+payment.ReceiptNumber, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	clientName := "—"
	if payment.Client != nil {
		clientName = payment.Client.FullName()
	}

	rows := [][2]string{
		{"Клиент", clientName},
		{"Дата оплаты", payment.PaidAt.Format("02.01.2006 15:04")},
		{"Тип платежа", payment.PaymentType.Label()},
		{"Способ оплаты", payment.Method.Label()},
		{"Сумма", fmt.Sprintf("%.2f", payment.Amount)},
	}
	if payment.Notes != nil && *payment.Notes != "" {
		rows = append(rows, [2]string{"Примечание", *payment.Notes})
	}
	kvTable(pdf, rows)

	if err := pdf.OutputFileAndClose(fullPath); err != nil {
		return "", err
	}
	return relative, nil
}

// TimeStampedName builds a timestamped file name for exports.
func TimeStampedName(prefix, ext string) string {
	return fmt.Sprintf("%s_%s%s", prefix, time.Now().Format("20060102_150405"), ext)
}
