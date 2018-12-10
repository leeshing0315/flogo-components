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

func GenDeviceErrorFromSinglePacket(singlePacket *SinglePacket, seqNo string, cntrNum string) *OperationModeChange {
	operationModeChange := &OperationModeChange{}

	operationModeChange.Seqno = seqNo
	operationModeChange.Logtime = singlePacket.Date
	operationModeChange.Revtime = time.Now().Format("2006-01-02 15:04:05.0")
	operationModeChange.Cntrnum = cntrNum
	operationModeChange.Opmode = singlePacket.InfoItem.OpMode
	operationModeChange.TableName = "Tblopmoderec"

	return operationModeChange
}
