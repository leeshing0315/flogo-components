package entity

type OperationModeChange struct {
	seqno     string
	cntrnum   string
	opmode    string
	logtime   string
	revtime   string
	tableName string // default: "Tblopmoderec"
}
