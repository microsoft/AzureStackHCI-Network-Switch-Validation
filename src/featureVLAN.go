package main

import (
	"fmt"
	"net"
	"sort"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	// Has to be PVST
	PVSTDestMAC = "01:00:0c:cc:cc:cd"
)

func decodePVSTPacket(packet gopacket.Packet) {
	// All VLANs
	EthernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if EthernetLayer != nil {
		EthernetType := EthernetLayer.(*layers.Ethernet)
		if net.HardwareAddr(EthernetType.DstMAC).String() == PVSTDestMAC {
			lenPayload := len(EthernetType.Payload)
			vlanID := bytesToDec(EthernetType.Payload[lenPayload-2 : lenPayload])
			if _, ok := VLANResult[int(vlanID)]; !ok {
				VLANResult[int(vlanID)] = struct{}{}
			}
		}
	}
}

func VLANResultValidation() {
	var restultFail []string
	var vlanList []int
	for k := range VLANResult {
		vlanList = append(vlanList, k)
	}

	sort.Slice(vlanList, func(i, j int) bool {
		return vlanList[i] < vlanList[j]
	})

	sort.Slice(INIObj.VlanIDs, func(i, j int) bool {
		return INIObj.VlanIDs[i] < INIObj.VlanIDs[j]
	})

	if len(VLANResult) == len(INIObj.VlanIDs) {
		for _, v := range INIObj.VlanIDs {
			if _, ok := VLANResult[v]; !ok {
				vlanError := fmt.Sprintf("%s - Input: %v, Found: %v", VLAN_NOT_MATCH, INIObj.VlanIDs, vlanList)
				restultFail = append(restultFail, vlanError)
			}
		}
	} else {
		vlanError := fmt.Sprintf("%s - Input: %v, Found: %v", VLAN_NOT_MATCH, INIObj.VlanIDs, vlanList)
		restultFail = append(restultFail, vlanError)
	}

	if len(restultFail) == 0 {
		ResultSummary["VLAN - PASS"] = restultFail
	} else {
		ResultSummary["VLAN - FAIL"] = restultFail
	}
}
