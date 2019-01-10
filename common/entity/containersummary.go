package entity

type ContainerSummary struct {
	Simno      string `json:"simno"`
	Carno      string `json:"carno"` // e.g. C04254
	Carid      string `json:"carid"` // e.g. CXRU1495240
	Commmode   string `json:"commmode"`
	Unitcode   string `json:"unitcode"`
	Cartype    string `json:"cartype"`
	Saveflag   string `json:"saveflag"`
	Calcflag   string `json:"calcflag"`
	Changeflag string `json:"changeflag"`
	Changetime string `json:"changetime"`
	Regtime    string `json:"regtime"`
	Devtype    string `json:"devtype"`
	Useacc     string `json:"useacc"`
	Groupname  string `json:"groupname"`
	Checkflag  string `json:"checkflag"`
	Boxtype    string `json:"boxtype"`
	Boxsize    string `json:"boxsize"`
	TableName  string `json:"tableName"` // default: 'Tblcarbaseinfo'
}
