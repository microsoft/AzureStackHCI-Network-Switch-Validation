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
			VLANIDList = append(VLANIDList, int(v.ID))
		}
		l.Subtype3_VLANList = RemoveSliceDup(VLANIDList)
		// Update Global Var for VLAN feature validation
		VLANIDList = l.Subtype3_VLANList

		// Subtype1_PortVLANID
		NativeVLANID = int(info8021.PVID)
		l.Subtype1_PortVLANID = NativeVLANID

		// Subtype7_LinkAgg
		l.Subtype7_LinkAggCap = info8021.LinkAggregation.Supported
	}
}

func (o *OutputType) LLDPResultValidation(l *LLDPResultType, i *INIType) {

	var LLDPSubtype1ReportType TypeResult
	LLDPSubtype1ReportType.TypeName = LLDP_Subtype1_PortVLANID
	if l.Subtype1_PortVLANID != i.NativeVlanID {
		errMsg := fmt.Sprintf("%s - Input: %d, Found: %d", INCORRECT_LLDP_Subtype1_PortVLANID, i.NativeVlanID, l.Subtype1_PortVLANID)
		LLDPSubtype1ReportType.TypePass = FAIL
		LLDPSubtype1ReportType.TypeLog = errMsg
	} else {
		LLDPSubtype1ReportType.TypePass = PASS
	}
	LLDPSubtype1ReportType.TypeRoles = []string{MANAGEMENT}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPSubtype1ReportType)

	var LLDPSubtype3ReportType TypeResult
	LLDPSubtype3ReportType.TypeName = LLDP_Subtype3_VLANList
	if len(l.Subtype3_VLANList) != len(i.AllVlanIDs) {
		errMsg := fmt.Sprintf("%s - Input: %d, Found: %d", INCORRECT_LLDP_Subtype3_VLANList, i.AllVlanIDs, l.Subtype3_VLANList)
		LLDPSubtype3ReportType.TypePass = FAIL
		LLDPSubtype3ReportType.TypeLog = errMsg
	} else {
		LLDPSubtype3ReportType.TypePass = PASS
	}
	LLDPSubtype3ReportType.TypeRoles = []string{COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPSubtype3ReportType)

	var LLDPSubtype4ReportType TypeResult
	LLDPSubtype4ReportType.TypeName = LLDP_MAXIMUM_FRAME_SIZE
	if l.Subtype4_MaxFrameSize != i.MTUSize {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", INCORRECT_LLDP_MAXIMUM_FRAME_SIZE, i.MTUSize, l.Subtype4_MaxFrameSize)
		LLDPSubtype4ReportType.TypePass = FAIL
		LLDPSubtype4ReportType.TypeLog = errMsg
	} else {
		LLDPSubtype4ReportType.TypePass = PASS
	}
	LLDPSubtype4ReportType.TypeRoles = []string{COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPSubtype4ReportType)

	var LLDPSubtype7ReportType TypeResult
	LLDPSubtype7ReportType.TypeName = LLDP_LINK_AGGREGATION
	if !l.Subtype7_LinkAggCap {
		errMsg := fmt.Sprint(UNSUPPORT_LLDP_LINK_AGGREGATION)
		LLDPSubtype7ReportType.TypePass = FAIL
		LLDPSubtype7ReportType.TypeLog = errMsg
	} else {
		LLDPSubtype7ReportType.TypePass = PASS
	}
	LLDPSubtype7ReportType.TypeRoles = []string{MANAGEMENT, COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPSubtype7ReportType)

	var LLDPETSMaxClassReportType TypeResult
	LLDPETSMaxClassReportType.TypeName = LLDP_ETS_MAX_CLASSES
	if l.Subtype9_ETS.ETSTotalPG != i.ETSMaxClass {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", INCORRECT_LLDP_ETS_MAX_CLASSES, i.ETSMaxClass, l.Subtype9_ETS.ETSTotalPG)
		LLDPETSMaxClassReportType.TypePass = FAIL
		LLDPETSMaxClassReportType.TypeLog = errMsg
	} else {
		LLDPETSMaxClassReportType.TypePass = PASS
	}
	LLDPETSMaxClassReportType.TypeRoles = []string{STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPETSMaxClassReportType)

	var LLDPETSBWReportType TypeResult
	LLDPETSBWReportType.TypeName = LLDP_ETS_BW
	etsBWString := mapintToSlicestring(l.Subtype9_ETS.ETSBWbyPGID)
	if etsBWString != i.ETSBWbyPG {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", INCORRECT_LLDP_ETS_BW, i.ETSBWbyPG, etsBWString)
		LLDPETSBWReportType.TypePass = FAIL
		LLDPETSBWReportType.TypeLog = errMsg
	} else {
		LLDPETSBWReportType.TypePass = PASS
	}
	LLDPETSBWReportType.TypeRoles = []string{STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPETSBWReportType)

	var LLDPPFCMaxClassReportType TypeResult
	LLDPPFCMaxClassReportType.TypeName = LLDP_PFC_MAX_CLASSES
	if l.SubtypeB_PFC.PFCMaxClasses != i.PFCMaxClass {
		errMsg := fmt.Sprintf("%s - Input:%d, Found: %d", INCORRECT_LLDP_PFC_MAX_CLASSES, i.PFCMaxClass, l.SubtypeB_PFC.PFCMaxClasses)
		LLDPPFCMaxClassReportType.TypePass = FAIL
		LLDPPFCMaxClassReportType.TypeLog = errMsg
	} else {
		LLDPPFCMaxClassReportType.TypePass = PASS
	}
	LLDPPFCMaxClassReportType.TypeRoles = []string{STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPPFCMaxClassReportType)

	var LLDPPFCEnableReportType TypeResult
	LLDPPFCEnableReportType.TypeName = LLDP_PFC_ENABLE
	pfcEnableString := mapintToSlicestring(l.SubtypeB_PFC.PFCConfig)
	if pfcEnableString != i.PFCPriorityEnabled {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", INCORRECT_LLDP_PFC_ENABLE, i.PFCPriorityEnabled, pfcEnableString)
		LLDPPFCEnableReportType.TypePass = FAIL
		LLDPPFCEnableReportType.TypeLog = errMsg
	} else {
		LLDPPFCEnableReportType.TypePass = PASS
	}
	LLDPPFCEnableReportType.TypeRoles = []string{STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, LLDPPFCEnableReportType)

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
