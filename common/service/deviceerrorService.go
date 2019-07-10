package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/leeshing0315/flogo-components/common/entity"
)

func GenDeviceErrorsFromSinglePacket(singlePacket *entity.SinglePacket, seqNo string, devId string) []string {
	deviceErrors := make([]string, 6)

	positioningModuleFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	positioningModuleFailure.Faulttype = "1"
	positioningModuleFailure.Status = ternaryOperator(singlePacket.PositioningModuleFailure, "1", "0")
	positioningModuleFailureStr, _ := json.Marshal(positioningModuleFailure)
	deviceErrors[0] = string(positioningModuleFailureStr)

	serialCommunicationFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	serialCommunicationFailure.Faulttype = "2"
	serialCommunicationFailure.Status = ternaryOperator(singlePacket.SerialCommunicationFailure, "1", "0")
	serialCommunicationFailureStr, _ := json.Marshal(serialCommunicationFailure)
	deviceErrors[1] = string(serialCommunicationFailureStr)

	communicationModuleFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	communicationModuleFailure.Faulttype = "3"
	communicationModuleFailure.Status = ternaryOperator(singlePacket.CommunicationModuleFailure, "1", "0")
	communicationModuleFailureStr, _ := json.Marshal(communicationModuleFailure)
	deviceErrors[2] = string(communicationModuleFailureStr)

	powerSupplyFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	powerSupplyFailure.Faulttype = "4"
	powerSupplyFailure.Status = ternaryOperator(singlePacket.PowerSupplyFailure, "1", "0")
	powerSupplyFailureStr, _ := json.Marshal(powerSupplyFailure)
	deviceErrors[3] = string(powerSupplyFailureStr)

	batteryChargingFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	batteryChargingFailure.Faulttype = "5"
	batteryChargingFailure.Status = ternaryOperator(singlePacket.BatteryChargingFailure, "1", "0")
	batteryChargingFailureStr, _ := json.Marshal(batteryChargingFailure)
	deviceErrors[4] = string(batteryChargingFailureStr)

	clockModuleFailure := genCommonDeviceError(singlePacket, seqNo, devId)
	clockModuleFailure.Faulttype = "6"
	clockModuleFailure.Status = ternaryOperator(singlePacket.ClockModuleFailure, "1", "0")
	clockModuleFailureStr, _ := json.Marshal(clockModuleFailure)
	deviceErrors[5] = string(clockModuleFailureStr)

	return deviceErrors
}

var UTC8 *time.Location = getUTC8()

func getUTC8() *time.Location {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	return location
}

func genCommonDeviceError(singlePacket *entity.SinglePacket, seqNo string, devId string) *entity.DeviceError {
	deviceError := &entity.DeviceError{}

	deviceError.Devid = devId
	deviceError.Logtime = changeDateFormatFromECMA(singlePacket.Date)
	deviceError.Revtime = time.Now().In(UTC8).Format("2006-01-02 15:04:05.0")
	deviceError.Seqno = seqNo
	deviceError.TableName = "Tbldevicefault"

	deviceError.Source = "TCP_SERVER"

	return deviceError
}

func ternaryOperator(b bool, trueVal string, falseVal string) string {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

func changeDateFormatFromECMA(input string) string {
	return strings.Replace(strings.Replace(input, "T", " ", 1), "+08:00", ".0", 1)
}
