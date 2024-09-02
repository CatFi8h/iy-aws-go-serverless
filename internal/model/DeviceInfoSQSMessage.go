package model

type DeviceInfoSQSMessage struct {
	DeviceId string `json:"deviceId"`
	HomeId   string `json:"homeId"`
}
