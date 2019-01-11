package smugetcntrnum

import (
	"encoding/json"
	"sort"

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
	containersummaryInterface := context.GetInput("containersummaries")

	if containersummaryInterface != nil {
		var containersummaries []entity.ContainerSummary
		err = json.Unmarshal(([]byte)(containersummaryInterface.(string)), &containersummaries)
		if err != nil {
			return true, nil
		}

		sort.Slice(containersummaries, func(i, j int) bool {
			return containersummaries[i].Regtime > containersummaries[j].Regtime
		})
		latestContainerSummary := containersummaries[0]

		println("**********devId", latestContainerSummary.Carno, "**********")
		println("**********cntrNum", latestContainerSummary.Carid, "**********")
		context.SetOutput("devId", latestContainerSummary.Carno)
		context.SetOutput("cntrNum", latestContainerSummary.Carid)
	}
	return true, nil
}
