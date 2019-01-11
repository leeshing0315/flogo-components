package entity

type ContainerSummary struct {
	Simno      string `json:"simno"`
	Carno      string `json:"carno"` // e.g. C04254
	Carid      string `json:"carid"` // e.g. CXRU1495240
	Commmode   string `json:"commmode"`
	Unitcode   int    `json:"unitcode"`
	Cartype    string `json:"cartype"`
	Saveflag   int    `json:"saveflag"`
	Calcflag   int    `json:"calcflag"`
	Changeflag int    `json:"changeflag"`
	Changetime string `json:"changetime"`
	Regtime    string `json:"regtime"`
	Devtype    string `json:"devtype"`
	Useacc     int    `json:"useacc"`
	Groupname  string `json:"groupname"`
	Checkflag  int    `json:"checkflag"`
	Boxtype    string `json:"boxtype"`
	Boxsize    string `json:"boxsize"`
	TableName  string `json:"tableName"` // default: 'Tblcarbaseinfo'
}
