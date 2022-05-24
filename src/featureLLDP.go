package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type LLDPResultType struct {
	SysDes           string
	PortName         string
	ChasisID         string
	ChasisIDType     string
	VLANID           int
	IEEE8021Subtype3 []uint16
	MTU              int
	ETS              ETSType
	PFC              PFCType
}

type PFCType struct {
	PFCMaxClasses      int
	PFCPriorityEnabled map[int]int
}

type ETSType struct {
	ETSTotalPG  int
	ETSBWbyPGID map[int]int
}

func (l *LLDPResultType) decodeLLDPPacket(packet gopacket.Packet) {
	LLDPLayer := packet.Layer(layers.LayerTypeLinkLayerDiscovery)
	if LLDPLayer != nil {
		LLDPType := LLDPLayer.(*layers.LinkLayerDiscovery)
		ChassisIDHex := hex.EncodeToString(LLDPType.ChassisID.ID)
		l.ChasisID = ChassisIDHex
		l.ChasisIDType = LLDPType.ChassisID.Subtype.String()
		l.PortName = string(LLDPType.PortID.ID)

		for _, v := range LLDPType.Values {
			// Subtype: Priority Flow Control Configuration 0x0b
			if v.Length == 6 && v.Value[3] == 11 {
				l.PFC.PFCMaxClasses = bytesToDec(v.Value[4:5])
				PFCStatusDec := bytesToDec(v.Value[5:])
				l.PFC.PFCPriorityEnabled = postProPFCStatus(PFCStatusDec)
			}

			//Subtype: Port VLAN ID 0x01
			if v.Length == 6 && v.Value[3] == 1 {
				l.VLANID = bytesToDec(v.Value[4:])
			}

			//Subtype: ETS Recommendation 0x0a
			if v.Length == 25 && v.Value[3] == 10 {
				PGIDs := hex.EncodeToString(v.Value[5:9])
				BWbyPGID := make(map[int]int)
				if len(PGIDs) != 0 {
					l.ETS.ETSTotalPG = len(PGIDs)
					for i := 0; i < 8; i++ {
						BWbyPGID[i] = int(v.Value[9+i])
					}
					l.ETS.ETSBWbyPGID = BWbyPGID
				}
			}
		}
	}
}

func (l *LLDPResultType) decodeLLDPInfoPacket(packet gopacket.Packet) {
	LLDPInfoLayer := packet.Layer(layers.LayerTypeLinkLayerDiscoveryInfo)
	if LLDPInfoLayer != nil {
		LLDPInfoType := LLDPInfoLayer.(*layers.LinkLayerDiscoveryInfo)
		l.SysDes = LLDPInfoType.SysDescription

		info8023, err := LLDPInfoType.Decode8023()
		if err != nil {
			log.Println(err)
		}
		l.MTU = int(info8023.MTU)

		info8021, err := LLDPInfoType.Decode8021()
		if err != nil {
			log.Println(err)
		}
		for _, v := range info8021.VLANNames {
			l.IEEE8021Subtype3 = append(l.IEEE8021Subtype3, v.ID)
		}
	}
}

func (o *OutputType) LLDPResultValidation(l *LLDPResultType, i *INIType) {
	var restultFail []string

	if len(l.SysDes) == 0 {
		restultFail = append(restultFail, NO_LLDP_SYS_DSC)
	}
	if l.ChasisIDType != CHASIS_ID_TYPE {
		restultFail = append(restultFail, NO_LLDP_CHASSIS_SUBTYPE)
	}
	if len(l.PortName) == 0 {
		restultFail = append(restultFail, NO_LLDP_PORT_SUBTYPE)
	}
	if len(l.IEEE8021Subtype3) == 0 {
		restultFail = append(restultFail, NO_LLDP_IEEE_8021_Subtype3)
	}

	if l.MTU != i.MTUSize {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_MAXIMUM_FRAME_SIZE, i.MTUSize, l.MTU)
		restultFail = append(restultFail, errMsg)
	}

	if !sliceContains(o.VLANResult.VLANIDs, l.VLANID) {
		errMsg := fmt.Sprintf("%s - VLANList:%v, Found: %d", WRONG_LLDP_VLAN_ID, o.VLANResult.VLANIDs, l.VLANID)
		restultFail = append(restultFail, errMsg)
	}

	if l.ETS.ETSTotalPG != i.ETSMaxClass {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_ETS_MAX_CLASSES, i.ETSMaxClass, l.ETS.ETSTotalPG)
		restultFail = append(restultFail, errMsg)
	}
	etsBWString := mapintToSlicestring(l.ETS.ETSBWbyPGID)
	if etsBWString != i.ETSBWbyPG {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_ETS_BW, i.ETSBWbyPG, etsBWString)
		restultFail = append(restultFail, errMsg)
	}
	if l.PFC.PFCMaxClasses != i.PFCMaxClass {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_PFC_MAX_CLASSES, i.PFCMaxClass, l.PFC.PFCMaxClasses)
		restultFail = append(restultFail, errMsg)
	}
	pfcEnableString := mapintToSlicestring(l.PFC.PFCPriorityEnabled)
	if pfcEnableString != i.PFCPriorityEnabled {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_PFC_ENABLE, i.PFCPriorityEnabled, pfcEnableString)
		restultFail = append(restultFail, errMsg)
	}

	if len(restultFail) == 0 {
		o.ResultSummary["LLDP - PASS"] = restultFail
	} else {
		o.ResultSummary["LLDP - FAIL"] = restultFail
	}
}

func postProPFCStatus(decNum int) map[int]int {
	binary := 0
	bit := 1
	remainder := 0
	pfcStatus := make(map[int]int, 8)

	for i := 0; i < 8; i++ {
		remainder = decNum % 2
		decNum = decNum / 2
		binary += remainder * bit
		bit *= 10
		pfcStatus[i] = remainder
	}

	return pfcStatus
}

func mapintToSlicestring(mapint map[int]int) string {
	var outputStringSlice []string

	for i := 0; i < 8; i++ {
		outputStringSlice = append(outputStringSlice, fmt.Sprintf("%d:%d", i, mapint[i]))
	}

	return strings.Join(outputStringSlice, ",")
}
