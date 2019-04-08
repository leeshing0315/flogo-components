package entity

type EventLog struct {
	Seq       string  `json:"seqno,omitempty"`
	CntrNum   string  `json:"cntrnum,omitempty"`
	RevTime   string  `json:"revtime,omitempty"`
	LogTime   string  `json:"logtime,omitempty"`
	Sp        float64 `json:"sp,omitempty"`
	Isc       int32   `json:"isc,omitempty"`
	Ss        float64 `json:"ss,omitempty"`
	Rs        float64 `json:"rs,omitempty"`
	Dss       float64 `json:"dss,omitempty"`
	Drs       float64 `json:"drs,omitempty"`
	Ambs      float64 `json:"ambs,omitempty"`
	Hus       string  `json:"hus,omitempty"`
	Sh        string  `json:"sh,omitempty"`
	Usda1     string  `json:"usda1,omitempty"`
	Usda2     string  `json:"usda2,omitempty"`
	Usda3     string  `json:"usda3,omitempty"`
	Cts       string  `json:"cts,omitempty"`
	Smode     string  `json:"smode,omitempty"`
	Isa       int32   `json:"isa,omitempty"`
	TableName string  `json:"tableName,omitempty"`

	Source string `json:"source,omitempty"`
}
