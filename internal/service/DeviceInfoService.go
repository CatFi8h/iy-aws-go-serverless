package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
)

type deviceInfoService struct {
	repo repository.IDeviceInfoRepository
}

func NewDeviceInfoService(repository repository.IDeviceInfoRepository) IDeviceInfoService {
	return &deviceInfoService{
		repo: repository,
	}
}

func (service *deviceInfoService) CreateDeviceInfo(ctx context.Context, deviceInfoReq model.DeviceInfo) (*model.DeviceInfo, error) {

	currentTime := time.Now().UnixMilli()

	deviceInfoReq.CreateAt = currentTime
	deviceInfoReq.UpdatedAt = currentTime
	log.Printf(deviceInfoReq.DeviceId + " " + deviceInfoReq.DeviceName)
	resp, err := service.repo.CreateDeviceInfo(ctx, &deviceInfoReq)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func (service *deviceInfoService) GetDeviceInfo(ctx context.Context, deviceId string) (*model.DeviceInfo, error) {

	if deviceId == "" {
		return nil, errors.New("device ID is empty")
	}
	deviceInfo := model.DeviceInfo{DeviceId: deviceId}

	deviceInfoResp, err := service.repo.GetDeviceInfo(ctx, &deviceInfo)
	if err != nil {
		return nil, err
	}
	if deviceInfoResp == nil {
		return nil, nil
	}
	return deviceInfoResp, nil
}

func (service *deviceInfoService) UpdateDeviceInfo(ctx context.Context, deviceId string, deviceInfoReq model.DeviceInfo) (*model.DeviceInfo, error) {

	deviceInfoReq.DeviceId = deviceId
	deviceInfoReq.UpdatedAt = time.Now().UnixMilli()

	resp, err := service.repo.UpdateDeviceInfo(ctx, &deviceInfoReq)

	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return nil, err
	}

	return resp, nil
}

func (service *deviceInfoService) DeleteDeviceInfo(ctx context.Context, deviceId string) (*model.DeviceInfo, error) {
	deviceInfo := model.DeviceInfo{DeviceId: deviceId}

	resp, err := service.repo.DeleteDeviceInfoByDeviceId(ctx, &deviceInfo)

	if err != nil {
		log.Printf("Could not delete Device Info : [%s]", err.Error())
		return nil, err
	}

	return resp, nil
}
