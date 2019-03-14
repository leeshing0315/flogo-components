package smuversionupgrade

import (
	"encoding/binary"

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
	reqDataSegmentBytes := context.GetInput("reqDataSegment").([]byte)

	serialNumber := reqDataSegmentBytes[2]
	contentLength := binary.BigEndian.Uint16(reqDataSegmentBytes[3:5])
	identifier := reqDataSegmentBytes[5:13]

	context.SetOutput("serialNumber", int(serialNumber))
	context.SetOutput("contentLength", int(contentLength))
	context.SetOutput("identifier", string(identifier))
	return true, nil
}
