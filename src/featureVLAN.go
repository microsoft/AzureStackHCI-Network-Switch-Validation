package main

import (
	"fmt"
	"sort"
)

type VLANResultType struct {
	NativeVlanID int
	AllVlanIDs   []int
}

func (o *OutputType) VLANResultValidation(v *VLANResultType, i *InputType) {
	v.NativeVlanID = NativeVLANID
	v.AllVlanIDs = VLANIDList

	var VLANReportType FeatureResultType

	VLANReportType.FeatureName = VLAN
	if len(v.AllVlanIDs) == 0 {
		errMsg := VLAN_NOT_DETECT
		VLANReportType.FeaturePass = FAIL
		VLANReportType.FeatureLog = errMsg
	} else if len(v.AllVlanIDs) < 10 {
		errMsg := fmt.Sprintf("%s - Detect: %d", VLAN_MINIMUM_10_ERROR, v.AllVlanIDs)
		VLANReportType.FeaturePass = FAIL
		VLANReportType.FeatureLog = errMsg
	} else if EqualArray(v.AllVlanIDs, i.AllVlanIDs) {
		errMsg := fmt.Sprintf("%s - Detect: %d, but Input: %d", VLAN_MISMATCH, v.AllVlanIDs, i.AllVlanIDs)
		VLANReportType.FeaturePass = FAIL
		VLANReportType.FeatureLog = errMsg
	} else {
		VLANReportType.FeaturePass = PASS
	}
	VLANReportType.FeatureRoles = []string{MANAGEMENT, COMPUTEBASIC, COMPUTESDN, STORAGE}
	o.FeatureResultList = append(o.FeatureResultList, VLANReportType)
}

func EqualArray(a, b []int) bool {
	sort.Ints(a)
	sort.Ints(b)
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
