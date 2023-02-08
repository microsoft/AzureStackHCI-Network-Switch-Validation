package main

const (
	//Role Types
	MANAGEMENT            = "Management"
	COMPUTEBASIC          = "Compute(Basic)"
	COMPUTESDN            = "Compute(SDN)"
	STORAGE               = "Storage"
	PASS                  = "Pass"
	FAIL                  = "Fail"
	REPORT_SUMMARY_TITTLE = "Supported Role-Types"
	TYPE_SUMMARY_TITTLE   = "Type Result with Log"
	ALL_LOGS              = "ALL Detail Logs"
	INPUT_VARIABLES       = "Input Variables"

	// Output
	CONSOLE_OUTPUT      = "### Validation Summary Result ###"
	GENERATE_PDF_OUTPUT = "Result PDF File Generated, Please Check in Path Folder"

	// Host
	NO_INTF        = "!!! No Interfaces being detected"
	INTF_NOT_MATCH = "[Error] No matched host interface found by IP, please check all host interfaces in menu and update input.ini accordingly."

	// BGP
	BGP = "BGP"

	// VLAN
	VLAN                     = "VLAN"
	INCORRECT_NATIVE_VLAN_ID = "Incorrect Native VLAN ID"
	INCORRECT_VLAN_ID_LIST   = "Incorrect VLAN ID List"

	// DHCPRelay
	DHCPRelay                    = "DHCP Relay Agent IP"
	DHCPPacket_NOT_Detect        = "DHCP Packet Not Detected from switch"
	DHCPRelay_AgentIP_Not_Detect = "DHCP Relay Agent IP Not Detected from switch"

	// DHCPRelay
	BGPPacket_NOT_Detect = "TCP 179 Packet Not Detected from switch"

	// LLDP
	LLDP_Subtype1_PortVLANID = "LLDP-Port VLAN ID (Subtype = 1)"
	LLDP_Subtype3_VLANList   = "LLDP-VLAN Name (Subtype = 3)"
	LLDP_MAXIMUM_FRAME_SIZE  = "LLDP-Maximum Frame Size"
	LLDP_ETS_MAX_CLASSES     = "LLDP-ETS Maximum Number of Traffic Classes"
	LLDP_ETS_BW              = "LLDP-ETS Class Bandwidth Configured"
	LLDP_PFC_MAX_CLASSES     = "LLDP-PFC Maximum Number of Traffic Classes"
	LLDP_PFC_ENABLE          = "LLDP-PFC Priority Class Enabled"
	LLDP_LINK_AGGREGATION    = "LLDP-Link Aggregation"

	NO_LLDP_PACKET                     = "!!! No LLDP Packets Detected from switch via the Interface"
	CHASIS_ID_TYPE                     = "MAC Address"
	NO_LLDP_SYS_DSC                    = "No System Description detected from switch"
	NO_LLDP_CHASSIS_SUBTYPE            = "No Chassis Subtype detected from switch"
	NO_LLDP_PORT_SUBTYPE               = "No Port Subtype detected from switch"
	INCORRECT_LLDP_MAXIMUM_FRAME_SIZE  = "Incorrect Maximum Frame Size"
	INCORRECT_LLDP_Subtype1_PortVLANID = "Incorrect Subtype1 PortVLANID"
	INCORRECT_LLDP_Subtype3_VLANList   = "Incorrect LLDP Subtype3 VLANList"
	INCORRECT_LLDP_ETS_MAX_CLASSES     = "Incorrect ETS Maximum Number of Traffic Classes"
	INCORRECT_LLDP_ETS_BW              = "Incorrect ETS Class Bandwidth Configured"
	INCORRECT_LLDP_PFC_MAX_CLASSES     = "Incorrect PFC Maximum Number of Traffic Classes"
	INCORRECT_LLDP_PFC_ENABLE          = "Incorrect PFC Priority Class Enabled"
	UNSUPPORT_LLDP_LINK_AGGREGATION    = "No Link Aggregation Support"

	// MTU
	INCORRECT_MTU_SIZE = "Incorrect MTU Size"
)
