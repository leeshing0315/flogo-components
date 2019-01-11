package entity

import (
	"strconv"
	"strings"

	"github.com/leeshing0315/flogo-components/common/util"
)

// DeviceConfigCmd entity
type DeviceConfigCmd struct {
	DeviceID       string `json:"devid"`
	SeqNo          string `json:"seqno"`
	Subcmd         string `json:"subcmd"`
	Value          string `json:"value"`
	SendFlag       string `json:"sendflag"`
	SendTime       string `json:"sendtime"`
	LastUpdateTime string `json:"lastupdatetime"`
}

// GenSetConfigCommand for generating Command for setting config
func GenSetConfigCommand(cmdVal map[string]string) (setConfigCommand []byte, err error) {

	powerOnFrequency := make([]byte, 5)
	if val, found := cmdVal["1"]; found {
		powerOnFrequency[0] = 1
		copy(powerOnFrequency[1:], util.GetEndBytes(util.FromDecStrToHexStr(val), 4))
	}

	powerOffFrequency := make([]byte, 5)
	if val, found := cmdVal["2"]; found {
		powerOffFrequency[0] = 1
		copy(powerOffFrequency[1:], util.GetEndBytes(util.FromDecStrToHexStr(val), 4))
	}

	collectFrequency := make([]byte, 5)
	if val, found := cmdVal["3"]; found {
		collectFrequency[0] = 1
		copy(collectFrequency[1:], util.GetEndBytes(util.FromDecStrToHexStr(val), 4))
	}

	serverIPAndPort := make([]byte, 13)
	if val, found := cmdVal["4"]; found {
		serverIPAndPort[0] = 1
		array := strings.Split(val, ":")
		serverIPStr := array[0]
		portStr := array[1]
		ipSegmentStrs := strings.Split(serverIPStr, ".")
		i := 1
		for _, ipSegmentStr := range ipSegmentStrs {
			copy(serverIPAndPort[i:i+2], util.GetEndBytes(util.FromDecStrToHexStr(ipSegmentStr), 2))
			i += 2
		}
		copy(serverIPAndPort[i:i+4], util.GetEndBytes(util.FromDecStrToHexStr(portStr), 4))
	}

	sleepMode := make([]byte, 2)
	if val, found := cmdVal["5"]; found {
		sleepMode[0] = 1
		copy(sleepMode[1:], util.FromStrToUint32(val))
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

// DecodeReadConfigAck for decoding the Read Config Cmd Ack
func DecodeReadConfigAck(original []byte) string {
	var str strings.Builder
	str.WriteString("*L")

	// powerOnFrequency
	str.WriteString(strconv.FormatUint(uint64(original[4]), 10))
	for i := 5; i < 9; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// powerOffFrequency
	str.WriteString(strconv.FormatUint(uint64(original[9]), 10))
	for i := 10; i < 14; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// collectFrequency
	str.WriteString(strconv.FormatUint(uint64(original[14]), 10))
	for i := 15; i < 19; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	// serverIPAndPort
	str.WriteString(strconv.FormatUint(uint64(original[19]), 10))
	for i := 20; i < 32; i++ {
		if strconv.FormatUint(uint64(original[i]), 10) == "0" {
			str.WriteString("0")
		} else {
			current := string(original[i])
			str.WriteString(current)
		}
	}

	//sleepMode
	str.WriteString(strconv.FormatUint(uint64(original[32]), 10))
	str.WriteString(strconv.FormatUint(uint64(original[33]), 10))

	str.WriteString("#")

	return str.String()
}
