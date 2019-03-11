package entity

type OperationModeChange struct {
	Seqno     string `json:"seqno,omitempty"`
	Cntrnum   string `json:"cntrnum,omitempty"`
	Opmode    string `json:"opmode,omitempty"`
	Logtime   string `json:"logtime,omitempty"`
	Revtime   string `json:"revtime,omitempty"`
	TableName string `json:"tableName,omitempty"` // default: "Tblopmoderec"
}
