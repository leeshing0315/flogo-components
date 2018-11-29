package smusingledata

import (
	"encoding/json"

	"github.com/leeshing0315/flogo-components/entity"
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

	positioningMask                    byte = 0x80 // 1000 0000
	latitudeNorthSouthMask             byte = 0x40 // 0100 0000
	longitudeEastWestMask              byte = 0x20 // 0010 0000
	useGpsSatellitesForPositioningMask byte = 0x10 // 0001 0000
)

func handleData(data []byte) string {
	gpsEvent := &entity.GpsEvent{}

	resultMap := make(map[string]string)
	handleLocationBasicInformation(data[0:23], gpsEvent, resultMap)

	resultStr, _ := json.Marshal(resultMap)
	return string(resultStr)
}

func handleLocationBasicInformation(data []byte, gpsEvent *entity.GpsEvent, resultMap map[string]string) {
	handleTriggerEvent(data, gpsEvent, resultMap)
}

func handleTriggerEvent(data []byte, gpsEvent *entity.GpsEvent, resultMap map[string]string) {
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
		resultMap["串口通訊故障"] = "1"
	} else {
		resultMap["串口通訊故障"] = "0"
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
