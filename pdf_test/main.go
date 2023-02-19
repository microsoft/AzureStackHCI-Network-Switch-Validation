package main

import (
	"github.com/go-pdf/fpdf"
)

func main() {
	// Define the data for each section
	data := map[string][2]string{
		"Storage": {"VLAN", "Fail"},
		"Compute": {"BGP", "Pass"},
	}

	// Create a new PDF document
	pdf := fpdf.New("P", "mm", "A4", "")

	addSection(pdf, data)

	// Save the PDF to a file
	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		panic(err)
	}
}

func addSection(pdf *fpdf.Fpdf, data map[string][2]string) {
	// Add a new page to the PDF document
	pdf.AddPage()
	// Set the font to Arial bold, size 16
	for key, value := range data {
		// Add the title to the page
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, key)
		pdf.Ln(10)

		// Set the font to Arial regular, size 12
		pdf.SetFont("Arial", "", 12)

		if value[1] == "Fail" {
			pdf.SetTextColor(255, 0, 0)
		} else if value[1] == "Pass" {
			pdf.SetTextColor(0, 255, 0)
		}

		// Add the cells to the table
		for _, cell := range value {
			pdf.CellFormat(20.0, 7, cell, "1", 0, "", false, 0, "")
		}
		pdf.CellFormat(60.0, 7, "", "1", 1, "", false, 0, "")

		pdf.Ln(10)
	}
}
