package interfaces

import (
	"context"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/models"
)

type IDeviceInfoService interface {
	CreateDeviceInfo(bodyStr string) error

	GetDeviceInfo(ctx context.Context, deviceId string) (models.DeviceInfoResponse, error)

	DeleteDeviceInfo(deviceId string) error

	UpdateDeviceInfo(deviceId string, bodyStr string) error
}
