package smugetsim

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
	simcardPinInterface := context.GetInput("simcardpin")

	if simcardPinInterface != nil {
		simcardPin := &entity.SimcardPin{}
		err = json.Unmarshal(([]byte)(simcardPinInterface.(string)), simcardPin)
		if err == nil {
			context.SetOutput("simno", simcardPin.Simno)
		}
	}
	return true, nil
}
