package entity

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/leeshing0315/flogo-components/common/util"
)

const (
	positioningModuleFailureMask   byte = 0x80 // 1000 0000
	serialCommunicationFailureMask byte = 0x40 // 0100 0000
	communicationModuleFailureMask byte = 0x20 // 0010 0000
	powerSupplyFailureMask         byte = 0x10 // 0001 0000
	batteryChargingFailureMask     byte = 0x08 // 0000 1000
	clockModuleFailureMask         byte = 0x04 // 0000 0100
	coldBoxFaultCodeChangeMask     byte = 0x02 // 0000 0010
	coldBoxOperationModeChangeMask byte = 0x01 // 0000 0001
	powerSupplyStatusChangeMask    byte = 0x80 // 1000 0000

	positioningMask                      byte = 0x40 // 0100 0000
	latitudeNorthSouthMask               byte = 0x20 // 0010 0000
	longitudeEastWestMask                byte = 0x10 // 0001 0000
	useGpsSatellitesForPositioningMask   byte = 0x08 // 0000 1000
	useBeidouSatelliteForPositioningMask byte = 0x04 // 0000 0100
	supplyByBatteryOrPowerMask           byte = 0x02 // 0000 0010
	inThePolygonAreaMask                 byte = 0x01 // 0000 0001
	positioningModuleStatusMask          byte = 0x80 // 1000 0000
	serialCommunicationStatusMask        byte = 0x40 // 0100 0000
	communicationModuleStatusMask        byte = 0x20 // 0010 0000
	powerSupplyStatusMask                byte = 0x10 // 0001 0000
	batteryChargingStatusMask            byte = 0x08 // 0000 1000
	clockModuleStatusMask                byte = 0x04 // 0000 0100
	timedUploadDataMask                  byte = 0x01 // 0000 0001

	Bit0Mask byte = 0x80 // 1000 0000
	Bit1Mask byte = 0x40 // 0100 0000
	Bit2Mask byte = 0x20 // 0010 0000
	Bit3Mask byte = 0x10 // 0001 0000
	Bit4Mask byte = 0x08 // 0000 1000
	Bit5Mask byte = 0x04 // 0000 0100
	Bit6Mask byte = 0x02 // 0000 0010
	Bit7Mask byte = 0x01 // 0000 0001
)

var BitMask = []byte{
	Bit0Mask,
	Bit1Mask,
	Bit2Mask,
	Bit3Mask,
	Bit4Mask,
	Bit5Mask,
	Bit6Mask,
	Bit7Mask,
}

type SinglePacket struct {
	PositioningModuleFailure   bool
	SerialCommunicationFailure bool
	CommunicationModuleFailure bool
	PowerSupplyFailure         bool
	BatteryChargingFailure     bool
	ClockModuleFailure         bool
	ColdBoxFaultCodeChange     bool
	ColdBoxOperationModeChange bool
	PowerSupplyStatusChange    bool

	Positioning                      bool
	LatitudeNorthSouth               bool
	LongitudeEastWest                bool
	UseGpsSatellitesForPositioning   bool
	UseBeidouSatelliteForPositioning bool
	SupplyByBatteryOrPower           bool
	InThePolygonArea                 bool
	PositioningModuleStatus          bool
	SerialCommunicationStatus        bool
	CommunicationModuleStatus        bool
	PowerSupplyStatus                bool
	BatteryChargingStatus            bool
	ClockModuleStatus                bool
	TimedUploadData                  bool

	Date      string
	Lat       string
	Lng       string
	Speed     string
	Direction string
	BatLevel  string

	LoginItem              LoginItem
	InfoItem               InfoItem
	DebugTextItem          DebugTextItem
	NumberOfSatellitesItem NumberOfSatellitesItem
	OpModeItem             OpModeItem
	FaultCodeItem          FaultCodeItem
	ColdBoxTimeItem        ColdBoxTimeItem
}

type LoginItem struct {
	Pin             string
	DeviceID        string
	ContainerNumber string
	ReeferType      string
}

