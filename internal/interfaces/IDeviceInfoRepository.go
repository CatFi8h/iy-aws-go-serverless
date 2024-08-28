package interfaces

import (
	"context"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/models"
)

type IDeviceInfoRepository interface {
	CreateDeviceInfo(ctx context.Context, deviceInfo models.DeviceInfo) error

	GetDeviceInfo(ctx context.Context, deviceInfo models.DeviceInfo) (models.DeviceInfo, error)

	UpdateDeviceInfo(ctx context.Context, deviceInfo models.DeviceInfo) error

	DeleteDeviceInfoByDeviceId(ctx context.Context, deviceInfo models.DeviceInfo) error
}
