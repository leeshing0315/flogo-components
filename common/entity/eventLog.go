package entity

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type EventLog struct {
	Seq     string `json:"SEQ"`
	CntrNum string `json:"CNTR_NUM"`
	RevTime string `json:"REV_TIME"`
	LogTime string `json:"LOG_TIME"`
	Sp      string `json:"SP"`
	Isc     string `json:"ISC"`
	Ss      string `json:"SS"`
	Rs      string `json:"RS"`
	Dss     string `json:"DSS"`
	Drs     string `json:"DRS"`
	Ambs    string `json:"AMBS"`
	Hus     string `json:"HUS"`
	Sh      string `json:"SH"`
	Usda1   string `json:"USDA1"`
	Usda2   string `json:"USDA2"`
	Usda3   string `json:"USDA3"`
	Cts     string `json:"CTS"`
	Smode   string `json:"SMODE"`
	Isa     int32  `json:"ISA"`
}

var smodeMapping = map[int]string{
	0x2001: "POWER OFF",
	0x2002: "UNIT OFF",
	0x2003: "G-SET OFF",
	0x2004: "Battery OFF",
	0x2005: "SPTI START(MODEM)",
	0x2006: "FPTI START(MODEM)",
	0xA001: "POWER ON",
	0xA002: "UNIT ON",
	0xA003: "G-SET ON",
	0xA004: "Battery ON",
	0xA005: "SPTI START(PANEL)",
	0xA006: "FPTI START(PANEL)",
	0xA007: "M.CHECK START(PANEL)",
	0xA008: "MANUAL DEFROST(PANEL)",
	0xA009: "Change Container ID",
	0xB010: "Defrost finish",
	0xB011: "Defrost start",
	0xB040: "Finish PTI with NG result",
	0xB041: "Finish PTI with OK result",
	0xB100: "Set DHU to OFF by PANEL",
	0xB101: "Set DHU to ON by PANEL",
	0xB102: "Set DHU to OFF by MODEM",
	0xB103: "Set DHU to ON by MODEM",
}

func ParseToEventLog(bytes []byte) *EventLog {
	if bytes[0] == 1 {
		return parseTemperatureLog(bytes[1:])
	} else if bytes[0] == 2 {
		return parseSmodeLog(bytes[1:])
	} else {
		return nil
	}
}

func parseTemperatureLog(bytes []byte) *EventLog {
	eventLog := &EventLog{}
	eventLog.LogTime = parseDateTime(bytes[0:4])
	eventLog.Smode = opModeMapping[bytes[4]>>4]
	eventLog.Ss = strconv.FormatFloat(float64((int(bytes[4]&0xf)<<6)+int(bytes[5]>>2))/10.0-40.0, 'f', 1, 64)
	eventLog.Rs = strconv.FormatFloat(float64(((int(bytes[5]&0x3)<<8)+int(bytes[6])))/10.0-40.0, 'f', 1, 64)
	eventLog.Sp = strconv.FormatFloat(float64((int(bytes[7])<<3)+int(bytes[8]>>5))/10.0-40.0, 'f', 1, 64)
	if (bytes[8] & 0x10) == 0 {
		eventLog.Isc = "1"
	} else {
		eventLog.Isc = "0"
	}
	eventLog.Dss = strconv.FormatFloat(float64((int(bytes[8]&0xf)<<6)+int(bytes[9]>>2))/10.0-40.0, 'f', 1, 64)
	eventLog.Drs = strconv.FormatFloat(float64(((int(bytes[9]&0x3)<<8)+int(bytes[10])))/10.0-40.0, 'f', 1, 64)
	eventLog.Ambs = strconv.FormatFloat(float64(bytes[11]&0x7f)-40.0, 'f', 1, 64)
	if (bytes[13] & 0x80) == 0x80 {
		eventLog.Isa = 1
		eventLog.Hus = strconv.FormatInt(int64(bytes[12]&0x7f), 10)
		eventLog.Sh = strconv.FormatInt(int64(bytes[13]&0x7f), 10)
	} else {
		eventLog.Isa = 0
	}
	eventLog.Usda1 = strconv.FormatFloat(float64((int(bytes[18]>>6)<<8)+int(bytes[14]))/10.0-40.0, 'f', 1, 64)
	eventLog.Usda2 = strconv.FormatFloat(float64((int(bytes[18]&0x30)<<4)+int(bytes[15]))/10.0-40.0, 'f', 1, 64)
	eventLog.Usda3 = strconv.FormatFloat(float64((int(bytes[18]&0xc)<<6)+int(bytes[16]))/10.0-40.0, 'f', 1, 64)
	eventLog.Cts = strconv.FormatFloat(float64((int(bytes[18]&0x3)<<8)+int(bytes[17]))/10.0-40.0, 'f', 1, 64)
	return eventLog
}

