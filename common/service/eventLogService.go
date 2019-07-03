package service

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/leeshing0315/flogo-components/common/entity"
)

func ParseToEventLog(bytes []byte, now time.Time, cntrNum string, seqNo int) *entity.EventLog {
	var eventLog *entity.EventLog
	if bytes[0] == 1 {
		eventLog = parseTemperatureLog(bytes[1:])
	} else if bytes[0] == 2 {
		eventLog = parseSmodeLog(bytes[1:])
	} else {
		return nil
	}
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	eventLog.RevTime = now.In(loc).Format("2006-01-02T15:04:05+08:00")
	eventLog.CntrNum = cntrNum
	eventLog.Seq = strconv.FormatInt(int64(seqNo), 10)
	eventLog.Source = "TCP_SERVER"
	return eventLog
}

func parseTemperatureLog(bytes []byte) *entity.EventLog {
	eventLog := &entity.EventLog{}
	eventLog.LogTime = parseDateTime(bytes[0:4])
	eventLog.Smode = opModeMapping[bytes[4]>>4]
	eventLog.Ss = roundFloat(float64((int(bytes[4]&0xf)<<6)+int(bytes[5]>>2))/10.0 - 40.0)
	eventLog.Rs = roundFloat(float64(((int(bytes[5]&0x3)<<8)+int(bytes[6])))/10.0 - 40.0)
	eventLog.Sp = roundFloat(float64((int(bytes[7])<<3)+int(bytes[8]>>5))/10.0 - 40.0)
	if (bytes[8] & 0x10) == 0 {
		eventLog.Isc = 1
	} else {
		eventLog.Isc = 0
	}
	eventLog.Dss = roundFloat(float64((int(bytes[8]&0xf)<<6)+int(bytes[9]>>2))/10.0 - 40.0)
	eventLog.Drs = roundFloat(float64(((int(bytes[9]&0x3)<<8)+int(bytes[10])))/10.0 - 40.0)
	eventLog.Ambs = roundFloat(float64(bytes[11]&0x7f) - 40.0)
	if (bytes[13] & 0x80) == 0x80 {
		// eventLog.Isa = 1
		eventLog.Hus = strconv.FormatInt(int64(bytes[12]&0x7f), 10)
		eventLog.Sh = strconv.FormatInt(int64(bytes[13]&0x7f), 10)
	} else {
		// eventLog.Isa = 0
	}
	eventLog.Isa = 1
	eventLog.Usda1 = strconv.FormatFloat(float64((int(bytes[18]>>6)<<8)+int(bytes[14]))/10.0-40.0, 'f', 1, 64)
	eventLog.Usda2 = strconv.FormatFloat(float64((int(bytes[18]&0x30)<<4)+int(bytes[15]))/10.0-40.0, 'f', 1, 64)
	eventLog.Usda3 = strconv.FormatFloat(float64((int(bytes[18]&0xc)<<6)+int(bytes[16]))/10.0-40.0, 'f', 1, 64)
	eventLog.Cts = strconv.FormatFloat(float64((int(bytes[18]&0x3)<<8)+int(bytes[17]))/10.0-40.0, 'f', 1, 64)
	eventLog.IsTemperatureLog = true
	return eventLog
}

func parseSmodeLog(bytes []byte) *entity.EventLog {
	eventLog := &entity.EventLog{}
	eventLog.LogTime = parseDateTime(bytes[0:4])
	smode := smodeMapping[int(binary.BigEndian.Uint16(bytes[4:6]))]
	if smode == "" {
		var isSmodeParsingFail bool
		smode, isSmodeParsingFail = getSmodeByCal(bytes[4], bytes[5])
		if isSmodeParsingFail {
			eventLog.Isa = 0
		} else {
			eventLog.Isa = 1
		}
	} else {
		eventLog.Isa = 1
	}
	eventLog.Smode = smode
	eventLog.IsTemperatureLog = false
	return eventLog
}

func getSmodeByCal(hiByte, loByte byte) (opMode string, isSmodeParsingFail bool) {
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
		return fmt.Sprintf("Set Point %s by %s", sp, mode), false
	}
	// Change Defrost Interval
	if ((hiByte >> 4) == 3) && (hiByte != 0x31) {
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
		return fmt.Sprintf("Change Defrost Interval %s hour by %s", value, mode), false
	}
	// Change Set Humidity
	if (hiByte >> 4) == 0x8 {
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
		return fmt.Sprintf("Change Set Humidity %s%%RH by %s", value, mode), false
	}
	// Change Set Time
	if hiByte == 0x65 {
		var value string = strconv.FormatInt(int64(loByte)+2000, 10)
		return fmt.Sprintf("Change Time Setting to %s (YEAR)", value), false
	} else if hiByte == 0x66 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (MONTH)", value), false
	} else if hiByte == 0x67 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (DAY)", value), false
	} else if hiByte == 0x68 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (HOUR)", value), false
	} else if hiByte == 0x69 {
		var value string = strconv.FormatInt(int64(loByte), 10)
		return fmt.Sprintf("Change Time Setting to %s (MINUTE)", value), false
	}
	return strings.ToUpper(strconv.FormatUint((uint64(hiByte)<<8)+uint64(loByte), 16)), true
}

