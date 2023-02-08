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
		"success_lldp_subtype3": {
			input:    testFolder + "input.ini",
			pcapFile: testFolder + "success_lldp.pcap",
			want:     OutputType{TestDate: time.Date(2022, time.November, 6, 5, 22, 6, 701322000, time.Local), RoleReportSummary: map[string]string{"Compute(Basic)": "Pass", "Compute(SDN)": "Fail", "Management": "Fail", "Storage": "Pass"}, TypeReportSummary: []TypeResult{TypeResult{TypeName: "BGP", TypePass: "Fail", TypeLog: "TCP 179 Packet Not Detected from switch", TypeRoles: []string{"Compute(SDN)"}}, TypeResult{TypeName: "VLAN", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Management", "Compute(Basic)", "Compute(SDN)", "Storage"}}, TypeResult{TypeName: "DHCP Relay Agent IP", TypePass: "Fail", TypeLog: "DHCP Relay Agent IP Not Detected from switch", TypeRoles: []string{"Management"}}, TypeResult{TypeName: "LLDP-Port VLAN ID (Subtype = 1)", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Management"}}, TypeResult{TypeName: "LLDP-VLAN Name (Subtype = 3)", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Compute(Basic)", "Compute(SDN)", "Storage"}}, TypeResult{TypeName: "LLDP-Maximum Frame Size", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Compute(Basic)", "Compute(SDN)", "Storage"}}, TypeResult{TypeName: "LLDP-Link Aggregation", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Management", "Compute(Basic)", "Compute(SDN)", "Storage"}}, TypeResult{TypeName: "LLDP-ETS Maximum Number of Traffic Classes", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Storage"}}, TypeResult{TypeName: "LLDP-ETS Class Bandwidth Configured", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Storage"}}, TypeResult{TypeName: "LLDP-PFC Maximum Number of Traffic Classes", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Storage"}}, TypeResult{TypeName: "LLDP-PFC Priority Class Enabled", TypePass: "Pass", TypeLog: "", TypeRoles: []string{"Storage"}}}, VLANResult: VLANResultType{NativeVlanID: 1, AllVlanIDs: []int{1, 710, 711, 712}}, LLDPResult: LLDPResultType{SysDes: "Cumulus Linux version 5.2.1 running on Mellanox Technologies Ltd. MSN2100", PortName: "swp1", ChasisID: "98039b5cbb20", ChasisIDType: "MAC Address", Subtype1_PortVLANID: 1, Subtype3_VLANList: []int{1, 710, 711, 712}, Subtype4_MaxFrameSize: 9214, Subtype7_LinkAggCap: true, Subtype9_ETS: ETSType{ETSTotalPG: 8, ETSBWbyPGID: map[int]int{0: 48, 1: 0, 2: 0, 3: 50, 4: 0, 5: 2, 6: 0, 7: 0}}, SubtypeB_PFC: PFCType{PFCMaxClasses: 8, PFCConfig: map[int]int{0: 0, 1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 0, 7: 0}}}, DHCPResult: DHCPResultType{DHCPPacketDetected: true, RelayAgentIP: net.IP(nil)}, BGPResult: BGPResultType{BGPTCPPacketDetected: false, SwitchInterfaceIP: "", SwitchInterfaceMAC: "", HostInterfaceIP: "", HostInterfaceMAC: ""}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			iniObj := &INIType{}
			got := &OutputType{}
			iniObj.loadIniFile(tc.input)
			got.resultAnalysis(tc.pcapFile, iniObj)
			// fmt.Printf("%s - %#v\n", name, got)
			if !reflect.DeepEqual(tc.want, *got) {
				t.Errorf("name: %s failed \n WANT: %#v \n GOT: %#v", name, tc.want, *got)
			}
		})
	}
}
