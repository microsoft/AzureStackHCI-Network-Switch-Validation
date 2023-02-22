package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"
	"gopkg.in/yaml.v3"
)

func (o *OutputType) resultAnalysis(pcapFilePath string, i *InputType) {
	o.decodePacketLayer(pcapFilePath)
	o.FeatureSummary = []FeatureResult{}
	o.BGPResultValidation(&o.BGPResult)
	o.VLANResultValidation(&o.VLANResult, i)
	o.DHCPResultValidation(&o.DHCPResult)
	o.LLDPResultValidation(&o.LLDPResult, i)
	o.RoleTypeResult()
}

func (o *OutputType) RoleTypeResult() {

	o.RoleSummary = map[string]string{
		MANAGEMENT:   PASS,
		COMPUTEBASIC: PASS,
		COMPUTESDN:   PASS,
		STORAGE:      PASS,
	}

	for _, v := range o.FeatureSummary {
		for _, role := range v.FeatureRoles {
			if v.FeaturePass == FAIL {
				if o.RoleSummary[role] != FAIL {
					o.RoleSummary[role] = FAIL
				}
			}
		}
	}
}

func (o *OutputType) outputPDFFile(pdfFilePath string) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	reportDate := o.TestDate.Format("2006-01-02 15:04:05")
	titleName := fmt.Sprintf("%s - Validation Report\n", reportDate)
	pdf.Cell(40, 10, titleName)
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(100, 10, ROLE_SUMMARY_TITTLE)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 16)
	for key, value := range o.RoleSummary {
		pdf.SetX(20)
		pdf.SetFont("Arial", "B", 14)
		if value == FAIL {
			pdf.SetTextColor(255, 0, 0)
		} else {
			pdf.SetTextColor(0, 255, 0)
		}
		roleTitle := fmt.Sprintf("%s - %s", key, value)
		pdf.Cell(40, 10, roleTitle)
		pdf.Ln(10)
	}
	pdf.Ln(10)

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(100, 10, FEATURE_SUMMARY_TITTLE)
	pdf.Ln(10)
	// Feature Summary List
	for _, featureObj := range o.FeatureSummary {
		titleWidth, contentWidth := 20.0, 150.0
		// Set the font to Arial regular, size 12
		pdf.SetFont("Arial", "", 8)

		if featureObj.FeaturePass == "Fail" {
			pdf.SetTextColor(255, 0, 0)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
		pdf.CellFormat(titleWidth, 7, "Feature", "1", 0, "", false, 0, "")
		pdf.CellFormat(contentWidth, 7, featureObj.FeatureName, "1", 0, "", false, 0, "")
		pdf.Ln(7)
		pdf.CellFormat(titleWidth, 7, "Result", "1", 0, "", false, 0, "")
		pdf.CellFormat(contentWidth, 7, featureObj.FeaturePass, "1", 0, "", false, 0, "")
		pdf.Ln(7)
		pdf.CellFormat(titleWidth, 7, "Log", "1", 0, "", false, 0, "")
		pdf.CellFormat(contentWidth, 7, featureObj.FeatureLog, "1", 0, "", false, 0, "")
		pdf.Ln(7)
		pdf.CellFormat(titleWidth, 7, "RoleType", "1", 0, "", false, 0, "")
		pdf.CellFormat(contentWidth, 7, strings.Join(featureObj.FeatureRoles, ", "), "1", 0, "", false, 0, "")
		// Line break
		pdf.Ln(10)
	}

	// Logs
	pdf.AddPage()
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(100, 10, ALL_LOGS)
	pdf.Ln(10)
	resultDetailBytes, err := yaml.Marshal(o)
	if err != nil {
		log.Fatalln("YAML marshal failed, err:", err)
	}
	pdf.SetFont("Arial", "", 8)
	pdf.MultiCell(0, 5, string(resultDetailBytes), "", "", false)

	// Input
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(100, 10, INPUT_VARIABLES)
	pdf.Ln(10)
	iniBytes, err := yaml.Marshal(inputObj)
	if err != nil {
		log.Fatalln("YAML marshal failed, err:", err)
	}
	pdf.SetFont("Arial", "", 8)
	pdf.MultiCell(0, 5, string(iniBytes), "", "", false)

	err = pdf.OutputFileAndClose(pdfFilePath)
	if err != nil {
		log.Fatalln("pdf creation failed: ", err)
	}
}

func (o *OutputType) outputYAMLFile(yamlFilePath string) {
	out, err := yaml.Marshal(o)
	if err != nil {
		log.Fatalln(err)
		return
	}

	f, err := os.Create(yamlFilePath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer f.Close()

	_, err = f.Write(out)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func (o *OutputType) outputJSONFile(jsonFilePath string) {
	out, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Fatalln(err)
		return
	}

	f, err := os.Create(jsonFilePath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer f.Close()

	_, err = f.Write(out)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
