package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
)

type deviceInfoService struct{}

var (
	repo repository.IDeviceInfoRepository
)

func NewDeviceInfoService(repository repository.IDeviceInfoRepository) IDeviceInfoService {
	repo = repository
	return &deviceInfoService{}
}

func (service *deviceInfoService) CreateDeviceInfo(ctx context.Context, deviceInfoReq model.DeviceInfo) (*model.DeviceInfo, error) {

	currentTime := time.Now().UnixMilli()

	if deviceInfoReq.DeviceId == "" {
		log.Printf("Can not read JSON")
		return nil, errors.New("can not read JSON")
	}

	deviceInfoReq.CreateAt = currentTime
	deviceInfoReq.UpdatedAt = currentTime

	resp, err := repo.CreateDeviceInfo(ctx, deviceInfoReq)
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

	deviceInfoResp, err := repo.GetDeviceInfo(ctx, &deviceInfo)
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

	resp, err := repo.UpdateDeviceInfo(ctx, deviceInfoReq)

	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return nil, err
	}

	return resp, nil
}

func (service *deviceInfoService) DeleteDeviceInfo(ctx context.Context, deviceId string) error {
	deviceInfo := model.DeviceInfo{DeviceId: deviceId}

	err := repo.DeleteDeviceInfoByDeviceId(ctx, deviceInfo)

	if err != nil {
		log.Printf("Could not delete Device Info : [%s]", err.Error())
		return err
	}

	return nil
}
