package smucntrdevmapping

import (
	"encoding/json"

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
	cntrDevMappingStr := context.GetInput("cntrDevMapping").(string)

	// do eval
	var cntrDevMapping = make(map[string]interface{})
	err = json.Unmarshal([]byte(cntrDevMappingStr), &cntrDevMapping)
	if err != nil {
		return false, err
	}

	context.SetOutput("sim", cntrDevMapping["simno"])
	context.SetOutput("devId", cntrDevMapping["carno"])
	context.SetOutput("cntrNum", cntrDevMapping["carid"])

	println("**********login*sim", cntrDevMapping["simno"].(string), "**********")
	println("**********login*devId", cntrDevMapping["carno"].(string), "**********")
	println("**********login*cntrNum", cntrDevMapping["carid"].(string), "**********")

	return true, nil
}
