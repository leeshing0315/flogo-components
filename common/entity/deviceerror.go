package entity

type DeviceError struct {
	Seqno     string `json:"seqno"`
	Devid     string `json:"devid"`
	Faulttype string `json:"faulttype"`
	Status    string `json:"status"`
	Logtime   string `json:"logtime"`
	Revtime   string `json:"revtime"`
	TableName string `json:"tableName"` // default: "Tbldevicefault"
}
