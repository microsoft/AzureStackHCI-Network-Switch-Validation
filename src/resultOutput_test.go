package main

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestResultAnalysis(t *testing.T) {
	type test struct {
		input    string
		pcapFile string
		want     OutputType
	}
	testFolder := "./test/"
	testCases := map[string]test{
		"fail_lldp_subtype3": {
			testFolder + "input.ini",
			testFolder + "fail_lldp.pcap",
			OutputType{TestDate: time.Date(2022, time.April, 27, 22, 51, 42, 982572000, time.Local), ResultSummary: map[string][]string{"BGP - PASS": []string(nil), "DHCPRelay - PASS": []string(nil), "LLDP - FAIL": []string{"No LLDP IEEE 802.1 VLAN Name (Subtype 3) Founded"}, "VLAN - PASS": []string(nil)}, VLANResult: VLANResultType{VLANIDs: []int{710, 711, 712}}, LLDPResult: LLDPResultType{SysDes: "Dell EMC Networking OS10 Enterprise.\r\nCopyright (c) 1999-2021 by Dell Inc. All Rights Reserved.\r\nSystem Description: OS10 Enterprise.\r\nOS Version: 10.5.3.0.\r\nSystem Type: S5248F-VM", PortName: "ethernet1/1/1", ChasisID: "0cc23e6c0000", ChasisIDType: "MAC Address", VLANID: 710, IEEE8021Subtype3: []uint16(nil), MTU: 9214, ETS: ETSType{ETSTotalPG: 0x8, ETSBWbyPGID: map[uint8]uint8{0x0: 0x30, 0x1: 0x0, 0x2: 0x0, 0x3: 0x32, 0x4: 0x0, 0x5: 0x2, 0x6: 0x0, 0x7: 0x0}}, PFC: PFCType{PFCMaxClasses: 0x8, PFCPriorityEnabled: map[uint8]uint8{0x0: 0x0, 0x1: 0x0, 0x2: 0x0, 0x3: 0x1, 0x4: 0x0, 0x5: 0x0, 0x6: 0x0, 0x7: 0x0}}}, DHCPResult: DHCPResultType{DHCPPacketDetected: true, RelayAgentIP: net.IP{0xa, 0xa, 0xc, 0x1}}, BGPResult: BGPResultType{BGPTCPPacketDetected: true, SwitchInterfaceIP: "10.10.10.1", SwitchInterfaceMAC: "0c:c2:3e:6c:00:a1", HostInterfaceIP: "10.10.10.11", HostInterfaceMAC: "0c:11:4a:d6:00:01"}},
		},
		"fail_mtu_ets": {
			testFolder + "input.ini",
			testFolder + "fail_mtu.pcap",
			OutputType{TestDate: time.Date(2022, time.April, 14, 22, 38, 35, 99334000, time.Local), ResultSummary: map[string][]string{"BGP - PASS": []string(nil), "DHCPRelay - PASS": []string(nil), "LLDP - FAIL": []string{"No LLDP IEEE 802.1 VLAN Name (Subtype 3) Founded", "Incorrect Maximum Frame Size - Input:9214, Found: 9216", "Incorrect ETS Class Bandwidth Configured:\n \t\tInput:0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0\n \t\tFound: 0:46,1:1,2:1,3:48,4:1,5:1,6:1,7:1"}, "VLAN - PASS": []string(nil)}, VLANResult: VLANResultType{VLANIDs: []int{710, 711, 712}}, LLDPResult: LLDPResultType{SysDes: "Dell EMC Networking OS10 Enterprise.\r\nCopyright (c) 1999-2021 by Dell Inc. All Rights Reserved.\r\nSystem Description: OS10 Enterprise.\r\nOS Version: 10.5.3.0.\r\nSystem Type: S5248F-VM", PortName: "ethernet1/1/1", ChasisID: "0cc23e6c0000", ChasisIDType: "MAC Address", VLANID: 710, IEEE8021Subtype3: []uint16(nil), MTU: 9216, ETS: ETSType{ETSTotalPG: 0x8, ETSBWbyPGID: map[uint8]uint8{0x0: 0x2e, 0x1: 0x1, 0x2: 0x1, 0x3: 0x30, 0x4: 0x1, 0x5: 0x1, 0x6: 0x1, 0x7: 0x1}}, PFC: PFCType{PFCMaxClasses: 0x8, PFCPriorityEnabled: map[uint8]uint8{0x0: 0x0, 0x1: 0x0, 0x2: 0x0, 0x3: 0x1, 0x4: 0x0, 0x5: 0x0, 0x6: 0x0, 0x7: 0x0}}}, DHCPResult: DHCPResultType{DHCPPacketDetected: true, RelayAgentIP: net.IP{0xa, 0xa, 0xc, 0x1}}, BGPResult: BGPResultType{BGPTCPPacketDetected: true, SwitchInterfaceIP: "10.10.10.1", SwitchInterfaceMAC: "0c:c2:3e:6c:00:a1", HostInterfaceIP: "10.10.10.100", HostInterfaceMAC: "0c:b0:99:4a:00:01"}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			iniObj := &INIType{}
			got := &OutputType{}
			iniObj.loadIniFile(tc.input)
			got.resultAnalysis(tc.pcapFile, iniObj)
			// fmt.Printf("%s - %#v\n", name, got)
			if !reflect.DeepEqual(*got, tc.want) {
				t.Errorf("name: %s failed \n want: %v \n got: %v", name, tc.want, *got)
			}
		})
	}
}
