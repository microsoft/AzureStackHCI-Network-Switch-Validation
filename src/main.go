package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type OutputType struct {
	TestDate      time.Time           `yaml:"TimeDate"`
	ResultSummary map[string][]string `yaml:"ResultSummary"`
	VLANResult    VLANResultType      `yaml:"VLANResult"`
	LLDPResult    LLDPResultType      `yaml:"LLDPResult"`
	DHCPResult    DHCPResultType      `yaml:"DHCPResult"`
	BGPResult     BGPResultType       `yaml:"BGPResult"`
}

type INIType struct {
	InterfaceName      string
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

	INIObj       = &INIType{}
	OutputObj    = &OutputType{}
	VLANIDList   []int
	NativeVLANID int
)

func init() {
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
	var iniFilePath string
	flag.StringVar(&iniFilePath, "iniFilePath", "./input.ini", "Please input INI file path.")
	flag.Parse()
	fileIsExist(iniFilePath)
	INIObj.loadIniFile(iniFilePath)
	fmt.Println(INIObj)

	pcapFilePath := fmt.Sprintf("./%s.pcap", INIObj.InterfaceName)
	writePcapFile(INIObj.InterfaceName, pcapFilePath)
	fileIsExist(pcapFilePath)
	OutputObj.resultAnalysis(pcapFilePath, INIObj)
	pdfFilePath := fmt.Sprintf("./%s.pdf", INIObj.InterfaceName)
	OutputObj.outputPDFbyFile(pdfFilePath)
}

func fileIsExist(filepath string) {
	_, err := os.Stat(filepath)
	if err != nil {
		log.Fatalf("[Error] Fail founding %s: %v\n", filepath, err)
	}
	log.Println(filepath, "founded.")
}

func (i *INIType) loadIniFile(filePath string) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Fatalf("Fail to read file: %v\n", err)
	}
	i.InterfaceName = cfg.Section("host").Key("interfaceName").String()
	i.NativeVlanID = cfg.Section("vlan").Key("nativeVlanID").MustInt()
	i.AllVlanIDs = cfg.Section("vlan").Key("vlanIDs").ValidInts(",")
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
