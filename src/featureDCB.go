package main

import (
	"fmt"
	"strings"
)

type PFCType struct {
	PFCMaxClasses      uint8
	PFCPriorityEnabled map[uint8]uint8
}

type ETSType struct {
	ETSTotalPG  uint8
	ETSBWbyPGID map[uint8]uint8
}

func DCBResultValidation() {
	var restultFail []string

	etsBWString := mapintToSlicestring(LLDPResult.ETS.ETSBWbyPGID)
	if etsBWString != INIObj.ETSBWbyPG {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_ETS_BW, INIObj.ETSBWbyPG, etsBWString)
		restultFail = append(restultFail, errMsg)
	}

	pfcEnableString := mapintToSlicestring(LLDPResult.PFC.PFCPriorityEnabled)
	if pfcEnableString != INIObj.PFCPriorityEnabled {
		errMsg := fmt.Sprintf("%s:\n \t\tInput:%s\n \t\tFound: %s", WRONG_LLDP_PFC_ENABLE, INIObj.PFCPriorityEnabled, pfcEnableString)
		restultFail = append(restultFail, errMsg)
	}

	if len(restultFail) == 0 {
		ResultSummary["DCB - PASS"] = restultFail
	} else {
		ResultSummary["DCB - FAIL"] = restultFail
	}
}

func mapintToSlicestring(mapint map[uint8]uint8) string {
	var outputStringSlice []string

	for i := 0; i < 8; i++ {
		outputStringSlice = append(outputStringSlice, fmt.Sprintf("%d:%d", i, mapint[uint8(i)]))
	}

	return strings.Join(outputStringSlice, ",")
}
