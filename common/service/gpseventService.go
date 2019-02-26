package service

import "github.com/leeshing0315/flogo-components/common/entity"

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

func GenGpsEventFromSinglePacket(singlePacket *entity.SinglePacket, seqNo string, cntrNum string, nowDateStr string) *entity.GpsEvent {
	gpsEvent := &entity.GpsEvent{}

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
