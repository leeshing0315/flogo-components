package entity

import (
	"encoding/json"
)

// DeviceConfigCmd entity
type DeviceConfigCmd struct {
	DeviceID       string      `json:"devid,omitempty"`
	SeqNo          json.Number `json:"seqno"`
	Subcmd         string      `json:"subcmd,omitempty"`
	Value          string      `json:"value,omitempty"`
	SendFlag       string      `json:"sendflag,omitempty"`
	SendTime       string      `json:"sendtime,omitempty"`
	LastUpdateTime string      `json:"lastupdatetime,omitempty"`
}
