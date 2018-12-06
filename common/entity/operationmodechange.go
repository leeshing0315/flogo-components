package entity

import "time"

type OperationModeChange struct {
	Seqno     string
	Cntrnum   string
	Opmode    string
	Logtime   string
	Revtime   string
	TableName string // default: "Tblopmoderec"
}

func GenOpModeChangeFromSinglePacket(singlePacket *SinglePacket, seqNo string, cntrNum string) *OperationModeChange {
	operationModeChange := &OperationModeChange{}

	operationModeChange.Seqno = seqNo
	operationModeChange.Logtime = singlePacket.Date
	operationModeChange.Revtime = time.Now().Format("2018-12-03 09:29:21.0")
	operationModeChange.Cntrnum = cntrNum
	operationModeChange.Opmode = singlePacket.InfoItem.OpMode
	operationModeChange.TableName = "Tblopmoderec"

	return operationModeChange
}
