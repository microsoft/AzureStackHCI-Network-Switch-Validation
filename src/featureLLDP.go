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
	SysDes                string
	PortName              string
	ChasisID              string
	ChasisIDType          string
	Subtype1_PortVLANID   int
	Subtype3_VLANList     []int
	Subtype4_MaxFrameSize int
	Subtype7_LinkAggCap   bool
	Subtype9_ETS          ETSType
	SubtypeB_PFC          PFCType
}

type PFCType struct {
	PFCMaxClasses int
	PFCConfig     map[int]int
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
				l.SubtypeB_PFC.PFCMaxClasses = bytesToDec(v.Value[4:5])
				PFCStatusDec := bytesToDec(v.Value[5:])
				l.SubtypeB_PFC.PFCConfig = postProPFCStatus(PFCStatusDec)
			}

			// //Subtype: ETS Recommendation 0x0a
			// if v.Length == 25 && v.Value[3] == 10 {
			// 	PGIDs := hex.EncodeToString(v.Value[5:9])
			// 	BWbyPGID := make(map[int]int)
			// 	if len(PGIDs) != 0 {
			// 		l.Subtype9_ETS.ETSTotalPG = len(PGIDs)
			// 		for i := 0; i < 8; i++ {
			// 			BWbyPGID[i] = int(v.Value[9+i])
			// 		}
			// 		l.Subtype9_ETS.ETSBWbyPGID = BWbyPGID
			// 	}
			// }

			//Subtype: ETS Configuration 0x09
			if v.Length == 25 && v.Value[3] == 9 {
				PGIDs := hex.EncodeToString(v.Value[5:9])
				BWbyPGID := make(map[int]int)
				if len(PGIDs) != 0 {
					l.Subtype9_ETS.ETSTotalPG = len(PGIDs)
					for i := 0; i < 8; i++ {
						BWbyPGID[i] = int(v.Value[9+i])
					}
					l.Subtype9_ETS.ETSBWbyPGID = BWbyPGID
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
		l.Subtype4_MaxFrameSize = int(info8023.MTU)

		info8021, err := LLDPInfoType.Decode8021()
		if err != nil {
			log.Println(err)
		}
		// Subtype3_VLANList
		for _, v := range info8021.VLANNames {
			VLANNameList = append(VLANNameList, int(v.ID))
		}
		l.Subtype3_VLANList = RemoveSliceDup(VLANNameList)

		// Subtype1_PortVLANID
		NativeVLANID := int(info8021.PVID)
		l.Subtype1_PortVLANID = NativeVLANID

		// Subtype7_LinkAgg
		l.Subtype7_LinkAggCap = info8021.LinkAggregation.Supported
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
	if len(l.Subtype3_VLANList) != len(i.TrunkVlanList) {
		restultFail = append(restultFail, INCORRECT_LLDP_Subtype3_VLANList)
	}

	if l.Subtype1_PortVLANID != i.NativeVlanID {
		errMsg := fmt.Sprintf("%s - Input: %d, Found: %d", INCORRECT_LLDP_Subtype1_PortVLANID, i.NativeVlanID, l.Subtype1_PortVLANID)
		restultFail = append(restultFail, errMsg)
	}

	if l.Subtype4_MaxFrameSize != i.MTUSize {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", INCORRECT_LLDP_MAXIMUM_FRAME_SIZE, i.MTUSize, l.Subtype4_MaxFrameSize)
		restultFail = append(restultFail, errMsg)
	}

	if !l.Subtype7_LinkAggCap {
		errMsg := fmt.Sprint(UNSUPPORT_LINK_AGGREGATION)
		restultFail = append(restultFail, errMsg)
	}

	if l.Subtype9_ETS.ETSTotalPG != i.ETSMaxClass {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", INCORRECT_LLDP_ETS_MAX_CLASSES, i.ETSMaxClass, l.Subtype9_ETS.ETSTotalPG)
		restultFail = append(restultFail, errMsg)
	}
	etsBWString := mapintToSlicestring(l.Subtype9_ETS.ETSBWbyPGID)
	if etsBWString != i.ETSBWbyPG {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", INCORRECT_LLDP_ETS_BW, i.ETSBWbyPG, etsBWString)
		restultFail = append(restultFail, errMsg)
	}
	if l.SubtypeB_PFC.PFCMaxClasses != i.PFCMaxClass {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", INCORRECT_LLDP_PFC_MAX_CLASSES, i.PFCMaxClass, l.SubtypeB_PFC.PFCMaxClasses)
		restultFail = append(restultFail, errMsg)
	}
	pfcEnableString := mapintToSlicestring(l.SubtypeB_PFC.PFCConfig)
	if pfcEnableString != i.PFCPriorityEnabled {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", INCORRECT_LLDP_PFC_ENABLE, i.PFCPriorityEnabled, pfcEnableString)
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
