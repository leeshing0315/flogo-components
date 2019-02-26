package entity

// GpsEvent entity gpsEvent
type GpsEvent struct {
	Seqno       string          `json:"seqno,omitempty"`
	CntrNum     string          `json:"cntrNum,omitempty"`
	RevTime     string          `json:"revTime,omitempty"`
	CltTime     string          `json:"cltTime,omitempty"`
	LocateTime  string          `json:"locateTime,omitempty"`
	EleState    string          `json:"eleState,omitempty"`
	BatLevel    string          `json:"batLevel,omitempty"`
	OpMode      string          `json:"opMode,omitempty"`
	SetTem      string          `json:"setTem,omitempty"`
	SupTem      string          `json:"supTem,omitempty"`
	RetTem      string          `json:"retTem,omitempty"`
	Hum         string          `json:"hum,omitempty"`
	Lng         string          `json:"lng,omitempty"`
	Lat         string          `json:"lat,omitempty"`
	Speed       string          `json:"speed,omitempty"`
	Direction   string          `json:"direction,omitempty"`
	PosFlag     string          `json:"posFlag,omitempty"`
	GpsNum      string          `json:"gpsNum,omitempty"`
	BdNum       string          `json:"bdNum,omitempty"`
	Source      string          `json:"source,omitempty"`
	Address     GpsEventAddress `json:"address,omitempty"`
	DisplayName string          `json:"displayName,omitempty"`
	Ambs        string          `json:"ambs,omitempty"`
	Hs          string          `json:"hs,omitempty"`
	Usda1       string          `json:"usda1,omitempty"`
	Usda2       string          `json:"usda2,omitempty"`
	Usda3       string          `json:"usda3,omitempty"`

	Hpt       string `json:"hpt,omitempty"`
	FaultCode string `json:"faultCode,omitempty"`
	Ism       string `json:"ism,omitempty"`
	FromDate  string `json:"fromDate,omitempty"`
	ToDate    string `json:"toDate,omitempty"`
	Carrier   string `json:"carrier,omitempty"`
	Lpt       string `json:"lpt,omitempty"`
	Pt        string `json:"pt,omitempty"`
	Ct1       string `json:"ct1,omitempty"`
	Ct2       string `json:"ct2,omitempty"`
	Eis       string `json:"eis,omitempty"`
	Eos       string `json:"eos,omitempty"`
	Dchs      string `json:"dchs,omitempty"`
	Sgs       string `json:"sgs,omitempty"`
	Smv       string `json:"smv,omitempty"`
	Ev        string `json:"ev,omitempty"`
	Dss       string `json:"dss,omitempty"`
	Drs       string `json:"drs,omitempty"`

	Isc        string `json:"isc,omitempty"`
	Isa        string `json:"isa,omitempty"`
	Cts        string `json:"cts,omitempty"`
	IsEventLog bool   `json:"isEventLog"`
	CreatedAt  string `json:"createdAt,omitempty"`
	Smode      string `json:"smode,omitempty"`
	EventLog   string `json:"eventLog,omitempty"`
}

// GpsEventAddress entity gpsEventAddress
type GpsEventAddress struct {
	Distance  float64 `json:"distance"`
	Longitude float64
	Latitude  float64

	Code        string `json:"code"`
	Name        string `json:"name"`
	City        string `json:"city"`
	RegionCode  string `json:"region_code"`
	Region      string `json:"region"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
}
