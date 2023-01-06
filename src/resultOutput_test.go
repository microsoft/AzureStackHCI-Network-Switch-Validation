package main

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestResultAnalysis(t *testing.T) {
	type test struct {
		input    *InputType
		pcapFile string
		want     OutputType
	}
	testFolder := "./test/"
	testCases := map[string]test{
		"fail_lldp_subtype3": {
			input:    &InputType{InterfaceGUID: "\\Device\\NPF_{0217D729-CED0-4D06-9C66-592E032A37A8}", InterfaceAlias: "Ethernet", NativeVlanID: 710, AllVlanIDs: []int{710, 711, 712}, MTUSize: 9214, ETSMaxClass: 8, ETSBWbyPG: "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0", PFCMaxClass: 8, PFCPriorityEnabled: "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0"},
			pcapFile: testFolder + "fail_lldp.pcap",
			want:     OutputType{TestDate: time.Date(2022, time.April, 27, 22, 51, 42, 982572000, time.Local), ResultSummary: map[string][]string{"BGP - PASS": []string(nil), "DHCPRelay - PASS": []string(nil), "LLDP - FAIL": []string{"Incorrect LLDP Subtype3 VLANList - Input: [710 711 712], Found: []"}, "VLAN - FAIL": []string{"Incorrect VLAN ID List - Input: [710 711 712], Found: []"}}, VLANResult: VLANResultType{NativeVlanID: 710, AllVlanIDs: []int{}}, LLDPResult: LLDPResultType{SysDes: "Dell EMC Networking OS10 Enterprise.\r\nCopyright (c) 1999-2021 by Dell Inc. All Rights Reserved.\r\nSystem Description: OS10 Enterprise.\r\nOS Version: 10.5.3.0.\r\nSystem Type: S5248F-VM", PortName: "ethernet1/1/1", ChasisID: "0cc23e6c0000", ChasisIDType: "MAC Address", Subtype1_PortVLANID: 710, Subtype3_VLANList: []int{}, Subtype4_MaxFrameSize: 9214, Subtype7_LinkAggCap: true, Subtype9_ETS: ETSType{ETSTotalPG: 8, ETSBWbyPGID: map[int]int{0: 48, 1: 0, 2: 0, 3: 50, 4: 0, 5: 2, 6: 0, 7: 0}}, SubtypeB_PFC: PFCType{PFCMaxClasses: 8, PFCConfig: map[int]int{0: 0, 1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 0, 7: 0}}}, DHCPResult: DHCPResultType{DHCPPacketDetected: true, RelayAgentIP: net.IP{0xa, 0xa, 0xc, 0x1}}, BGPResult: BGPResultType{BGPTCPPacketDetected: true, SwitchInterfaceIP: "10.10.10.1", SwitchInterfaceMAC: "0c:c2:3e:6c:00:a1", HostInterfaceIP: "10.10.10.11", HostInterfaceMAC: "0c:11:4a:d6:00:01"}},
		},
		"fail_mtu_ets": {
			input:    &InputType{InterfaceGUID: "\\Device\\NPF_{0217D729-CED0-4D06-9C66-592E032A37A8}", InterfaceAlias: "Ethernet", NativeVlanID: 710, AllVlanIDs: []int{710, 711, 712}, MTUSize: 9214, ETSMaxClass: 8, ETSBWbyPG: "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0", PFCMaxClass: 8, PFCPriorityEnabled: "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0"},
			pcapFile: testFolder + "fail_mtu.pcap",
			want:     OutputType{TestDate: time.Date(2022, time.April, 14, 22, 38, 35, 99334000, time.Local), ResultSummary: map[string][]string{"BGP - PASS": []string(nil), "DHCPRelay - PASS": []string(nil), "LLDP - FAIL": []string{"Incorrect LLDP Subtype3 VLANList - Input: [710 711 712], Found: []", "Incorrect Maximum Frame Size - Input:9214, Found: 9216", "Incorrect ETS Class Bandwidth Configured:\n \t\tInput:0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0\n \t\tFound: 0:46,1:1,2:1,3:48,4:1,5:1,6:1,7:1"}, "VLAN - FAIL": []string{"Incorrect VLAN ID List - Input: [710 711 712], Found: []"}}, VLANResult: VLANResultType{NativeVlanID: 710, AllVlanIDs: []int{}}, LLDPResult: LLDPResultType{SysDes: "Dell EMC Networking OS10 Enterprise.\r\nCopyright (c) 1999-2021 by Dell Inc. All Rights Reserved.\r\nSystem Description: OS10 Enterprise.\r\nOS Version: 10.5.3.0.\r\nSystem Type: S5248F-VM", PortName: "ethernet1/1/1", ChasisID: "0cc23e6c0000", ChasisIDType: "MAC Address", Subtype1_PortVLANID: 710, Subtype3_VLANList: []int{}, Subtype4_MaxFrameSize: 9216, Subtype7_LinkAggCap: true, Subtype9_ETS: ETSType{ETSTotalPG: 8, ETSBWbyPGID: map[int]int{0: 46, 1: 1, 2: 1, 3: 48, 4: 1, 5: 1, 6: 1, 7: 1}}, SubtypeB_PFC: PFCType{PFCMaxClasses: 8, PFCConfig: map[int]int{0: 0, 1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 0, 7: 0}}}, DHCPResult: DHCPResultType{DHCPPacketDetected: true, RelayAgentIP: net.IP{0xa, 0xa, 0xc, 0x1}}, BGPResult: BGPResultType{BGPTCPPacketDetected: true, SwitchInterfaceIP: "10.10.10.1", SwitchInterfaceMAC: "0c:c2:3e:6c:00:a1", HostInterfaceIP: "10.10.10.100", HostInterfaceMAC: "0c:b0:99:4a:00:01"}},
		},
		"success_lldp_subtype3": {
			input:    &InputType{InterfaceGUID: "\\Device\\NPF_{0217D729-CED0-4D06-9C66-592E032A37A8}", InterfaceAlias: "Ethernet", NativeVlanID: 1, AllVlanIDs: []int{1, 710, 711, 712}, MTUSize: 9214, ETSMaxClass: 8, ETSBWbyPG: "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0", PFCMaxClass: 8, PFCPriorityEnabled: "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0"},
			pcapFile: testFolder + "success_lldp.pcap",
			want:     OutputType{TestDate: time.Date(2022, time.November, 6, 5, 22, 6, 701322000, time.Local), ResultSummary: map[string][]string{"BGP - FAIL": []string{"TCP 179 Packet Not Detected from switch"}, "DHCPRelay - FAIL": []string{"DHCP Relay Agent IP Not Detected from switch"}, "LLDP - PASS": []string(nil), "VLAN - PASS": []string(nil)}, VLANResult: VLANResultType{NativeVlanID: 1, AllVlanIDs: []int{1, 710, 711, 712}}, LLDPResult: LLDPResultType{SysDes: "Cumulus Linux version 5.2.1 running on Mellanox Technologies Ltd. MSN2100", PortName: "swp1", ChasisID: "98039b5cbb20", ChasisIDType: "MAC Address", Subtype1_PortVLANID: 1, Subtype3_VLANList: []int{1, 710, 711, 712}, Subtype4_MaxFrameSize: 9214, Subtype7_LinkAggCap: true, Subtype9_ETS: ETSType{ETSTotalPG: 8, ETSBWbyPGID: map[int]int{0: 48, 1: 0, 2: 0, 3: 50, 4: 0, 5: 2, 6: 0, 7: 0}}, SubtypeB_PFC: PFCType{PFCMaxClasses: 8, PFCConfig: map[int]int{0: 0, 1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 0, 7: 0}}}, DHCPResult: DHCPResultType{DHCPPacketDetected: true, RelayAgentIP: net.IP(nil)}, BGPResult: BGPResultType{BGPTCPPacketDetected: false, SwitchInterfaceIP: "", SwitchInterfaceMAC: "", HostInterfaceIP: "", HostInterfaceMAC: ""}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := &OutputType{}
			got.resultAnalysis(tc.pcapFile, tc.input)
			// fmt.Printf("%s - %#v\n", name, got)
			if !reflect.DeepEqual(tc.want, *got) {
				t.Errorf("name: %s failed \n want: %#v \n got: %#v", name, tc.want, *got)
			}
		})
	}
}
