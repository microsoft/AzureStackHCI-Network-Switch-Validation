package main

import (
	"reflect"
	"testing"
)

func TestResultAnalysis(t *testing.T) {
	type test struct {
		input    *InputType
		pcapFile string
		want     OutputType
	}
	// testFolder := "./test/"
	testCases := map[string]test{}

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
