package main

import (
	"context"
	"log"

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

	log.Println(request.PathParameters["deviceId"])
	deviceId := request.PathParameters["deviceId"]
	if deviceId == "" {
		return transport.SendValidationError(400, "ID is empty")
	}

	deviceInfo, error := deviceInfoService.GetDeviceInfo(ctx, deviceId)

	if error != nil {
		return transport.SendError(400, error.Error())
	} else if deviceInfo == "" {
		return transport.SendError(404, "Device Info not found")
	}

	return transport.Send(200, deviceInfo)

}

func main() {
	lambda.Start(handler)
}
