package main

import (
	"context"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

func writePcapFile(intfName string) {
	var (
		// False will only collect packages to the interface
		promiscuous    bool          = true
		timeout        time.Duration = -1 * time.Second
		handle         *pcap.Handle
		packetCountMax int = 300
		// Duration default unit is nanosecond
		sessionTimeout time.Duration = 90000000000
	)
	// Open output pcap file and write header
	f, _ := os.Create(pcapFilePath)
	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(uint32(packetMaxSize), layers.LinkTypeEthernet)
	defer f.Close()

	// Open the device for capturing
	handle, err := pcap.OpenLive(intfName, packetMaxSize, promiscuous, timeout)
	if err != nil {
		log.Fatalf("Error opening interface %s, %v", intfName, err)
	}
	// Create timeout for Live Session
	go OpenLiveTimeout(handle, sessionTimeout)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Refresh interface to collect dhcp packet
	triggerWinDHCP()
	triggerLinuxDHCP()
	start := 0
	for packet := range packetSource.Packets() {
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		start++
		log.Printf("Collecting Network Packages: [%d / %d (Max)]\n", start, packetCountMax)
		// Set maximum packets to collect
		if start > packetCountMax {
			break
		}
	}
}

func OpenLiveTimeout(handle *pcap.Handle, sessionTimeout time.Duration) {
	time.Sleep(sessionTimeout)
	handle.Close()
	log.Printf("Reach preset max session time %v, close live collection.\n", sessionTimeout)
}

func (o *OutputType) decodePacketLayer(pcapFilePath string) {
	handle, err := pcap.OpenOffline(pcapFilePath)
	if err != nil {
		log.Fatalf("Error opening %s, error:%v", pcapFilePath, err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {

		o.VLANResult.decodePVSTPacket(packet)

		if decodeDHCPPacket(packet) {
			o.DHCPResult.DHCPPacketDetected = true
			o.DHCPResult.decodeDHCPRelayPacket(packet)
		}

		o.LLDPResult.decodeLLDPPacket(packet)
		o.LLDPResult.decodeLLDPInfoPacket(packet)

		if decodeBGPPacket(packet) {
			o.BGPResult.BGPTCPPacketDetected = true
			o.BGPResult.SwitchInterfaceMAC, o.BGPResult.HostInterfaceMAC = getPacketMACs(packet)
			o.BGPResult.SwitchInterfaceIP, o.BGPResult.HostInterfaceIP = getPacketIPv4s(packet)
		}

		o.TestDate = packet.Metadata().Timestamp
	}
}

func bytesToDec(bytes []byte) int {
	hexNum := hex.EncodeToString(bytes)
	decNum, err := strconv.ParseInt(hexNum, 16, 32)
	if err != nil {
		log.Fatalln(err)
	}
	return int(decNum)
}

func getPacketMACs(packet gopacket.Packet) (SrcMac, DstMac string) {
	EthernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if EthernetLayer != nil {
		EthernetType := EthernetLayer.(*layers.Ethernet)
		return EthernetType.SrcMAC.String(), EthernetType.DstMAC.String()
	}
	log.Fatalf("Not able to decode the network packet: %#v\n", packet)
	return "", ""
}

func getPacketIPv4s(packet gopacket.Packet) (SrcIP, DstIP string) {
	IPv4Layer := packet.Layer(layers.LayerTypeIPv4)
	if IPv4Layer != nil {
		IPv4Type := IPv4Layer.(*layers.IPv4)
		return IPv4Type.SrcIP.String(), IPv4Type.DstIP.String()
	}
	log.Fatalf("Not able to decode the network packet: %#v\n", packet)
	return "", ""
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

func sliceContains(elems []int, v int) bool {
	for _, i := range elems {
		if v == i {
			return true
		}
	}
	return false
}
