package main

import (
	"fmt"
	"net"
	"reflect"
	"sort"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	// Has to be PVST
	PVSTDestMAC = "01:00:0c:cc:cc:cd"
)

type VLANResultType struct {
	VLANIDs []int
}

func (v *VLANResultType) decodePVSTPacket(packet gopacket.Packet) {
	EthernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if EthernetLayer != nil {
		EthernetType := EthernetLayer.(*layers.Ethernet)
		if net.HardwareAddr(EthernetType.DstMAC).String() == PVSTDestMAC {
			lenPayload := len(EthernetType.Payload)
			vlanID := bytesToDec(EthernetType.Payload[lenPayload-2 : lenPayload])

			if !sliceContains(v.VLANIDs, int(vlanID)) {
				v.VLANIDs = append(v.VLANIDs, int(vlanID))
			}
		}
	}
}

func (o *OutputType) VLANResultValidation(v *VLANResultType, i *INIType) {
	var restultFail []string
	var vlanList []int

	for k := range v.VLANIDs {
		vlanList = append(vlanList, k)
	}

	sort.Slice(vlanList, func(i, j int) bool {
		return vlanList[i] < vlanList[j]
	})

	sort.Slice(i.VlanIDs, func(m, n int) bool {
		return i.VlanIDs[m] < i.VlanIDs[n]
	})

	if !reflect.DeepEqual(v.VLANIDs, i.VlanIDs) {
		vlanError := fmt.Sprintf("%s - Input: %v, Found: %v", VLAN_NOT_MATCH, i.VlanIDs, vlanList)
		restultFail = append(restultFail, vlanError)
	}

	if len(restultFail) == 0 {
		o.ResultSummary["VLAN - PASS"] = restultFail
	} else {
		o.ResultSummary["VLAN - FAIL"] = restultFail
	}
}
