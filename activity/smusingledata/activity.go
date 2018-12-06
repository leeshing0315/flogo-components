package smusingledata

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/common/entity"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	seqNo, _ := context.GetInput("seqNo").(string)
	cntrNum, _ := context.GetInput("cntrNum").(string)
	reqDataSegment, _ := context.GetInput("reqDataSegment").([]byte)

	singlePacket := entity.ParseToSinglePacket(reqDataSegment)
	gpsEvent := entity.GenGpsEventFromSinglePacket(singlePacket, seqNo, cntrNum)
	opModeChange := entity.GenOpModeChangeFromSinglePacket(singlePacket, seqNo, cntrNum)

	println(singlePacket)

	gpsEventStr, _ := json.Marshal(gpsEvent)
	opModeChangeStr, _ := json.Marshal(opModeChange)

	context.SetOutput("resDataSegment", []byte{})
	context.SetOutput("gpsEvent", gpsEventStr)
	context.SetOutput("opModeChange", opModeChangeStr)

	return true, nil
}
