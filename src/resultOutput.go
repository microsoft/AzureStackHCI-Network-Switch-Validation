package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"gopkg.in/yaml.v3"
)

func (o *OutputType) resultAnalysis(pcapFilePath string, i *INIType) {
	o.decodePacketLayer(pcapFilePath)
	o.ResultSummary = map[string][]string{}
	o.VLANResultValidation(&o.VLANResult, i)
	o.DHCPResultValidation(&o.DHCPResult)
	o.LLDPResultValidation(&o.LLDPResult, i)
	o.BGPResultValidation(&o.BGPResult)
}

func (o *OutputType) outputPDFbyFile() {

	var resultTemplate = `
	{{range $key, $value := .}}
	{{$key -}}
		{{range $value}}
		- {{. -}}
		{{end}}
	{{end}}
	`
	ret := OutputObj.ResultSummary

	t := template.New("resultTemplate")
	t, err := t.Parse(resultTemplate)
	if err != nil {
		log.Fatalln("parse file: ", err)
		return
	}

	fmt.Print(CONSOLE_OUTPUT)
	err = t.Execute(os.Stdout, ret)
	if err != nil {
		log.Fatalln("execute: ", err)
		return
	}

	var retSummary bytes.Buffer
	err = t.Execute(&retSummary, ret)
	if err != nil {
		log.Fatalln("execute: ", err)
		return
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	var titleResult string
	if OutputObj.resultPass() {
		pdf.SetTextColor(0, 220, 0)
		titleResult = "PASS"
	} else {
		pdf.SetTextColor(220, 0, 0)
		titleResult = "FAIL"
	}

	OutputObj.resultPass()
	reportDate := time.Now().Format("2006-01-02 15:04:05")
	titleName := fmt.Sprintf("%s - Validation Report:  %s\n", reportDate, titleResult)

	pdf.Cell(0, 10, titleName)
	pdf.Ln(10.0)

	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.MultiCell(0, 8, retSummary.String(), "", "", false)
	err = pdf.OutputFileAndClose("result.pdf")
	if err != nil {
		log.Fatalln("pdf creation failed: ", err)
	}

	pdf.AddPage()
	yamlBytes, err := yaml.Marshal(OutputObj)
	if err != nil {
		log.Fatalln("YAML marshal failed, err:", err)
	}

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, string(yamlBytes), "", "", false)
	err = pdf.OutputFileAndClose("result.pdf")
	if err != nil {
		log.Fatalln("pdf creation failed: ", err)
	}
	fmt.Println(GENERATE_PDF_OUTPUT)
}

func (o *OutputType) resultPass() bool {
	for k := range o.ResultSummary {
		if strings.Contains(k, "FAIL") {
			return false
		}
	}
	return true
}
