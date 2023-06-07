package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type RoleResultType struct {
	RoleName       string
	RolePass       string
	FeaturesByRole []FeatureResultType
}

type FeatureResultType struct {
	FeatureName       string
	FeaturePass       string
	FeatureLogSubject string
	FeatureLogDetail  string
	FeatureRoles      []string
}

type OutputType struct {
	TestDate          time.Time           `yaml:"TestDate"`
	ToolBuildVersion  string              `yaml:"ToolBuildVersion"`
	RoleResultList    []RoleResultType    `yaml:"RoleResultList"`
	FeatureResultList []FeatureResultType `yaml:"-" json:"-"`
	VLANResult        VLANResultType      `yaml:"VLANResult"`
	LLDPResult        LLDPResultType      `yaml:"LLDPResult"`
	DHCPResult        DHCPResultType      `yaml:"DHCPResult"`
	BGPResult         BGPResultType       `yaml:"BGPResult"`
}

type InputType struct {
	InterfaceName      string
	InterfaceAlias     string
	InterfaceIndex     int
	NativeVlanID       int
	AllVlanIDs         []int
	MTUSize            int
	ETSMaxClass        int
	ETSBWbyPG          string
	PFCMaxClass        int
	PFCPriorityEnabled string
}

type WinNetAdapter struct {
	InterfaceAlias string `json:"InterfaceAlias"`
	InterfaceIndex int    `json:"InterfaceIndex"`
	InterfaceGuid  string `json:"InterfaceGuid"`
}

var (
	logFilePath      = "./result.log"
	ToolBuildVersion = "1.2305.01"

	inputObj                                = &InputType{}
	OutputObj                               = &OutputType{}
	VLANIDList                              []int
	NativeVLANID                            int
	pdfFilePath, yamlFilePath, jsonFilePath string
)

func init() {
	// set up log format
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	// output both logfile and console
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(mw)
}

func main() {
	// parse input variables
	inputObj.loadInputVariable()
	if len(inputObj.AllVlanIDs) < 10 {
		log.Fatalln(VLAN_MINIMUM_10_ERROR)
	}
	inputObj.UpdateInterfaceValues()
	log.Printf("## Interface Name %s is selected, start collecting packages[Maximum 90s or 300 packages]: ##", inputObj.InterfaceAlias)
	log.Println()
	// Scan and collect traffic data to pcap file
	pcapFilePath := fmt.Sprintf("./%s.pcap", inputObj.InterfaceAlias)
	writePcapFile(inputObj.InterfaceName, pcapFilePath)
	// Analyst traffic packages
	fileIsExist(pcapFilePath)
	OutputObj.resultAnalysis(pcapFilePath, inputObj)
	// Write result to format outputs
	pdfFilePath = fmt.Sprintf("./Result_%s.pdf", inputObj.InterfaceAlias)
	yamlFilePath = fmt.Sprintf("./Result_%s.yml", inputObj.InterfaceAlias)
	// jsonFilePath = fmt.Sprintf("./Result_%s.json", inputObj.InterfaceAlias)
	// OutputObj.outputJSONFile(jsonFilePath)
	OutputObj.outputYAMLFile(yamlFilePath)
	OutputObj.outputPDFFile(pdfFilePath, inputObj)

	fmt.Println("---------------------")
	fmt.Println(GENERATE_REPORT_FILES)
}

func (i *InputType) loadInputVariable() {
	var allVlanIDs string
	flag.StringVar(&i.InterfaceName, "InterfaceName", "eth0", "[Linux Only - ifconfig] name of interface connects with network device to be validated")
	flag.IntVar(&i.InterfaceIndex, "InterfaceIndex", 15, "[Windows Only - Get-NetAdapter] index of interface connects with network device to be validated")
	flag.IntVar(&i.NativeVlanID, "nativeVlanID", 710, "native vlan id")
	flag.StringVar(&allVlanIDs, "allVlanIDs", "710,711,712,713,714,715,716,717,718,719,720", "vlan list string separate with comma. Minimum 10 vlans required.")
	flag.IntVar(&i.MTUSize, "mtu", 9214, "mtu value configured on the switch interface")
	flag.IntVar(&i.ETSMaxClass, "etsMaxClass", 8, "maximum number of traffic classes in ETS configuration")
	flag.StringVar(&i.ETSBWbyPG, "etsBWbyPG", "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0", "bandwidth for PGID in ETS configuration")
	flag.IntVar(&i.PFCMaxClass, "pfcMaxClass", 8, "maximum PFC enabled traffic classes in PFC configuration")
	flag.StringVar(&i.PFCPriorityEnabled, "pfcPriorityEnabled", "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0", "PFC for priority in PFC configuration")

	flag.Parse()
	res := strings.Split(allVlanIDs, ",")
	for _, vlan := range res {
		vlanid, err := strconv.Atoi(vlan)
		if err != nil {
			log.Fatalln(err)
		}
		i.AllVlanIDs = append(i.AllVlanIDs, vlanid)
	}
}

func (i *InputType) UpdateInterfaceValues() {
	if runtime.GOOS == "windows" {
		// Define the modified PowerShell command to execute
		command := exec.Command("powershell.exe", "-Command", "Get-NetAdapter | Select-Object InterfaceAlias,InterfaceIndex,InterfaceGuid | ConvertTo-Json")
		// Run the command and collect the output
		output, err := command.Output()
		if err != nil {
			log.Fatal(err)
		}
		// Parse Get-NetAdapter JSON output
		var winNetAdapters []WinNetAdapter
		err = json.Unmarshal(output, &winNetAdapters)
		if err != nil {
			log.Println(output)
			log.Fatal(err)
		}
		// Update Interface Name based on Interface Index
		for _, adapter := range winNetAdapters {
			// fmt.Println("Name:", adapter.InterfaceAlias)
			// fmt.Println("Interface Index:", adapter.InterfaceIndex)
			// fmt.Println("Status:", adapter.InterfaceGuid)
			// fmt.Println()
			if i.InterfaceIndex == adapter.InterfaceIndex {
				i.InterfaceAlias = adapter.InterfaceAlias
				// Customized the interface name to fit gopackage
				// "{89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221}" -> "\Device\NPF_{89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221}"
				winIntfName := fmt.Sprintf("NPF_%s", adapter.InterfaceGuid)
				i.InterfaceName = filepath.Join("\\", "Device", winIntfName)
			}
		}
	} else if runtime.GOOS == "linux" {
		i.InterfaceAlias = i.InterfaceName
	} else {
		log.Fatalln(runtime.GOOS, "Not Support")
	}
}

func fileIsExist(filepath string) {
	_, err := os.Stat(filepath)
	if err != nil {
		log.Fatalf("[Error] Fail founding %s: %v\n", filepath, err)
	}
	log.Println(filepath, "founded.")
}

func RemoveSliceDup(intSlice []int) []int {
	set := make(map[int]bool)
	nonDupSlice := []int{}
	for _, item := range intSlice {
		if _, value := set[item]; !value {
			set[item] = true
			nonDupSlice = append(nonDupSlice, item)
		}
	}
	return nonDupSlice
}
