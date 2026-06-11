package media

import _ "embed"

// Embedded Cyrillic-capable TrueType fonts used for PDF generation so the
// binary is fully self-contained and renders Russian text correctly.

//go:embed fonts/AppFont-Regular.ttf
var fontRegular []byte

//go:embed fonts/AppFont-Bold.ttf
var fontBold []byte
