package main

import (
	"context"
	"errors"
	"log"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

var respository = repository.NewDeviceInfoRepository()
var deviceInfoService = service.NewDeviceInfoService(respository)

func handler(ctx context.Context, request Request) (Response, error) {

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{StatusCode: 404}, errors.New("request body is empty")
	}
	deviceId := request.PathParameters["deviceId"]

	err := deviceInfoService.UpdateDeviceInfo(ctx, deviceId, request.Body)

	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400}, nil
	}

	log.Printf("Entity Updated")

	a := events.APIGatewayProxyResponse{StatusCode: 200}
	return a, err

}

func main() {
	lambda.Start(handler)
}
