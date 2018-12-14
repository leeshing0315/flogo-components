package entity

import "time"

type DeviceError struct {
	Seqno     string
	Devid     string
	Faulttype string
	Status    string
	Logtime   string
	Revtime   string
	TableName string // default: "Tbldevicefault"
}

func GenDeviceErrorsFromSinglePacket(singlePacket *SinglePacket, seqNo string, devId string) []*DeviceError {
	deviceErrors := make([]*DeviceError, 6)

	positioningModuleFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	positioningModuleFailure.Faulttype = "1"
	positioningModuleFailure.Status = ternaryOperator(singlePacket.PositioningModuleFailure, "1", "0")
	deviceErrors[0] = positioningModuleFailure

	serialCommunicationFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	serialCommunicationFailure.Faulttype = "2"
	serialCommunicationFailure.Status = ternaryOperator(singlePacket.SerialCommunicationFailure, "1", "0")
	deviceErrors[1] = serialCommunicationFailure

	communicationModuleFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	communicationModuleFailure.Faulttype = "3"
	communicationModuleFailure.Status = ternaryOperator(singlePacket.CommunicationModuleFailure, "1", "0")
	deviceErrors[2] = communicationModuleFailure

	powerSupplyFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	powerSupplyFailure.Faulttype = "4"
	powerSupplyFailure.Status = ternaryOperator(singlePacket.PowerSupplyFailure, "1", "0")
	deviceErrors[3] = powerSupplyFailure

	batteryChargingFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	batteryChargingFailure.Faulttype = "5"
	batteryChargingFailure.Status = ternaryOperator(singlePacket.BatteryChargingFailure, "1", "0")
	deviceErrors[4] = batteryChargingFailure

	clockModuleFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	clockModuleFailure.Faulttype = "6"
	clockModuleFailure.Status = ternaryOperator(singlePacket.ClockModuleFailure, "1", "0")
	deviceErrors[5] = clockModuleFailure

	return deviceErrors
}

func genCommonDeviceError(singlePacket *SinglePacket, seqNo string, devId string) *DeviceError {
	deviceError := &DeviceError{}

	deviceError.Devid = devId
	deviceError.Logtime = singlePacket.Date
	deviceError.Revtime = time.Now().Format("2006-01-02 15:04:05.0")
	deviceError.Seqno = seqNo
	deviceError.TableName = "Tbldevicefault"

	return deviceError
}

func ternaryOperator(b bool, trueVal string, falseVal string) string {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}