type InfoItem struct {
	ReeferType     string
	OpModeValid    bool
	OpMode         string
	SetTemValid    bool
	SetTem         string
	SupTemValid    bool
	SupTem         string
	RetTemValid    bool
	RetTem         string
	HumValid       bool
	Hum            string
	HptValid       bool
	Hpt            string
	Usda1Valid     bool
	Usda1          string
	Usda2Valid     bool
	Usda2          string
	Usda3Valid     bool
	Usda3          string
	CtlTypeValid   bool
	CtlType        string
	LptValid       bool
	Lpt            string
	PtValid        bool
	Pt             string
	Ct1Valid       bool
	Ct1            string
	Ct2Valid       bool
	Ct2            string
	AmbsValid      bool
	Ambs           string
	EisValid       bool
	Eis            string
	EosValid       bool
	Eos            string
	DchsValid      bool
	Dchs           string
	SgsValid       bool
	Sgs            string
	SmvValid       bool
	Smv            string
	EvValid        bool
	Ev             string
	DssValid       bool
	Dss            string
	DrsValid       bool
	Drs            string
	HsValid        bool
	Hs             string
	FaultCodeValid bool
	FaultCode      string
}

type DebugTextItem struct {
	debugText string
}

type NumberOfSatellitesItem struct {
	GpsSatelliteNumber    string
	BeidouSatelliteNumber string
}

type OpModeItem struct {
	OpMode string
}

type FaultCodeItem struct {
	FaultCode string
}

type ColdBoxTimeItem struct {
	CntrTime           string
	CollectColdBoxTime string
}

func ParseToSinglePacket(data []byte) *SinglePacket {
	singlePacket := &SinglePacket{}

	handleLocationBasicInformation(data[0:23], singlePacket)

	return singlePacket
}

func handleLocationBasicInformation(data []byte, singlePacket *SinglePacket) {
	handleTriggerEvent(data[0:2], singlePacket)
	handleStatus(data[2:4], singlePacket)
	handleDate(data[4:10], singlePacket)
	handleRemainingInfo(data, singlePacket)
}

func handleTriggerEvent(data []byte, singlePacket *SinglePacket) {
	positioningModuleFailure := data[0] & positioningModuleFailureMask
	if positioningModuleFailure == positioningModuleFailureMask {
		singlePacket.PositioningModuleFailure = true
	}

	serialCommunicationFailure := data[0] & serialCommunicationFailureMask
	if serialCommunicationFailure == serialCommunicationFailureMask {
		singlePacket.SerialCommunicationFailure = true
	}

	communicationModuleFailure := data[0] & communicationModuleFailureMask
	if communicationModuleFailure == communicationModuleFailureMask {
		singlePacket.CommunicationModuleFailure = true
	}

	powerSupplyFailure := data[0] & powerSupplyFailureMask
	if powerSupplyFailure == powerSupplyFailureMask {
		singlePacket.PowerSupplyFailure = true
	}

	batteryChargingFailure := data[0] & batteryChargingFailureMask
	if batteryChargingFailure == batteryChargingFailureMask {
		singlePacket.BatteryChargingFailure = true
	}

	clockModuleFailure := data[0] & clockModuleFailureMask
	if clockModuleFailure == clockModuleFailureMask {
		singlePacket.ClockModuleFailure = true
	}

	coldBoxFaultCodeChange := data[0] & coldBoxFaultCodeChangeMask
	if coldBoxFaultCodeChange == coldBoxFaultCodeChangeMask {
		singlePacket.ColdBoxFaultCodeChange = true
	}

	coldBoxOperationModeChange := data[0] & coldBoxOperationModeChangeMask
	if coldBoxOperationModeChange == coldBoxOperationModeChangeMask {
		singlePacket.ColdBoxOperationModeChange = true
	}

	powerSupplyStatusChange := data[1] & powerSupplyStatusChangeMask
	if powerSupplyStatusChange == powerSupplyStatusChangeMask {
		singlePacket.PowerSupplyStatusChange = true
	}
}

