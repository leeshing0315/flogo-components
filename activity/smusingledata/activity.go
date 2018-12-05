package smusingledata

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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
	seqNo, _ := context.GetInput("reqDataSegment").(int)
	reqDataSegment, _ := context.GetInput("reqDataSegment").([]byte)

	gpsEvent, deviceError, operationModeChange, cntrNum, devid := handleData(seqNo, reqDataSegment)

	context.SetOutput("resDataSegment", []byte{})
	context.SetOutput("output", "")

	return true, nil
}
