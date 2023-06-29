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
				// Interate TSA for Traffic Class to make sure ETS is configured.
				ETSEnableClassList := []int{}
				for idx, ets2 := range v.Value[17:] {
					// ETS enabled is 0x02
					if ets2 == 2 {
						ETSEnableClassList = append(ETSEnableClassList, idx)
					}
				}
				BWbyPGID := make(map[int]int)
				if len(ETSEnableClassList) > 0 {
					ETSTotalPG := bytesToDec(v.Value[4:5])
					if ETSTotalPG == 0 {
						ETSTotalPG = 8
					}
					l.Subtype9_ETS.ETSTotalPG = ETSTotalPG
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

		if NativeVLANID != 0 {
			l.Subtype1_PortVLANID = NativeVLANID
		}

		// Subtype7_LinkAgg
		l.Subtype7_LinkAggCap = info8021.LinkAggregation.Supported
	}
}

func (o *OutputType) LLDPResultValidation(l *LLDPResultType, i *InputType) {

	var LLDPSubtype1ReportType FeatureResultType
	LLDPSubtype1ReportType.FeatureName = LLDP_Subtype1_PortVLANID
	if l.Subtype1_PortVLANID == 0 {
		LLDPSubtype1ReportType.FeaturePass = FAIL
		LLDPSubtype1ReportType.FeatureLogSubject = LLDP_Subtype1_NOT_DETECT
	} else if l.Subtype1_PortVLANID != i.NativeVlanID {
		LLDPSubtype1ReportType.FeaturePass = FAIL
		LLDPSubtype1ReportType.FeatureLogSubject = LLDP_Subtype1_MISMATCH
		errMsg := fmt.Sprintf("Detect: %d, but Input: %d", l.Subtype1_PortVLANID, i.NativeVlanID)
		LLDPSubtype1ReportType.FeatureLogDetail = errMsg
	} else {
		LLDPSubtype1ReportType.FeaturePass = PASS
	}
	LLDPSubtype1ReportType.FeatureRoles = []string{MANAGEMENT}
	o.FeatureResultList = append(o.FeatureResultList, LLDPSubtype1ReportType)

	var LLDPSubtype3ReportType FeatureResultType
	LLDPSubtype3ReportType.FeatureName = LLDP_Subtype3_VLANList
	if len(l.Subtype3_VLANList) == 0 {
		LLDPSubtype3ReportType.FeaturePass = FAIL
		LLDPSubtype3ReportType.FeatureLogSubject = LLDP_Subtype3_NOT_DETECT
	} else if len(l.Subtype3_VLANList) != len(i.AllVlanIDs) {
		LLDPSubtype3ReportType.FeaturePass = FAIL
		LLDPSubtype3ReportType.FeatureLogSubject = LLDP_Subtype3_MISMATCH
		errMsg := fmt.Sprintf("Detect: %d, but Input: %d", l.Subtype3_VLANList, i.AllVlanIDs)
		LLDPSubtype3ReportType.FeatureLogDetail = errMsg
	} else {
		LLDPSubtype3ReportType.FeaturePass = PASS
	}
	LLDPSubtype3ReportType.FeatureRoles = []string{COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, LLDPSubtype3ReportType)

	var LLDPSubtype4ReportType FeatureResultType
	LLDPSubtype4ReportType.FeatureName = LLDP_Subtype4_MAX_FRAME_SIZE
	if l.Subtype4_MaxFrameSize == 0 {
		LLDPSubtype4ReportType.FeaturePass = FAIL
		LLDPSubtype4ReportType.FeatureLogSubject = LLDP_Subtype4_NOT_DETECT
	} else if l.Subtype4_MaxFrameSize != i.MTUSize {
		LLDPSubtype4ReportType.FeaturePass = FAIL
		LLDPSubtype4ReportType.FeatureLogSubject = LLDP_Subtype4_MISMATCH
		errMsg := fmt.Sprintf("Detect: %d, but Input: %d", l.Subtype4_MaxFrameSize, i.MTUSize)
		LLDPSubtype4ReportType.FeatureLogDetail = errMsg
	} else {
		LLDPSubtype4ReportType.FeaturePass = PASS
	}
	LLDPSubtype4ReportType.FeatureRoles = []string{COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, LLDPSubtype4ReportType)

	var LLDPSubtype7ReportType FeatureResultType
	LLDPSubtype7ReportType.FeatureName = LLDP_Subtype7_LINK_AGGREGATION
	if !l.Subtype7_LinkAggCap {
		LLDPSubtype7ReportType.FeaturePass = FAIL
		LLDPSubtype7ReportType.FeatureLogSubject = LLDP_Subtype7_NOT_DETECT
	} else {
		LLDPSubtype7ReportType.FeaturePass = PASS
	}
	LLDPSubtype7ReportType.FeatureRoles = []string{MANAGEMENT, COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, LLDPSubtype7ReportType)

	var LLDPETSMaxClassReportType FeatureResultType
	LLDPETSMaxClassReportType.FeatureName = LLDP_Subtype9_ETS_MAX_CLASSES
	if l.Subtype9_ETS.ETSTotalPG == 0 {
		LLDPETSMaxClassReportType.FeaturePass = FAIL
		LLDPETSMaxClassReportType.FeatureLogSubject = LLDP_Subtype9_ETS_MAX_CLASSES_NOT_DETECT
	} else if l.Subtype9_ETS.ETSTotalPG != i.ETSMaxClass {
		LLDPETSMaxClassReportType.FeaturePass = FAIL
		LLDPETSMaxClassReportType.FeatureLogSubject = LLDP_Subtype9_ETS_MAX_CLASSES_MISMATCH
		errMsg := fmt.Sprintf("Detect: %d, but Input: %d", l.Subtype9_ETS.ETSTotalPG, i.ETSMaxClass)
		LLDPETSMaxClassReportType.FeatureLogDetail = errMsg
	} else {
		LLDPETSMaxClassReportType.FeaturePass = PASS
	}
	LLDPETSMaxClassReportType.FeatureRoles = []string{STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, LLDPETSMaxClassReportType)

	var LLDPETSBWReportType FeatureResultType
	LLDPETSBWReportType.FeatureName = LLDP_Subtype9_ETS_BW
	// etsBWMap := stringToMap(i.ETSBWbyPG)
	// bwLogs := comparePriorityMap(l.Subtype9_ETS.ETSBWbyPGID, etsBWMap)
	detectEtsBW := mapToString(l.Subtype9_ETS.ETSBWbyPGID)
	if detectEtsBW != i.ETSBWbyPG {
		LLDPETSBWReportType.FeaturePass = FAIL
		LLDPETSBWReportType.FeatureLogSubject = LLDP_Subtype9_ETS_BW_MISMATCH
		errMsg := fmt.Sprintf("Detect %s, but should be %s", detectEtsBW, i.ETSBWbyPG)
		LLDPETSBWReportType.FeatureLogDetail = errMsg
	} else {
		LLDPETSBWReportType.FeaturePass = PASS
	}
	LLDPETSBWReportType.FeatureRoles = []string{STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, LLDPETSBWReportType)

	var LLDPPFCMaxClassReportType FeatureResultType
	LLDPPFCMaxClassReportType.FeatureName = LLDP_SubtypeB_PFC_MAX_CLASSES
	if l.SubtypeB_PFC.PFCMaxClasses == 0 {
		LLDPPFCMaxClassReportType.FeaturePass = FAIL
		LLDPPFCMaxClassReportType.FeatureLogSubject = LLDP_SubtypeB_PFC_MAX_CLASSES_NOT_DETECT
	} else if l.SubtypeB_PFC.PFCMaxClasses != i.PFCMaxClass {
		LLDPPFCMaxClassReportType.FeaturePass = FAIL
		LLDPPFCMaxClassReportType.FeatureLogSubject = LLDP_SubtypeB_PFC_MAX_CLASSES_MISMATCH
		errMsg := fmt.Sprintf("Detect: %d, but Input: %d", l.SubtypeB_PFC.PFCMaxClasses, i.PFCMaxClass)
		LLDPPFCMaxClassReportType.FeatureLogDetail = errMsg
	} else {
		LLDPPFCMaxClassReportType.FeaturePass = PASS
	}
	LLDPPFCMaxClassReportType.FeatureRoles = []string{STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, LLDPPFCMaxClassReportType)

	var LLDPPFCEnableReportType FeatureResultType
	LLDPPFCEnableReportType.FeatureName = LLDP_SubtypeB_PFC_ENABLE
	// pfcBWMap := stringToMap(i.PFCPriorityEnabled)
	// pfcLogs := comparePriorityMap(l.SubtypeB_PFC.PFCConfig, pfcBWMap)
	detectPfcConfig := mapToString(l.SubtypeB_PFC.PFCConfig)
	if detectPfcConfig != i.PFCPriorityEnabled {
		LLDPPFCEnableReportType.FeaturePass = FAIL
		LLDPPFCEnableReportType.FeatureLogSubject = LLDP_SubtypeB_PFC_ENABLE_MISMATCH
		errMsg := fmt.Sprintf("Detect %s, but should be %s", detectPfcConfig, i.PFCPriorityEnabled)
		LLDPPFCEnableReportType.FeatureLogDetail = errMsg
	} else {
		LLDPPFCEnableReportType.FeaturePass = PASS
	}
	LLDPPFCEnableReportType.FeatureRoles = []string{STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, LLDPPFCEnableReportType)

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
