package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/leeshing0315/flogo-components/common/util"
)

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
	sleepModeVal, sleepModeFound := cmdVal["05"]
	opModeFilterVal, opModeFilterFound := cmdVal["05-2"]
	if sleepModeFound && !opModeFilterFound {
		if sleepModeVal == "1" {
			// SLEEP MODE ON
			sleepMode[0] = '1'
			sleepMode[1] = '0'
		} else {
			// SLEEP MODE OFF
			sleepMode[0] = '1'
			sleepMode[1] = '1'
		}
	} else if opModeFilterFound && !sleepModeFound {
		if opModeFilterVal == "1" {
			// OPMODE FILTER ON
			sleepMode[0] = '1'
			sleepMode[1] = '2'
		} else {
			// OPMODE FILTER OFF
			sleepMode[0] = '1'
			sleepMode[1] = '3'
		}
	} else if sleepModeFound && opModeFilterFound {
		if sleepModeVal == "1" && opModeFilterVal == "1" {
			// SLEEP MODE ON & OPMODE FILTER ON
			sleepMode[0] = '1'
			sleepMode[1] = '4'
		} else if sleepModeVal == "1" && opModeFilterVal == "0" {
			// SLEEP MODE ON & OPMODE FILTER OFF
			sleepMode[0] = '1'
			sleepMode[1] = '5'
		} else if sleepModeVal == "0" && opModeFilterVal == "1" {
			// SLEEP MODE OFF & OPMODE FILTER ON
			sleepMode[0] = '1'
			sleepMode[1] = '6'
		} else if sleepModeVal == "0" && opModeFilterVal == "0" {
			// SLEEP MODE OFF & OPMODE FILTER OFF
			sleepMode[0] = '1'
			sleepMode[1] = '7'
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
	if original[2] == 49 {
		str.WriteString("1")
	} else {
		str.WriteString("1")
	}
	for i := 3; i < 7; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// powerOffFrequency
	if original[7] == 49 {
		str.WriteString("1")
	} else {
		str.WriteString("1")
	}
	for i := 8; i < 12; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// collectFrequency
	if original[12] == 49 {
		str.WriteString("1")
	} else {
		str.WriteString("1")
	}
	for i := 13; i < 17; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// serverIPAndPort
	if original[17] == 49 {
		str.WriteString("1")
	} else {
		str.WriteString("1")
	}
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
