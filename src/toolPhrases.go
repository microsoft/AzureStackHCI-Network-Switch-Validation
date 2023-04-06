package main

const (
	//Role Types
	MANAGEMENT             = "Management"
	COMPUTEBASIC           = "Compute (Standard)"
	COMPUTESDN             = "Compute (SDN)"
	STORAGE                = "Storage"
	PASS                   = "Pass"
	FAIL                   = "Fail"
	ROLE_SUMMARY_TITTLE    = "Role-type Support for Azure Stack HCI:"
	FEATURE_SUMMARY_TITTLE = "Role-Type Test Result List:"
	LOG_DETAIL_BELOW       = "For a more detailed log, please refer to the next section below"
	ALL_LOGS               = "ALL Detail Logs"
	INPUT_VARIABLES        = "Input Variables"
	GENERATE_REPORT_FILES  = "Report Files have been generated."

	// BGP
	BGP                  = "BGP"
	BGPPacket_NOT_Detect = "TCP 179 Packet Not Detected from switch, please check switch BGP configuration."

	// VLAN
	VLAN                  = "VLAN"
	VLAN_NOT_DETECT       = "No VLAN detected, please check the VLAN configuration on switch"
	VLAN_MISMATCH         = "VLAN Mismatch"
	VLAN_MINIMUM_10_ERROR = "Minimum of 10 VLANs need to be input and configured"

	// DHCPRelay
	DHCPRelay                    = "DHCP - Relay Agent IP"
	DHCPRelay_AgentIP_Not_Detect = "DHCP Relay Agent IP Not Detected from switch, please check switch dhcp configuration."

	// LLDP
	LLDP_Subtype1_PortVLANID = "LLDP - Port VLAN ID (Subtype = 1)"
	LLDP_Subtype1_NOT_DETECT = "LLDP Subtype 1 not detected from switch"
	LLDP_Subtype1_MISMATCH   = "LLDP Subtype 1 Mismatch, please check switch VLAN configuration"

	LLDP_Subtype3_VLANList   = "LLDP - VLAN Name (Subtype = 3)"
	LLDP_Subtype3_NOT_DETECT = "LLDP Subtype 3 not detected from switch"
	LLDP_Subtype3_MISMATCH   = "LLDP Subtype 3 Mismatch, please check switch VLAN configuration"

	LLDP_Subtype4_MAX_FRAME_SIZE = "LLDP - Maximum Frame Size (Subtype = 4)"
	LLDP_Subtype4_NOT_DETECT     = "LLDP Subtype 4 not detected from switch"
	LLDP_Subtype4_MISMATCH       = "LLDP Subtype 4 Mismatch, please check switch MTU configuration"

	LLDP_Subtype7_LINK_AGGREGATION = "LLDP - Link Aggregation (Subtype = 7)"
	LLDP_Subtype7_NOT_DETECT       = "LLDP Subtype 7 not detected from switch"

	LLDP_Subtype9_ETS_MAX_CLASSES            = "LLDP - ETS Maximum Number of Traffic Classes (Subtype = 9)"
	LLDP_Subtype9_ETS_MAX_CLASSES_NOT_DETECT = "ETS Maximum Number of Traffic Classes not detected from switch"
	LLDP_Subtype9_ETS_MAX_CLASSES_MISMATCH   = "ETS Maximum Number of Traffic Classes Mismatch"

	LLDP_Subtype9_ETS_BW          = "LLDP - ETS Class Bandwidth Configuration (Subtype = 9)"
	LLDP_Subtype9_ETS_BW_MISMATCH = "Priority 0~7 Mismatch"

	LLDP_SubtypeB_PFC_MAX_CLASSES            = "LLDP - PFC Maximum Number of Traffic Classes (Subtype = B)"
	LLDP_SubtypeB_PFC_MAX_CLASSES_NOT_DETECT = "PFC Maximum Number of Traffic Classes not detected from switch"
	LLDP_SubtypeB_PFC_MAX_CLASSES_MISMATCH   = "PFC Maximum Number of Traffic Classes Mismatch"

	LLDP_SubtypeB_PFC_ENABLE          = "LLDP - PFC Priority Class Enabled (Subtype = B)"
	LLDP_SubtypeB_PFC_ENABLE_MISMATCH = "Priority 0~7 Mismatch"
)
