package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

var (
	inputVariables = &InputType{InterfaceName: "eth0", InterfaceGUID: "\\Device\\NPF_{0217D729-CED0-4D06-9C66-592E032A37A8}", InterfaceAlias: "Ethernet", NativeVlanID: 710, AllVlanIDs: []int{1, 710, 711, 712}, MTUSize: 9214, ETSMaxClass: 8, ETSBWbyPG: "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0", PFCMaxClass: 8, PFCPriorityEnabled: "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0"}
)

func TestResultOutput(t *testing.T) {
	type test struct {
		inputFileName string
	}

	testCases := map[string]test{
		"storage_pass": {
			inputFileName: "storage_pass",
		},
	}

	srcFolder, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testFolder := filepath.Join(srcFolder, "test")
	testInputFolder := filepath.Join(testFolder, "testInput")
	testOutputFolder := filepath.Join(testFolder, "testOutput")
	testGoldenFolder := filepath.Join(testFolder, "goldenConfig")

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := &OutputType{}
			pcapFile := filepath.Join(testInputFolder, tc.inputFileName+".pcap")
			got.resultAnalysis(pcapFile, inputVariables)
			// Parse Golden Yaml to Go Object to compare
			goldenYamlFile := filepath.Join(testGoldenFolder, tc.inputFileName+".yml")
			want := parseYamlToGo(goldenYamlFile)
			if !reflect.DeepEqual(want.RoleResultList, got.RoleResultList) {
				t.Errorf("%s - RoleResultList Failed \n want: %v \n got: %v", name, want.RoleResultList, got.RoleResultList)
			}
			if !reflect.DeepEqual(want.LLDPResult, got.LLDPResult) {
				t.Errorf("%s - LLDPResult Failed \n want: %v \n got: %v", name, want.LLDPResult, got.LLDPResult)
			}
			// Generate Yaml JSON PDF files for view
			pdfFileName := filepath.Join(testOutputFolder, tc.inputFileName+".pdf")
			got.outputPDFFile(pdfFileName)
			yamlFileName := filepath.Join(testOutputFolder, tc.inputFileName+".yml")
			got.outputYAMLFile(yamlFileName)
			jsonFileName := filepath.Join(testOutputFolder, tc.inputFileName+".json")
			got.outputJSONFile(jsonFileName)
		})
	}
}

// Parse Golden Yaml to Go Object
func parseYamlToGo(yamlFileName string) *OutputType {
	outputObj := &OutputType{}
	bytes, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(bytes, outputObj)
	if err != nil {
		log.Fatalln(err)
	}
	return outputObj
}
