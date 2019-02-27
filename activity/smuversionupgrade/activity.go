package smuversionupgrade

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"encoding/binary"
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
func (a *MyActivity) Eval(context activity.Context) (done bool, err error)  {

	// do eval
	reqDataSegmentBytes := context.GetInput("reqDataSegment").([]byte)

	serialNumber := reqDataSegmentBytes[0]
	contentLength := binary.BigEndian.Uint16(reqDataSegmentBytes[1:3])
	identifier := reqDataSegmentBytes[3:11]

	context.SetOutput("serialNumber", int(serialNumber))
	context.SetOutput("contentLength", int(contentLength))
	context.SetOutput("identifier", string(identifier))
	return true, nil
}