func handleStatus(data []byte, singlePacket *SinglePacket) {
	positioning := data[0] & positioningMask
	if positioning == positioningMask {
		singlePacket.Positioning = true
	}

	latitudeNorthSouth := data[0] & latitudeNorthSouthMask
	if latitudeNorthSouth == latitudeNorthSouthMask {
		singlePacket.LatitudeNorthSouth = true
	}

	longitudeEastWest := data[0] & longitudeEastWestMask
	if longitudeEastWest == longitudeEastWestMask {
		singlePacket.LongitudeEastWest = true
	}

	useGpsSatellitesForPositioning := data[0] & useGpsSatellitesForPositioningMask
	if useGpsSatellitesForPositioning == useGpsSatellitesForPositioningMask {
		singlePacket.UseGpsSatellitesForPositioning = true
	}

	useBeidouSatelliteForPositioning := data[0] & useBeidouSatelliteForPositioningMask
	if useBeidouSatelliteForPositioning == useBeidouSatelliteForPositioningMask {
		singlePacket.UseBeidouSatelliteForPositioning = true
	}

	supplyByBatteryOrPower := data[0] & supplyByBatteryOrPowerMask
	if supplyByBatteryOrPower == supplyByBatteryOrPowerMask {
		singlePacket.SupplyByBatteryOrPower = true
	}

	inThePolygonArea := data[0] & inThePolygonAreaMask
	if inThePolygonArea == inThePolygonAreaMask {
		singlePacket.InThePolygonArea = true
	}

	positioningModuleStatus := data[1] & positioningModuleStatusMask
	if positioningModuleStatus == positioningModuleStatusMask {
		singlePacket.PositioningModuleStatus = true
	}

	serialCommunicationStatus := data[1] & serialCommunicationStatusMask
	if serialCommunicationStatus == serialCommunicationStatusMask {
		singlePacket.SerialCommunicationStatus = true
	}

	communicationModuleStatus := data[1] & communicationModuleStatusMask
	if communicationModuleStatus == communicationModuleStatusMask {
		singlePacket.CommunicationModuleStatus = true
	}

	powerSupplyStatus := data[1] & powerSupplyStatusMask
	if powerSupplyStatus == powerSupplyStatusMask {
		singlePacket.PowerSupplyStatus = true
	}

	batteryChargingStatus := data[1] & batteryChargingStatusMask
	if batteryChargingStatus == batteryChargingStatusMask {
		singlePacket.BatteryChargingStatus = true
	}

	clockModuleStatus := data[1] & clockModuleStatusMask
	if clockModuleStatus == clockModuleStatusMask {
		singlePacket.ClockModuleStatus = true
	}

	timedUploadData := data[1] & timedUploadDataMask
	if timedUploadData == timedUploadDataMask {
		singlePacket.TimedUploadData = true
	}
}

func handleDate(data []byte, singlePacket *SinglePacket) {
	singlePacket.Date = util.ParseDateStrFromBCD6(data)
}

func handleRemainingInfo(data []byte, singlePacket *SinglePacket) {
	singlePacket.Lat = strconv.FormatFloat(float64(binary.BigEndian.Uint32(data[10:14]))/math.Pow10(6), 'f', -1, 64)
	singlePacket.Lng = strconv.FormatFloat(float64(binary.BigEndian.Uint32(data[14:18]))/math.Pow10(6), 'f', -1, 64)
	singlePacket.Speed = strconv.FormatFloat(float64(binary.BigEndian.Uint16(data[18:20]))/10, 'f', -1, 64)
	singlePacket.Direction = strconv.FormatUint(uint64(binary.BigEndian.Uint16(data[20:22])), 10)
	singlePacket.BatLevel = strconv.FormatUint(uint64(data[22]), 10)
}

func handleAdditionalInformationItems(data []byte, singlePacket *SinglePacket) {
	items := splitItems(data)
	for _, item := range items {
		switch item[0] {
		case 0x01:
			handleLoginItem(item, singlePacket)
		case 0x02:
			handleInfoItem(item, singlePacket)
		case 0x03:
			handleDebugTextItem(item, singlePacket)
		case 0x04:
			handleNumberOfSatellitesItem(item, singlePacket)
		case 0x05:
			handleOpModeItem(item, singlePacket)
		case 0x06:
			handleFaultCodeItem(item, singlePacket)
		case 0x07:
			handleColdBoxTimeItem(item, singlePacket)
		}
	}
}

func splitItems(data []byte) [][]byte {
	var items = [][]byte{}
	// 附加信息项：附加信息ID(BYTE) + 附加信息长度(BYTE) + 附加信息
	for i := data; len(i) > 1; {
		itemLen := i[1]
		item := i[0 : itemLen+2]
		items = append(items, item)

		i = i[itemLen+2:]
	}
	return items
}

