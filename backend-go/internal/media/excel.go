package media

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelReport describes the data to render into a styled worksheet.
type ExcelReport struct {
	SheetName   string
	Title       string
	CompanyName string
	Headers     []string
	Rows        [][]interface{}
}

// ExportReportExcel writes a styled Excel report (green headers, auto-sized
// columns, frozen header) and returns the absolute path of the saved file.
func (g *Generator) ExportReportExcel(rep ExcelReport) (string, error) {
	dir := filepath.Join(g.ExportDir, "reports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheet := rep.SheetName
	if sheet == "" {
		sheet = "Отчёт"
	}
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return "", err
	}

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16, Color: "1B5E20"},
	})
	subStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Italic: true, Size: 9, Color: "666666"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1B5E20"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "C8E6C9", Style: 1},
			{Type: "right", Color: "C8E6C9", Style: 1},
			{Type: "top", Color: "C8E6C9", Style: 1},
			{Type: "bottom", Color: "C8E6C9", Style: 1},
		},
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "C8E6C9", Style: 1},
			{Type: "right", Color: "C8E6C9", Style: 1},
			{Type: "top", Color: "C8E6C9", Style: 1},
			{Type: "bottom", Color: "C8E6C9", Style: 1},
		},
	})

	company := rep.CompanyName
	if company == "" {
		company = "Rentacar CRM"
	}
	_ = f.SetCellValue(sheet, "A1", company)
	_ = f.SetCellStyle(sheet, "A1", "A1", titleStyle)
	_ = f.SetCellValue(sheet, "A2", rep.Title)
	_ = f.SetCellValue(sheet, "A3", "Сформировано: "+time.Now().Format("02.01.2006 15:04"))
	_ = f.SetCellStyle(sheet, "A3", "A3", subStyle)

	headerRow := 5
	for col, h := range rep.Headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, headerRow)
		_ = f.SetCellValue(sheet, cell, h)
		_ = f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for r, row := range rep.Rows {
		for c, value := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, headerRow+1+r)
			_ = f.SetCellValue(sheet, cell, value)
			_ = f.SetCellStyle(sheet, cell, cell, cellStyle)
		}
	}

	// Auto-size columns based on the longest value.
	for col := 0; col < len(rep.Headers); col++ {
		maxLen := len([]rune(rep.Headers[col]))
		for _, row := range rep.Rows {
			if col < len(row) {
				l := len([]rune(fmt.Sprintf("%v", row[col])))
				if l > maxLen {
					maxLen = l
				}
			}
		}
		width := float64(maxLen + 4)
		if width > 50 {
			width = 50
		}
		colName, _ := excelize.ColumnNumberToName(col + 1)
		_ = f.SetColWidth(sheet, colName, colName, width)
	}

	// Freeze the header row.
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      headerRow,
		TopLeftCell: "A6",
		ActivePane:  "bottomLeft",
	})

	name := fmt.Sprintf("%s_%s.xlsx", rep.SheetName, time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(dir, name)
	if err := f.SaveAs(fullPath); err != nil {
		return "", err
	}
	return fullPath, nil
}
