package entity

type DeviceConfig struct {
	CntrNum                        string `json:"cntrNum"`
	DeviceId                       string `json:"devid"`
	Status                         string `json:"status"`
	UpdateTime                     string `json:"updateTime"`
	PowerOnCommunicationFrequency  string `json:"powerOnCommunicationFrequency"`
	PowerOffCommunicationFrequency string `json:"powerOffCommunicationFrequency"`
	CollectFrequency               string `json:"collectFrequency"`
	ServerIpAndPort                string `json:"serverIpAndPort"`
	SleepMode                      string `json:"sleepMode"`
	ReadDeviceConfig               string `json:"readDeviceConfig"`
	SeqNo                          string `json:"seqNo"`
}
