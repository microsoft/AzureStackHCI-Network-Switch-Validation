package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

func writePcapFile(intfName string) {
	var (
		// false will only collect packages to the interface
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
		// Process packet here
		// fmt.Println(packet)
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		start++
		fmt.Printf("Collecting Network Packages: [%d / %d (Max)]\n", start, packetCountMax)
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
