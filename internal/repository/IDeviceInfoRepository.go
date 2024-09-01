package repository

import (
	"context"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
)

type IDeviceInfoRepository interface {
	CreateDeviceInfo(ctx context.Context, deviceInfo model.DeviceInfo) (*model.DeviceInfo, error)

	UpdateDeviceInfo(ctx context.Context, deviceInfo model.DeviceInfo) (*model.DeviceInfo, error)

	GetDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error)

	DeleteDeviceInfoByDeviceId(ctx context.Context, deviceInfo model.DeviceInfo) error
}
