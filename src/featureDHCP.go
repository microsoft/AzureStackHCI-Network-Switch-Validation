package main

import (
	"context"
	"log"
	"net"
	"os/exec"
	"time"

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

func decodeDHCPRelayPacket(packet gopacket.Packet) {
	DHCPLayer := packet.Layer(layers.LayerTypeDHCPv4)
	if DHCPLayer != nil {
		DHCPType := DHCPLayer.(*layers.DHCPv4)
		_, dstMac := getPacketMACs(packet)
		// Exclude broadcast dhcp
		if (string(DHCPType.RelayAgentIP) != "0.0.0.0") && dstMac.String() != "ff:ff:ff:ff:ff:ff" {
			// fmt.Println(DHCPType.Contents)
			// fmt.Println(DHCPType.RelayAgentIP)
			// fmt.Println(dstMac.String())
			DHCPResult.RelayAgentIP = DHCPType.RelayAgentIP
			DHCPResult.DHCPPacketDetected = true
		}
	}
}

func triggerWinDHCP() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	cmd := exec.CommandContext(ctx, "ipconfig", "/renew")
	err := cmd.Run()
	if err != nil {
		log.Println("cmd exec error:", err)
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Command timed out")
		return
	}
}

func triggerLinuxDHCP() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	cmd := exec.CommandContext(ctx, "dhclient")
	err := cmd.Run()
	if err != nil {
		log.Println("cmd exec error:", err)
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Command timed out")
		return
	}
}

func DHCPResultValidation() {
	var restultFail []string

	if !DHCPResult.DHCPPacketDetected {
		restultFail = append(restultFail, DHCPPacket_NOT_Detect)
	}
	if DHCPResult.RelayAgentIP == nil {
		restultFail = append(restultFail, DHCPRelay_AgentIP_Not_Detect)
	}

	if len(restultFail) == 0 {
		ResultSummary["DHCPRelay - PASS"] = restultFail
	} else {
		ResultSummary["DHCPRelay - FAIL"] = restultFail
	}
}
