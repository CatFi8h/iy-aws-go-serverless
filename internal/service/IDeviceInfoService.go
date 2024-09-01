package service

import (
	"context"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
)

type IDeviceInfoService interface {
	CreateDeviceInfo(ctx context.Context, deviceInfoReq model.DeviceInfo) (*model.DeviceInfo, error)

	GetDeviceInfo(ctx context.Context, deviceId string) (*model.DeviceInfo, error)

	UpdateDeviceInfo(ctx context.Context, deviceId string, deviceInfoReq model.DeviceInfo) (*model.DeviceInfo, error)

	DeleteDeviceInfo(ctx context.Context, deviceId string) error
}
