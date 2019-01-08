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
	devId, _ := context.GetInput("devid").(string)
	reqDataSegment, _ := context.GetInput("reqDataSegment").([]byte)
	eventTime := context.GetInput("eventTime").(string)

	singlePacket := entity.ParseToSinglePacket(reqDataSegment)
	if singlePacket.LoginItem.ContainerNumber != "" {
		cntrNum = singlePacket.LoginItem.ContainerNumber
	}
	if singlePacket.LoginItem.DeviceID != "" {
		devId = singlePacket.LoginItem.DeviceID
	}

	context.SetOutput("cntrNum", cntrNum)
	context.SetOutput("devId", devId)
	context.SetOutput("resDataSegment", []byte{})

	gpsEvent := entity.GenGpsEventFromSinglePacket(singlePacket, seqNo, cntrNum, eventTime)
	gpsEventBytes, _ := json.Marshal(gpsEvent)
	context.SetOutput("gpsEvent", string(gpsEventBytes))

	opModeChange := entity.GenOpModeChangeFromSinglePacket(singlePacket, seqNo, cntrNum)
	if opModeChange != nil {
		opModeChangeBytes, _ := json.Marshal(opModeChange)
		context.SetOutput("opModeChange", string(opModeChangeBytes))
	}

	deviceErrors := entity.GenDeviceErrorsFromSinglePacket(singlePacket, seqNo, devId)
	deviceErrorsBytes, _ := json.Marshal(deviceErrors)
	context.SetOutput("deviceError", string(deviceErrorsBytes))

	return true, nil
}
