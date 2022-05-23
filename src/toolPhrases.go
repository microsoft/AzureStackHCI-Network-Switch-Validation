package main

const (
	// Output
	CONSOLE_OUTPUT      = "### Validation Summary Result ###"
	GENERATE_PDF_OUTPUT = "Result PDF File Generated"

	// Host
	INTF_NOT_MATCH = "[Error] No matched host interface founded by IP, please check all host interfaces in menu and update input.ini accordingly."

	// VLAN
	VLAN_NOT_MATCH = "VLAN Not Match"

	// DHCPRelay
	DHCPPacket_NOT_Detect        = "DHCP Packet Not Detected"
	DHCPRelay_AgentIP_Not_Detect = "DHCP Relay Agent IP Not Detected"

	// DHCPRelay
	BGPPacket_NOT_Detect = "TCP 179 Packet Not Detected"

	// LLDP
	CHASIS_ID_TYPE                = "MAC Address"
	NO_LLDP_SYS_DSC               = "No System Description Founded"
	NO_LLDP_CHASSIS_SUBTYPE       = "No Chassis Subtype Founded"
	NO_LLDP_PORT_SUBTYPE          = "No Port Subtype Founded"
	WRONG_LLDP_MAXIMUM_FRAME_SIZE = "Incorrect Maximum Frame Size"
	WRONG_LLDP_VLAN_ID            = "Incorrect Port VLAN ID"
	NO_LLDP_IEEE_8021_Subtype3    = "No LLDP IEEE 802.1 VLAN Name (Subtype 3) Founded"
	WRONG_LLDP_ETS_MAX_CLASSES    = "Incorrect ETS Maximum Number of Traffic Classes"
	WRONG_LLDP_ETS_BW             = "Incorrect ETS Class Bandwidth Configured"
	WRONG_LLDP_PFC_MAX_CLASSES    = "Incorrect PFC Maximum Number of Traffic Classes"
	WRONG_LLDP_PFC_ENABLE         = "Incorrect PFC Priority Class Enabled"

	// MTU
	WRONG_MTU_SIZE = "Incorrect MTU Size"
)
