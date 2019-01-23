package smusendcommand

import (
	"encoding/json"
	"time"

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
	commands := context.GetInput("commands")
	commandNum := context.GetInput("commandNum")

	readCommandSegment := make([]byte, 0)
	readCommandSeqNo := ""
	setCommandSegment := make([]byte, 0)
	setCommandSeqNo := ""

	if int(commandNum.(float64)) > 0 {
		var commandArr []entity.DeviceConfigCmd
		err = json.Unmarshal([]byte(commands.(string)), &commandArr)
		if err != nil {
			return false, err
		}

		setCmdVal := make(map[string]string)

		for _, command := range commandArr {
			if command.Subcmd == "FF" {
				readCommandSegment = []byte{0x2A, 0x4C, 0x32, 0x23}
				readCommandSeqNo = string(command.SeqNo)
			} else {
				setCmdVal[command.Subcmd] = command.Value
				setCommandSeqNo = string(command.SeqNo)
			}
		}
		if len(setCmdVal) > 0 {
			setCommand, err := entity.GenSetConfigCommand(setCmdVal)
			if err != nil {
				return false, err
			}
			setCommandSegment = setCommand
		}

		valueMap := make(map[string]string)
		valueMap["sendflag"] = "1"
		now := time.Now().Format("2006-01-02 15:04:05")
		valueMap["sendtime"] = now
		valueMap["lastupdatetime"] = now

		jsonBytes, _ := json.Marshal(valueMap)
		context.SetOutput("updateMongoVal", string(jsonBytes))
	}
	context.SetOutput("readCommandSegment", readCommandSegment)
	context.SetOutput("readCommandSeqNo", readCommandSeqNo)
	context.SetOutput("setCommandSegment", setCommandSegment)
	context.SetOutput("setCommandSeqNo", setCommandSeqNo)

	return true, nil
}
