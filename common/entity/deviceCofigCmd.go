package entity

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/leeshing0315/flogo-components/common/util"
)

// DeviceConfigCmd entity
type DeviceConfigCmd struct {
	DeviceID       string      `json:"devid"`
	SeqNo          json.Number `json:"seqno"`
	Subcmd         string      `json:"subcmd"`
	Value          string      `json:"value"`
	SendFlag       string      `json:"sendflag"`
	SendTime       string      `json:"sendtime"`
	LastUpdateTime string      `json:"lastupdatetime"`
}

// GenSetConfigCommand for generating Command for setting config
func GenSetConfigCommand(cmdVal map[string]string) (setConfigCommand []byte, err error) {

	powerOnFrequency := make([]byte, 5)
	powerOnFrequency[0] = '0'
	if val, found := cmdVal["01"]; found {
		powerOnFrequency[0] = '1'
		// uintVal, _ := strconv.ParseUint(val, 10, 32)
		// binary.BigEndian.PutUint32(powerOnFrequency[1:], uint32(uintVal))
		copy(powerOnFrequency[1:], util.GetEndBytes(util.FromDecStrToHexStr(val), 4))
	}

	powerOffFrequency := make([]byte, 5)
	powerOffFrequency[0] = '0'
	if val, found := cmdVal["02"]; found {
		powerOffFrequency[0] = '1'
		copy(powerOffFrequency[1:], util.GetEndBytes(util.FromDecStrToHexStr(val), 4))
	}

	collectFrequency := make([]byte, 5)
	collectFrequency[0] = '0'
	if val, found := cmdVal["03"]; found {
		collectFrequency[0] = '1'
		copy(collectFrequency[1:], util.GetEndBytes(util.FromDecStrToHexStr(val), 4))
	}

	serverIPAndPort := make([]byte, 13)
	serverIPAndPort[0] = '0'
	if val, found := cmdVal["04"]; found {
		serverIPAndPort[0] = '1'
		array := strings.Split(val, ":")
		serverIPStr := array[0]
		portStr := array[1]
		ipSegmentStrs := strings.Split(serverIPStr, ".")
		if len(ipSegmentStrs) != 4 {
			return nil, errors.New("ipSegment length wrong")
		}
		i := 1
		for _, ipSegmentStr := range ipSegmentStrs {
			copy(serverIPAndPort[i:i+2], util.GetEndBytes(util.FromDecStrToHexStr(ipSegmentStr), 2))
			i += 2
		}
		copy(serverIPAndPort[i:i+4], util.GetEndBytes(util.FromDecStrToHexStr(portStr), 4))
	}

	sleepMode := make([]byte, 2)
	sleepMode[0] = '0'
	if val, found := cmdVal["05"]; found {
		sleepMode[0] = '1'
		if val == "1" {
			sleepMode[1] = '1'
		} else {
			sleepMode[1] = '0'
		}
	}

	content := make([]byte, 33)
	content[0] = 0x2A
	content[1] = 0x4C
	copy(content[2:7], powerOnFrequency)
	copy(content[7:12], powerOffFrequency)
	copy(content[12:17], collectFrequency)
	copy(content[17:30], serverIPAndPort)
	copy(content[30:32], sleepMode)
	content[32] = 0x23

	return content, nil
}

// DecodeReadConfigAck for decoding the Read Config Cmd Ack
func DecodeReadConfigAck(original []byte) string {
	var str strings.Builder
	str.WriteString("*L")

	// powerOnFrequency
	str.WriteString(strconv.FormatUint(uint64(original[2]), 10))
	for i := 3; i < 7; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// powerOffFrequency
	str.WriteString(strconv.FormatUint(uint64(original[7]), 10))
	for i := 8; i < 12; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// collectFrequency
	str.WriteString(strconv.FormatUint(uint64(original[12]), 10))
	for i := 13; i < 17; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// serverIPAndPort
	str.WriteString(strconv.FormatUint(uint64(original[17]), 10))
	for i := 18; i < 30; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	//sleepMode
	str.WriteString(strconv.FormatUint(uint64(original[30]), 10))
	str.WriteString(strconv.FormatUint(uint64(original[31]), 10))

	str.WriteString("#")

	return str.String()
}
