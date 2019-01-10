package entity

type SimcardPin struct {
	Pinno     string `json:"pinno"`
	Simno     string `json:"simno"`
	Regtype   string `json:"regtype"`
	Regtime   string `json:"regtime"`
	TableName string `json:"tableName"` // default: 'TblPin2SimNo'
}
