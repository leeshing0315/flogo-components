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

	devtype, pin, terminalNum, firmwareVersion := parseDataSegment(reqDataSegment)

	context.SetOutput("resDataSegment", []byte{})
	context.SetOutput("devtype", devtype)
	context.SetOutput("pin", pin)
	context.SetOutput("terminalNum", terminalNum)
	context.SetOutput("firmwareVersion", firmwareVersion)

	return true, nil
}

func parseDataSegment(data []byte) (devtype, pin, terminalNum, firmwareVersion string) {
	var cursor int
	devtype = strconv.FormatUint(uint64(data[cursor]), 16)
	println("SMU Type: " + "0x" + devtype)
	cursor++
	pinLen := int(data[cursor : cursor+1][0])
	cursor++
	pin = string(data[cursor : cursor+pinLen])
	println("Pin: " + pin)
	cursor += pinLen
	terminalNumLen := int(data[cursor : cursor+1][0])
	cursor++
	terminalNum = string(data[cursor : cursor+terminalNumLen])
	println("TerminalNum: " + terminalNum)
	cursor += terminalNumLen
	hardwareVerLen := int(data[cursor : cursor+1][0])
	cursor++
	firmwareVersion = string(data[cursor : cursor+hardwareVerLen])
	println("HardwareVer: " + firmwareVersion)

	return devtype, pin, terminalNum, firmwareVersion
}
