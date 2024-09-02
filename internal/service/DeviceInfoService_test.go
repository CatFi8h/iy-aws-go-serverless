package service

import (
	"context"
	"testing"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDeviceInfoRepository struct {
	mock.Mock
}

func (mock *MockDeviceInfoRepository) CreateDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*model.DeviceInfo), args.Error(1)
}

func (mock *MockDeviceInfoRepository) GetDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*model.DeviceInfo), args.Error(1)
}

func (mock *MockDeviceInfoRepository) UpdateDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*model.DeviceInfo), args.Error(1)
}

func (mock *MockDeviceInfoRepository) DeleteDeviceInfoByDeviceId(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*model.DeviceInfo), args.Error(1)
}

var deviceInfo = model.DeviceInfo{
	DeviceId:   "1",
	DeviceName: "My Device 1",
	Mac:        "aaa-aaa-aaa",
	DeviceType: "Phone",
	HomeId:     "1",
	CreateAt:   1,
	UpdatedAt:  1,
}

func TestGetDeviceInfo_valid(t *testing.T) {
	mockRepo := new(MockDeviceInfoRepository)
	mockRepo.On("GetDeviceInfo").Return(&deviceInfo, nil).Once()
	service := NewDeviceInfoService(mockRepo)

	res, _ := service.GetDeviceInfo(context.TODO(), "1")

	mockRepo.AssertExpectations(t)
	assert.NotNil(t, res)
	assert.Equal(t, deviceInfo.DeviceId, res.DeviceId)
	assert.Equal(t, deviceInfo.DeviceName, res.DeviceName)
	assert.Equal(t, deviceInfo.DeviceType, res.DeviceType)
	assert.Equal(t, deviceInfo.Mac, res.Mac)
	assert.Equal(t, deviceInfo.HomeId, res.HomeId)
}

func TestGetDeviceInfo_ErrorDeviceIdEmpty(t *testing.T) {

	service := NewDeviceInfoService(nil)

	_, err := service.GetDeviceInfo(context.TODO(), "")

	assert.Error(t, err)
	assert.Equal(t, "device ID is empty", err.Error())

}

func TestGetDeviceInfo_NotFoundDeviceInfo(t *testing.T) {
	mockRepo := new(MockDeviceInfoRepository)
	service := NewDeviceInfoService(mockRepo)
	mockRepo.On("GetDeviceInfo").Return((*model.DeviceInfo)(nil), nil)

	resp, err := service.GetDeviceInfo(context.TODO(), "2")

	assert.NoError(t, err)
	assert.Nil(t, resp)
}

func TestCreateDeviceInfo_valid(t *testing.T) {
	mockRepo := new(MockDeviceInfoRepository)
	service := NewDeviceInfoService(mockRepo)
	mockRepo.On("CreateDeviceInfo").Return(&model.DeviceInfo{DeviceId: deviceInfo.DeviceId}, nil)

	resp, err := service.CreateDeviceInfo(context.TODO(), deviceInfo)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, deviceInfo.DeviceId, resp.DeviceId)
	assert.Equal(t, resp.CreateAt, resp.UpdatedAt)
}

func TestCreateDeviceInfo_noRequiredIdValue(t *testing.T) {
	service := NewDeviceInfoService(nil)

	_, err := service.CreateDeviceInfo(context.TODO(), model.DeviceInfo{DeviceName: "Name test"})

	assert.Error(t, err)
}

func TestUpdateDeviceInfo_valid(t *testing.T) {
	const newName = "New Name"
	mockRepo := new(MockDeviceInfoRepository)
	service := NewDeviceInfoService(mockRepo)
	mockRepo.On("UpdateDeviceInfo").Return(&model.DeviceInfo{DeviceId: deviceInfo.DeviceId, DeviceName: newName, CreateAt: deviceInfo.CreateAt}, nil)
	resp, err := service.UpdateDeviceInfo(context.TODO(), "1", deviceInfo)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, newName, resp.DeviceName)
	assert.NotEqual(t, resp.CreateAt, resp.UpdatedAt)
}

func TestDeleteDeviceInfo_valid(t *testing.T) {
	mockRepo := new(MockDeviceInfoRepository)
	service := NewDeviceInfoService(mockRepo)
	deviceId := "1"
	mockRepo.On("DeleteDeviceInfoByDeviceId").Return(&model.DeviceInfo{DeviceId: deviceId}, nil)
	resp, err := service.DeleteDeviceInfo(context.TODO(), deviceId)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, deviceId, resp.DeviceId)
}
