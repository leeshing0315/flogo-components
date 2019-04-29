package service

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/leeshing0315/flogo-components/common/entity"
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

	Bit0Mask byte = 0x01 // 0000 0001
	Bit1Mask byte = 0x02 // 0000 0010
	Bit2Mask byte = 0x04 // 0000 0100
	Bit3Mask byte = 0x08 // 0000 1000
	Bit4Mask byte = 0x10 // 0001 0000
	Bit5Mask byte = 0x20 // 0010 0000
	Bit6Mask byte = 0x40 // 0100 0000
	Bit7Mask byte = 0x80 // 1000 0000

	LOW_4_BITS_MASK byte = 0x0F // 0000 1111
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

func ParseToSinglePacket(data []byte) (*entity.SinglePacket, error) {
	singlePacket := &entity.SinglePacket{}

	err := handleLocationBasicInformation(data[0:23], singlePacket)
	if err != nil {
		return nil, err
	}
	err = handleAdditionalInformationItems(data[23:], singlePacket)
	if err != nil {
		return nil, err
	}

	return singlePacket, nil
}

func handleLocationBasicInformation(data []byte, singlePacket *entity.SinglePacket) error {
	handleTriggerEvent(data[0:2], singlePacket)

	handleStatus(data[2:4], singlePacket)

	err := handleDate(data[4:10], singlePacket)
	if err != nil {
		return err
	}

	handleRemainingInfo(data, singlePacket)

	return nil
}

func handleTriggerEvent(data []byte, singlePacket *entity.SinglePacket) {
	// positioningModuleFailure := data[0] & positioningModuleFailureMask
	positioningModuleFailure := data[1] & Bit0Mask
	if positioningModuleFailure == Bit0Mask {
		singlePacket.PositioningModuleFailure = true
	}

	// serialCommunicationFailure := data[0] & serialCommunicationFailureMask
	serialCommunicationFailure := data[1] & Bit1Mask
	if serialCommunicationFailure == Bit1Mask {
		singlePacket.SerialCommunicationFailure = true
	}

	// communicationModuleFailure := data[0] & communicationModuleFailureMask
	communicationModuleFailure := data[1] & Bit2Mask
	if communicationModuleFailure == Bit2Mask {
		singlePacket.CommunicationModuleFailure = true
	}

	// powerSupplyFailure := data[0] & powerSupplyFailureMask
	powerSupplyFailure := data[1] & Bit3Mask
	if powerSupplyFailure == Bit3Mask {
		singlePacket.PowerSupplyFailure = true
	}

	// batteryChargingFailure := data[0] & batteryChargingFailureMask
	batteryChargingFailure := data[1] & Bit4Mask
	if batteryChargingFailure == Bit4Mask {
		singlePacket.BatteryChargingFailure = true
	}

	// clockModuleFailure := data[0] & clockModuleFailureMask
	clockModuleFailure := data[1] & Bit5Mask
	if clockModuleFailure == Bit5Mask {
		singlePacket.ClockModuleFailure = true
	}

	// coldBoxFaultCodeChange := data[0] & coldBoxFaultCodeChangeMask
	coldBoxFaultCodeChange := data[1] & Bit6Mask
	if coldBoxFaultCodeChange == Bit6Mask {
		singlePacket.ColdBoxFaultCodeChange = true
	}

	// coldBoxOperationModeChange := data[0] & coldBoxOperationModeChangeMask
	coldBoxOperationModeChange := data[1] & Bit7Mask
	if coldBoxOperationModeChange == Bit7Mask {
		singlePacket.ColdBoxOperationModeChange = true
	}

	// powerSupplyStatusChange := data[1] & powerSupplyStatusChangeMask
	powerSupplyStatusChange := data[0] & Bit0Mask
	if powerSupplyStatusChange == Bit0Mask {
		singlePacket.PowerSupplyStatusChange = true
	}
}

