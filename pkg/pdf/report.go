package pdf

import (
	"bytes"

	"github.com/jung-kurt/gofpdf"
	"github.com/whiterage/webserver_go/pkg/models"
)

func BuildReport(tasks []*models.Task) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitMillimeter, gofpdf.PageSizeA4, "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Hello, World!")
	return buffer.Bytes(), nil
}
