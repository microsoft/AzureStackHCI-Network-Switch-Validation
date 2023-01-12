package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type OutputType struct {
	TestDate      time.Time           `yaml:"TimeDate"`
	ResultSummary map[string][]string `yaml:"ResultSummary"`
	VLANResult    VLANResultType      `yaml:"VLANResult"`
	LLDPResult    LLDPResultType      `yaml:"LLDPResult"`
	DHCPResult    DHCPResultType      `yaml:"DHCPResult"`
	BGPResult     BGPResultType       `yaml:"BGPResult"`
}

type InputType struct {
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

	inputObj     = &InputType{}
	OutputObj    = &OutputType{}
	VLANIDList   []int
	NativeVLANID int
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

	// Scan and collect traffic data to pcap file
	pcapFilePath := fmt.Sprintf("./%s.pcap", inputObj.InterfaceAlias)
	writePcapFile(inputObj, pcapFilePath)
	// pcapFilePath := "./test/success_lldp.pcap"
	fileIsExist(pcapFilePath)
	OutputObj.resultAnalysis(pcapFilePath, inputObj)
	// log.Println(OutputObj)
	pdfFilePath := fmt.Sprintf("./Report_%s.pdf", inputObj.InterfaceAlias)
	OutputObj.outputPDFbyFile(pdfFilePath)
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
	flag.StringVar(&allVlanIDs, "allVlanIDs", "710,711,712", "vlan list string separate with comma")
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
