package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-pdf/fpdf"
	"gopkg.in/yaml.v3"
)

type TemplateResult struct {
	TasksWithResult map[string][]string
}

func resultAnalysis(pcapFilePath string) {
	decodePacketLayer(pcapFilePath)
	// VLANResultValidation()
	OutputObj.ResultSummary = make(map[string][]string, 100)
	OutputObj.DHCPResultValidation(&OutputObj.DHCPResult)
	OutputObj.LLDPResultValidation(&OutputObj.LLDPResult)
	OutputObj.BGPResultValidation(&OutputObj.BGPResult)

	testDate := time.Now().Format("2006-01-02 15:04:05")
	OutputObj.TimeDate = testDate
	writeToYAML(OutputObj, yamlFilePath)
	outputResultByTemplate()
	outputPDFbyFile()
}

func outputResultByTemplate() {

	var resultTemplate = `
{{range $key, $value := .TasksWithResult}}
{{$key -}}
	{{range $value}}
  	- {{. -}}
	{{end}}
{{end}}
	`

	ret := TemplateResult{
		TasksWithResult: OutputObj.ResultSummary,
	}

	t := template.New("resultTemplate")
	t, err := t.Parse(resultTemplate)
	if err != nil {
		log.Fatalln("parse file: ", err)
		return
	}

	err = t.Execute(os.Stdout, ret)
	if err != nil {
		log.Fatalln("execute: ", err)
		return
	}

	f, err := os.OpenFile(resultSummaryFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln("create file: ", err)
		return
	}
	err = t.Execute(f, ret)
	if err != nil {
		log.Fatalln("execute: ", err)
		return
	}

	f.Close()
}

func outputPDFbyFile() {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	// gap := 10.0
	testDate := time.Now().Format("2006-01-02 15:04:05")
	titleName := testDate + " Validation Result Summary: "
	pdf.Cell(0, 10, titleName)

	pdf.Ln(10.0)
	content1, err := ioutil.ReadFile(resultSummaryFile)
	if err != nil {
		log.Fatalln("file not found:", resultSummaryFile)
	}
	pdf.SetFont("Arial", "", 14)
	pdf.MultiCell(0, 8, string(content1), "", "", false)
	err = pdf.OutputFileAndClose("result.pdf")
	if err != nil {
		log.Fatalln("pdf creation failed: ", err)
	}

	pdf.AddPage()
	content2, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		log.Fatalln("file not found:", yamlFilePath)
	}
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, string(content2), "", "", false)
	err = pdf.OutputFileAndClose("result.pdf")
	if err != nil {
		log.Fatalln("pdf creation failed: ", err)
	}
	fmt.Println(GENERATE_PDF_OUTPUT)

	delFile(resultSummaryFile)
	delFile(yamlFilePath)
}

func writeToYAML(Results interface{}, yamlFileName string) {
	yamlBytes, err := yaml.Marshal(Results)
	if err != nil {
		log.Fatalln("YAML marshal failed, err:", err)
	}

	err = ioutil.WriteFile(yamlFileName, yamlBytes, 0666)
	if err != nil {
		log.Fatalln("WriteFile failed, err:", err)
	}
}
