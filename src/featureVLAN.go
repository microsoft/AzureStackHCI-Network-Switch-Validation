package main

import (
	"fmt"
)

// var (
// 	// Has to be PVST
// 	PVSTDestMAC = "01:00:0c:cc:cc:cd"
// )

type VLANResultType struct {
	NativeVlanID int
	AllVlanIDs   []int
}

// func (v *VLANResultType) decodePVSTPacket(packet gopacket.Packet) {
// 	EthernetLayer := packet.Layer(layers.LayerTypeEthernet)
// 	if EthernetLayer != nil {
// 		EthernetType := EthernetLayer.(*layers.Ethernet)
// 		if net.HardwareAddr(EthernetType.DstMAC).String() == PVSTDestMAC {
// 			lenPayload := len(EthernetType.Payload)
// 			vlanID := bytesToDec(EthernetType.Payload[lenPayload-2 : lenPayload])

// 			if !sliceContains(v.VLANIDs, int(vlanID)) {
// 				v.VLANIDs = append(v.VLANIDs, int(vlanID))
// 			}
// 		}
// 	}
// }

func (o *OutputType) VLANResultValidation(v *VLANResultType, i *InputType) {
	v.NativeVlanID = NativeVLANID
	v.AllVlanIDs = VLANIDList

	var VLANReportType TypeResult

	VLANReportType.TypeName = VLAN
	if len(v.AllVlanIDs) != len(i.AllVlanIDs) {
		errMsg := fmt.Sprintf("%s - Input: %d, Found: %d", INCORRECT_VLAN_ID_LIST, i.AllVlanIDs, v.AllVlanIDs)
		VLANReportType.TypePass = FAIL
		VLANReportType.TypeLog = errMsg
	} else {
		VLANReportType.TypePass = PASS
	}
	VLANReportType.TypeRoles = []string{MANAGEMENT, COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.TypeReportSummary = append(o.TypeReportSummary, VLANReportType)
}
