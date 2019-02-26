package entity

import (
	"fmt"
	"strconv"

	crg "github.com/leeshing0315/go-city-reverse-geocoder"
)

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

const (
	GPSEVENT_SOURCE_TCPSERVER = "TCP_SERVER"
	GPSEVENT_CARRIER_COSU     = "COSU"
)

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

func GenGpsEventFromSinglePacket(singlePacket *SinglePacket, seqNo string, cntrNum string, nowDateStr string) *GpsEvent {
	gpsEvent := &GpsEvent{}

	gpsEvent.Seqno = seqNo
	gpsEvent.CntrNum = cntrNum
	gpsEvent.RevTime = nowDateStr
	gpsEvent.CltTime = singlePacket.Date
	gpsEvent.LocateTime = singlePacket.Date
	gpsEvent.EleState = returnValueByCondition(singlePacket.PowerSupplyStatus, "1", "0").(string)
	gpsEvent.BatLevel = singlePacket.BatLevel
	gpsEvent.OpMode = returnValueByCondition(singlePacket.InfoItem.OpModeValid, singlePacket.InfoItem.OpMode, nil).(string)
	gpsEvent.SetTem = returnValueByCondition(singlePacket.InfoItem.SetTemValid, singlePacket.InfoItem.SetTem, nil).(string)
	gpsEvent.SupTem = returnValueByCondition(singlePacket.InfoItem.SupTemValid, singlePacket.InfoItem.SupTem, nil).(string)
	gpsEvent.RetTem = returnValueByCondition(singlePacket.InfoItem.RetTemValid, singlePacket.InfoItem.RetTem, nil).(string)
	gpsEvent.Hum = returnValueByCondition(singlePacket.InfoItem.HumValid, singlePacket.InfoItem.Hum, nil).(string)
	gpsEvent.Lat = singlePacket.Lat
	gpsEvent.Lng = singlePacket.Lng
	gpsEvent.Speed = singlePacket.Speed
	gpsEvent.Direction = singlePacket.Direction
	gpsEvent.Hpt = returnValueByCondition(singlePacket.InfoItem.HptValid, singlePacket.InfoItem.Hpt, nil).(string)
	gpsEvent.Usda1 = returnValueByCondition(singlePacket.InfoItem.Usda1Valid, singlePacket.InfoItem.Usda1, nil).(string)
	gpsEvent.Usda2 = returnValueByCondition(singlePacket.InfoItem.Usda2Valid, singlePacket.InfoItem.Usda2, nil).(string)
	gpsEvent.Usda3 = returnValueByCondition(singlePacket.InfoItem.Usda3Valid, singlePacket.InfoItem.Usda3, nil).(string)
	gpsEvent.FaultCode = returnValueByCondition(singlePacket.InfoItem.FaultCodeValid, singlePacket.InfoItem.FaultCode, nil).(string)
	gpsEvent.PosFlag = returnValueByCondition(singlePacket.PositioningModuleStatus, "1", "0").(string)
	gpsEvent.Ism = "0"
	gpsEvent.GpsNum = singlePacket.NumberOfSatellitesItem.GpsSatelliteNumber
	gpsEvent.BdNum = singlePacket.NumberOfSatellitesItem.BeidouSatelliteNumber
	gpsEvent.Lpt = returnValueByCondition(singlePacket.InfoItem.LptValid, singlePacket.InfoItem.Lpt, nil).(string)
	gpsEvent.Pt = returnValueByCondition(singlePacket.InfoItem.PtValid, singlePacket.InfoItem.Pt, nil).(string)
	gpsEvent.Ct1 = returnValueByCondition(singlePacket.InfoItem.Ct1Valid, singlePacket.InfoItem.Ct1, nil).(string)
	gpsEvent.Ct2 = returnValueByCondition(singlePacket.InfoItem.Ct2Valid, singlePacket.InfoItem.Ct2, nil).(string)
	gpsEvent.Ambs = returnValueByCondition(singlePacket.InfoItem.AmbsValid, singlePacket.InfoItem.Ambs, nil).(string)
	gpsEvent.Eis = returnValueByCondition(singlePacket.InfoItem.EisValid, singlePacket.InfoItem.Eis, nil).(string)
	gpsEvent.Eos = returnValueByCondition(singlePacket.InfoItem.EosValid, singlePacket.InfoItem.Eos, nil).(string)
	gpsEvent.Dchs = returnValueByCondition(singlePacket.InfoItem.DchsValid, singlePacket.InfoItem.Dchs, nil).(string)
	gpsEvent.Sgs = returnValueByCondition(singlePacket.InfoItem.SgsValid, singlePacket.InfoItem.Sgs, nil).(string)
	gpsEvent.Smv = returnValueByCondition(singlePacket.InfoItem.SmvValid, singlePacket.InfoItem.Smv, nil).(string)
	gpsEvent.Ev = returnValueByCondition(singlePacket.InfoItem.EvValid, singlePacket.InfoItem.Ev, nil).(string)
	gpsEvent.Dss = returnValueByCondition(singlePacket.InfoItem.DssValid, singlePacket.InfoItem.Dss, nil).(string)
	gpsEvent.Drs = returnValueByCondition(singlePacket.InfoItem.DrsValid, singlePacket.InfoItem.Drs, nil).(string)
	gpsEvent.Hs = returnValueByCondition(singlePacket.InfoItem.HsValid, singlePacket.InfoItem.Hs, nil).(string)
	gpsEvent.Source = GPSEVENT_SOURCE_TCPSERVER
	gpsEvent.Carrier = GPSEVENT_CARRIER_COSU
	// address, displayName := getAddress(singlePacket.Lat, singlePacket.Lng)
	// gpsEvent.Address = address
	// gpsEvent.DisplayName = displayName.(string)
	// service.AttachLocation(gpsEvent)

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
