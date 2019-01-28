package smucntrdevmapping

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
	cntrDevMappingStr := context.GetInput("cntrDevMapping").(string)

	// do eval
	var cntrDevMapping = entity.ContainerDeviceMapping{}
	err = json.Unmarshal([]byte(cntrDevMappingStr), &cntrDevMapping)
	if err != nil {
		return false, err
	}

	context.SetOutput("sim", cntrDevMapping.Simno)
	context.SetOutput("devId", cntrDevMapping.DeviceId)
	context.SetOutput("cntrNum", cntrDevMapping.ContainerNumber)

	return true, nil
}
