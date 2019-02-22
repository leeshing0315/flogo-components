package entity

import (
	"fmt"
	"strconv"

	crg "github.com/leeshing0315/go-city-reverse-geocoder"
)

// GpsEvent entity gpsEvent
type GpsEvent struct {
	Seqno       interface{}     `json:"seqno,omitempty"`
	CntrNum     interface{}     `json:"cntrNum,omitempty"`
	RevTime     interface{}     `json:"revTime,omitempty"`
	CltTime     interface{}     `json:"cltTime,omitempty"`
	LocateTime  interface{}     `json:"locateTime,omitempty"`
	EleState    interface{}     `json:"eleState,omitempty"`
	BatLevel    interface{}     `json:"batLevel,omitempty"`
	OpMode      interface{}     `json:"opMode,omitempty"`
	SetTem      interface{}     `json:"setTem,omitempty"`
	SupTem      interface{}     `json:"supTem,omitempty"`
	RetTem      interface{}     `json:"retTem,omitempty"`
	Hum         interface{}     `json:"hum,omitempty"`
	Lng         interface{}     `json:"lng,omitempty"`
	Lat         interface{}     `json:"lat,omitempty"`
	Speed       interface{}     `json:"speed,omitempty"`
	Direction   interface{}     `json:"direction,omitempty"`
	PosFlag     interface{}     `json:"posFlag,omitempty"`
	GpsNum      interface{}     `json:"gpsNum,omitempty"`
	BdNum       interface{}     `json:"bdNum,omitempty"`
	Source      interface{}     `json:"source,omitempty"`
	Address     GpsEventAddress `json:"address,omitempty"`
	DisplayName interface{}     `json:"displayName,omitempty"`
	Ambs        interface{}     `json:"ambs,omitempty"`
	Hs          interface{}     `json:"hs,omitempty"`
	Usda1       interface{}     `json:"usda1,omitempty"`
	Usda2       interface{}     `json:"usda2,omitempty"`
	Usda3       interface{}     `json:"usda3,omitempty"`

	Hpt       interface{} `json:"hpt,omitempty"`
	FaultCode interface{} `json:"faultCode,omitempty"`
	Ism       interface{} `json:"ism,omitempty"`
	FromDate  interface{} `json:"fromDate,omitempty"`
	ToDate    interface{} `json:"toDate,omitempty"`
	Carrier   interface{} `json:"carrier,omitempty"`
	Lpt       interface{} `json:"lpt,omitempty"`
	Pt        interface{} `json:"pt,omitempty"`
	Ct1       interface{} `json:"ct1,omitempty"`
	Ct2       interface{} `json:"ct2,omitempty"`
	Eis       interface{} `json:"eis,omitempty"`
	Eos       interface{} `json:"eos,omitempty"`
	Dchs      interface{} `json:"dchs,omitempty"`
	Sgs       interface{} `json:"sgs,omitempty"`
	Smv       interface{} `json:"smv,omitempty"`
	Ev        interface{} `json:"ev,omitempty"`
	Dss       interface{} `json:"dss,omitempty"`
	Drs       interface{} `json:"drs,omitempty"`

	Isc        interface{} `json:"isc,omitempty"`
	Isa        interface{} `json:"isa,omitempty"`
	Smode      interface{} `json:"smode,omitempty"`
	Cts        interface{} `json:"cts,omitempty"`
	IsEventLog bool        `json:"isEventLog"`
	CreatedAt  interface{} `json:"createdAt,omitempty"`
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
	gpsEvent.EleState = returnValueByCondition(singlePacket.PowerSupplyStatus, "1", "0")
	gpsEvent.BatLevel = singlePacket.BatLevel
	gpsEvent.OpMode = returnValueByCondition(singlePacket.InfoItem.OpModeValid, singlePacket.InfoItem.OpMode, nil)
	gpsEvent.SetTem = returnValueByCondition(singlePacket.InfoItem.SetTemValid, singlePacket.InfoItem.SetTem, nil)
	gpsEvent.SupTem = returnValueByCondition(singlePacket.InfoItem.SupTemValid, singlePacket.InfoItem.SupTem, nil)
	gpsEvent.RetTem = returnValueByCondition(singlePacket.InfoItem.RetTemValid, singlePacket.InfoItem.RetTem, nil)
	gpsEvent.Hum = returnValueByCondition(singlePacket.InfoItem.HumValid, singlePacket.InfoItem.Hum, nil)
	gpsEvent.Lat = singlePacket.Lat
	gpsEvent.Lng = singlePacket.Lng
	gpsEvent.Speed = singlePacket.Speed
	gpsEvent.Direction = singlePacket.Direction
	gpsEvent.Hpt = returnValueByCondition(singlePacket.InfoItem.HptValid, singlePacket.InfoItem.Hpt, nil)
	gpsEvent.Usda1 = returnValueByCondition(singlePacket.InfoItem.Usda1Valid, singlePacket.InfoItem.Usda1, nil)
	gpsEvent.Usda2 = returnValueByCondition(singlePacket.InfoItem.Usda2Valid, singlePacket.InfoItem.Usda2, nil)
	gpsEvent.Usda3 = returnValueByCondition(singlePacket.InfoItem.Usda3Valid, singlePacket.InfoItem.Usda3, nil)
	gpsEvent.FaultCode = returnValueByCondition(singlePacket.InfoItem.FaultCodeValid, singlePacket.InfoItem.FaultCode, nil)
	gpsEvent.PosFlag = returnValueByCondition(singlePacket.PositioningModuleStatus, "1", "0")
	gpsEvent.Ism = "0"
	gpsEvent.GpsNum = singlePacket.NumberOfSatellitesItem.GpsSatelliteNumber
	gpsEvent.BdNum = singlePacket.NumberOfSatellitesItem.BeidouSatelliteNumber
	gpsEvent.Lpt = returnValueByCondition(singlePacket.InfoItem.LptValid, singlePacket.InfoItem.Lpt, nil)
	gpsEvent.Pt = returnValueByCondition(singlePacket.InfoItem.PtValid, singlePacket.InfoItem.Pt, nil)
	gpsEvent.Ct1 = returnValueByCondition(singlePacket.InfoItem.Ct1Valid, singlePacket.InfoItem.Ct1, nil)
	gpsEvent.Ct2 = returnValueByCondition(singlePacket.InfoItem.Ct2Valid, singlePacket.InfoItem.Ct2, nil)
	gpsEvent.Ambs = returnValueByCondition(singlePacket.InfoItem.AmbsValid, singlePacket.InfoItem.Ambs, nil)
	gpsEvent.Eis = returnValueByCondition(singlePacket.InfoItem.EisValid, singlePacket.InfoItem.Eis, nil)
	gpsEvent.Eos = returnValueByCondition(singlePacket.InfoItem.EosValid, singlePacket.InfoItem.Eos, nil)
	gpsEvent.Dchs = returnValueByCondition(singlePacket.InfoItem.DchsValid, singlePacket.InfoItem.Dchs, nil)
	gpsEvent.Sgs = returnValueByCondition(singlePacket.InfoItem.SgsValid, singlePacket.InfoItem.Sgs, nil)
	gpsEvent.Smv = returnValueByCondition(singlePacket.InfoItem.SmvValid, singlePacket.InfoItem.Smv, nil)
	gpsEvent.Ev = returnValueByCondition(singlePacket.InfoItem.EvValid, singlePacket.InfoItem.Ev, nil)
	gpsEvent.Dss = returnValueByCondition(singlePacket.InfoItem.DssValid, singlePacket.InfoItem.Dss, nil)
	gpsEvent.Drs = returnValueByCondition(singlePacket.InfoItem.DrsValid, singlePacket.InfoItem.Drs, nil)
	gpsEvent.Hs = returnValueByCondition(singlePacket.InfoItem.HsValid, singlePacket.InfoItem.Hs, nil)
	gpsEvent.Source = GPSEVENT_SOURCE_TCPSERVER
	gpsEvent.Carrier = GPSEVENT_CARRIER_COSU
	address, displayName := getAddress(singlePacket.Lat, singlePacket.Lng)
	gpsEvent.Address = address
	gpsEvent.DisplayName = displayName

	return gpsEvent
}

func returnValueByCondition(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition && (trueVal.(string) != "") {
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
