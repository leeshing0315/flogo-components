package entity

type ContainerSummary struct {
	Simno    string `json:"simno,omitempty"`
	Carno    string `json:"carno,omitempty"` // e.g. C04254
	Carid    string `json:"carid,omitempty"` // e.g. CXRU1495240
	Commmode string `json:"commmode,omitempty"`
	// Unitcode   int    `json:"unitcode,omitempty"`
	Cartype string `json:"cartype,omitempty"`
	// Saveflag   int    `json:"saveflag,omitempty"`
	// Calcflag   int    `json:"calcflag,omitempty"`
	// Changeflag int    `json:"changeflag,omitempty"`
	Changetime string `json:"changetime,omitempty"`
	Regtime    string `json:"regtime,omitempty"`
	Devtype    string `json:"devtype,omitempty"`
	// Useacc     int    `json:"useacc,omitempty"`
	Groupname string `json:"groupname,omitempty"`
	// Checkflag  int    `json:"checkflag,omitempty"`
	Boxtype   string `json:"boxtype,omitempty"`
	Boxsize   string `json:"boxsize,omitempty"`
	TableName string `json:"tableName,omitempty"` // default: 'Tblcarbaseinfo'
}
