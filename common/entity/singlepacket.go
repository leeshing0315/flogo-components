package entity

const (
	positioningModuleFailureMask   byte = 0x80 // 1000 0000
	serialCommunicationFailureMask byte = 0x40 // 0100 0000
	communicationModuleFailureMask byte = 0x20 // 0010 0000
	powerSupplyFailureMask         byte = 0x10 // 0001 0000
	batteryChargingFailureMask     byte = 0x08 // 0000 1000
	clockModuleFailureMask         byte = 0x04 // 0000 0100
	coldBoxFaultCodeChangeMask     byte = 0x02 // 0000 0010
	coldBoxOperationModeChangeMask byte = 0x01 // 0000 0001
	powerSupplyStatusChangeMask    byte = 0x80 // 1000 0000

	positioningMask                      byte = 0x40 // 0100 0000
	latitudeNorthSouthMask               byte = 0x20 // 0010 0000
	longitudeEastWestMask                byte = 0x10 // 0001 0000
	useGpsSatellitesForPositioningMask   byte = 0x08 // 0000 1000
	useBeidouSatelliteForPositioningMask byte = 0x04 // 0000 0100
	supplyByBatteryOrPowerMask           byte = 0x02 // 0000 0010
	inThePolygonAreaMask                 byte = 0x01 // 0000 0001
	positioningModuleStatusMask          byte = 0x80 // 1000 0000
	serialCommunicationStatusMask        byte = 0x40 // 0100 0000
	communicationModuleStatusMask        byte = 0x20 // 0010 0000
	powerSupplyStatusMask                byte = 0x10 // 0001 0000
	batteryChargingStatusMask            byte = 0x08 // 0000 1000
	clockModuleStatusMask                byte = 0x04 // 0000 0100
	timedUploadDataMask                  byte = 0x01 // 0000 0001

	Bit0Mask byte = 0x01 // 0000 0001
	Bit1Mask byte = 0x02 // 0000 0010
	Bit2Mask byte = 0x04 // 0000 0100
	Bit3Mask byte = 0x08 // 0000 1000
	Bit4Mask byte = 0x10 // 0001 0000
	Bit5Mask byte = 0x20 // 0010 0000
	Bit6Mask byte = 0x40 // 0100 0000
	Bit7Mask byte = 0x80 // 1000 0000

	LOW_4_BITS_MASK byte = 0x0F // 0000 1111
)

var BitMask = []byte{
	Bit0Mask,
	Bit1Mask,
	Bit2Mask,
	Bit3Mask,
	Bit4Mask,
	Bit5Mask,
	Bit6Mask,
	Bit7Mask,
}

type SinglePacket struct {
	PositioningModuleFailure   bool
	SerialCommunicationFailure bool
	CommunicationModuleFailure bool
	PowerSupplyFailure         bool
	BatteryChargingFailure     bool
	ClockModuleFailure         bool
	ColdBoxFaultCodeChange     bool
	ColdBoxOperationModeChange bool
	PowerSupplyStatusChange    bool

	Positioning                      bool
	LatitudeNorthSouth               bool
	LongitudeEastWest                bool
	UseGpsSatellitesForPositioning   bool
	UseBeidouSatelliteForPositioning bool
	SupplyByBatteryOrPower           bool
	InThePolygonArea                 bool
	PositioningModuleStatus          bool
	SerialCommunicationStatus        bool
	CommunicationModuleStatus        bool
	PowerSupplyStatus                bool
	BatteryChargingStatus            bool
	ClockModuleStatus                bool
	TimedUploadData                  bool

	Date      string
	Lat       string
	Lng       string
	Speed     string
	Direction string
	BatLevel  string

	LoginItem              LoginItem
	InfoItem               InfoItem
	DebugTextItem          DebugTextItem
	NumberOfSatellitesItem NumberOfSatellitesItem
	OpModeItem             OpModeItem
	FaultCodeItem          FaultCodeItem
	ColdBoxTimeItem        ColdBoxTimeItem
}

type LoginItem struct {
	Pin             string
	DeviceID        string
	ContainerNumber string
	ReeferType      string
}

type InfoItem struct {
	ReeferType     string
	OpModeValid    bool
	OpMode         string
	SetTemValid    bool
	SetTem         string
	SupTemValid    bool
	SupTem         string
	RetTemValid    bool
	RetTem         string
	HumValid       bool
	Hum            string
	HptValid       bool
	Hpt            string
	Usda1Valid     bool
	Usda1          string
	Usda2Valid     bool
	Usda2          string
	Usda3Valid     bool
	Usda3          string
	CtlTypeValid   bool
	CtlType        string
	LptValid       bool
	Lpt            string
	PtValid        bool
	Pt             string
	Ct1Valid       bool
	Ct1            string
	Ct2Valid       bool
	Ct2            string
	AmbsValid      bool
	Ambs           string
	EisValid       bool
	Eis            string
	EosValid       bool
	Eos            string
	DchsValid      bool
	Dchs           string
	SgsValid       bool
	Sgs            string
	SmvValid       bool
	Smv            string
	EvValid        bool
	Ev             string
	DssValid       bool
	Dss            string
	DrsValid       bool
	Drs            string
	HsValid        bool
	Hs             string
	FaultCodeValid bool
	FaultCode      string
}

type DebugTextItem struct {
	DebugText string
}

type NumberOfSatellitesItem struct {
	GpsSatelliteNumber    string
	BeidouSatelliteNumber string
}

type OpModeItem struct {
	OpMode string
}

type FaultCodeItem struct {
	FaultCode string
}

type ColdBoxTimeItem struct {
	CntrTime           string
	CollectColdBoxTime string
}
