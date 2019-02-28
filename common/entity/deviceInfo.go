package entity

type DeviceInfo struct {
	Simno     string `json:"simno"`
	Devtype   string `json:"devtype"`
	Ip        string `json:"ip"`
	Remark    string `json:"remark"`
	Savetime  string `json:"savetime"`
	Setaddr   string `json:"setaddr"`
	TableName string `json:"tableName"` // default: Tbldevinfo
}
