package entity

type SimcardPin struct {
	Pinno     string `json:"pinno,omitempty"`
	Simno     string `json:"simno,omitempty"`
	Regtype   string `json:"regtype,omitempty"`
	Regtime   string `json:"regtime,omitempty"`
	TableName string `json:"tableName,omitempty"` // default: 'TblPin2SimNo'
}
