package smusingledata

import (
	"encoding/json"
	"strconv"

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
	seqNo := strconv.FormatUint(uint64(context.GetInput("seqNo").(int)), 10)
	cntrNum, _ := context.GetInput("cntrNum").(string)
	// devId, _ := context.GetInput("devId").(string)
	reqDataSegment, _ := context.GetInput("reqDataSegment").([]byte)

	singlePacket := entity.ParseToSinglePacket(reqDataSegment)
	if singlePacket.LoginItem.ContainerNumber != "" {
		cntrNum = singlePacket.LoginItem.ContainerNumber
	}
	// if singlePacket.LoginItem.DeviceID != "" {
	// 	devId = singlePacket.LoginItem.DeviceID
	// }

	gpsEvent := entity.GenGpsEventFromSinglePacket(singlePacket, seqNo, cntrNum)
	opModeChange := entity.GenOpModeChangeFromSinglePacket(singlePacket, seqNo, cntrNum)

	gpsEventStr, _ := json.Marshal(gpsEvent)
	opModeChangeStr, _ := json.Marshal(opModeChange)

	context.SetOutput("cntrNum", cntrNum)
	// context.SetOutput("devId", devId)
	context.SetOutput("resDataSegment", []byte{})
	context.SetOutput("gpsEvent", gpsEventStr)
	context.SetOutput("opModeChange", opModeChangeStr)

	return true, nil
}
