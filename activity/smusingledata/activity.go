package smusingledata

import (
	"encoding/json"
	"strconv"

	"github.com/leeshing0315/flogo-components/common/util"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/common/service"
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
	pin := context.GetInput("pin").(string)
	originalPacket := context.GetInput("originalPacket").([]byte)
	carrier := context.GetInput("carrier").(string)

	defer func() {
		if r := recover(); r != nil {
			util.LogDownException(eventTime, pin, originalPacket, r)
		}
	}()

	autoReg := "false"

	singlePacket, err := service.ParseToSinglePacket(reqDataSegment)
	if err != nil {
		// save parsing error to DB
		util.LogDownException(eventTime, pin, originalPacket, err)
		context.SetOutput("cntrNum", cntrNum)
		context.SetOutput("devId", devId)
		context.SetOutput("resDataSegment", []byte{})
		context.SetOutput("autoReg", autoReg)
		return true, nil
	}
	if singlePacket.LoginItem.ContainerNumber != "" {
		cntrNum = singlePacket.LoginItem.ContainerNumber
		autoReg = "true"
	}
	if singlePacket.LoginItem.DeviceID != "" {
		devId = singlePacket.LoginItem.DeviceID
		autoReg = "true"
	}

	context.SetOutput("cntrNum", cntrNum)
	context.SetOutput("devId", devId)
	context.SetOutput("resDataSegment", []byte{})
	context.SetOutput("autoReg", autoReg)
	println("**********singleData*cntrNum", cntrNum, "**********")
	println("**********singleData*devId", devId, "**********")
	println("**********singleData*autoReg", autoReg, "**********")

	gpsEvent := service.GenGpsEventFromSinglePacket(singlePacket, seqNo, cntrNum, eventTime, carrier)
	gpsEventBytes, _ := json.Marshal(gpsEvent)
	context.SetOutput("gpsEvent", string(gpsEventBytes))

	opModeChange := service.GenOpModeChangeFromSinglePacket(singlePacket, seqNo, cntrNum)
	if opModeChange != nil {
		opModeChangeBytes, _ := json.Marshal(opModeChange)
		context.SetOutput("opModeChange", string(opModeChangeBytes))
	}

	deviceErrors := service.GenDeviceErrorsFromSinglePacket(singlePacket, seqNo, devId)
	deviceErrorsBytes, _ := json.Marshal(deviceErrors)
	context.SetOutput("deviceError", string(deviceErrorsBytes))

	return true, nil
}
