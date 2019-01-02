package smucommandack

import (
	"encoding/json"
	"strings"
	"time"

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
	reqDataSegment := context.GetInput("reqDataSegment").([]byte)
	devid := context.GetInput("devid").(string)
	seqNo := context.GetInput("seqNo").(string)

	context.SetOutput("keyName", strings.Join([]string{"devid", "seqNo"}, ","))
	context.SetOutput("keyValue", strings.Join([]string{devid, seqNo}, ","))

	valueMap := make(map[string]string)
	valueMap["sendflag"] = "2"
	valueMap["lastupdatetime"] = time.Now().Format("2006-01-02 15:04:05")

	if len(reqDataSegment) > 1 {
		// read config
		valueMap["value"] = string(reqDataSegment)
	}

	jsonBytes, _ := json.Marshal(valueMap)
	context.SetOutput("updateVal", string(jsonBytes))

	return true, nil
}
