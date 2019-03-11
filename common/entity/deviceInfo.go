package entity

type DeviceInfo struct {
	Simno     string `json:"simno,omitempty"`
	Devtype   string `json:"devtype,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Remark    string `json:"remark,omitempty"`
	Savetime  string `json:"savetime,omitempty"`
	Setaddr   string `json:"setaddr,omitempty"`
	TableName string `json:"tableName,omitempty"` // default: Tbldevinfo
}
