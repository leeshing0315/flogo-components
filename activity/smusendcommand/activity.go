package smusendcommand

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
	deviceCommand := &entity.DeviceCommand{}
	commands := context.GetInput("commands")

	var commandArr []deviceCommand
	err = json.Unmarshal([]byte(commands.(string)), &commandArr)
	if err != nil {
		return false, err
	}

	for _, command := range commandArr {
		commandType := command.CommandType
		if commandType == "READ" {
			context.SetOutput("readCommandSegment", []byte{0x32, 0x41, 0x34, 0x43, 0x33, 0x32})
			context.SetOutput("readCommandSeqNo", command.seqNum)

		} else if commandType == "WRITE" {
			setCommand, err := entity.getSetConfigCommand(command)
			context.SetOutput("setCommandSegment", setCommand)
			context.SetOutput("setCommandSeqNo", command.seqNum)
		}
	}

	return true, nil
}
