package smumultidata

import (
	"encoding/binary"
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
	eventTime, _ := context.GetInput("eventTime").(string)

	gpsEventStrs := []string{}
	opModeChangeStrs := []string{}

	packets := splitPackets(reqDataSegment)
	for _, dateSegment := range packets {
		singlePacket := entity.ParseToSinglePacket(dateSegment)
		if singlePacket.LoginItem.ContainerNumber != "" {
			cntrNum = singlePacket.LoginItem.ContainerNumber
		}
		// if singlePacket.LoginItem.DeviceID != "" {
		// 	devId = singlePacket.LoginItem.DeviceID
		// }

		gpsEvent := entity.GenGpsEventFromSinglePacket(singlePacket, seqNo, cntrNum, eventTime)
		opModeChange := entity.GenOpModeChangeFromSinglePacket(singlePacket, seqNo, cntrNum)

		gpsEventBytes, _ := json.Marshal(gpsEvent)
		gpsEventStrs = append(gpsEventStrs, string(gpsEventBytes))
		opModeChangeBytes, _ := json.Marshal(opModeChange)
		opModeChangeStrs = append(opModeChangeStrs, string(opModeChangeBytes))
	}

	context.SetOutput("cntrNum", cntrNum)
	// context.SetOutput("devId", devId)
	context.SetOutput("resDataSegment", []byte{})
	if len(gpsEventStrs) > 0 {
		gpsEventsOutput, _ := json.Marshal(gpsEventStrs)
		context.SetOutput("gpsEvents", gpsEventsOutput)
	}
	if len(opModeChangeStrs) > 0 {
		opModeChangesOutput, _ := json.Marshal(opModeChangeStrs)
		context.SetOutput("opModeChanges", opModeChangesOutput)
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
