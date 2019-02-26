package entity

type OperationModeChange struct {
	Seqno     string `json:"seqno"`
	Cntrnum   string `json:"cntrnum"`
	Opmode    string `json:"opmode"`
	Logtime   string `json:"logtime"`
	Revtime   string `json:"revtime"`
	TableName string `json:"tableName"` // default: "Tblopmoderec"
}
