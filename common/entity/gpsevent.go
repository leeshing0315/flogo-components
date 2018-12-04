package entity

// GpsEvent entity gpsEvent
type GpsEvent struct {
	ID string `json:"id"`

	Seqno       string          `json:"seqno"`
	CntrNum     string          `json:"cntrNum"`
	RevTime     string          `json:"revTime"`
	CltTime     string          `json:"cltTime"`
	LocateTime  string          `json:"locateTime"`
	EleState    string          `json:"eleState"`
	BatLevel    string          `json:"batLevel"`
	OpMode      string          `json:"opMode"`
	SetTem      string          `json:"setTem"`
	SupTem      string          `json:"supTem"`
	RetTem      string          `json:"retTem"`
	Hum         string          `json:"hum"`
	Lng         string          `json:"lng"`
	Lat         string          `json:"lat"`
	Speed       string          `json:"speed"`
	Direction   string          `json:"direction"`
	PosFlag     string          `json:"posFlag"`
	GpsNum      string          `json:"gpsNum"`
	BdNum       string          `json:"bdNum"`
	Source      string          `json:"source"`
	Address     GpsEventAddress `json:"address"`
	DisplayName string          `json:"displayName"`
	Ambs        string          `json:"ambs"`
	Hs          string          `json:"hs"`
	Usda1       string          `json:"usda1"`
	Usda2       string          `json:"usda2"`
	Usda3       string          `json:"usda3"`

	Hpt       string `json:"hpt"`
	FaultCode string `json:"faultCode"`
	Ism       string `json:"ism"`
	FromDate  string `json:"fromDate"`
	ToDate    string `json:"toDate"`
	Carrier   string `json:"carrier"`
	Lpt       string `json:"lpt"`
	Pt        string `json:"pt"`
	Ct1       string `json:"ct1"`
	Ct2       string `json:"ct2"`
	Eis       string `json:"eis"`
	Eos       string `json:"eos"`
	Dchs      string `json:"dchs"`
	Sgs       string `json:"sgs"`
	Smv       string `json:"smv"`
	Ev        string `json:"ev"`
	Dss       string `json:"dss"`
	Drs       string `json:"drs"`
}

// GpsEventAddress entity gpsEventAddress
type GpsEventAddress struct {
	distance  float64
	longitude float64
	latitude  float64

	city        string
	RegionCode  string `json:"region_code"`
	region      string
	CountryCode string `json:"country_code"`
	country     string
}
