package entity

type DeviceError struct {
	Seqno     string
	Devid     string
	Faulttype string
	Status    string
	Logtime   string
	Revtime   string
	TableName string // default: "Tbldevicefault"
}
