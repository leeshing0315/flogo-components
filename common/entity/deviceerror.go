package entity

type DeviceError struct {
	seqno     string
	devid     string
	faulttype string
	status    string
	logtime   string
	revtime   string
	tableName string // default: "Tbldevicefault"
}
