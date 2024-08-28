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

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

var respository = repository.NewDeviceInfoRepository()
var deviceInfoService = service.NewDeviceInfoService(respository)

func handler(ctx context.Context, request Request) (Response, error) {

	deviceId := request.PathParameters["deviceId"]
	if deviceId == "" {
		return transport.SendValidationError(400, "ID is empty")
	}

	err := deviceInfoService.DeleteDeviceInfo(ctx, deviceId)

	if err != nil {
		transport.SendError(400, err.Error())
	}
	log.Println("Entity Removed. Starting parsing")

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil

}

func main() {
	lambda.Start(handler)
}
