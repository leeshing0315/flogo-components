package smucommandack

import (
	"encoding/json"
	"strings"

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
	eventTime := context.GetInput("eventTime").(string)
	reqDataSegment := context.GetInput("reqDataSegment").([]byte)
	cntrNum := context.GetInput("cntrNum").(string)
	devid := context.GetInput("devid").(string)
	seqNo := context.GetInput("seqNo").(string)

	context.SetOutput("keyName", strings.Join([]string{"cntrNum", "devid", "seqNo"}, ","))
	context.SetOutput("keyValue", strings.Join([]string{cntrNum, devid, seqNo}, ","))

	if len(reqDataSegment) > 1 {
		context.SetOutput("collectionName", "deviceConfig")
		context.SetOutput("method", "INSERT")
	} else {
		context.SetOutput("collectionName", "deviceCommand")
		context.SetOutput("method", "UPDATE")

		valueMap := make(map[string]string)
		valueMap["status"] = "Received"
		valueMap["updateTime"] = eventTime
		jsonBytes, _ := json.Marshal(valueMap)
		context.SetOutput("value", string(jsonBytes))
	}

	return true, nil
}
