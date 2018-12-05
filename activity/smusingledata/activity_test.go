package smusingledata

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	act.Eval(tc)

	//check result attr
}

func TestHandleData(t *testing.T) {
	dataSegment := []byte{0, 0, 0, 114, 24, 17, 33, 23, 67, 81, 1, 215, 31, 42, 7, 68, 30, 6, 0, 1, 0, 40, 100, 4, 2, 15, 10, 1, 33, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 67, 48, 48, 48, 48, 49, 83, 77, 85, 84, 48, 48, 48, 48, 48, 48, 49, 68, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80}
	println(handleData(dataSegment))
}
