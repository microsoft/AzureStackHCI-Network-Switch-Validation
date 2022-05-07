package main

import (
	"encoding/hex"
	"fmt"

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

func decodeLLDPPacket(packet gopacket.Packet) {
	LLDPLayer := packet.Layer(layers.LayerTypeLinkLayerDiscovery)
	if LLDPLayer != nil {
		LLDPType := LLDPLayer.(*layers.LinkLayerDiscovery)
		ChassisIDHex := hex.EncodeToString(LLDPType.ChassisID.ID)
		LLDPResult.ChasisID = ChassisIDHex
		LLDPResult.PortName = string(LLDPType.PortID.ID)

		for _, v := range LLDPType.Values {
			// Subtype: Maximum Frame Size 0x04
			if v.Length == 6 && v.Value[3] == 4 {
				LLDPResult.MTU = int(bytesToDec(v.Value[4:6]))
			}

			// Subtype: Priority Flow Control Configuration 0x0b
			if v.Length == 6 && v.Value[3] == 11 {
				LLDPResult.PFC.PFCMaxClasses = uint8(bytesToDec(v.Value[4:5]))
				PFCStatusDec := int(bytesToDec(v.Value[5:]))
				LLDPResult.PFC.PFCPriorityEnabled = postProPFCStatus(PFCStatusDec)
			}

			//Subtype: Port VLAN ID 0x01
			if v.Length == 6 && v.Value[3] == 1 {
				LLDPResult.VLANID = int(bytesToDec(v.Value[4:]))
			}

			//Subtype: ETS Recommendation 0x0a
			if v.Length == 25 && v.Value[3] == 10 {
				PGIDs := hex.EncodeToString(v.Value[5:9])
				BWbyPGID := make(map[uint8]uint8)
				if len(PGIDs) != 0 {
					LLDPResult.ETS.ETSTotalPG = uint8(len(PGIDs))
					for i := 0; i < 8; i++ {
						BWbyPGID[uint8(i)] = v.Value[9+i]
					}
					LLDPResult.ETS.ETSBWbyPGID = BWbyPGID
				}
			}

		}
		// fmt.Println(LLDPType.Contents)
		// Organisation Specific 6 [0 128 194 11 8 8] PFC
		// {Organisation Specific 25 [0 128 194 10 0 0 3 5 0 46 1 1 48 1 1 1 1 2 2 2 2 2 2 2 2]} ETS
		// {Organisation Specific 6 [0 18 15 4 36 0]} Maximum Frame Size
		// {Organisation Specific 6 [0 128 194 1 0 7]} Port VLAN ID
	}
}

func decodeLLDPInfoPacket(packet gopacket.Packet) {
	LLDPInfoLayer := packet.Layer(layers.LayerTypeLinkLayerDiscoveryInfo)
	if LLDPInfoLayer != nil {
		LLDPInfoType := LLDPInfoLayer.(*layers.LinkLayerDiscoveryInfo)
		LLDPResult.SysDes = LLDPInfoType.SysDescription
	}
}

func LLDPResultValidation() {
	var restultFail []string

	if len(LLDPResult.SysDes) == 0 {
		restultFail = append(restultFail, NO_LLDP_SYS_DSC)
	}
	if len(LLDPResult.ChasisID) != 12 {
		restultFail = append(restultFail, NO_LLDP_CHASSIS_SUBTYPE)
	}
	if len(LLDPResult.PortName) == 0 {
		restultFail = append(restultFail, NO_LLDP_PORT_SUBTYPE)
	}
	if LLDPResult.MTU != INIObj.MTUSize {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_MAXIMUM_FRAME_SIZE, INIObj.MTUSize, LLDPResult.MTU)
		restultFail = append(restultFail, errMsg)
	}
	if _, ok := VLANResult[LLDPResult.VLANID]; !ok {
		var vlanList []int
		for k := range VLANResult {
			vlanList = append(vlanList, k)
		}
		errMsg := fmt.Sprintf("%s - VLANLIST:%v, Found: %d", WRONG_LLDP_VLAN_ID, vlanList, LLDPResult.VLANID)
		restultFail = append(restultFail, errMsg)
	}
	if LLDPResult.ETS.ETSTotalPG != uint8(INIObj.ETSMaxClass) {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_ETS_MAX_CLASSES, INIObj.ETSMaxClass, LLDPResult.ETS.ETSTotalPG)
		restultFail = append(restultFail, errMsg)
	}
	etsBWString := mapintToSlicestring(LLDPResult.ETS.ETSBWbyPGID)
	if etsBWString != INIObj.ETSBWbyPG {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_ETS_BW, INIObj.ETSBWbyPG, etsBWString)
		restultFail = append(restultFail, errMsg)
	}
	if LLDPResult.PFC.PFCMaxClasses != uint8(INIObj.PFCMaxClass) {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", WRONG_LLDP_PFC_MAX_CLASSES, INIObj.PFCMaxClass, LLDPResult.PFC.PFCMaxClasses)
		restultFail = append(restultFail, errMsg)
	}
	pfcEnableString := mapintToSlicestring(LLDPResult.PFC.PFCPriorityEnabled)
	if pfcEnableString != INIObj.PFCPriorityEnabled {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_PFC_ENABLE, INIObj.PFCPriorityEnabled, pfcEnableString)
		restultFail = append(restultFail, errMsg)
	}

	if len(restultFail) == 0 {
		ResultSummary["LLDP - PASS"] = restultFail
	} else {
		ResultSummary["LLDP - FAIL"] = restultFail
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
