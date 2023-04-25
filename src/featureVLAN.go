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
		VLANReportType.FeaturePass = FAIL
		VLANReportType.FeatureLogSubject = VLAN_NOT_DETECT
	} else if len(v.AllVlanIDs) < 10 {
		VLANReportType.FeaturePass = FAIL
		VLANReportType.FeatureLogSubject = VLAN_MINIMUM_10_ERROR
		errMsg := fmt.Sprintf("Detect: %d", v.AllVlanIDs)
		VLANReportType.FeatureLogDetail = errMsg
	} else if !EqualArray(v.AllVlanIDs, i.AllVlanIDs) {
		VLANReportType.FeaturePass = FAIL
		VLANReportType.FeatureLogSubject = VLAN_MISMATCH
		errMsg := fmt.Sprintf("Detect: %d, but Input: %d", v.AllVlanIDs, i.AllVlanIDs)
		VLANReportType.FeatureLogDetail = errMsg
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
	} else {
		for i := range a {
			fmt.Println(a[i], b[i])
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}