func handleLoginItem(item []byte, singlePacket *SinglePacket) {
	dataSegment := item[2:]
	loginItem := LoginItem{}

	loginItem.Pin = string(dataSegment[0:15])
	loginItem.DeviceID = string(dataSegment[15:21])
	loginItem.ContainerNumber = string(dataSegment[21:32])
	loginItem.ReeferType = string(dataSegment[33])

	singlePacket.LoginItem = loginItem
}

func handleInfoItem(item []byte, singlePacket *SinglePacket) {
	dataSegmentLen := item[1]
	dataSegment := item[2:]
	infoItem := InfoItem{}

	// VALIDSTATE: 4byte 下面冷箱数据段有效位 0：无效  1：有效，bit0对应OP_MODE……bit24对应FAULT_CODE  现在有25个数据段
	validState := dataSegment[1:5]

	// 冷箱类型: 1byte 目前都是大金的，ascall码‘D’
	infoItem.ReeferType = string(dataSegment[0])

	// OP_MODE: 1byte
	infoItem.OpModeValid = isValid(validState, 0)
	infoItem.OpMode = opModeMapping[dataSegment[5]]

	// SET_TEM: 2byte 多字节时，注意低位在右，参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.SetTemValid = isValid(validState, 1)
	infoItem.SetTem = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[6:8]))*0.0625, 'f', -1, 64)

	// SUP_TEM: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.SupTemValid = isValid(validState, 2)
	infoItem.SupTem = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[8:10]))*0.0625, 'f', -1, 64)

	// RET_TEM: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.RetTemValid = isValid(validState, 3)
	infoItem.RetTem = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[10:12]))*0.0625, 'f', -1, 64)

	// HUM: 1byte 参与计算的区间：b0-b7，除以2.54，其中LSB：b0，MSB：b7，单位是%RH，范围0~100%。
	infoItem.HumValid = isValid(validState, 4)
	infoItem.Hum = strconv.FormatFloat(float64(dataSegment[12])/2.54, 'f', -1, 64)

	// HPT: 1byte 参与计算的区间：b0-b7，乘以10，单位是Kpa，如果是FF则为无效值，范围0~2500Kpa。
	infoItem.HptValid = isValid(validState, 5)
	if dataSegment[13] != 0xFF {
		infoItem.Hpt = strconv.FormatInt(int64(dataSegment[13])*10, 10)
	}

	// USDA1: 1byte 参与计算的区间：b0-b7，乘以0.5，然后再减去5.0，单位是C，如果是FF则为无效值，范围-5.0~+20.0。
	infoItem.Usda1Valid = isValid(validState, 6)
	if dataSegment[14] != 0xFF {
		infoItem.Usda1 = strconv.FormatFloat(float64(dataSegment[14])*0.5-5.0, 'f', -1, 64)
	}

	// USDA2: 1byte 参与计算的区间：b0-b7，乘以0.5，然后再减去5.0，单位是C，如果是FF则为无效值，范围-5.0~+20.0。
	infoItem.Usda2Valid = isValid(validState, 7)
	if dataSegment[15] != 0xFF {
		infoItem.Usda2 = strconv.FormatFloat(float64(dataSegment[15])*0.5-5.0, 'f', -1, 64)
	}

	// USDA3: 1byte 参与计算的区间：b0-b7，乘以0.5，然后再减去5.0，单位是C，如果是FF则为无效值，范围-5.0~+20.0。
	infoItem.Usda3Valid = isValid(validState, 8)
	if dataSegment[16] != 0xFF {
		infoItem.Usda3 = strconv.FormatFloat(float64(dataSegment[16])*0.5-5.0, 'f', -1, 64)
	}

	// CTLTYPE: 1byte 这个是冷机软件版本，直接记录，如25H、26H等。
	infoItem.CtlTypeValid = isValid(validState, 9)
	infoItem.CtlType = strconv.FormatInt(int64(dataSegment[17]), 10)

	// 	LPT: 1byte 参与计算的区间：b0-b7，乘以10，然后再减去70，单位是Kpa，如果是FF则为无效值，范围0~2500Kpa。
	infoItem.LptValid = isValid(validState, 10)
	infoItem.Lpt = strconv.FormatInt(int64(dataSegment[18])*10-70, 10)

	// PT: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是V，范围20~654。
	infoItem.PtValid = isValid(validState, 11)
	infoItem.Pt = strconv.FormatFloat(float64(binary.BigEndian.Uint16(dataSegment[19:21]))*0.1, 'f', -1, 64)

	// CT1: 2byte 参与计算的区间：b0-b15，乘以0.01，单位是A，范围-1.11~54。
	infoItem.Ct1Valid = isValid(validState, 12)
	infoItem.Ct1 = strconv.FormatFloat(float64(int16(binary.BigEndian.Uint16(dataSegment[21:23])))*0.01, 'f', -1, 64)

	// CT2: 2byte 参与计算的区间：b0-b15，乘以0.01，单位是A，范围-1.01~50.96。
	infoItem.Ct2Valid = isValid(validState, 13)
	infoItem.Ct2 = strconv.FormatFloat(float64(int16(binary.BigEndian.Uint16(dataSegment[23:25])))*0.01, 'f', -1, 64)

	// AMBS: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.AmbsValid = isValid(validState, 14)
	infoItem.Ambs = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[25:27]))*0.0625, 'f', -1, 64)

	// EIS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.EisValid = isValid(validState, 15)
	infoItem.Eis = strconv.FormatFloat(float64(int16(binary.BigEndian.Uint16(dataSegment[27:29])))*0.1, 'f', -1, 64)

	// EOS: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.EosValid = isValid(validState, 16)
	infoItem.Eos = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[29:31]))*0.0625, 'f', -1, 64)

	// DCHS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围3~187。
	infoItem.DchsValid = isValid(validState, 17)
	infoItem.Dchs = strconv.FormatFloat(float64(binary.BigEndian.Uint16(dataSegment[31:33]))*0.1, 'f', -1, 64)

	// SGS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.SgsValid = isValid(validState, 18)
	infoItem.Sgs = strconv.FormatFloat(float64(int16(binary.BigEndian.Uint16(dataSegment[33:35])))*0.1, 'f', -1, 64)

	// SMV: 2byte 参与计算的区间：b0-b15，除以328，单位是PLS，范围0~500。
	infoItem.SmvValid = isValid(validState, 19)
	infoItem.Smv = strconv.FormatFloat(float64(binary.BigEndian.Uint16(dataSegment[35:37]))/328, 'f', -1, 64)

	// EV: 2byte 参与计算的区间：b0-b15，单位是%，范围0~100。
	infoItem.EvValid = isValid(validState, 20)
	infoItem.Ev = strconv.FormatFloat(float64(binary.BigEndian.Uint16(dataSegment[37:39])), 'f', -1, 64)

	// DSS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.DssValid = isValid(validState, 21)
	infoItem.Dss = strconv.FormatFloat(float64(int16(binary.BigEndian.Uint16(dataSegment[39:41])))*0.1, 'f', -1, 64)

	// DRS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.DrsValid = isValid(validState, 22)
	infoItem.Drs = strconv.FormatFloat(float64(int16(binary.BigEndian.Uint16(dataSegment[41:43])))*0.1, 'f', -1, 64)

	// HS: 1byte 参与计算的区间：b0-b6，其中LSB：b0，MSB：b6，单位是%RH，范围是30~95。
	infoItem.HsValid = isValid(validState, 23)
	infoItem.Hs = strconv.FormatInt(util.BigEndianFromBit0ToBit6(dataSegment[43]), 10)

	if dataSegmentLen > 44 {
		faultCodes := make([]string, dataSegmentLen-44)
		for _, faultCode := range dataSegment[44:] {
			faultCodes = append(faultCodes, faultCodeMapping[faultCode])
		}
		infoItem.FaultCode = strings.Join(faultCodes, "|")
	}

	singlePacket.InfoItem = infoItem
}

func isValid(validState []byte, number int) bool {
	index := number / 8
	offset := number % 8
	return (validState[index] & BitMask[offset]) == BitMask[offset]
}

func handleDebugTextItem(item []byte, singlePacket *SinglePacket) {
}

func handleNumberOfSatellitesItem(item []byte, singlePacket *SinglePacket) {
}

func handleOpModeItem(item []byte, singlePacket *SinglePacket) {
}

func handleFaultCodeItem(item []byte, singlePacket *SinglePacket) {
}

func handleColdBoxTimeItem(item []byte, singlePacket *SinglePacket) {
}
