package entity

import "time"

type OperationModeChange struct {
	Seqno     string `json:"seqno"`
	Cntrnum   string `json:"cntrnum"`
	Opmode    string `json:"opmode"`
	Logtime   string `json:"logtime"`
	Revtime   string `json:"revtime"`
	TableName string `json:"tableName"` // default: "Tblopmoderec"
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
