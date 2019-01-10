package smugetcntrnum

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
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
	tc.SetInput("containersummary", `{ 
		"_id" : "5c35ae47d5adca004faa4054", 
		"tableName" : "Tblcarbaseinfo", 
		"simno" : "14540694003", 
		"carno" : "C04254", 
		"carid" : "CXRU1495240", 
		"commmode" : "GPRS", 
		"unitcode" : "00000", 
		"cartype" : "船舶", 
		"saveflag" : "1", 
		"calcflag" : "3", 
		"changeflag" : "1", 
		"changetime" : "2018-08-27 15:08:27.0", 
		"regtime" : "2018-08-27 15:08:27.0", 
		"useacc" : "0", 
		"checkflag" : "1", 
		"boxtype" : "D"
	}`)

	act.Eval(tc)

	//check result attr
	assert.Equal(t, "C04254", tc.GetOutput("devId"))
	assert.Equal(t, "CXRU1495240", tc.GetOutput("cntrNum"))
}

func TestEvalWhenEmpty(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	done, err := act.Eval(tc)

	//check result attr
	assert.Equal(t, true, done)
	assert.Equal(t, nil, err)
}
