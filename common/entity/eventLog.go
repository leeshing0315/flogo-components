package entity

type EventLog struct {
	Seq     string `json:"SEQ"`
	CntrNum string `json:"CNTR_NUM"`
	RevTime string `json:"REV_TIME"`
	LogTime string `json:"LOG_TIME"`
	Sp      string `json:"SP"`
	Isc     string `json:"ISC"`
	Ss      string `json:"SS"`
	Rs      string `json:"RS"`
	Dss     string `json:"DSS"`
	Drs     string `json:"DRS"`
	Ambs    string `json:"AMBS"`
	Hus     string `json:"HUS"`
	Sh      string `json:"SH"`
	Usda1   string `json:"USDA1"`
	Usda2   string `json:"USDA2"`
	Usda3   string `json:"USDA3"`
	Cts     string `json:"CTS"`
	Smode   string `json:"SMODE"`
	Isa     int32  `json:"ISA"`
}
