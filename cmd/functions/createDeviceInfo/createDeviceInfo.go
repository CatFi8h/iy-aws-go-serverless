package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

var deviceInfoRepository = repository.NewDeviceInfoRepository()
var deviceInfoService = service.NewDeviceInfoService(deviceInfoRepository)

func handler(ctx context.Context, request Request) (*Response, error) {

	if len(request.Body) < 1 {
		return &events.APIGatewayProxyResponse{StatusCode: 400}, errors.New("request body is empty")
	}
	var detailsStucture model.DeviceInfo
	var requestBodyStr string = request.Body

	err := json.Unmarshal([]byte(requestBodyStr), &detailsStucture)
	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return nil, err
	}

	resp, err := deviceInfoService.CreateDeviceInfo(ctx, detailsStucture)

	if err != nil {
		log.Printf("Couldn't insert an item. Here's why: %v\n", err)
		return &events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}
	log.Println("Insert success!")

	return &events.APIGatewayProxyResponse{StatusCode: 200, Body: resp.DeviceId}, nil
}

func main() {
	lambda.Start(handler)
}