func handleStatus(data []byte, singlePacket *entity.SinglePacket) {
	// positioning := data[0] & positioningMask
	positioning := data[1] & Bit1Mask
	if positioning == Bit1Mask {
		singlePacket.Positioning = true
	}

	// latitudeNorthSouth := data[0] & latitudeNorthSouthMask
	latitudeNorthSouth := data[1] & Bit2Mask
	if latitudeNorthSouth == Bit2Mask {
		singlePacket.LatitudeNorthSouth = true
	}

	// longitudeEastWest := data[0] & longitudeEastWestMask
	longitudeEastWest := data[1] & Bit3Mask
	if longitudeEastWest == Bit3Mask {
		singlePacket.LongitudeEastWest = true
	}

	// useGpsSatellitesForPositioning := data[0] & useGpsSatellitesForPositioningMask
	useGpsSatellitesForPositioning := data[1] & Bit4Mask
	if useGpsSatellitesForPositioning == Bit4Mask {
		singlePacket.UseGpsSatellitesForPositioning = true
	}

	// useBeidouSatelliteForPositioning := data[0] & useBeidouSatelliteForPositioningMask
	useBeidouSatelliteForPositioning := data[1] & Bit5Mask
	if useBeidouSatelliteForPositioning == Bit5Mask {
		singlePacket.UseBeidouSatelliteForPositioning = true
	}

	// supplyByBatteryOrPower := data[0] & supplyByBatteryOrPowerMask
	supplyByBatteryOrPower := data[1] & Bit6Mask
	if supplyByBatteryOrPower == Bit6Mask {
		singlePacket.SupplyByBatteryOrPower = true
	}

	// inThePolygonArea := data[0] & inThePolygonAreaMask
	inThePolygonArea := data[1] & Bit7Mask
	if inThePolygonArea == Bit7Mask {
		singlePacket.InThePolygonArea = true
	}

	// positioningModuleStatus := data[1] & positioningModuleStatusMask
	positioningModuleStatus := data[0] & Bit0Mask
	if positioningModuleStatus == Bit0Mask {
		singlePacket.PositioningModuleStatus = true
	}

	// serialCommunicationStatus := data[1] & serialCommunicationStatusMask
	serialCommunicationStatus := data[0] & Bit1Mask
	if serialCommunicationStatus == Bit1Mask {
		singlePacket.SerialCommunicationStatus = true
	}

	// communicationModuleStatus := data[1] & communicationModuleStatusMask
	communicationModuleStatus := data[0] & Bit2Mask
	if communicationModuleStatus == Bit2Mask {
		singlePacket.CommunicationModuleStatus = true
	}

	// powerSupplyStatus := data[1] & powerSupplyStatusMask
	powerSupplyStatus := data[0] & Bit3Mask
	if powerSupplyStatus == Bit3Mask {
		singlePacket.PowerSupplyStatus = true
	}

	// batteryChargingStatus := data[1] & batteryChargingStatusMask
	batteryChargingStatus := data[0] & Bit4Mask
	if batteryChargingStatus == Bit4Mask {
		singlePacket.BatteryChargingStatus = true
	}

	// clockModuleStatus := data[1] & clockModuleStatusMask
	clockModuleStatus := data[0] & Bit5Mask
	if clockModuleStatus == Bit5Mask {
		singlePacket.ClockModuleStatus = true
	}

	// timedUploadData := data[1] & timedUploadDataMask
	timedUploadData := data[0] & Bit7Mask
	if timedUploadData == Bit7Mask {
		singlePacket.TimedUploadData = true
	}
}

func handleDate(data []byte, singlePacket *entity.SinglePacket) error {
	dateStr, err := util.ParseDateStrFromBCD6(data)
	if err != nil {
		return err
	}
	singlePacket.Date = dateStr
	return nil
}

