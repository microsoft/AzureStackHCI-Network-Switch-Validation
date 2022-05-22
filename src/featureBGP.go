package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type BGPResultType struct {
	BGPTCPPacketDetected bool
	SwitchInterfaceIP    string
	SwitchInterfaceMAC   string
	HostInterfaceIP      string
	HostInterfaceMAC     string
}

func decodeBGPPacket(packet gopacket.Packet) bool {
	// Incoming Packet TCP DstPort=179
	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	if tcpLayer != nil {
		tcpType := tcpLayer.(*layers.TCP)
		if tcpType.DstPort == 179 {
			// fmt.Println(tcpType.Contents)
			return true
		}
	}
	return false
}

func (o *OutputType) BGPResultValidation(b *BGPResultType) {
	var restultFail []string

	if !b.BGPTCPPacketDetected {
		restultFail = append(restultFail, BGPPacket_NOT_Detect)
	}

	if len(restultFail) == 0 {
		o.ResultSummary["BGP - PASS"] = restultFail
	} else {
		o.ResultSummary["BGP - FAIL"] = restultFail
	}
}