func parseDateTime(bytes []byte) string {
	var year int = int(bytes[0]) + 2000
	var month int = int(bytes[1]&0xf) + 1
	var day int = int((bytes[2]>>3)&0x1f) + 1
	var hour int = int(((bytes[2] & 0x7) << 2) + ((bytes[3] >> 6) & 0x3))
	var minute int = int(bytes[3] & 0x3f)
	yearStr := strconv.FormatInt(int64(year), 10)
	var monthStr string
	if month < 10 {
		monthStr = "0" + strconv.FormatInt(int64(month), 10)
	} else {
		monthStr = strconv.FormatInt(int64(month), 10)
	}
	var dayStr string
	if day < 10 {
		dayStr = "0" + strconv.FormatInt(int64(day), 10)
	} else {
		dayStr = strconv.FormatInt(int64(day), 10)
	}
	var hourStr string
	if hour < 10 {
		hourStr = "0" + strconv.FormatInt(int64(hour), 10)
	} else {
		hourStr = strconv.FormatInt(int64(hour), 10)
	}
	var minuteStr string
	if minute < 10 {
		minuteStr = "0" + strconv.FormatInt(int64(minute), 10)
	} else {
		minuteStr = strconv.FormatInt(int64(minute), 10)
	}
	// 1/26/2019 6:00
	// 1/28/2019 15:44
	// 2019-1-28T15:44:00+8:00
	var builder strings.Builder
	builder.WriteString(yearStr)
	builder.WriteString("-")
	builder.WriteString(monthStr)
	builder.WriteString("-")
	builder.WriteString(dayStr)
	builder.WriteString("T")
	builder.WriteString(hourStr)
	builder.WriteString(":")
	builder.WriteString(minuteStr)
	builder.WriteString(":")
	builder.WriteString("00Z")
	t, _ := time.Parse(time.RFC3339, builder.String())
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return t.In(loc).Format("2006-01-02T15:04:05+08:00")
}

func ConvertEventLogToGPSEvent(eventLog *entity.EventLog) *entity.GpsEvent {
	gpsEvent := &entity.GpsEvent{}
	gpsEvent.Seqno = eventLog.Seq
	gpsEvent.CntrNum = eventLog.CntrNum
	gpsEvent.RevTime = eventLog.RevTime
	gpsEvent.CltTime = eventLog.LogTime
	gpsEvent.LocateTime = eventLog.LogTime
	if eventLog.Smode == "Electric Power Shut Off" || eventLog.Smode == "POWER OFF" || eventLog.Smode == "Battery ON" {
		gpsEvent.EleState = "0"
	} else {
		// Reefer datalog
		if eventLog.Isa == 1 && eventLog.Isc == 1 {
			if eventLog.Ss == 0 {
				gpsEvent.EleState = "0"
			} else {
				gpsEvent.EleState = "1"
			}
		} else {
			// eventLog
			gpsEvent.EleState = "1"
		}
	}
	if eventLog.Isa == 1 && eventLog.Isc == 1 {
		gpsEvent.OpMode = eventLog.Smode
	} else {
		gpsEvent.EventLog = eventLog.Smode
	}
	gpsEvent.PosFlag = "0"
	gpsEvent.Source = "TCP_SERVER"
	gpsEvent.Carrier = "COSU"
	gpsEvent.IsEventLog = true
	gpsEvent.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	if eventLog.IsTemperatureLog {
		gpsEvent.SetTem = strconv.FormatFloat(eventLog.Ss, 'f', 1, 64)
		gpsEvent.SupTem = strconv.FormatFloat(eventLog.Sp, 'f', 1, 64)
		gpsEvent.RetTem = strconv.FormatFloat(eventLog.Rs, 'f', 1, 64)
		gpsEvent.Hum = eventLog.Hus
		gpsEvent.Ambs = strconv.FormatFloat(eventLog.Ambs, 'f', 1, 64)
		gpsEvent.Hs = eventLog.Sh
		gpsEvent.Usda1 = eventLog.Usda1
		gpsEvent.Usda2 = eventLog.Usda2
		gpsEvent.Usda3 = eventLog.Usda3
		gpsEvent.Drs = strconv.FormatFloat(eventLog.Drs, 'f', 1, 64)
		gpsEvent.Dss = strconv.FormatFloat(eventLog.Dss, 'f', 1, 64)
		gpsEvent.Cts = eventLog.Cts
		gpsEvent.Isc = strconv.FormatInt(int64(eventLog.Isc), 10)
	}
	gpsEvent.Isa = strconv.FormatInt(int64(eventLog.Isa), 10)

	return gpsEvent
}

func roundFloat(input float64) float64 {
	result, _ := strconv.ParseFloat(strconv.FormatFloat(input, 'f', 1, 64), 64)
	return result
}
