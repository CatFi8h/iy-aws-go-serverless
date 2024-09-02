package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDeviceInfoService struct {
	mock.Mock
}

func (mock *MockDeviceInfoService) CreateDeviceInfo(ctx context.Context, deviceInfoReq model.DeviceInfo) (*model.DeviceInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*model.DeviceInfo), args.Error(1)
}

func (mock *MockDeviceInfoService) GetDeviceInfo(ctx context.Context, deviceId string) (*model.DeviceInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*model.DeviceInfo), args.Error(1)
}

func (mock *MockDeviceInfoService) UpdateDeviceInfo(ctx context.Context, deviceId string, deviceInfoReq model.DeviceInfo) (*model.DeviceInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*model.DeviceInfo), args.Error(1)
}

func (mock *MockDeviceInfoService) DeleteDeviceInfo(ctx context.Context, deviceId string) (*model.DeviceInfo, error) {
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

func getMockApiGetDeviceInfo() *ApiGatewayHandler {
	mockService := new(MockDeviceInfoService)
	mockService.On("GetDeviceInfo").Return(&deviceInfo, nil)
	return NewApiGatewayHandler(mockService)
}

func TestApiGetDeviceInfo_valid(t *testing.T) {

	apiGateway := getMockApiGetDeviceInfo()
	resp, err := apiGateway.GetHandler(context.TODO(), events.APIGatewayProxyRequest{PathParameters: map[string]string{"deviceId": "1"}})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	fmt.Print(resp.Body)
	var respDevInfo model.DeviceInfo
	err = json.Unmarshal([]byte(resp.Body), &respDevInfo)
	assert.NoError(t, err)
	validateDeviceInfoEqueal(t, deviceInfo, respDevInfo)
}

func validateDeviceInfoEqueal(t *testing.T, d1 model.DeviceInfo, d2 model.DeviceInfo) {
	assert.Equal(t, d1.DeviceId, d2.DeviceId)
	assert.Equal(t, d1.DeviceName, d2.DeviceName)
	assert.Equal(t, d1.DeviceType, d2.DeviceType)
	assert.Equal(t, d1.Mac, d2.Mac)
	assert.Equal(t, d1.HomeId, d2.HomeId)
	assert.Equal(t, d1.CreateAt, d2.CreateAt)
	assert.Equal(t, d1.UpdatedAt, d2.UpdatedAt)
}

func TestApiGetDeviceInfo_NoPathParam(t *testing.T) {
	apiGateway := getMockApiGetDeviceInfo()
	resp, err := apiGateway.GetHandler(context.TODO(), events.APIGatewayProxyRequest{})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestApiGetDeviceInfo_NotFound(t *testing.T) {
	mockService := new(MockDeviceInfoService)
	mockService.On("GetDeviceInfo").Return((*model.DeviceInfo)(nil), nil)
	apiGateway := NewApiGatewayHandler(mockService)
	resp, err := apiGateway.GetHandler(context.TODO(), events.APIGatewayProxyRequest{PathParameters: map[string]string{"deviceId": "2"}})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestApiCreateDeviceInfo_valid(t *testing.T) {
	jsonStr, err := json.Marshal(deviceInfo)
	assert.NoError(t, err)
	mockService := new(MockDeviceInfoService)
	mockService.On("CreateDeviceInfo").Return(&deviceInfo, nil)
	apiGateway := NewApiGatewayHandler(mockService)
	resp, err := apiGateway.CreateHandler(context.TODO(), events.APIGatewayProxyRequest{Body: string(jsonStr)})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestApiCreateDeviceInfo_NoDeviceId(t *testing.T) {
	deviceInfo.DeviceId = ""
	jsonStr, err := json.Marshal(deviceInfo)
	assert.NoError(t, err)
	apiGateway := NewApiGatewayHandler(nil)
	resp, err := apiGateway.CreateHandler(context.TODO(), events.APIGatewayProxyRequest{Body: string(jsonStr)})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestApiUpdateDeviceInfo_valid(t *testing.T) {
	jsonStr, err := json.Marshal(deviceInfo)
	assert.NoError(t, err)
	mockService := new(MockDeviceInfoService)
	mockService.On("UpdateDeviceInfo").Return(&deviceInfo, nil)
	apiGateway := NewApiGatewayHandler(mockService)
	resp, err := apiGateway.UpdateHandler(context.TODO(), events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"deviceId": "1"},
		Body:           string(jsonStr),
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestApiDeleteDeviceInfo_valid(t *testing.T) {
	mockService := new(MockDeviceInfoService)
	mockService.On("DeleteDeviceInfo").Return(&deviceInfo, nil)
	apiGateway := NewApiGatewayHandler(mockService)
	resp, err := apiGateway.DeleteHandler(context.TODO(), events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"deviceId": "1"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSqsUpdateDeviceInfo_valid(t *testing.T) {
	devMessage := model.DeviceInfoSQSMessage{
		DeviceId: "1",
		HomeId:   "My New Home ID",
	}
	jsonStr, err := json.Marshal(devMessage)
	assert.NoError(t, err)
	mockService := new(MockDeviceInfoService)
	mockService.On("UpdateDeviceInfo").Return(&deviceInfo, nil)
	apiGateway := NewApiGatewayHandler(mockService)
	err = apiGateway.SQSHandler(context.TODO(), events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "id1",
				Body:      string(jsonStr),
			},
		},
	})
	assert.NoError(t, err)
}
