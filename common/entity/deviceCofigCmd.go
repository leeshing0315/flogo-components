package entity

import (
	"encoding/json"
)

// DeviceConfigCmd entity
type DeviceConfigCmd struct {
	DeviceID       string      `json:"devid"`
	SeqNo          json.Number `json:"seqno"`
	Subcmd         string      `json:"subcmd"`
	Value          string      `json:"value"`
	SendFlag       string      `json:"sendflag"`
	SendTime       string      `json:"sendtime"`
	LastUpdateTime string      `json:"lastupdatetime"`
}
