package entity

type DeviceError struct {
	Seqno     string `json:"seqno,omitempty"`
	Devid     string `json:"devid,omitempty"`
	Faulttype string `json:"faulttype,omitempty"`
	Status    string `json:"status,omitempty"`
	Logtime   string `json:"logtime,omitempty"`
	Revtime   string `json:"revtime,omitempty"`
	TableName string `json:"tableName,omitempty"` // default: "Tbldevicefault"

	Source string `json:"source,omitempty"`
}
