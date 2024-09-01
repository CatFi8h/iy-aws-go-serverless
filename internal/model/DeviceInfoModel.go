package model

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DeviceInfo struct {
	DeviceId   string `dynamodbav:"deviceid" json:"deviceid"`
	DeviceName string `dynamodbav:"deviceName" json:"deviceName"`
	Mac        string `dynamodbav:"mac" json:"mac"`
	DeviceType string `dynamodbav:"deviceType" json:"deviceType"`
	HomeId     string `dynamodbav:"homeId" json:"homeId"`
	CreateAt   int64  `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt  int64  `dynamodbav:"updatedAt" json:"updateAt"`
}

type DeviceInfoResponse struct {
	Data string
}

func (deviceInfo DeviceInfo) GetKey() map[string]types.AttributeValue {
	deviceId, err := attributevalue.Marshal(deviceInfo.DeviceId)
	if err != nil {
		panic("DeviceId Not parsed " + err.Error())
	}

	return map[string]types.AttributeValue{"deviceid": deviceId}
}
