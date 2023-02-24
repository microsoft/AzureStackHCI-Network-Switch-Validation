package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type DHCPResultType struct {
	DHCPPacketDetected bool
	RelayAgentIP       net.IP
}

func decodeDHCPPacket(packet gopacket.Packet) bool {
	UDPLayer := packet.Layer(layers.LayerTypeUDP)
	if UDPLayer != nil {
		UDPType := UDPLayer.(*layers.UDP)
		if UDPType.DstPort == 67 {
			return true
		}
	}
	return false
}

func (d *DHCPResultType) decodeDHCPRelayPacket(packet gopacket.Packet) {
	DHCPLayer := packet.Layer(layers.LayerTypeDHCPv4)
	if DHCPLayer != nil {
		DHCPType := DHCPLayer.(*layers.DHCPv4)
		_, dstMac := getPacketMACs(packet)
		// Exclude broadcast dhcp
		if (string(DHCPType.RelayAgentIP) != "0.0.0.0") && dstMac != "ff:ff:ff:ff:ff:ff" {
			// fmt.Println(DHCPType.Contents)
			// fmt.Println(DHCPType.RelayAgentIP)
			// fmt.Println(dstMac.String())
			d.RelayAgentIP = DHCPType.RelayAgentIP
			d.DHCPPacketDetected = true
		}
	}
}

func (o *OutputType) DHCPResultValidation(d *DHCPResultType) {

	var DHCPRelayReportType FeatureResultType

	DHCPRelayReportType.FeatureName = DHCPRelay
	if d.RelayAgentIP == nil {
		DHCPRelayReportType.FeaturePass = FAIL
		DHCPRelayReportType.FeatureLog = DHCPRelay_AgentIP_Not_Detect
	} else {
		DHCPRelayReportType.FeaturePass = PASS
	}
	DHCPRelayReportType.FeatureRoles = []string{MANAGEMENT}
	o.FeatureResultList = append(o.FeatureResultList, DHCPRelayReportType)
}
