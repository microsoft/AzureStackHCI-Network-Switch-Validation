package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func decodePacketLayer(pcapFilePath string) {
	handle, err := pcap.OpenOffline(pcapFilePath)
	if err != nil {
		fmt.Printf("Error opening %s, error:%v", pcapFilePath, err)
		log.Fatalf("Error opening %s, error:%v", pcapFilePath, err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		decodePVSTPacket(packet)
		decodeLLDPPacket(packet)
		decodeLLDPInfoPacket(packet)
		if decodeBGPPacket(packet) {
			SwitchInterfaceMAC, HostInterfaceMAC = getPacketMACs(packet)
			SwitchInterfaceIP, HostInterfaceIP = getPacketIPv4s(packet)
			BGPResult.BGPTCPPacketDetected = true
			BGPResult.SwitchInterfaceMAC = net.HardwareAddr(SwitchInterfaceMAC).String()
			BGPResult.SwitchInterfaceIP = SwitchInterfaceIP.String()
		}

		if decodeDHCPPacket(packet) {
			decodeDHCPRelayPacket(packet)
		}
	}
}

func bytesToDec(bytes []byte) int64 {
	hexNum := hex.EncodeToString(bytes)
	decNum, err := strconv.ParseInt(hexNum, 16, 32)
	if err != nil {
		log.Fatalln(err)
	}
	return decNum
}

func delFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Fatalln("remove file:", err)
	}
}
