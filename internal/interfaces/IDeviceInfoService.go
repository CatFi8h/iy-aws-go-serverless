package interfaces

import (
	"context"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
)

type IDeviceInfoService interface {
	CreateDeviceInfo(bodyStr string) error

	GetDeviceInfo(ctx context.Context, deviceId string) (model.DeviceInfoResponse, error)

	DeleteDeviceInfo(deviceId string) error

	UpdateDeviceInfo(deviceId string, bodyStr string) error
}
