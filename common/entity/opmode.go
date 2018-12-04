package entity

var opModeMapping = map[byte]string{
	0x0: "Defrost",
	0x1: "Full Cool",
	0x2: "Modulation",
	0x3: "Thermo Off",
	0x4: "Remote Unit Off",
	0x5: "Unit Off",
	0x6: "Emergency Stop (Automatically reset)",
	0x7: "Emergency Shut Down",
	0x8: "Manual check",
	0x9: "Short PTI",
	0xa: "Full PTI",
	0xb: "Emergency Stop in PTI",
	0xc: "Stop in PTI",
	0xd: "Heating",
	0xe: "Non Control",
	0xf: "Electric Power Shut Off",
}
