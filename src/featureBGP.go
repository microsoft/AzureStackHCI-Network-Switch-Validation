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

	var BGPReportType FeatureResult

	BGPReportType.FeatureName = BGP

	if !b.BGPTCPPacketDetected {
		BGPReportType.FeaturePass = FAIL
		BGPReportType.FeatureLog = BGPPacket_NOT_Detect
	} else {
		BGPReportType.FeaturePass = PASS
	}

	BGPReportType.FeatureRoles = []string{COMPUTESDN}
	o.FeatureSummary = append(o.FeatureSummary, BGPReportType)
}
