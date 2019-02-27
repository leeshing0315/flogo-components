package smueventlog

import (
	"encoding/json"
	"time"

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
	cntrNum, _ := context.GetInput("cntrNum").(string)
	seqNo, _ := context.GetInput("seqNo").(int)
	reqDataSegment, _ := context.GetInput("reqDataSegment").([]byte)

	eventLogStrs := []string{}
	gpsEventStrs := []string{}

	packets := splitPackets(reqDataSegment)
	for _, dataSegment := range packets {
		eventLog := service.ParseToEventLog(dataSegment, time.Now(), cntrNum, seqNo)
		if eventLog == nil {
			break
		}
		eventLogBytes, err := json.Marshal(eventLog)
		if err != nil {
			break
		}

		gpsEvent := service.ConvertEventLogToGPSEvent(eventLog)
		gpsEventBytes, err := json.Marshal(gpsEvent)
		if err != nil {
			break
		}

		eventLogStrs = append(eventLogStrs, string(eventLogBytes))
		gpsEventStrs = append(gpsEventStrs, string(gpsEventBytes))
	}

	context.SetOutput("resDataSegment", []byte{})
	if len(eventLogStrs) > 0 {
		eventLogsOutput, _ := json.Marshal(eventLogStrs)
		context.SetOutput("eventLogs", string(eventLogsOutput))

		gpsEventsOutput, _ := json.Marshal(gpsEventStrs)
		context.SetOutput("gpsEvents", string(gpsEventsOutput))
	}

	return true, nil
}

func splitPackets(data []byte) [][]byte {
	var result = [][]byte{}
	for i := data; len(i) > 1; {
		var dataSegmentLen int
		if i[0] == 1 {
			dataSegmentLen = 19
		} else if i[0] == 2 {
			dataSegmentLen = 6
		}
		dataSegment := i[0 : dataSegmentLen+1]
		result = append(result, dataSegment)

		i = i[dataSegmentLen+1:]
	}
	return result
}
