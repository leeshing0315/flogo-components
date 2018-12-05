package smusingledata

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/leeshing0315/flogo-components/common/entity"
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
)

func handleData(seqNo int, data []byte) (gpsEvent *entity.GpsEvent, deviceError *entity.DeviceError, operationModeChange *entity.OperationModeChange, cntrNum string, devid string) {
	gpsEvent = &entity.GpsEvent{}
	deviceError = &entity.DeviceError{}
	operationModeChange = &entity.OperationModeChange{}

	println(len(data))
	handleLocationBasicInformation(data[0:23], gpsEvent)
	cntrNum, devid = handleAdditionalInformationItems(data[23:], gpsEvent)

	return gpsEvent, deviceError, operationModeChange, cntrNum, devid
}

func handleLocationBasicInformation(data []byte, gpsEvent *entity.GpsEvent) {
	handleTriggerEvent(data[0:2], gpsEvent)
	handleStatus(data[2:4], gpsEvent)
	handleDate(data[4:10], gpsEvent)
	handleRemainingInfo(data, gpsEvent)
}

func handleTriggerEvent(data []byte, gpsEvent *entity.GpsEvent) {
	positioningModuleFailure := data[0] & positioningModuleFailureMask
	if positioningModuleFailure == positioningModuleFailureMask {
		gpsEvent.PosFlag = "1"
		resultMap["定位模块故障"] = "1"
	} else {
		gpsEvent.PosFlag = "0"
		resultMap["定位模块故障"] = "0"
	}

	serialCommunicationFailure := data[0] & serialCommunicationFailureMask
	if serialCommunicationFailure == serialCommunicationFailureMask {
		resultMap["串口通讯故障"] = "1"
	} else {
		resultMap["串口通讯故障"] = "0"
	}

	communicationModuleFailure := data[0] & communicationModuleFailureMask
	if communicationModuleFailure == communicationModuleFailureMask {
		resultMap["通讯模块故障"] = "1"
	} else {
		resultMap["通讯模块故障"] = "0"
	}

	powerSupplyFailure := data[0] & powerSupplyFailureMask
	if powerSupplyFailure == powerSupplyFailureMask {
		resultMap["电源供电故障"] = "1"
	} else {
		resultMap["电源供电故障"] = "0"
	}

	batteryChargingFailure := data[0] & batteryChargingFailureMask
	if batteryChargingFailure == batteryChargingFailureMask {
		resultMap["电池充电故障"] = "1"
	} else {
		resultMap["电池充电故障"] = "0"
	}

	clockModuleFailure := data[0] & clockModuleFailureMask
	if clockModuleFailure == clockModuleFailureMask {
		resultMap["时钟模块故障"] = "1"
	} else {
		resultMap["时钟模块故障"] = "0"
	}

	coldBoxFaultCodeChange := data[0] & coldBoxFaultCodeChangeMask
	if coldBoxFaultCodeChange == coldBoxFaultCodeChangeMask {
		resultMap["冷箱故障码变化"] = "1"
	} else {
		resultMap["冷箱故障码变化"] = "0"
	}

	coldBoxOperationModeChange := data[0] & coldBoxOperationModeChangeMask
	if coldBoxOperationModeChange == coldBoxOperationModeChangeMask {
		resultMap["冷箱操作模式变化"] = "1"
	} else {
		resultMap["冷箱操作模式变化"] = "0"
	}

	powerSupplyStatusChange := data[1] & powerSupplyStatusChangeMask
	if powerSupplyStatusChange == powerSupplyStatusChangeMask {
		resultMap["供电状态变化"] = "1"
	} else {
		resultMap["供电状态变化"] = "0"
	}
}

