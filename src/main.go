package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type RoleResultType struct {
	RolePass       string
	FeaturesByRole []FeatureResultType
}

type FeatureResultType struct {
	FeatureName  string
	FeaturePass  string
	FeatureLog   string
	FeatureRoles []string
}

type OutputType struct {
	TestDate          time.Time                 `yaml:"TestDate"`
	RoleResultList    map[string]RoleResultType `yaml:"RoleResultList"`
	FeatureResultList []FeatureResultType       `yaml:"FeatureResultList"`
	VLANResult        VLANResultType            `yaml:"VLANResult"`
	LLDPResult        LLDPResultType            `yaml:"LLDPResult"`
	DHCPResult        DHCPResultType            `yaml:"DHCPResult"`
	BGPResult         BGPResultType             `yaml:"BGPResult"`
}

type InputType struct {
	InterfaceName      string
	InterfaceGUID      string
	InterfaceAlias     string
	NativeVlanID       int
	AllVlanIDs         []int
	MTUSize            int
	ETSMaxClass        int
	ETSBWbyPG          string
	PFCMaxClass        int
	PFCPriorityEnabled string
}

var (
	logFilePath = "./result.log"

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
	if runtime.GOOS == "windows" {
		// fmt.Println("Running on Windows")
		inputObj.loadInputVariable()
		// Scan and collect traffic data to pcap file
		pcapFilePath := fmt.Sprintf("./%s.pcap", inputObj.InterfaceAlias)
		writePcapFile(inputObj.InterfaceGUID, pcapFilePath)
		// srcFolder, err := os.Getwd()
		// if err != nil {
		// 	panic(err)
		// }
		// testFolder := filepath.Join(srcFolder, "test")
		// testInputFolder := filepath.Join(testFolder, "testInput")
		// pcapFilePath := filepath.Join(testInputFolder, "storage_pass.pcap")
		fileIsExist(pcapFilePath)
		OutputObj.resultAnalysis(pcapFilePath, inputObj)
		// log.Println(OutputObj)
		pdfFilePath = fmt.Sprintf("./Report_%s.pdf", inputObj.InterfaceAlias)
		yamlFilePath = fmt.Sprintf("./Report_%s.yml", inputObj.InterfaceAlias)
		jsonFilePath = fmt.Sprintf("./Report_%s.json", inputObj.InterfaceAlias)
	} else if runtime.GOOS == "linux" {
		// fmt.Println("Running on Linux")
		var iniFilePath string
		// Parse iniFile to Input Object
		flag.StringVar(&iniFilePath, "iniFilePath", "./input.ini", "Please input INI file path.")
		flag.Parse()
		fileIsExist(iniFilePath)
		inputObj.loadIniFile(iniFilePath)

		// Scan and collect traffic data to pcap file
		pcapFilePath := fmt.Sprintf("./%s.pcap", inputObj.InterfaceName)
		writePcapFile(inputObj.InterfaceName, pcapFilePath)
		fileIsExist(pcapFilePath)
		OutputObj.resultAnalysis(pcapFilePath, inputObj)
		pdfFilePath = fmt.Sprintf("./%s.pdf", inputObj.InterfaceName)
		yamlFilePath = fmt.Sprintf("./%s.yml", inputObj.InterfaceName)
		jsonFilePath = fmt.Sprintf("./%s.json", inputObj.InterfaceName)
	} else {
		fmt.Println(runtime.GOOS, "Not Support")
	}

	OutputObj.outputPDFFile(pdfFilePath)
	OutputObj.outputYAMLFile(yamlFilePath)
	OutputObj.outputJSONFile(jsonFilePath)
	fmt.Println("---------------------")
	fmt.Println(GENERATE_REPORT_FILES)
}

func fileIsExist(filepath string) {
	_, err := os.Stat(filepath)
	if err != nil {
		log.Fatalf("[Error] Fail founding %s: %v\n", filepath, err)
	}
	log.Println(filepath, "founded.")
}

func (i *InputType) loadInputVariable() {
	var allVlanIDs string
	flag.IntVar(&i.NativeVlanID, "nativeVlanID", 1, "native vlan id")
	flag.StringVar(&allVlanIDs, "allVlanIDs", "1,710,711,712", "vlan list string separate with comma")
	flag.IntVar(&i.MTUSize, "mtu", 9214, "mtu value configured on the switch interface")
	flag.IntVar(&i.ETSMaxClass, "etsMaxClass", 8, "maximum number of traffic classes in ETS configuration")
	flag.StringVar(&i.ETSBWbyPG, "etsBWbyPG", "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0", "bandwidth for PGID in ETS configuration")
	flag.IntVar(&i.PFCMaxClass, "pfcMaxClass", 8, "maximum PFC enabled traffic classes in PFC configuration")
	flag.StringVar(&i.PFCPriorityEnabled, "pfcPriorityEnabled", "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0", "PFC for priority in PFC configuration")
	flag.StringVar(&i.InterfaceGUID, "interfaceGUID", "", "Powershell: Get-NetAdapter | Select-Object InterfaceAlias,InterfaceGuid")
	flag.StringVar(&i.InterfaceAlias, "interfaceAlias", "", "Powershell: Get-NetAdapter | Select-Object InterfaceAlias,InterfaceGuid")

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

func (i *InputType) loadIniFile(filePath string) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Fatalf("Fail to read file: %v\n", err)
	}
	i.InterfaceName = cfg.Section("host").Key("interfaceName").String()
	i.NativeVlanID = cfg.Section("vlan").Key("nativeVlanID").MustInt()
	i.AllVlanIDs = cfg.Section("vlan").Key("allVlanIDs").ValidInts(",")
	i.MTUSize = cfg.Section("mtu").Key("mtuSize").MustInt(9174)
	i.ETSMaxClass = cfg.Section("ets").Key("ETSMaxClass").MustInt(8)
	i.ETSBWbyPG = cfg.Section("ets").Key("ETSBWbyPG").MustString("0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0")
	i.PFCMaxClass = cfg.Section("pfc").Key("PFCMaxClass").MustInt(8)
	i.PFCPriorityEnabled = cfg.Section("pfc").Key("PFCPriorityEnabled").MustString("0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0")
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
