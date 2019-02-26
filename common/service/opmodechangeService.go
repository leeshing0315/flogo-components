package service

import (
	"time"

	"github.com/leeshing0315/flogo-components/common/entity"
)

func GenOpModeChangeFromSinglePacket(singlePacket *entity.SinglePacket, seqNo string, cntrNum string) *entity.OperationModeChange {
	if singlePacket.ColdBoxOperationModeChange && singlePacket.InfoItem.OpMode != "" {
		operationModeChange := &entity.OperationModeChange{}

		operationModeChange.Seqno = seqNo
		operationModeChange.Logtime = changeDateFormatFromECMA(singlePacket.Date)
		operationModeChange.Revtime = time.Now().Format("2006-01-02 15:04:05.0")
		operationModeChange.Cntrnum = cntrNum
		operationModeChange.Opmode = singlePacket.InfoItem.OpMode
		operationModeChange.TableName = "Tblopmoderec"

		return operationModeChange
	} else {
		return nil
	}
}