func handleStatus(data []byte, gpsEvent *entity.GpsEvent) {
	positioning := data[0] & positioningMask
	if positioning == positioningMask {
		resultMap["定位"] = "1"
	} else {
		resultMap["定位"] = "0"
	}

	latitudeNorthSouth := data[0] & latitudeNorthSouthMask
	if latitudeNorthSouth == latitudeNorthSouthMask {
		resultMap["北纬南纬"] = "1"
	} else {
		resultMap["北纬南纬"] = "0"
	}

	longitudeEastWest := data[0] & longitudeEastWestMask
	if longitudeEastWest == longitudeEastWestMask {
		resultMap["东经西经"] = "1"
	} else {
		resultMap["东经西经"] = "0"
	}

	useGpsSatellitesForPositioning := data[0] & useGpsSatellitesForPositioningMask
	if useGpsSatellitesForPositioning == useGpsSatellitesForPositioningMask {
		resultMap["使用GPS卫星进行定位"] = "1"
	} else {
		resultMap["使用GPS卫星进行定位"] = "0"
	}

	useBeidouSatelliteForPositioning := data[0] & useBeidouSatelliteForPositioningMask
	if useBeidouSatelliteForPositioning == useBeidouSatelliteForPositioningMask {
		resultMap["使用北斗卫星进行定位"] = "1"
	} else {
		resultMap["使用北斗卫星进行定位"] = "0"
	}

	supplyByBatteryOrPower := data[0] & supplyByBatteryOrPowerMask
	if supplyByBatteryOrPower == supplyByBatteryOrPowerMask {
		resultMap["电池电源供电"] = "1"
	} else {
		resultMap["电池电源供电"] = "0"
	}

	inThePolygonArea := data[0] & inThePolygonAreaMask
	if inThePolygonArea == inThePolygonAreaMask {
		resultMap["在多边形区域内"] = "1"
	} else {
		resultMap["在多边形区域内"] = "0"
	}

	positioningModuleStatus := data[1] & positioningModuleStatusMask
	if positioningModuleStatus == positioningModuleStatusMask {
		resultMap["定位模块状态"] = "1"
	} else {
		resultMap["定位模块状态"] = "0"
	}

	serialCommunicationStatus := data[1] & serialCommunicationStatusMask
	if serialCommunicationStatus == serialCommunicationStatusMask {
		resultMap["串口通讯状态"] = "1"
	} else {
		resultMap["串口通讯状态"] = "0"
	}

	communicationModuleStatus := data[1] & communicationModuleStatusMask
	if communicationModuleStatus == communicationModuleStatusMask {
		resultMap["通讯模块状态"] = "1"
	} else {
		resultMap["通讯模块状态"] = "0"
	}

	powerSupplyStatus := data[1] & powerSupplyStatusMask
	if powerSupplyStatus == powerSupplyStatusMask {
		resultMap["电源供电状态"] = "1"
	} else {
		resultMap["电源供电状态"] = "0"
	}

	batteryChargingStatus := data[1] & batteryChargingStatusMask
	if batteryChargingStatus == batteryChargingStatusMask {
		resultMap["电池充电状态"] = "1"
	} else {
		resultMap["电池充电状态"] = "0"
	}

	clockModuleStatus := data[1] & clockModuleStatusMask
	if clockModuleStatus == clockModuleStatusMask {
		resultMap["时钟模块状态"] = "1"
	} else {
		resultMap["时钟模块状态"] = "0"
	}

	timedUploadData := data[1] & timedUploadDataMask
	if timedUploadData == timedUploadDataMask {
		resultMap["定时上传数据"] = "1"
	} else {
		resultMap["定时上传数据"] = "0"
	}
}

func handleDate(data []byte, gpsEvent *entity.GpsEvent) {
	year := FormatBCD(data[0])
	month := FormatBCD(data[1])
	day := FormatBCD(data[2])
	hour := FormatBCD(data[3])
	minute := FormatBCD(data[4])
	second := FormatBCD(data[5])

	resultMap["时间"] = strings.Join([]string{year, month, day, hour, minute, second}, "-")
}

func handleRemainingInfo(data []byte, gpsEvent *entity.GpsEvent) {
	gpsEvent.Lat = strconv.FormatFloat(float64(binary.BigEndian.Uint32(data[10:14]))/math.Pow10(6), 'f', -1, 64)
	gpsEvent.Lng = strconv.FormatFloat(float64(binary.BigEndian.Uint32(data[14:18]))/math.Pow10(6), 'f', -1, 64)
	gpsEvent.Speed = strconv.FormatFloat(float64(binary.BigEndian.Uint16(data[18:20]))/10, 'f', -1, 64)
	gpsEvent.Direction = strconv.FormatUint(uint64(binary.BigEndian.Uint16(data[20:22])), 10)
	gpsEvent.BatLevel = strconv.FormatUint(uint64(data[22]), 10)
}

func handleAdditionalInformationItems(data []byte, gpsEvent *entity.GpsEvent, resultMap map[string]string) (cntrNum string, devid string) {
	items := splitItems(data)
	cntrNum, devid = "", ""
	for _, item := range items {
		switch item[0] {
		case 0x01:
			cntrNum, devid = handleLoginItem(item, gpsEvent)
		case 0x02:
			handleInfoItem(item, gpsEvent)
		case 0x03:
			handleDebugTextItem(item, gpsEvent)
		case 0x04:
			handleNumberOfSatellitesItem(item, gpsEvent)
		case 0x05:
			handleOpModeItem(item, gpsEvent)
		case 0x06:
			handleFaultCodeItem(item, gpsEvent)
		case 0x07:
			handleColdBoxTimeItem(item, gpsEvent)
		}
	}
	return cntrNum, devid
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

func handleLoginItem(item []byte, gpsEvent *entity.GpsEvent) (cntrNum string, devid string) {
	pin := string(item[0:15])
	devid = string(item[15:21])
	cntrNum = string(item[21:32])
	reeferType := string(item[33])
	return cntrNum, devid
}

func handleInfoItem(item []byte, gpsEvent *entity.GpsEvent) {
}

func handleDebugTextItem(item []byte, gpsEvent *entity.GpsEvent) {
}

func handleNumberOfSatellitesItem(item []byte, gpsEvent *entity.GpsEvent) {
}

func handleOpModeItem(item []byte, gpsEvent *entity.GpsEvent) {
}

func handleFaultCodeItem(item []byte, gpsEvent *entity.GpsEvent) {
}

func handleColdBoxTimeItem(item []byte, gpsEvent *entity.GpsEvent) {
}
