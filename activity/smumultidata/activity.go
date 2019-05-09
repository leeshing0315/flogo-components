package smumultidata

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

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
	eventTime, _ := context.GetInput("eventTime").(string)

	gpsEventStrs := []string{}
	opModeChangeStrs := []string{}
	deviceErrorsStrs := []string{}
	autoReg := "false"

	packets := splitPackets(reqDataSegment)
	for _, dateSegment := range packets {
		singlePacket, err := service.ParseToSinglePacket(dateSegment)
		if err != nil {
			// TODO save parsing error to DB
			continue
		}
		if singlePacket.LoginItem.ContainerNumber != "" {
			cntrNum = singlePacket.LoginItem.ContainerNumber
			autoReg = "true"
		}
		if singlePacket.LoginItem.DeviceID != "" {
			devId = singlePacket.LoginItem.DeviceID
			autoReg = "true"
		}

		gpsEvent := service.GenGpsEventFromSinglePacket(singlePacket, seqNo, cntrNum, eventTime)
		gpsEventBytes, _ := json.Marshal(gpsEvent)
		gpsEventStrs = append(gpsEventStrs, string(gpsEventBytes))

		opModeChange := service.GenOpModeChangeFromSinglePacket(singlePacket, seqNo, cntrNum)
		if opModeChange != nil {
			opModeChangeBytes, _ := json.Marshal(opModeChange)
			opModeChangeStrs = append(opModeChangeStrs, string(opModeChangeBytes))
		}

		deviceErrors := service.GenDeviceErrorsFromSinglePacket(singlePacket, seqNo, devId)
		for _, val := range deviceErrors {
			deviceErrorsStrs = append(deviceErrorsStrs, val)
		}
	}

	context.SetOutput("cntrNum", cntrNum)
	context.SetOutput("devId", devId)
	context.SetOutput("resDataSegment", []byte{})
	context.SetOutput("autoReg", autoReg)
	println("**********multiData*cntrNum", cntrNum, "**********")
	println("**********multiData*devId", devId, "**********")
	println("**********multiData*autoReg", autoReg, "**********")
	if len(gpsEventStrs) > 0 {
		gpsEventsOutput, _ := json.Marshal(gpsEventStrs)
		context.SetOutput("gpsEvents", string(gpsEventsOutput))
	}
	if len(opModeChangeStrs) > 0 {
		opModeChangesOutput, _ := json.Marshal(opModeChangeStrs)
		context.SetOutput("opModeChanges", string(opModeChangesOutput))
	}
	if len(deviceErrorsStrs) > 0 {
		deviceErrorsOutput, _ := json.Marshal(deviceErrorsStrs)
		context.SetOutput("deviceErrors", string(deviceErrorsOutput))
	}

	return true, nil
}

func splitPackets(data []byte) [][]byte {
	var result = [][]byte{}
	for i := data; len(i) > 1; {
		dataSegmentLen := binary.BigEndian.Uint16(i[3:5])
		dataSegment := i[5 : dataSegmentLen+5]
		result = append(result, dataSegment)

		i = i[dataSegmentLen+7:]
	}
	return result
}
