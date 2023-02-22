package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-pdf/fpdf"
)

type RoleReportType struct {
	RoleName         string                `json:"RoleName"`
	RoleResult       string                `json:"RoleResult"`
	RoleTestFeatures []RoleTestFeatureType `json:"RoleTestFeatures"`
}

type RoleTestFeatureType struct {
	FeatureName string `json:"FeatureName"`
	FeaturePass string `json:"FeaturePass"`
	FeatureLog  string `json:"FeatureLog"`
}

// Load input.json to RoleReportType
func loadInputJSON(inputJsonFile string) []RoleReportType {
	data, err := ioutil.ReadFile(inputJsonFile)
	if err != nil {
		fmt.Println("Read input.json error:", err)
	}
	var roleReport []RoleReportType
	err = json.Unmarshal(data, &roleReport)
	if err != nil {
		fmt.Println("Unmarshal input.json error:", err)
	}
	return roleReport
}

func main() {
	// Define the data for each section
	roleReportObj := loadInputJSON("input.json")
	// Create a new PDF document
	pdf := fpdf.New("P", "mm", "A4", "")

	// addRolebySection(pdf, roleReportObj)
	addRolebySummary(pdf, roleReportObj)

	// Save the PDF to a file
	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		panic(err)
	}
}

func addRolebySection(pdf *fpdf.Fpdf, roleReportObj []RoleReportType) {
	// Add a new page to the PDF document
	pdf.AddPage()
	// Set the font to Arial bold, size 16
	for _, roleObj := range roleReportObj {
		// Add the title to the page
		pdf.SetFont("Arial", "B", 16)
		if roleObj.RoleResult == "Fail" {
			pdf.SetTextColor(255, 0, 0)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
		roleTitle := fmt.Sprintf("%s - %s", roleObj.RoleName, roleObj.RoleResult)
		pdf.Cell(40, 10, roleTitle)
		pdf.Ln(10)

		// Set the font to Arial regular, size 12
		pdf.SetFont("Arial", "", 10)

		for _, featureObj := range roleObj.RoleTestFeatures {
			if featureObj.FeaturePass == "Fail" {
				pdf.SetTextColor(255, 0, 0)
			} else {
				pdf.SetTextColor(0, 0, 0)
			}
			pdf.CellFormat(20.0, 7, "Feature", "1", 0, "", false, 0, "")
			pdf.CellFormat(100.0, 7, featureObj.FeatureName, "1", 0, "", false, 0, "")
			pdf.Ln(7)
			pdf.CellFormat(20.0, 7, "Result", "1", 0, "", false, 0, "")
			pdf.CellFormat(100.0, 7, featureObj.FeaturePass, "1", 0, "", false, 0, "")
			pdf.Ln(7)
			pdf.CellFormat(20.0, 7, "Log", "1", 0, "", false, 0, "")
			pdf.CellFormat(100.0, 7, featureObj.FeatureLog, "1", 0, "", false, 0, "")
			pdf.Ln(10)
		}

		pdf.Ln(10)
	}
}

func addRolebySummary(pdf *fpdf.Fpdf, roleReportObj []RoleReportType) {
	// Add a new page to the PDF document
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(40, 10, "2022-08-31 11:02:00 Validation Report")
	pdf.Ln(20)
	// Set the font to Arial bold, size 16
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Supported Role-Types:")
	pdf.Ln(10)

	// Set the font to Arial bold, size 16
	for _, roleObj := range roleReportObj {
		// Indent the role name
		pdf.SetX(20)
		// Add the title to the page
		pdf.SetFont("Arial", "B", 14)
		if roleObj.RoleResult == "Fail" {
			pdf.SetTextColor(255, 0, 0)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
		roleTitle := fmt.Sprintf("%s - %s", roleObj.RoleName, roleObj.RoleResult)
		pdf.Cell(40, 10, roleTitle)
		pdf.Ln(10)
	}
	pdf.Ln(10)

	for _, roleObj := range roleReportObj {
		// Set the font to Arial regular, size 12
		pdf.SetFont("Arial", "", 8)

		for _, featureObj := range roleObj.RoleTestFeatures {
			if featureObj.FeaturePass == "Fail" {
				pdf.SetTextColor(255, 0, 0)
			} else {
				pdf.SetTextColor(0, 0, 0)
			}
			pdf.CellFormat(20.0, 7, "Feature", "1", 0, "", false, 0, "")
			pdf.CellFormat(140.0, 7, featureObj.FeatureName, "1", 0, "", false, 0, "")
			pdf.Ln(7)
			pdf.CellFormat(20.0, 7, "Result", "1", 0, "", false, 0, "")
			pdf.CellFormat(140.0, 7, featureObj.FeaturePass, "1", 0, "", false, 0, "")
			pdf.Ln(7)
			pdf.CellFormat(20.0, 7, "Log", "1", 0, "", false, 0, "")
			pdf.CellFormat(140.0, 7, featureObj.FeatureLog, "1", 0, "", false, 0, "")
			pdf.Ln(10)
		}
	}
}
