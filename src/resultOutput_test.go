package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

var (
	testFolder       = "/workspaces/AzureStackHCI-Network-Switch-Validation/src/test"
	testInputFolder  = testFolder + "/testInput/"
	testOutputFolder = testFolder + "/testOutput/"
	testGoldenFolder = testFolder + "/goldenConfig/"
	inputVariables   = &InputType{InterfaceGUID: "\\Device\\NPF_{0217D729-CED0-4D06-9C66-592E032A37A8}", InterfaceAlias: "Ethernet", NativeVlanID: 710, AllVlanIDs: []int{710, 711, 712}, MTUSize: 9214, ETSMaxClass: 8, ETSBWbyPG: "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0", PFCMaxClass: 8, PFCPriorityEnabled: "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0"}
)

func TestResultOutput(t *testing.T) {
	type test struct {
		inputFileName string
	}

	testCases := map[string]test{
		"all_fail": {
			inputFileName: "all_fail",
		},
		"storage_pass": {
			inputFileName: "storage_pass",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := &OutputType{}
			pcapFile := testInputFolder + tc.inputFileName + ".pcap"
			got.resultAnalysis(pcapFile, inputVariables)
			// Generate Yaml and PDF files for view
			outputPdfFile := testOutputFolder + tc.inputFileName + ".pdf"
			got.outputPDFFile(outputPdfFile)
			outputYamlFile := testOutputFolder + tc.inputFileName + ".yml"
			got.outputYAMLFile(outputYamlFile)
			// Parse Golden Yaml to Go Object to compare
			goldenYamlFile := testGoldenFolder + tc.inputFileName + ".yml"
			want := parseYamltoGo(goldenYamlFile)
			if !reflect.DeepEqual(*want, *got) {
				t.Errorf("name: %s failed \n want: %v \n got: %v", name, *want, *got)
			}
		})
	}
}

// Parse Golden Yaml to Go Object
func parseYamltoGo(yamlFileName string) *OutputType {
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
