package entity

import "time"

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

const (
	GPSEVENT_SOURCE_TCPSERVER = "TCP_SERVER"
	GPSEVENT_CARRIER_COSU     = "COSU"
)

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

func GenGpsEventFromSinglePacket(singlePacket *SinglePacket, seqNo string, cntrNum string) *GpsEvent {
	gpsEvent := &GpsEvent{}

	gpsEvent.Seqno = seqNo
	gpsEvent.CntrNum = cntrNum
	gpsEvent.RevTime = time.Now().Format("2018-12-03T09:29:21+08:00")
	gpsEvent.CltTime = singlePacket.Date
	gpsEvent.LocateTime = singlePacket.Date
	gpsEvent.EleState = ReturnValueByCondition(singlePacket.PowerSupplyStatus, "1", "0")
	gpsEvent.BatLevel = singlePacket.BatLevel
	gpsEvent.OpMode = ReturnValueByCondition(singlePacket.InfoItem.OpModeValid, singlePacket.InfoItem.OpMode, "")
	gpsEvent.SetTem = ReturnValueByCondition(singlePacket.InfoItem.SetTemValid, singlePacket.InfoItem.SetTem, "")
	gpsEvent.SupTem = ReturnValueByCondition(singlePacket.InfoItem.SupTemValid, singlePacket.InfoItem.SupTem, "")
	gpsEvent.RetTem = ReturnValueByCondition(singlePacket.InfoItem.RetTemValid, singlePacket.InfoItem.RetTem, "")
	gpsEvent.Hum = ReturnValueByCondition(singlePacket.InfoItem.HumValid, singlePacket.InfoItem.Hum, "")
	gpsEvent.Lat = singlePacket.Lat
	gpsEvent.Lng = singlePacket.Lng
	gpsEvent.Speed = singlePacket.Speed
	gpsEvent.Direction = singlePacket.Direction
	gpsEvent.Hpt = ReturnValueByCondition(singlePacket.InfoItem.HptValid, singlePacket.InfoItem.Hpt, "")
	gpsEvent.Usda1 = ReturnValueByCondition(singlePacket.InfoItem.Usda1Valid, singlePacket.InfoItem.Usda1, "")
	gpsEvent.Usda2 = ReturnValueByCondition(singlePacket.InfoItem.Usda2Valid, singlePacket.InfoItem.Usda2, "")
	gpsEvent.Usda3 = ReturnValueByCondition(singlePacket.InfoItem.Usda3Valid, singlePacket.InfoItem.Usda3, "")
	gpsEvent.FaultCode = ReturnValueByCondition(singlePacket.InfoItem.FaultCodeValid, singlePacket.InfoItem.FaultCode, "")
	gpsEvent.PosFlag = ReturnValueByCondition(singlePacket.PositioningModuleStatus, "1", "0")
	gpsEvent.Ism = "0"
	gpsEvent.GpsNum = singlePacket.NumberOfSatellitesItem.GpsSatelliteNumber
	gpsEvent.BdNum = singlePacket.NumberOfSatellitesItem.BeidouSatelliteNumber
	gpsEvent.Lpt = ReturnValueByCondition(singlePacket.InfoItem.LptValid, singlePacket.InfoItem.Lpt, "")
	gpsEvent.Pt = ReturnValueByCondition(singlePacket.InfoItem.PtValid, singlePacket.InfoItem.Pt, "")
	gpsEvent.Ct1 = ReturnValueByCondition(singlePacket.InfoItem.Ct1Valid, singlePacket.InfoItem.Ct1, "")
	gpsEvent.Ct2 = ReturnValueByCondition(singlePacket.InfoItem.Ct2Valid, singlePacket.InfoItem.Ct2, "")
	gpsEvent.Ambs = ReturnValueByCondition(singlePacket.InfoItem.AmbsValid, singlePacket.InfoItem.Ambs, "")
	gpsEvent.Eis = ReturnValueByCondition(singlePacket.InfoItem.EisValid, singlePacket.InfoItem.Eis, "")
	gpsEvent.Eos = ReturnValueByCondition(singlePacket.InfoItem.EosValid, singlePacket.InfoItem.Eos, "")
	gpsEvent.Dchs = ReturnValueByCondition(singlePacket.InfoItem.DchsValid, singlePacket.InfoItem.Dchs, "")
	gpsEvent.Sgs = ReturnValueByCondition(singlePacket.InfoItem.SgsValid, singlePacket.InfoItem.Sgs, "")
	gpsEvent.Smv = ReturnValueByCondition(singlePacket.InfoItem.SmvValid, singlePacket.InfoItem.Smv, "")
	gpsEvent.Ev = ReturnValueByCondition(singlePacket.InfoItem.EvValid, singlePacket.InfoItem.Ev, "")
	gpsEvent.Dss = ReturnValueByCondition(singlePacket.InfoItem.DssValid, singlePacket.InfoItem.Dss, "")
	gpsEvent.Drs = ReturnValueByCondition(singlePacket.InfoItem.DrsValid, singlePacket.InfoItem.Drs, "")
	gpsEvent.Hs = ReturnValueByCondition(singlePacket.InfoItem.HsValid, singlePacket.InfoItem.Hs, "")
	gpsEvent.Source = GPSEVENT_SOURCE_TCPSERVER
	gpsEvent.Carrier = GPSEVENT_CARRIER_COSU

	return gpsEvent
}

func ReturnValueByCondition(condition bool, trueVal string, falseVal string) string {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
