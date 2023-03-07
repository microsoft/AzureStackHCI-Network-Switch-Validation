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
	o.FeatureResultList = []FeatureResultType{}
	o.BGPResultValidation(&o.BGPResult)
	o.VLANResultValidation(&o.VLANResult, i)
	o.DHCPResultValidation(&o.DHCPResult)
	o.LLDPResultValidation(&o.LLDPResult, i)
	o.RoleTypeResult()
}

func (o *OutputType) RoleTypeResult() {
	o.RoleResultList = append(o.RoleResultList, RoleResultType{RoleName: MANAGEMENT, RolePass: PASS, FeaturesByRole: []FeatureResultType{}})
	o.RoleResultList = append(o.RoleResultList, RoleResultType{RoleName: STORAGE, RolePass: PASS, FeaturesByRole: []FeatureResultType{}})
	o.RoleResultList = append(o.RoleResultList, RoleResultType{RoleName: COMPUTEBASIC, RolePass: PASS, FeaturesByRole: []FeatureResultType{}})
	o.RoleResultList = append(o.RoleResultList, RoleResultType{RoleName: COMPUTESDN, RolePass: PASS, FeaturesByRole: []FeatureResultType{}})

	// Create Tmp Map to store role based object
	roleResultMap := map[string]RoleResultType{
		MANAGEMENT:   {RolePass: PASS, FeaturesByRole: []FeatureResultType{}},
		COMPUTEBASIC: {RolePass: PASS, FeaturesByRole: []FeatureResultType{}},
		COMPUTESDN:   {RolePass: PASS, FeaturesByRole: []FeatureResultType{}},
		STORAGE:      {RolePass: PASS, FeaturesByRole: []FeatureResultType{}},
	}
	for _, featureResultObj := range o.FeatureResultList {
		for _, role := range featureResultObj.FeatureRoles {
			tmpMap := roleResultMap[role]
			tmpMap.FeaturesByRole = append(tmpMap.FeaturesByRole, featureResultObj)
			if featureResultObj.FeaturePass == FAIL {
				if tmpMap.RolePass != FAIL {
					tmpMap.RolePass = FAIL
				}
			}
			roleResultMap[role] = tmpMap
		}
	}

	for idx, roleObj := range o.RoleResultList {
		o.RoleResultList[idx].RolePass = roleResultMap[roleObj.RoleName].RolePass
		o.RoleResultList[idx].FeaturesByRole = roleResultMap[roleObj.RoleName].FeaturesByRole
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

	for _, roleObj := range o.RoleResultList {
		// pdf.SetX(20)
		pdf.SetFont("Arial", "B", 14)
		if roleObj.RolePass == FAIL {
			pdf.SetTextColor(255, 0, 0)
		} else {
			pdf.SetTextColor(0, 255, 0)
		}
		roleTitle := fmt.Sprintf("%s - %s", roleObj.RoleName, roleObj.RolePass)
		pdf.Cell(40, 10, roleTitle)
		pdf.Ln(10)
		// Fail Feature Summary based on Role
		for _, featureObj := range roleObj.FeaturesByRole {
			titleWidth, contentWidth := 20.0, 150.0
			pdf.SetFont("Arial", "", 8)
			if featureObj.FeaturePass == FAIL {
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
		}
	}
	pdf.Ln(10)

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(100, 10, FEATURE_SUMMARY_TITTLE)
	pdf.Ln(10)
	// Feature Summary List
	for _, featureObj := range o.FeatureResultList {
		titleWidth, contentWidth := 20.0, 150.0
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
