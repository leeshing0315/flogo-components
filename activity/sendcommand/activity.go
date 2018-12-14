package sendcommand

import (
	"encoding/json"
	"strings"

	"github.com/leeshing0315/flogo-components/common/util"

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
	setDeviceConfigStr := context.GetInput("setDeviceConfig").(string)

	var content []byte
	if setDeviceConfigStr == "" {
		content = []byte{}
	} else {
		setDeviceConfig := &entity.DeviceConfig{}
		err = json.Unmarshal([]byte(setDeviceConfigStr), setDeviceConfig)
		if err != nil {
			return false, err
		}
		powerOnCommunicationFrequency := make([]byte, 5)
		if setDeviceConfig.PowerOnCommunicationFrequency != "" {
			powerOnCommunicationFrequency[0] = 1
			copy(powerOnCommunicationFrequency[1:], util.FromStrToUint32(setDeviceConfig.PowerOnCommunicationFrequency))
		}

		powerOffCommunicationFrequency := make([]byte, 5)
		if setDeviceConfig.PowerOffCommunicationFrequency != "" {
			powerOffCommunicationFrequency[0] = 1
			copy(powerOffCommunicationFrequency[1:], util.FromStrToUint32(setDeviceConfig.PowerOffCommunicationFrequency))
		}

		collectFrequency := make([]byte, 5)
		if setDeviceConfig.CollectFrequency != "" {
			collectFrequency[0] = 1
			copy(collectFrequency[1:], util.FromStrToUint32(setDeviceConfig.CollectFrequency))
		}

		serverIpAndPort := make([]byte, 13)
		if setDeviceConfig.ServerIpAndPort != "" {
			serverIpAndPort[0] = 1
			array := strings.Split(setDeviceConfig.ServerIpAndPort, ":")
			serverIpStr := array[0]
			portStr := array[1]
			ipSegmentStrs := strings.Split(serverIpStr, ".")
			i := 1
			for _, ipSegmentStr := range ipSegmentStrs {
				copy(serverIpAndPort[i:i+2], util.FromDecStrToHexStr(ipSegmentStr))
				i += 2
			}
			copy(serverIpAndPort[i:i+2], util.FromDecStrToHexStr(portStr))
		}

		sleepMode := make([]byte, 2)
		if setDeviceConfig.SleepMode != "" {
			sleepMode[0] = 1
			if setDeviceConfig.SleepMode == "ON" {
				sleepMode[1] = 1
			} else if setDeviceConfig.SleepMode == "OFF" {
				sleepMode[1] = 0
			}
		}

		content = make([]byte, 36)
		copy(content[0:2], []byte{0x32, 0x41}) // * => 0x2A => (2=>0x32, A=>0x41)
		copy(content[2:4], []byte{0x34, 0x43}) // L => 0x4C => (4=>0X34, C=>0x43)
		copy(content[4:9], powerOnCommunicationFrequency)
		copy(content[9:14], powerOffCommunicationFrequency)
		copy(content[14:19], collectFrequency)
		copy(content[19:32], serverIpAndPort)
		copy(content[32:34], sleepMode)
		copy(content[34:36], []byte{0x32, 0x33}) // # => 0x23 => (2=>0x32, 3=>0x33)

		context.SetOutput("setCommandSegment", content)
		context.SetOutput("setCommandSeqNo", setDeviceConfig.SeqNo)
	}

	readDeviceConfigStr := context.GetInput("readDeviceConfig").(string)

	if readDeviceConfigStr == "" {
		context.SetOutput("readCommandSegment", []byte{})
	} else {
		readDeviceConfig := &entity.DeviceConfig{}
		err = json.Unmarshal([]byte(readDeviceConfigStr), readDeviceConfig)
		if err != nil {
			return false, err
		}
		context.SetOutput("readCommandSegment", []byte{0x32, 0x41, 0x34, 0x43, 0x33, 0x32}) // *L2
		context.SetOutput("readCommandSeqNo", readDeviceConfig.SeqNo)
	}
	return true, nil
}
