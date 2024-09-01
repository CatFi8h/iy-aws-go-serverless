package main

import (
	"context"
	"encoding/json"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/service"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/transport"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// aliasing the types to keep lines short
type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

var respository = repository.NewDeviceInfoRepository()
var deviceInfoService = service.NewDeviceInfoService(respository)

func handler(ctx context.Context, request Request) (Response, error) {

	deviceId := request.PathParameters["deviceId"]
	if deviceId == "" {
		return transport.SendValidationError(400, "ID is empty")
	}

	deviceInfo, err := deviceInfoService.GetDeviceInfo(ctx, deviceId)

	if err != nil {
		return transport.SendError(400, error.Error())
	}
	if deviceInfo == nil {
		return transport.SendError(404, "Device Info not found")
	}
	resultByteArr, err := json.Marshal(deviceInfo)
	if err != nil {
		return transport.SendError(400, error.Error())
	}
	resultStr := string(resultByteArr)

	return transport.Send(200, resultStr)

}

func main() {
	lambda.Start(handler)
}
