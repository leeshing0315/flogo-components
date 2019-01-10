package smugetcntrnum

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
	containersummaryInterface := context.GetInput("containersummary")

	if containersummaryInterface != nil {
		containersummary := &entity.ContainerSummary{}
		err = json.Unmarshal(([]byte)(containersummaryInterface.(string)), containersummary)
		if err == nil {
			context.SetOutput("devId", containersummary.Carno)
			context.SetOutput("cntrNum", containersummary.Carid)
		}
	}
	return true, nil
}
