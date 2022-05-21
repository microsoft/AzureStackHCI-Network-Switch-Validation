package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type LLDPResultType struct {
	SysDes   string
	PortName string
	ChasisID string
	VLANID   int
	MTU      int
	ETS      ETSType
	PFC      PFCType
}

type PFCType struct {
	PFCMaxClasses      uint8
	PFCPriorityEnabled map[uint8]uint8
}

type ETSType struct {
	ETSTotalPG  uint8
	ETSBWbyPGID map[uint8]uint8
}

func (l *LLDPResultType) decodeLLDPPacket(packet gopacket.Packet) {
	LLDPLayer := packet.Layer(layers.LayerTypeLinkLayerDiscovery)
	if LLDPLayer != nil {
		LLDPType := LLDPLayer.(*layers.LinkLayerDiscovery)
		ChassisIDHex := hex.EncodeToString(LLDPType.ChassisID.ID)
		l.ChasisID = ChassisIDHex
		l.PortName = string(LLDPType.PortID.ID)

		for _, v := range LLDPType.Values {
			// Subtype: Maximum Frame Size 0x04
			if v.Length == 6 && v.Value[3] == 4 {
				l.MTU = int(bytesToDec(v.Value[4:6]))
			}

			// Subtype: Priority Flow Control Configuration 0x0b
			if v.Length == 6 && v.Value[3] == 11 {
				l.PFC.PFCMaxClasses = uint8(bytesToDec(v.Value[4:5]))
				PFCStatusDec := int(bytesToDec(v.Value[5:]))
				l.PFC.PFCPriorityEnabled = postProPFCStatus(PFCStatusDec)
			}

			//Subtype: Port VLAN ID 0x01
			if v.Length == 6 && v.Value[3] == 1 {
				l.VLANID = int(bytesToDec(v.Value[4:]))
			}

			//Subtype: ETS Recommendation 0x0a
			if v.Length == 25 && v.Value[3] == 10 {
				PGIDs := hex.EncodeToString(v.Value[5:9])
				BWbyPGID := make(map[uint8]uint8)
				if len(PGIDs) != 0 {
					l.ETS.ETSTotalPG = uint8(len(PGIDs))
					for i := 0; i < 8; i++ {
						BWbyPGID[uint8(i)] = v.Value[9+i]
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
	}
}

func (o *OutputType) LLDPResultValidation(l *LLDPResultType) {
	var restultFail []string

	if len(l.SysDes) == 0 {
		restultFail = append(restultFail, NO_LLDP_SYS_DSC)
	}
	if len(l.ChasisID) != 12 {
		restultFail = append(restultFail, NO_LLDP_CHASSIS_SUBTYPE)
	}
	if len(l.PortName) == 0 {
		restultFail = append(restultFail, NO_LLDP_PORT_SUBTYPE)
	}
	if l.MTU != INIObj.MTUSize {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_MAXIMUM_FRAME_SIZE, INIObj.MTUSize, l.MTU)
		restultFail = append(restultFail, errMsg)
	}

	if _, ok := o.VLANResult[l.VLANID]; !ok {
		var vlanList []int
		for k := range o.VLANResult {
			vlanList = append(vlanList, k)
		}
		errMsg := fmt.Sprintf("%s - VLANLIST:%v, Found: %d", WRONG_LLDP_VLAN_ID, vlanList, l.VLANID)
		restultFail = append(restultFail, errMsg)
	}
	if l.ETS.ETSTotalPG != uint8(INIObj.ETSMaxClass) {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_ETS_MAX_CLASSES, INIObj.ETSMaxClass, l.ETS.ETSTotalPG)
		restultFail = append(restultFail, errMsg)
	}
	etsBWString := mapintToSlicestring(l.ETS.ETSBWbyPGID)
	if etsBWString != INIObj.ETSBWbyPG {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_ETS_BW, INIObj.ETSBWbyPG, etsBWString)
		restultFail = append(restultFail, errMsg)
	}
	if l.PFC.PFCMaxClasses != uint8(INIObj.PFCMaxClass) {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_PFC_MAX_CLASSES, INIObj.PFCMaxClass, l.PFC.PFCMaxClasses)
		restultFail = append(restultFail, errMsg)
	}
	pfcEnableString := mapintToSlicestring(l.PFC.PFCPriorityEnabled)
	if pfcEnableString != INIObj.PFCPriorityEnabled {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_PFC_ENABLE, INIObj.PFCPriorityEnabled, pfcEnableString)
		restultFail = append(restultFail, errMsg)
	}

	if len(restultFail) == 0 {
		o.ResultSummary["LLDP - PASS"] = restultFail
	} else {
		o.ResultSummary["LLDP - FAIL"] = restultFail
	}
}

func postProPFCStatus(decNum int) map[uint8]uint8 {
	binary := 0
	bit := 1
	remainder := 0
	pfcStatus := make(map[uint8]uint8, 8)

	for i := 0; i < 8; i++ {
		remainder = decNum % 2
		decNum = decNum / 2
		binary += remainder * bit
		bit *= 10
		pfcStatus[uint8(i)] = uint8(remainder)
	}

	return pfcStatus
}

func mapintToSlicestring(mapint map[uint8]uint8) string {
	var outputStringSlice []string

	for i := 0; i < 8; i++ {
		outputStringSlice = append(outputStringSlice, fmt.Sprintf("%d:%d", i, mapint[uint8(i)]))
	}

	return strings.Join(outputStringSlice, ",")
}
