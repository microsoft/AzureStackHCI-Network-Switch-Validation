package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-pdf/fpdf"
	"gopkg.in/yaml.v3"
)

func (o *OutputType) resultAnalysis(pcapFilePath string, i *InputType) {
	o.decodePacketLayer(pcapFilePath)
	o.TypeReportSummary = []TypeResult{}
	o.BGPResultValidation(&o.BGPResult)
	o.VLANResultValidation(&o.VLANResult, i)
	o.DHCPResultValidation(&o.DHCPResult)
	o.LLDPResultValidation(&o.LLDPResult, i)
	o.RoleTypeResult()
}

func (o *OutputType) RoleTypeResult() {

	o.RoleReportSummary = map[string]string{
		MANAGEMENT:   PASS,
		COMPUTEBASIC: PASS,
		COMPUTESDN:   PASS,
		STORAGE:      PASS,
	}

	for _, v := range o.TypeReportSummary {
		for _, role := range v.TypeRoles {
			if v.TypePass == FAIL {
				if o.RoleReportSummary[role] != FAIL {
					o.RoleReportSummary[role] = FAIL
				}
			}
		}
	}
}

// func (o *OutputType) outputPDFbyFile(pdfFilePath string) {

// 	var resultTemplate = `
// 	{{range $key, $value := .}}
// 	{{$key -}}
// 		{{range $value}}
// 		- {{. -}}
// 		{{end}}
// 	{{end}}
// 	`
// 	ret := OutputObj.ResultSummary

// 	t := template.New("resultTemplate")
// 	t, err := t.Parse(resultTemplate)
// 	if err != nil {
// 		log.Fatalln("parse file: ", err)
// 		return
// 	}

// 	fmt.Print(CONSOLE_OUTPUT)
// 	err = t.Execute(os.Stdout, ret)
// 	if err != nil {
// 		log.Fatalln("execute: ", err)
// 		return
// 	}

// 	var retSummary bytes.Buffer
// 	err = t.Execute(&retSummary, ret)
// 	if err != nil {
// 		log.Fatalln("execute: ", err)
// 		return
// 	}

// 	pdf := fpdf.New("P", "mm", "A4", "")
// 	pdf.AddPage()
// 	pdf.SetFont("Arial", "B", 16)

// 	var titleResult string
// 	if OutputObj.resultPass() {
// 		pdf.SetTextColor(0, 220, 0)
// 		titleResult = "PASS"
// 	} else {
// 		pdf.SetTextColor(220, 0, 0)
// 		titleResult = "FAIL"
// 	}

// 	reportDate := time.Now().Format("2006-01-02 15:04:05")
// 	titleName := fmt.Sprintf("%s - Validation Report\n", reportDate)

// 	pdf.Cell(0, 10, titleName)
// 	pdf.Ln(10.0)

// 	pdf.SetFont("Arial", "", 14)
// 	pdf.SetTextColor(0, 0, 0)
// 	pdf.MultiCell(0, 8, retSummary.String(), "", "", false)

// 	pdf.AddPage()
// 	resultDetaillBytes, err := yaml.Marshal(OutputObj)
// 	if err != nil {
// 		log.Fatalln("YAML marshal failed, err:", err)
// 	}

// 	pdf.SetFont("Arial", "", 10)
// 	pdf.MultiCell(0, 5, string(resultDetaillBytes), "", "", false)
// 	err = pdf.OutputFileAndClose(pdfFilePath)
// 	if err != nil {
// 		log.Fatalln("pdf creation failed: ", err)
// 	}

// 	pdf.AddPage()
// 	iniBytes, err := yaml.Marshal(inputObj)
// 	if err != nil {
// 		log.Fatalln("YAML marshal failed, err:", err)
// 	}

// 	pdf.SetFont("Arial", "", 10)
// 	pdf.MultiCell(0, 5, string(iniBytes), "", "", false)
// 	err = pdf.OutputFileAndClose(pdfFilePath)
// 	if err != nil {
// 		log.Fatalln("pdf creation failed: ", err)
// 	}

// 	fmt.Println(GENERATE_PDF_OUTPUT)
// }

func (o *OutputType) outputPDFFile(pdfFilePath string) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	reportDate := o.TestDate.Format("2006-01-02 15:04:05")
	titleName := fmt.Sprintf("%s - Validation Report\n", reportDate)
	pdf.Cell(40, 10, titleName)
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(100, 10, REPORT_SUMMARY_TITTLE)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 16)
	for key, value := range o.RoleReportSummary {
		pdf.Cell(100, 10, key)
		if value == PASS {
			pdf.SetTextColor(0, 128, 0)
		} else {
			pdf.SetTextColor(255, 0, 0)
		}
		pdf.Cell(100, 10, value)
		pdf.Ln(10)
		pdf.SetTextColor(0, 0, 0)
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(100, 10, TYPE_SUMMARY_TITTLE)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	typeReportBytes, err := yaml.Marshal(o.TypeReportSummary)
	if err != nil {
		log.Fatalln("YAML marshal failed, err:", err)
	}
	pdf.MultiCell(0, 5, string(typeReportBytes), "", " ", false)

	// Logs
	pdf.AddPage()
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
	fmt.Println(GENERATE_PDF_OUTPUT)
}

func (o *OutputType) outputYAMLFile(yamlFilePath string) {
	out, err := yaml.Marshal(o)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(yamlFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = f.Write(out)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (o *OutputType) outputJSONFile(jsonFilePath string) {
	out, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(jsonFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = f.Write(out)
	if err != nil {
		fmt.Println(err)
		return
	}
}
