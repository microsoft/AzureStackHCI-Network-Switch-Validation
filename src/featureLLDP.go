package main

import (
	"encoding/hex"
	"fmt"
	"log"

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

func (o *OutputType) LLDPResultValidation(l *LLDPResultType, i *InputType) {

	var LLDPSubtype1ReportType FeatureResult
	LLDPSubtype1ReportType.FeatureName = LLDP_Subtype1_PortVLANID
	if l.Subtype1_PortVLANID == 0 {
		errMsg := LLDP_Subtype1_NOT_DETECT
		LLDPSubtype1ReportType.FeaturePass = FAIL
		LLDPSubtype1ReportType.FeatureLog = errMsg
	} else if l.Subtype1_PortVLANID != i.NativeVlanID {
		errMsg := fmt.Sprintf("%s - Detect: %d, but Input: %d", LLDP_Subtype1_MISMATCH, l.Subtype1_PortVLANID, i.NativeVlanID)
		LLDPSubtype1ReportType.FeaturePass = FAIL
		LLDPSubtype1ReportType.FeatureLog = errMsg
	} else {
		LLDPSubtype1ReportType.FeaturePass = PASS
	}
	LLDPSubtype1ReportType.FeatureRoles = []string{MANAGEMENT}
	o.FeatureSummary = append(o.FeatureSummary, LLDPSubtype1ReportType)

	var LLDPSubtype3ReportType FeatureResult
	LLDPSubtype3ReportType.FeatureName = LLDP_Subtype3_VLANList
	if len(l.Subtype3_VLANList) == 0 {
		errMsg := LLDP_Subtype3_NOT_DETECT
		LLDPSubtype3ReportType.FeaturePass = FAIL
		LLDPSubtype3ReportType.FeatureLog = errMsg
	} else if len(l.Subtype3_VLANList) != len(i.AllVlanIDs) {
		errMsg := fmt.Sprintf("%s - Detect: %d, but Input: %d", LLDP_Subtype3_MISMATCH, l.Subtype3_VLANList, i.AllVlanIDs)
		LLDPSubtype3ReportType.FeaturePass = FAIL
		LLDPSubtype3ReportType.FeatureLog = errMsg
	} else {
		LLDPSubtype3ReportType.FeaturePass = PASS
	}
	LLDPSubtype3ReportType.FeatureRoles = []string{COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.FeatureSummary = append(o.FeatureSummary, LLDPSubtype3ReportType)

	var LLDPSubtype4ReportType FeatureResult
	LLDPSubtype4ReportType.FeatureName = LLDP_Subtype4_MAX_FRAME_SIZE
	if l.Subtype4_MaxFrameSize == 0 {
		errMsg := LLDP_Subtype4_NOT_DETECT
		LLDPSubtype4ReportType.FeaturePass = FAIL
		LLDPSubtype4ReportType.FeatureLog = errMsg
	} else if l.Subtype4_MaxFrameSize != i.MTUSize {
		errMsg := fmt.Sprintf("%s - Detect: %d, but Input: %d", LLDP_Subtype4_MISMATCH, l.Subtype4_MaxFrameSize, i.MTUSize)
		LLDPSubtype4ReportType.FeaturePass = FAIL
		LLDPSubtype4ReportType.FeatureLog = errMsg
	} else {
		LLDPSubtype4ReportType.FeaturePass = PASS
	}
	LLDPSubtype4ReportType.FeatureRoles = []string{COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.FeatureSummary = append(o.FeatureSummary, LLDPSubtype4ReportType)

	var LLDPSubtype7ReportType FeatureResult
	LLDPSubtype7ReportType.FeatureName = LLDP_Subtype7_LINK_AGGREGATION
	if !l.Subtype7_LinkAggCap {
		errMsg := LLDP_Subtype7_NOT_DETECT
		LLDPSubtype7ReportType.FeaturePass = FAIL
		LLDPSubtype7ReportType.FeatureLog = errMsg
	} else {
		LLDPSubtype7ReportType.FeaturePass = PASS
	}
	LLDPSubtype7ReportType.FeatureRoles = []string{MANAGEMENT, COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.FeatureSummary = append(o.FeatureSummary, LLDPSubtype7ReportType)

	var LLDPETSMaxClassReportType FeatureResult
	LLDPETSMaxClassReportType.FeatureName = LLDP_Subtype9_ETS_MAX_CLASSES
	if l.Subtype9_ETS.ETSTotalPG == 0 {
		errMsg := LLDP_Subtype9_ETS_MAX_CLASSES_NOT_DETECT
		LLDPETSMaxClassReportType.FeaturePass = FAIL
		LLDPETSMaxClassReportType.FeatureLog = errMsg
	} else if l.Subtype9_ETS.ETSTotalPG != i.ETSMaxClass {
		errMsg := fmt.Sprintf("%s - Detect: %d, but Input: %d", LLDP_Subtype9_ETS_MAX_CLASSES_MISMATCH, l.Subtype9_ETS.ETSTotalPG, i.ETSMaxClass)
		LLDPETSMaxClassReportType.FeaturePass = FAIL
		LLDPETSMaxClassReportType.FeatureLog = errMsg
	} else {
		LLDPETSMaxClassReportType.FeaturePass = PASS
	}
	LLDPETSMaxClassReportType.FeatureRoles = []string{STORAGE}
	o.FeatureSummary = append(o.FeatureSummary, LLDPETSMaxClassReportType)

	var LLDPETSBWReportType FeatureResult
	LLDPETSBWReportType.FeatureName = LLDP_Subtype9_ETS_BW
	// etsBWMap := stringToMap(i.ETSBWbyPG)
	// bwLogs := comparePriorityMap(l.Subtype9_ETS.ETSBWbyPGID, etsBWMap)
	detectEtsBW := mapToString(l.Subtype9_ETS.ETSBWbyPGID)
	if detectEtsBW != i.ETSBWbyPG {
		errMsg := fmt.Sprintf("%s - Detect %s, but should be %s", LLDP_Subtype9_ETS_BW_MISMATCH, detectEtsBW, i.ETSBWbyPG)
		LLDPETSBWReportType.FeaturePass = FAIL
		LLDPETSBWReportType.FeatureLog = errMsg
	} else {
		LLDPETSBWReportType.FeaturePass = PASS
	}
	LLDPETSBWReportType.FeatureRoles = []string{STORAGE}
	o.FeatureSummary = append(o.FeatureSummary, LLDPETSBWReportType)

	var LLDPPFCMaxClassReportType FeatureResult
	LLDPPFCMaxClassReportType.FeatureName = LLDP_SubtypeB_PFC_MAX_CLASSES
	if l.SubtypeB_PFC.PFCMaxClasses == 0 {
		errMsg := LLDP_SubtypeB_PFC_MAX_CLASSES_NOT_DETECT
		LLDPPFCMaxClassReportType.FeaturePass = FAIL
		LLDPPFCMaxClassReportType.FeatureLog = errMsg
	} else if l.SubtypeB_PFC.PFCMaxClasses != i.PFCMaxClass {
		errMsg := fmt.Sprintf("%s - Detect: %d, but Input: %d", LLDP_SubtypeB_PFC_MAX_CLASSES_MISMATCH, l.SubtypeB_PFC.PFCMaxClasses, i.PFCMaxClass)
		LLDPPFCMaxClassReportType.FeaturePass = FAIL
		LLDPPFCMaxClassReportType.FeatureLog = errMsg
	} else {
		LLDPPFCMaxClassReportType.FeaturePass = PASS
	}
	LLDPPFCMaxClassReportType.FeatureRoles = []string{STORAGE}
	o.FeatureSummary = append(o.FeatureSummary, LLDPPFCMaxClassReportType)

	var LLDPPFCEnableReportType FeatureResult
	LLDPPFCEnableReportType.FeatureName = LLDP_SubtypeB_PFC_ENABLE
	// pfcBWMap := stringToMap(i.PFCPriorityEnabled)
	// pfcLogs := comparePriorityMap(l.SubtypeB_PFC.PFCConfig, pfcBWMap)
	detectPfcConfig := mapToString(l.SubtypeB_PFC.PFCConfig)
	if detectPfcConfig != i.PFCPriorityEnabled {
		errMsg := fmt.Sprintf("%s - Detect %s, but should be %s", LLDP_SubtypeB_PFC_ENABLE_MISMATCH, detectPfcConfig, i.PFCPriorityEnabled)
		LLDPPFCEnableReportType.FeaturePass = FAIL
		LLDPPFCEnableReportType.FeatureLog = errMsg
	} else {
		LLDPPFCEnableReportType.FeaturePass = PASS
	}
	LLDPPFCEnableReportType.FeatureRoles = []string{STORAGE}
	o.FeatureSummary = append(o.FeatureSummary, LLDPPFCEnableReportType)

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

// input: map[int]int{0:48, 1:0, 2:0, 3:50, 4:0, 5:2, 6:0, 7:0}
// output: string 0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0
func mapToString(inputMap map[int]int) string {
	outputString := ""
	for i := 0; i < 8; i++ {
		outputString += fmt.Sprintf("%d:%d,", i, inputMap[i])
	}
	return outputString[:len(outputString)-1]
}
