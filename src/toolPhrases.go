package main

const (
	// Output
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
	NO_LLDP_SYS_DSC               = "No System Description Founded"
	NO_LLDP_CHASSIS_SUBTYPE       = "No Chassis Subtype Founded"
	NO_LLDP_PORT_SUBTYPE          = "No Port Subtype Founded"
	WRONG_LLDP_MAXIMUM_FRAME_SIZE = "Incorrect Maximum Frame Size"
	WRONG_LLDP_VLAN_ID            = "Incorrect Port VLAN ID"
	WRONG_LLDP_ETS_MAX_CLASSES    = "Incorrect ETS Maximum Number of Traffic Classes"
	WRONG_LLDP_ETS_BW             = "Incorrect ETS Class Bandwidth Configured"
	WRONG_LLDP_PFC_MAX_CLASSES    = "Incorrect PFC Maximum Number of Traffic Classes"
	WRONG_LLDP_PFC_ENABLE         = "Incorrect PFC Priority Class Enabled"

	// MTU
	WRONG_MTU_SIZE = "Incorrect MTU Size"
)

// func printResultSummary() {
// 	fmt.Println("\n### Result Summary ###")
// 	var resultTemplate = `
// {{range $key, $value := .TasksWithResult}}
// {{- $key -}}
// 	{{range $value}}
//   	- {{. -}}
// 	{{end}}
// {{end}}
// 	`
// 	ret := TemplateResult{
// 		TasksWithResult: OutputObj.ResultSummary,
// 	}

// 	t := template.New("resultTemplate")
// 	t, err := t.Parse(resultTemplate)
// 	if err != nil {
// 		log.Fatalln("parse file: ", err)
// 		return
// 	}

// 	err = t.Execute(os.Stdout, ret)
// 	if err != nil {
// 		log.Fatalln("execute: ", err)
// 		return
// 	}
// }
