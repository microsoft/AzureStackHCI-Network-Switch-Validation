package main

const (
	// Output
	CONSOLE_OUTPUT      = "### Validation Summary Result ###"
	GENERATE_PDF_OUTPUT = "Result PDF File Generated"

	// Host
	NO_INTF        = "!!! No Interfaces being detected"
	INTF_NOT_MATCH = "[Error] No matched host interface found by IP, please check all host interfaces in menu and update input.ini accordingly."

	// VLAN
	INCORRECT_NATIVE_VLAN_ID = "Incorrect Native VLAN ID"
	INCORRECT_VLAN_ID_LIST   = "Incorrect VLAN ID List"

	// DHCPRelay
	DHCPPacket_NOT_Detect        = "DHCP Packet Not Detected from switch"
	DHCPRelay_AgentIP_Not_Detect = "DHCP Relay Agent IP Not Detected from switch"

	// DHCPRelay
	BGPPacket_NOT_Detect = "TCP 179 Packet Not Detected from switch"

	// LLDP
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
	UNSUPPORT_LINK_AGGREGATION         = "No Link Aggregation Support"

	// MTU
	INCORRECT_MTU_SIZE = "Incorrect MTU Size"
)
