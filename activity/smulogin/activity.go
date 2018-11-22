package smulogin

import (
	"strconv"

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
	// eventTime, _ := context.GetInput("eventTime").(string)
	// ip, _ := context.GetInput("ip").(string)
	reqDataSegment, _ := context.GetInput("reqDataSegment").([]byte)

	parseDataSegment(reqDataSegment)

	context.SetOutput("resDataSegment", []byte{})

	return true, nil
}

func parseDataSegment(data []byte) {
	var cursor int
	println("SMU Type: " + "0x" + strconv.FormatUint(uint64(data[cursor]), 16))
	cursor++
	pinLen := int(data[cursor : cursor+1][0])
	cursor++
	println("Pin: " + string(data[cursor:cursor+pinLen]))
	cursor += pinLen
	terminalNumLen := int(data[cursor : cursor+1][0])
	cursor++
	println("TerminalNum: " + string(data[cursor:cursor+terminalNumLen]))
	cursor += terminalNumLen
	hardwareVerLen := int(data[cursor : cursor+1][0])
	cursor++
	println("HardwareVer: " + string(data[cursor:cursor+hardwareVerLen]))
}