func parseSmodeLog(bytes []byte) *EventLog {
	eventLog := &EventLog{}
	eventLog.LogTime = parseDateTime(bytes[0:4])
	smode := smodeMapping[int(binary.BigEndian.Uint16(bytes[4:6]))]
	if smode == "" {
		smode = getSmodeByCal(bytes[4], bytes[5])
	}
	eventLog.Smode = smode
	return eventLog
}

func getSmodeByCal(hiByte, loByte byte) string {
	// Change set point
	if hiByte>>4 == 0 {
		var sp string = strconv.FormatFloat(float64((int(hiByte&0x3)<<8)+int(loByte))/10.0-40.0, 'f', 1, 64)
		modeValue := (hiByte & 0xc) >> 2
		var mode string
		if modeValue == 0 {
			mode = "Panel"
		} else if modeValue == 1 {
			mode = "Modem"
		} else {
			mode = "PC"
		}
		return fmt.Sprintf("Set Point %s by %s", sp, mode)
	}
	// Change Defrost Interval
	if (hiByte & 0x30) == 0x30 {
		modeValue := (hiByte & 0xc) >> 2
		var mode string
		if modeValue == 1 {
			mode = "Modem"
		} else if modeValue == 2 {
			mode = "PC"
		} else if modeValue == 3 {
			mode = "CNT"
		} else {
			mode = "PANEL"
		}
		var value string = strconv.FormatInt(int64((loByte&0xf)*3+3), 10)
		return fmt.Sprintf("Change Defrost Interval %s hour by %s", value, mode)
	}
	// Change Set Humidity
	if (hiByte & 0x80) == 0x80 {
		modeValue := (hiByte & 0xc) >> 2
		var mode string
		if modeValue == 0 {
			mode = "Panel"
		} else if modeValue == 1 {
			mode = "Modem"
		} else {
			mode = "PC"
		}
		var value string = strconv.FormatInt(int64(loByte&0x7f), 10)
		return fmt.Sprintf("Change Set Humidity %s%RH by %s", value, mode)
	}
	// Change Set Time
	if hiByte == 0x65 {
		var value string = strconv.FormatInt(int64(loByte)+2000, 10)
		return fmt.Sprintf("Change Time Setting to %s (YEAR)", value)
	} else if hiByte == 0x66 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (MONTH)", value)
	} else if hiByte == 0x67 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (DAY)", value)
	} else if hiByte == 0x68 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (HOUR)", value)
	} else if hiByte == 0x69 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (MINUTE)", value)
	}
	return strconv.FormatUint((uint64(hiByte)<<8)+uint64(loByte), 16)
}

func parseDateTime(bytes []byte) string {
	var year int = int(bytes[0]) + 2000
	var month int = int(bytes[1]&0xf) + 1
	var day int = int((bytes[2]>>3)&0x1f) + 1
	var hour int = int(((bytes[2] & 0x7) << 2) + ((bytes[3] >> 6) & 0x3))
	var minute int = int(bytes[3] & 0x3f)
	yearStr := strconv.FormatInt(int64(year), 10)
	monthStr := strconv.FormatInt(int64(month), 10)
	dayStr := strconv.FormatInt(int64(day), 10)
	hourStr := strconv.FormatInt(int64(hour), 10)
	var minuteStr string
	if minute < 10 {
		minuteStr = "0" + strconv.FormatInt(int64(minute), 10)
	} else {
		minuteStr = strconv.FormatInt(int64(minute), 10)
	}
	// 1/26/2019 6:00
	// 1/28/2019 15:44
	var builder strings.Builder
	builder.WriteString(monthStr)
	builder.WriteString("/")
	builder.WriteString(dayStr)
	builder.WriteString("/")
	builder.WriteString(yearStr)
	builder.WriteString(" ")
	builder.WriteString(hourStr)
	builder.WriteString(":")
	builder.WriteString(minuteStr)
	return builder.String()
}