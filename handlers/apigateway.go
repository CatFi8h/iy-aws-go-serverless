package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/service"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/transport"
	"github.com/aws/aws-lambda-go/events"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

const DEVICE_ID = "deviceId"

type ApiGatewayHandler struct {
	service service.IDeviceInfoService
}

func NewApiGatewayHandler(s service.IDeviceInfoService) *ApiGatewayHandler {
	return &ApiGatewayHandler{
		service: s,
	}
}

func (h *ApiGatewayHandler) GetHandler(ctx context.Context, request Request) (Response, error) {

	deviceId, ok := request.PathParameters[DEVICE_ID]
	if !ok {
		return transport.SendBadRequestError("missing [deviceId] parameter in path")
	}

	deviceInfo, err := h.service.GetDeviceInfo(ctx, deviceId)

	if err != nil {
		return transport.SendBadRequestError(err.Error())
	}
	if deviceInfo == nil {
		return transport.SendNotFoundError("Device Info not found")
	}
	resultByteArr, err := json.Marshal(deviceInfo)
	if err != nil {
		return transport.SendBadRequestError(err.Error())
	}

	return transport.Send(http.StatusOK, string(resultByteArr))

}

func (h *ApiGatewayHandler) CreateHandler(ctx context.Context, request Request) (Response, error) {

	if len(request.Body) < 1 {
		return transport.SendBadRequestError("request body is empty")
	}
	var detailsStucture model.DeviceInfo

	err := json.Unmarshal([]byte(request.Body), &detailsStucture)
	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return transport.SendBadRequestError(err.Error())
	}
	if detailsStucture.DeviceId == "" {
		log.Printf("Device ID not specified")
		return transport.SendBadRequestError("Device ID not specified")
	}

	resp, err := h.service.CreateDeviceInfo(ctx, detailsStucture)

	if err != nil {
		log.Printf("Couldn't insert an item. Here's why: %v\n", err)
		return transport.SendBadRequestError(err.Error())
	}
	log.Println("Insert success!")

	return transport.Send(http.StatusCreated, resp.DeviceId)
}

func (h *ApiGatewayHandler) UpdateHandler(ctx context.Context, request Request) (Response, error) {

	if len(request.Body) < 1 {
		return transport.SendBadRequestError("request body is empty")
	}
	deviceId, ok := request.PathParameters[DEVICE_ID]
	if !ok {
		return transport.SendBadRequestError("device ID not specified.")
	}

	var detailsStucture model.DeviceInfo

	err := json.Unmarshal([]byte(request.Body), &detailsStucture)
	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return transport.SendBadRequestError(fmt.Sprintf("Could not Unmarshal JSON : [%s]", err.Error()))
	}

	resp, err := h.service.UpdateDeviceInfo(ctx, deviceId, detailsStucture)

	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
		return transport.SendBadRequestError(fmt.Sprintf("Couldn't build expression for update. Here's why: %v\n", err))
	}

	log.Printf("Entity Updated.")

	return transport.Send(http.StatusOK, resp.DeviceId)
}

func (h *ApiGatewayHandler) DeleteHandler(ctx context.Context, request Request) (Response, error) {

	deviceId, ok := request.PathParameters[DEVICE_ID]
	if !ok {
		return transport.SendBadRequestError("device ID not specified.")
	}

	resp, err := h.service.DeleteDeviceInfo(ctx, deviceId)

	if err != nil {
		return transport.SendBadRequestError(err.Error())
	}
	log.Println("Entity Removed.")

	return transport.Send(http.StatusOK, resp.DeviceId)

}
