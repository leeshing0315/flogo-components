package entity

import (
	"strings"

	"github.com/leeshing0315/flogo-components/common/util"
)

// DeviceConfigCmd entity
type DeviceConfigCmd struct {
	CntrNum      string       `json:"cntrNum"`
	DeviceID     string       `json:"devid"`
	SeqNo        string       `json:"seqNo"`
	CommandValue CommandValue `json:"commandValue"`
	CommandType  string       `json:"commandType"`
	Status       string       `json:"status"`
	CreatedBy    string       `json:"createdBy"`
	UpdateBy     string       `json:"updateBy"`
	LastUpdated  string       `json:"lastUpdated"`
}

// CommandValue entity
type CommandValue struct {
	PowerOnFrequency  string `json:"powerOnFrequency"`
	PowerOffFrequency string `json:"powerOffFrequency"`
	CollectFrequency  string `json:"collectFrequency"`
	ServerIPAndPort   string `json:"serverIPAddress"`
	SleepMode         string `json:"sleepMode"`
}

// GenSetConfigCommand fot generating Command for setting config
func GenSetConfigCommand(command *DeviceConfigCmd) (setConfigCommand []byte, err error) {

	commandValue := command.CommandValue

	powerOnFrequency := make([]byte, 5)
	if commandValue.PowerOnFrequency != "" {
		powerOnFrequency[0] = 1
		copy(powerOnFrequency[1:], util.FromStrToUint32(commandValue.PowerOnFrequency))
	}

	powerOffFrequency := make([]byte, 5)
	if commandValue.PowerOffFrequency != "" {
		powerOffFrequency[0] = 1
		copy(powerOffFrequency[1:], util.FromStrToUint32(commandValue.PowerOffFrequency))
	}

	collectFrequency := make([]byte, 5)
	if commandValue.CollectFrequency != "" {
		collectFrequency[0] = 1
		copy(collectFrequency[1:], util.FromStrToUint32(commandValue.CollectFrequency))
	}

	serverIPAndPort := make([]byte, 13)
	if commandValue.ServerIPAndPort != "" {
		serverIPAndPort[0] = 1
		array := strings.Split(commandValue.ServerIPAndPort, ":")
		serverIPStr := array[0]
		portStr := array[1]
		ipSegmentStrs := strings.Split(serverIPStr, ".")
		i := 1
		for _, ipSegmentStr := range ipSegmentStrs {
			copy(serverIPAndPort[i:i+2], util.FromDecStrToHexStr(ipSegmentStr))
			i += 2
		}
		copy(serverIPAndPort[i:i+2], util.FromDecStrToHexStr(portStr))
	}

	sleepMode := make([]byte, 2)
	if commandValue.SleepMode != "" {
		sleepMode[0] = 1
		if commandValue.SleepMode == "ON" {
			sleepMode[1] = 1
		} else if commandValue.SleepMode == "OFF" {
			sleepMode[1] = 0
		}
	}

	content := make([]byte, 36)
	copy(content[0:2], []byte{0x32, 0x41}) // * => 0x2A => (2=>0x32, A=>0x41)
	copy(content[2:4], []byte{0x34, 0x43}) // L => 0x4C => (4=>0X34, C=>0x43)
	copy(content[4:9], powerOnFrequency)
	copy(content[9:14], powerOffFrequency)
	copy(content[14:19], collectFrequency)
	copy(content[19:32], serverIPAndPort)
	copy(content[32:34], sleepMode)
	copy(content[34:36], []byte{0x32, 0x33}) // # => 0x23 => (2=>0x32, 3=>0x33)

	return content, nil
}
