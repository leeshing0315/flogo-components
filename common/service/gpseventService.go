package service

import (
	"time"

	"github.com/leeshing0315/flogo-components/common/entity"
)

const (
	GPSEVENT_SOURCE_TCPSERVER = "TCP_SERVER"
	GPSEVENT_CARRIER_COSU     = "COSU"
)

// GpsEventAddress entity gpsEventAddress
type GpsEventAddress struct {
	Distance  float64 `json:"distance"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`

	Code        string `json:"code"`
	Name        string `json:"name"`
	City        string `json:"city"`
	RegionCode  string `json:"region_code"`
	Region      string `json:"region"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
}

func GenGpsEventFromSinglePacket(singlePacket *entity.SinglePacket, seqNo string, cntrNum string, nowDateStr string, carrier string) *entity.GpsEvent {
	gpsEvent := &entity.GpsEvent{}

	gpsEvent.Seqno = seqNo
	gpsEvent.CntrNum = cntrNum
	gpsEvent.RevTime = nowDateStr
	gpsEvent.CltTime = singlePacket.Date
	gpsEvent.LocateTime = singlePacket.Date
	// gpsEvent.EleState = returnValueByCondition(singlePacket.PowerSupplyStatus, "1", "0").(string)
	gpsEvent.EleState = returnValueByCondition(singlePacket.SupplyByBatteryOrPower, "1", "0").(string)
	gpsEvent.BatLevel = singlePacket.BatLevel
	gpsEvent.OpMode = returnValueByCondition(singlePacket.InfoItem.OpModeValid, singlePacket.InfoItem.OpMode, "").(string)
	gpsEvent.SetTem = returnValueByCondition(singlePacket.InfoItem.SetTemValid, singlePacket.InfoItem.SetTem, "").(string)
	gpsEvent.SupTem = returnValueByCondition(singlePacket.InfoItem.SupTemValid, singlePacket.InfoItem.SupTem, "").(string)
	gpsEvent.RetTem = returnValueByCondition(singlePacket.InfoItem.RetTemValid, singlePacket.InfoItem.RetTem, "").(string)
	gpsEvent.Hum = returnValueByCondition(singlePacket.InfoItem.HumValid, singlePacket.InfoItem.Hum, "").(string)

	// PosFlag
	gpsEvent.PosFlag = returnValueByCondition(singlePacket.Positioning, "1", "0").(string)

	if gpsEvent.PosFlag == "1" {
		if singlePacket.LatitudeNorthSouth == true {
			// 南纬
			gpsEvent.Lat = "-" + singlePacket.Lat
		} else {
			// 北纬
			gpsEvent.Lat = singlePacket.Lat
		}
		if singlePacket.LongitudeEastWest == true {
			// 西经
			gpsEvent.Lng = "-" + singlePacket.Lng
		} else {
			// 东经
			gpsEvent.Lng = singlePacket.Lng
		}
	} else {
		gpsEvent.Lat = "0"
		gpsEvent.Lng = "0"
	}

	gpsEvent.Speed = singlePacket.Speed
	gpsEvent.Direction = singlePacket.Direction
	gpsEvent.Hpt = returnValueByCondition(singlePacket.InfoItem.HptValid, singlePacket.InfoItem.Hpt, "").(string)
	gpsEvent.Usda1 = returnValueByCondition(singlePacket.InfoItem.Usda1Valid, singlePacket.InfoItem.Usda1, "").(string)
	gpsEvent.Usda2 = returnValueByCondition(singlePacket.InfoItem.Usda2Valid, singlePacket.InfoItem.Usda2, "").(string)
	gpsEvent.Usda3 = returnValueByCondition(singlePacket.InfoItem.Usda3Valid, singlePacket.InfoItem.Usda3, "").(string)
	gpsEvent.FaultCode = returnValueByCondition(singlePacket.InfoItem.FaultCodeValid, singlePacket.InfoItem.FaultCode, "").(string)

	gpsEvent.Ism = "0"
	gpsEvent.GpsNum = singlePacket.NumberOfSatellitesItem.GpsSatelliteNumber
	gpsEvent.BdNum = singlePacket.NumberOfSatellitesItem.BeidouSatelliteNumber
	gpsEvent.Lpt = returnValueByCondition(singlePacket.InfoItem.LptValid, singlePacket.InfoItem.Lpt, "").(string)
	gpsEvent.Pt = returnValueByCondition(singlePacket.InfoItem.PtValid, singlePacket.InfoItem.Pt, "").(string)
	gpsEvent.Ct1 = returnValueByCondition(singlePacket.InfoItem.Ct1Valid, singlePacket.InfoItem.Ct1, "").(string)
	gpsEvent.Ct2 = returnValueByCondition(singlePacket.InfoItem.Ct2Valid, singlePacket.InfoItem.Ct2, "").(string)
	gpsEvent.Ambs = returnValueByCondition(singlePacket.InfoItem.AmbsValid, singlePacket.InfoItem.Ambs, "").(string)
	gpsEvent.Eis = returnValueByCondition(singlePacket.InfoItem.EisValid, singlePacket.InfoItem.Eis, "").(string)
	gpsEvent.Eos = returnValueByCondition(singlePacket.InfoItem.EosValid, singlePacket.InfoItem.Eos, "").(string)
	gpsEvent.Dchs = returnValueByCondition(singlePacket.InfoItem.DchsValid, singlePacket.InfoItem.Dchs, "").(string)
	gpsEvent.Sgs = returnValueByCondition(singlePacket.InfoItem.SgsValid, singlePacket.InfoItem.Sgs, "").(string)
	gpsEvent.Smv = returnValueByCondition(singlePacket.InfoItem.SmvValid, singlePacket.InfoItem.Smv, "").(string)
	gpsEvent.Ev = returnValueByCondition(singlePacket.InfoItem.EvValid, singlePacket.InfoItem.Ev, "").(string)
	gpsEvent.Dss = returnValueByCondition(singlePacket.InfoItem.DssValid, singlePacket.InfoItem.Dss, "").(string)
	gpsEvent.Drs = returnValueByCondition(singlePacket.InfoItem.DrsValid, singlePacket.InfoItem.Drs, "").(string)
	gpsEvent.Hs = returnValueByCondition(singlePacket.InfoItem.HsValid, singlePacket.InfoItem.Hs, "").(string)
	gpsEvent.Source = GPSEVENT_SOURCE_TCPSERVER
	if carrier != "" {
		gpsEvent.Carrier = carrier
	} else {
		gpsEvent.Carrier = GPSEVENT_CARRIER_COSU
	}
	// address, displayName := getAddress(singlePacket.Lat, singlePacket.Lng)
	// gpsEvent.Address = address
	// gpsEvent.DisplayName = displayName.(string)
	if gpsEvent.PosFlag != "0" {
		AttachLocation(gpsEvent)
	}
	gpsEvent.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	return gpsEvent
}

func returnValueByCondition(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition && (trueVal.(string) != "") {
		return trueVal
	} else {
		return falseVal
	}
}
