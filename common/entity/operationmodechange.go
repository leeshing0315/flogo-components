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
	if singlePacket.ColdBoxOperationModeChange && singlePacket.InfoItem.OpMode != "" {
		operationModeChange := &OperationModeChange{}

		operationModeChange.Seqno = seqNo
		operationModeChange.Logtime = singlePacket.Date
		operationModeChange.Revtime = time.Now().Format("2006-01-02 15:04:05.0")
		operationModeChange.Cntrnum = cntrNum
		operationModeChange.Opmode = singlePacket.InfoItem.OpMode
		operationModeChange.TableName = "Tblopmoderec"

		return operationModeChange
	} else {
		return nil
	}
}
