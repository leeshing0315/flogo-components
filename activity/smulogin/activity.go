package smulogin

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/trigger/tcpreceiver"
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
	packet, _ := context.GetInput("packet").(*tcpreceiver.BinPacket)
	// ip, _ := ctx.GetInput("ip").(string)

	result := &tcpreceiver.BinPacket{
		Command:  packet.Command,
		Sequence: packet.Sequence,
	}

	var data = packet.DataSegment

	var cursor int
	println("SMU Type: " + string(data[cursor]))
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

	result.DataSegmentLength = make([]byte, 2)

	ctx.SetOutput("packet", &result)

	return true, nil
}