func handleRemainingInfo(data []byte, singlePacket *entity.SinglePacket) {
	singlePacket.Lat = strconv.FormatFloat(float64(binary.BigEndian.Uint32(data[10:14]))/math.Pow10(6), 'f', 6, 64)
	singlePacket.Lng = strconv.FormatFloat(float64(binary.BigEndian.Uint32(data[14:18]))/math.Pow10(6), 'f', 6, 64)
	singlePacket.Speed = strconv.FormatFloat(float64(binary.BigEndian.Uint16(data[18:20]))/10, 'f', 1, 64)
	singlePacket.Direction = strconv.FormatUint(uint64(binary.BigEndian.Uint16(data[20:22])), 10)
	singlePacket.BatLevel = strconv.FormatUint(uint64(data[22]), 10)
}

func handleAdditionalInformationItems(data []byte, singlePacket *entity.SinglePacket) error {
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
			err := handleColdBoxTimeItem(item, singlePacket)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

func handleLoginItem(item []byte, singlePacket *entity.SinglePacket) {
	dataSegment := item[2:]
	loginItem := entity.LoginItem{}

	loginItem.Pin = string(dataSegment[0:15])
	loginItem.DeviceID = string(dataSegment[15:21])
	loginItem.ContainerNumber = string(dataSegment[21:32])
	loginItem.ReeferType = string(dataSegment[32])

	singlePacket.LoginItem = loginItem
}

func handleInfoItem(item []byte, singlePacket *entity.SinglePacket) {
	dataSegmentLen := item[1]
	dataSegment := item[2:]
	infoItem := entity.InfoItem{}

	// VALIDSTATE: 4byte 下面冷箱数据段有效位 0：无效  1：有效，bit0对应OP_MODE……bit24对应FAULT_CODE  现在有25个数据段
	validState := dataSegment[1:5]

	// 冷箱类型: 1byte 目前都是大金的，ascall码‘D’
	infoItem.ReeferType = string(dataSegment[0])

	// OP_MODE: 1byte
	infoItem.OpModeValid = isValid(validState, 0)
	infoItem.OpMode = opModeMapping[dataSegment[5]&LOW_4_BITS_MASK]

	// SET_TEM: 2byte 多字节时，注意低位在右，参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.SetTemValid = isValid(validState, 1)
	infoItem.SetTem = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[6:8]))*0.0625, 'f', 1, 64)

	// SUP_TEM: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.SupTemValid = isValid(validState, 2)
	infoItem.SupTem = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[8:10]))*0.0625, 'f', 1, 64)

	// RET_TEM: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.RetTemValid = isValid(validState, 3)
	infoItem.RetTem = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[10:12]))*0.0625, 'f', 1, 64)

	// HUM: 1byte 参与计算的区间：b0-b7，除以2.54，其中LSB：b0，MSB：b7，单位是%RH，范围0~100%。
	infoItem.HumValid = isValid(validState, 4)
	infoItem.Hum = strconv.FormatFloat(intervalFloat64(float64(dataSegment[12])/2.54, 0, 100), 'f', 1, 64)

	// HPT: 1byte 参与计算的区间：b0-b7，乘以10，单位是Kpa，如果是FF则为无效值，范围0~2500Kpa。
	infoItem.HptValid = isValid(validState, 5)
	if dataSegment[13] != 0xFF {
		infoItem.Hpt = strconv.FormatInt(int64(dataSegment[13])*10, 10)
	}

	// USDA1: 1byte 参与计算的区间：b0-b7，乘以0.5，然后再减去5.0，单位是C，如果是FF则为无效值，范围-5.0~+20.0。
	infoItem.Usda1Valid = isValid(validState, 6)
	if dataSegment[14] != 0xFF {
		infoItem.Usda1 = strconv.FormatFloat(intervalFloat64(float64(dataSegment[14])*0.5-5.0, -5, 20), 'f', 1, 64)
	}

	// USDA2: 1byte 参与计算的区间：b0-b7，乘以0.5，然后再减去5.0，单位是C，如果是FF则为无效值，范围-5.0~+20.0。
	infoItem.Usda2Valid = isValid(validState, 7)
	if dataSegment[15] != 0xFF {
		infoItem.Usda2 = strconv.FormatFloat(intervalFloat64(float64(dataSegment[15])*0.5-5.0, -5, 20), 'f', 1, 64)
	}

	// USDA3: 1byte 参与计算的区间：b0-b7，乘以0.5，然后再减去5.0，单位是C，如果是FF则为无效值，范围-5.0~+20.0。
	infoItem.Usda3Valid = isValid(validState, 8)
	if dataSegment[16] != 0xFF {
		infoItem.Usda3 = strconv.FormatFloat(intervalFloat64(float64(dataSegment[16])*0.5-5.0, -5, 20), 'f', 1, 64)
	}

	// CTLTYPE: 1byte 这个是冷机软件版本，直接记录，如25H、26H等。
	infoItem.CtlTypeValid = isValid(validState, 9)
	infoItem.CtlType = strconv.FormatInt(int64(dataSegment[17]), 10)

	// 	LPT: 1byte 参与计算的区间：b0-b7，乘以10，然后再减去70，单位是Kpa，如果是FF则为无效值，范围0~2500Kpa。
	infoItem.LptValid = isValid(validState, 10)
	infoItem.Lpt = strconv.FormatInt(int64(dataSegment[18])*10-70, 10)

	// PT: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是V，范围20~654。
	infoItem.PtValid = isValid(validState, 11)
	infoItem.Pt = strconv.FormatFloat(intervalFloat64(float64(binary.BigEndian.Uint16(dataSegment[19:21]))*0.1, 20, 654), 'f', 1, 64)

	// CT1: 2byte 参与计算的区间：b0-b15，乘以0.01，单位是A，范围-1.11~54。
	infoItem.Ct1Valid = isValid(validState, 12)
	infoItem.Ct1 = strconv.FormatFloat(intervalFloat64(float64(int16(binary.BigEndian.Uint16(dataSegment[21:23])))*0.01, -1.11, 54), 'f', 1, 64)

	// CT2: 2byte 参与计算的区间：b0-b15，乘以0.01，单位是A，范围-1.01~50.96。
	infoItem.Ct2Valid = isValid(validState, 13)
	infoItem.Ct2 = strconv.FormatFloat(intervalFloat64(float64(int16(binary.BigEndian.Uint16(dataSegment[23:25])))*0.01, -1.01, 50.96), 'f', 1, 64)

	// AMBS: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.AmbsValid = isValid(validState, 14)
	infoItem.Ambs = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[25:27]))*0.0625, 'f', 1, 64)

	// EIS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.EisValid = isValid(validState, 15)
	infoItem.Eis = strconv.FormatFloat(intervalFloat64(float64(int16(binary.BigEndian.Uint16(dataSegment[27:29])))*0.1, -57, 100), 'f', 1, 64)

	// EOS: 2byte 参与计算的区间：b2-b13，乘以0.0625，其中LSB：b2，MSB：b13，单位是C。
	infoItem.EosValid = isValid(validState, 16)
	infoItem.Eos = strconv.FormatFloat(float64(util.BigEndianInt12(dataSegment[29:31]))*0.0625, 'f', 1, 64)

	// DCHS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围3~187。
	infoItem.DchsValid = isValid(validState, 17)
	infoItem.Dchs = strconv.FormatFloat(intervalFloat64(float64(binary.BigEndian.Uint16(dataSegment[31:33]))*0.1, 3, 187), 'f', 1, 64)

	// SGS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.SgsValid = isValid(validState, 18)
	infoItem.Sgs = strconv.FormatFloat(intervalFloat64(float64(int16(binary.BigEndian.Uint16(dataSegment[33:35])))*0.1, -57, 100), 'f', 1, 64)

	// SMV: 2byte 参与计算的区间：b0-b15，除以328，单位是PLS，范围0~500。
	infoItem.SmvValid = isValid(validState, 19)
	infoItem.Smv = strconv.FormatFloat(intervalFloat64(float64(binary.BigEndian.Uint16(dataSegment[35:37]))/328, 0, 500), 'f', 1, 64)

	// EV: 2byte 参与计算的区间：b0-b15，单位是%，范围0~100。
	infoItem.EvValid = isValid(validState, 20)
	infoItem.Ev = strconv.FormatInt(int64(binary.BigEndian.Uint16(dataSegment[37:39])), 10)

	// DSS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.DssValid = isValid(validState, 21)
	infoItem.Dss = strconv.FormatFloat(intervalFloat64(float64(int16(binary.BigEndian.Uint16(dataSegment[39:41])))*0.1, -57, 100), 'f', 1, 64)

	// DRS: 2byte 参与计算的区间：b0-b15，乘以0.1，单位是C，范围-57~100。
	infoItem.DrsValid = isValid(validState, 22)
	infoItem.Drs = strconv.FormatFloat(intervalFloat64(float64(int16(binary.BigEndian.Uint16(dataSegment[41:43])))*0.1, -57, 100), 'f', 1, 64)

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
	return (validState[len(validState)-1-index] & BitMask[offset]) == BitMask[offset]
}

