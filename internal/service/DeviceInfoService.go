package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/interfaces"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
)

type DeviceInfoService struct {
	repository interfaces.IDeviceInfoRepository
}

func NewDeviceInfoService(repo interfaces.IDeviceInfoRepository) DeviceInfoService {
	return DeviceInfoService{repository: repo}
}

func (service DeviceInfoService) GetDeviceInfo(ctx context.Context, deviceId string) (string, error) {

	if deviceId == "" {
		log.Fatal("Device ID is empty")
	}
	deviceInfo := model.DeviceInfo{DeviceId: deviceId}
	//database get data by ID
	deviceInfo, err := service.repository.GetDeviceInfo(ctx, deviceInfo)
	if err != nil {
		return "", err
	}
	resultByteArr, err := json.Marshal(deviceInfo)
	if err != nil {
		return "", err
	}

	return string(resultByteArr), err
}

func (service DeviceInfoService) CreateDeviceInfo(ctx context.Context, responseStr string) error {

	var detailsStucture model.DeviceInfo
	currentTime := time.Now().UnixMilli()

	err := json.Unmarshal([]byte(responseStr), &detailsStucture)
	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return err
	}
	if detailsStucture.DeviceId == "" {
		log.Printf("Can not read JSON")
		return errors.New("can not read JSON")
	}

	detailsStucture.CreateAt = currentTime
	detailsStucture.UpdatedAt = currentTime
	err = service.repository.CreateDeviceInfo(ctx, detailsStucture)
	if err != nil {
		return err
	}
	return nil

}

func (service DeviceInfoService) UpdateDeviceInfo(ctx context.Context, deviceId string, requestBody string) error {

	deviceInfo := model.DeviceInfo{DeviceId: deviceId}

	err := json.Unmarshal([]byte(requestBody), &deviceInfo)

	deviceInfo.UpdatedAt = time.Now().UnixMilli()

	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return err
	}

	err = repository.UpdateDeviceInfo(ctx, deviceInfo)

	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return err
	}

	return nil
}

func (service DeviceInfoService) DeleteDeviceInfo(ctx context.Context, deviceId string) error {
	deviceInfo := model.DeviceInfo{DeviceId: deviceId}

	err := repository.DeleteDeviceInfoByDeviceId(ctx, deviceInfo)

	if err != nil {
		log.Printf("Could not delete Device Info : [%s]", err.Error())
		return err
	}

	return nil
}
