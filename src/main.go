package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/gopacket/pcap"
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
	HostInterfaceIP    string
	VlanIDs            []int
	MTUSize            int
	ETSMaxClass        int
	ETSBWbyPG          string
	PFCMaxClass        int
	PFCPriorityEnabled string
}

var (
	logFilePath   = "./result/result.log"
	npcapFilePath = "npcap_install.ps1"

	INIObj    = &INIType{}
	OutputObj = &OutputType{}

	// intfName      = "\\Device\\NPF_{D63D7017-9A01-42F2-B2DB-C34055E3BDBD}"
	intfNameMap   = make(map[string]interface{}, 20)
	packetMaxSize = int32(9216)
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
	Check_Npcap(npcapFilePath)

	getInterfaceByIP()
	if len(intfNameMap) > 0 {
		log.Println(intfNameMap)
		tmpPCAPfolder := "tmpPcap"
		CreateFolder(tmpPCAPfolder)
		for intfName, intfDes := range intfNameMap {
			pcapFilePath := fmt.Sprintf("%s/%s.pcap", tmpPCAPfolder, intfDes)
			writePcapFile(intfName, pcapFilePath)
			fileIsExist(pcapFilePath)
			OutputObj.resultAnalysis(pcapFilePath, INIObj)
			log.Println(OutputObj)
			if len(OutputObj.LLDPResult.SysDes) != 0 {
				resultfolder := "result"
				CreateFolder(resultfolder)
				nbrport := OutputObj.LLDPResult.PortName
				fileName := CheckFileName(nbrport)
				pdfFilePath := fmt.Sprintf("%s/%s.pdf", resultfolder, fileName)
				OutputObj.outputPDFbyFile(pdfFilePath)
				pcapFileNewPath := fmt.Sprintf("%s/%s.pcap", resultfolder, fileName)
				err := os.Rename(pcapFilePath, pcapFileNewPath)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println(NO_LLDP_PACKET + intfName)
			}
		}
		return
	} else {
		log.Fatalln(NO_INTF)
	}
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
	i.HostInterfaceIP = cfg.Section("host").Key("hostInterfaceIP").String()
	i.VlanIDs = cfg.Section("vlan").Key("vlanIDs").ValidInts(",")
	i.MTUSize = cfg.Section("mtu").Key("mtuSize").MustInt(9174)
	i.ETSMaxClass = cfg.Section("ets").Key("ETSMaxClass").MustInt(8)
	i.ETSBWbyPG = cfg.Section("ets").Key("ETSBWbyPG").MustString("0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0")
	i.PFCMaxClass = cfg.Section("pfc").Key("PFCMaxClass").MustInt(8)
	i.PFCPriorityEnabled = cfg.Section("pfc").Key("PFCPriorityEnabled").MustString("0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0")
}

// func getInterfaceByIP() {
// 	interfaces, err := pcap.FindAllDevs()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	matchFlag := false
// 	for _, intf := range interfaces {
// 		for _, address := range intf.Addresses {
// 			fmt.Println(intf, address)
// 			maskNum, _ := address.Netmask.Size()
// 			lintfName := intf.Name
// 			intfIPMask := fmt.Sprintf("%v/%d", address.IP, maskNum)
// 			if intfIPMask == INIObj.HostInterfaceIP {
// 				intfName = lintfName
// 				log.Printf("Found matched host interface by IP: %s - %s\n", intfIPMask, intfName)
// 				matchFlag = true
// 				return
// 			}
// 		}
// 	}
// 	if !matchFlag {
// 		log.Printf("%s: %s - %s\n", INTF_NOT_MATCH, INIObj.HostInterfaceIP, intfName)
// 		return
// 	}
// }

func getInterfaceByIP() {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, intf := range interfaces {
		for _, address := range intf.Addresses {
			if address.Netmask != nil {
				// log.Println(intf.Name, address.IP, intf.Description)
				intfNameMap[intf.Name] = intf.Description
			}
		}
	}
}

func Check_Npcap(filePath string) {
	fileIsExist(filePath)
	abs, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	pscmd := fmt.Sprintf("'& %s'", abs)
	cmd := exec.Command("powershell", "-noexit", pscmd)
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Npcap is installed.")
}

func CheckFileName(filename string) string {
	filename = strings.Replace(filename, "/", "_", -1)
	filename = strings.Replace(filename, `\`, "_", -1)
	filename = strings.Replace(filename, `:`, "_", -1)
	return filename
}

func CreateFolder(folderName string) {
	if _, err := os.Stat(folderName); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
