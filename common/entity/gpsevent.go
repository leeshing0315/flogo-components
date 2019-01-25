package entity

import (
	"fmt"
	"strconv"

	crg "github.com/leeshing0315/go-city-reverse-geocoder"
)

// GpsEvent entity gpsEvent
type GpsEvent struct {
	Seqno       interface{}     `json:"seqno"`
	CntrNum     interface{}     `json:"cntrNum"`
	RevTime     interface{}     `json:"revTime"`
	CltTime     interface{}     `json:"cltTime"`
	LocateTime  interface{}     `json:"locateTime"`
	EleState    interface{}     `json:"eleState"`
	BatLevel    interface{}     `json:"batLevel"`
	OpMode      interface{}     `json:"opMode"`
	SetTem      interface{}     `json:"setTem"`
	SupTem      interface{}     `json:"supTem"`
	RetTem      interface{}     `json:"retTem"`
	Hum         interface{}     `json:"hum"`
	Lng         interface{}     `json:"lng"`
	Lat         interface{}     `json:"lat"`
	Speed       interface{}     `json:"speed"`
	Direction   interface{}     `json:"direction"`
	PosFlag     interface{}     `json:"posFlag"`
	GpsNum      interface{}     `json:"gpsNum"`
	BdNum       interface{}     `json:"bdNum"`
	Source      interface{}     `json:"source"`
	Address     GpsEventAddress `json:"address"`
	DisplayName interface{}     `json:"displayName"`
	Ambs        interface{}     `json:"ambs"`
	Hs          interface{}     `json:"hs"`
	Usda1       interface{}     `json:"usda1"`
	Usda2       interface{}     `json:"usda2"`
	Usda3       interface{}     `json:"usda3"`

	Hpt       interface{} `json:"hpt"`
	FaultCode interface{} `json:"faultCode"`
	Ism       interface{} `json:"ism"`
	FromDate  interface{} `json:"fromDate"`
	ToDate    interface{} `json:"toDate"`
	Carrier   interface{} `json:"carrier"`
	Lpt       interface{} `json:"lpt"`
	Pt        interface{} `json:"pt"`
	Ct1       interface{} `json:"ct1"`
	Ct2       interface{} `json:"ct2"`
	Eis       interface{} `json:"eis"`
	Eos       interface{} `json:"eos"`
	Dchs      interface{} `json:"dchs"`
	Sgs       interface{} `json:"sgs"`
	Smv       interface{} `json:"smv"`
	Ev        interface{} `json:"ev"`
	Dss       interface{} `json:"dss"`
	Drs       interface{} `json:"drs"`
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

func GenGpsEventFromSinglePacket(singlePacket *SinglePacket, seqNo string, cntrNum string, nowDateStr string) *GpsEvent {
	gpsEvent := &GpsEvent{}

	gpsEvent.Seqno = seqNo
	gpsEvent.CntrNum = cntrNum
	gpsEvent.RevTime = nowDateStr
	gpsEvent.CltTime = singlePacket.Date
	gpsEvent.LocateTime = singlePacket.Date
	gpsEvent.EleState = ReturnValueByCondition(singlePacket.PowerSupplyStatus, "1", "0")
	gpsEvent.BatLevel = singlePacket.BatLevel
	gpsEvent.OpMode = ReturnValueByCondition(singlePacket.InfoItem.OpModeValid, singlePacket.InfoItem.OpMode, nil)
	gpsEvent.SetTem = ReturnValueByCondition(singlePacket.InfoItem.SetTemValid, singlePacket.InfoItem.SetTem, nil)
	gpsEvent.SupTem = ReturnValueByCondition(singlePacket.InfoItem.SupTemValid, singlePacket.InfoItem.SupTem, nil)
	gpsEvent.RetTem = ReturnValueByCondition(singlePacket.InfoItem.RetTemValid, singlePacket.InfoItem.RetTem, nil)
	gpsEvent.Hum = ReturnValueByCondition(singlePacket.InfoItem.HumValid, singlePacket.InfoItem.Hum, nil)
	gpsEvent.Lat = singlePacket.Lat
	gpsEvent.Lng = singlePacket.Lng
	gpsEvent.Speed = singlePacket.Speed
	gpsEvent.Direction = singlePacket.Direction
	gpsEvent.Hpt = ReturnValueByCondition(singlePacket.InfoItem.HptValid, singlePacket.InfoItem.Hpt, nil)
	gpsEvent.Usda1 = ReturnValueByCondition(singlePacket.InfoItem.Usda1Valid, singlePacket.InfoItem.Usda1, nil)
	gpsEvent.Usda2 = ReturnValueByCondition(singlePacket.InfoItem.Usda2Valid, singlePacket.InfoItem.Usda2, nil)
	gpsEvent.Usda3 = ReturnValueByCondition(singlePacket.InfoItem.Usda3Valid, singlePacket.InfoItem.Usda3, nil)
	gpsEvent.FaultCode = ReturnValueByCondition(singlePacket.InfoItem.FaultCodeValid, singlePacket.InfoItem.FaultCode, nil)
	gpsEvent.PosFlag = ReturnValueByCondition(singlePacket.PositioningModuleStatus, "1", "0")
	gpsEvent.Ism = "0"
	gpsEvent.GpsNum = singlePacket.NumberOfSatellitesItem.GpsSatelliteNumber
	gpsEvent.BdNum = singlePacket.NumberOfSatellitesItem.BeidouSatelliteNumber
	gpsEvent.Lpt = ReturnValueByCondition(singlePacket.InfoItem.LptValid, singlePacket.InfoItem.Lpt, nil)
	gpsEvent.Pt = ReturnValueByCondition(singlePacket.InfoItem.PtValid, singlePacket.InfoItem.Pt, nil)
	gpsEvent.Ct1 = ReturnValueByCondition(singlePacket.InfoItem.Ct1Valid, singlePacket.InfoItem.Ct1, nil)
	gpsEvent.Ct2 = ReturnValueByCondition(singlePacket.InfoItem.Ct2Valid, singlePacket.InfoItem.Ct2, nil)
	gpsEvent.Ambs = ReturnValueByCondition(singlePacket.InfoItem.AmbsValid, singlePacket.InfoItem.Ambs, nil)
	gpsEvent.Eis = ReturnValueByCondition(singlePacket.InfoItem.EisValid, singlePacket.InfoItem.Eis, nil)
	gpsEvent.Eos = ReturnValueByCondition(singlePacket.InfoItem.EosValid, singlePacket.InfoItem.Eos, nil)
	gpsEvent.Dchs = ReturnValueByCondition(singlePacket.InfoItem.DchsValid, singlePacket.InfoItem.Dchs, nil)
	gpsEvent.Sgs = ReturnValueByCondition(singlePacket.InfoItem.SgsValid, singlePacket.InfoItem.Sgs, nil)
	gpsEvent.Smv = ReturnValueByCondition(singlePacket.InfoItem.SmvValid, singlePacket.InfoItem.Smv, nil)
	gpsEvent.Ev = ReturnValueByCondition(singlePacket.InfoItem.EvValid, singlePacket.InfoItem.Ev, nil)
	gpsEvent.Dss = ReturnValueByCondition(singlePacket.InfoItem.DssValid, singlePacket.InfoItem.Dss, nil)
	gpsEvent.Drs = ReturnValueByCondition(singlePacket.InfoItem.DrsValid, singlePacket.InfoItem.Drs, nil)
	gpsEvent.Hs = ReturnValueByCondition(singlePacket.InfoItem.HsValid, singlePacket.InfoItem.Hs, nil)
	gpsEvent.Source = GPSEVENT_SOURCE_TCPSERVER
	gpsEvent.Carrier = GPSEVENT_CARRIER_COSU
	address, displayName := getAddress(singlePacket.Lat, singlePacket.Lng)
	gpsEvent.Address = address
	gpsEvent.DisplayName = displayName

	return gpsEvent
}

func ReturnValueByCondition(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}

func getAddress(latStr string, lonStr string) (address GpsEventAddress, displayName interface{}) {
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return GpsEventAddress{}, nil
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return GpsEventAddress{}, nil
	}
	results, err := crg.GetNearestCities(lat, lon, 1, "mi")
	result := results[0]
	return GpsEventAddress{
		RegionCode:  result.Region_code,
		CountryCode: result.Country_code,
	}, fmt.Sprintf("%v, %v, %v", result.City, result.Region, result.Country)
}