func handleDebugTextItem(item []byte, singlePacket *entity.SinglePacket) {
	dataSegment := item[2:]
	debugTextItem := entity.DebugTextItem{}

	debugTextItem.DebugText = string(dataSegment)

	singlePacket.DebugTextItem = debugTextItem
}

func handleNumberOfSatellitesItem(item []byte, singlePacket *entity.SinglePacket) {
	dataSegment := item[2:]
	numberOfSatellitesItem := entity.NumberOfSatellitesItem{}

	numberOfSatellitesItem.GpsSatelliteNumber = strconv.FormatUint(uint64(dataSegment[0]), 10)
	numberOfSatellitesItem.BeidouSatelliteNumber = strconv.FormatUint(uint64(dataSegment[1]), 10)

	singlePacket.NumberOfSatellitesItem = numberOfSatellitesItem
}

func handleOpModeItem(item []byte, singlePacket *entity.SinglePacket) {
	dataSegment := item[2:]
	opModeItem := entity.OpModeItem{}

	opModeItem.OpMode = opModeMapping[dataSegment[0]]

	singlePacket.OpModeItem = opModeItem
}

func handleFaultCodeItem(item []byte, singlePacket *entity.SinglePacket) {
	dataSegmentLen := item[1]
	dataSegment := item[2:]
	faultCodeItem := entity.FaultCodeItem{}

	faultCodes := make([]string, dataSegmentLen)
	for _, faultCode := range dataSegment {
		faultCodes = append(faultCodes, faultCodeMapping[faultCode])
	}
	faultCodeItem.FaultCode = strings.Join(faultCodes, "|")

	singlePacket.FaultCodeItem = faultCodeItem
}

func handleColdBoxTimeItem(item []byte, singlePacket *entity.SinglePacket) error {
	dataSegment := item[2:]
	coldBoxTimeItem := entity.ColdBoxTimeItem{}

	cntrTime, err := util.ParseDateStrFromBCD6(dataSegment[0:6])
	if err != nil {
		return err
	}
	coldBoxTimeItem.CntrTime = cntrTime

	collectColdBoxTime, err := util.ParseDateStrFromBCD6(dataSegment[6:12])
	if err != nil {
		return err
	}
	coldBoxTimeItem.CollectColdBoxTime = collectColdBoxTime

	singlePacket.ColdBoxTimeItem = coldBoxTimeItem
	return nil
}

func intervalFloat64(number float64, min float64, max float64) float64 {
	if number < min {
		return min
	} else if number > max {
		return max
	}
	return number
}
