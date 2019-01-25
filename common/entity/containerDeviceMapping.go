package entity

type ContainerDeviceMapping struct {
	DeviceId        string `json:"carno,omitempty"`
	ContainerNumber string `json:"carid,omitempty"`
	Simno           string `json:"simno,omitempty"`
	Pin             string `json:"pin,omitempty"`
	Commmode        string `json:"commmode,omitempty"`
	Model           string `json:"model,omitempty"`
	Status          string `json:"status,omitempty"`
	Regtime         string `json:"regtime,omitempty"`
	ChangeTime      string `json:"changetime,omitempty"`
	LastUpdated     string `json:"lastUpdate,omitempty"`
}
