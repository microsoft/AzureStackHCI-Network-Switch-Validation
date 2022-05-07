package main

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type BGPResultType struct {
	BGPTCPPacketDetected bool
	SwitchInterfaceIP    string
	SwitchInterfaceMAC   string
}

func getPacketMACs(packet gopacket.Packet) (SrcMac, DstMac net.HardwareAddr) {
	EthernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if EthernetLayer != nil {
		EthernetType := EthernetLayer.(*layers.Ethernet)
		return EthernetType.SrcMAC, EthernetType.DstMAC
	}
	log.Fatalf("Not able to decode the network packet: %#v\n", packet)
	return nil, nil
}

func getPacketIPv4s(packet gopacket.Packet) (SrcIP, DstIP net.IP) {
	IPv4Layer := packet.Layer(layers.LayerTypeIPv4)
	if IPv4Layer != nil {
		IPv4Type := IPv4Layer.(*layers.IPv4)
		return IPv4Type.SrcIP, IPv4Type.DstIP
	}
	log.Fatalf("Not able to decode the network packet: %#v\n", packet)
	return nil, nil
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

func BGPResultValidation() {
	var restultFail []string

	if !BGPResult.BGPTCPPacketDetected {
		restultFail = append(restultFail, BGPPacket_NOT_Detect)
	}

	if len(restultFail) == 0 {
		ResultSummary["BGP - PASS"] = restultFail
	} else {
		ResultSummary["BGP - FAIL"] = restultFail
	}
}
