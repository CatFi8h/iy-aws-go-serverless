package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/service"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/transport"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

// var respository = repository.NewDeviceInfoRepository(initClient(context.TODO()))
var respository = repository.NewDeviceInfoRepository()
var deviceInfoService = service.NewDeviceInfoService(respository)

func handler(ctx context.Context, request Request) (Response, error) {

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{StatusCode: 404}, errors.New("request body is empty")
	}
	deviceId := request.PathParameters["deviceId"]

	if deviceId == "" {
		return transport.SendValidationError(400, "device ID not specified.")
	}

	var detailsStucture model.DeviceInfo
	var requestBodyStr string = request.Body

	err := json.Unmarshal([]byte(requestBodyStr), &detailsStucture)
	if err != nil {
		log.Printf("Could not Unmarshal JSON : [%s]", err.Error())
		return transport.SendError(400, fmt.Sprintf("Could not Unmarshal JSON : [%s]", err.Error()))
	}

	resp, err := deviceInfoService.UpdateDeviceInfo(ctx, deviceId, detailsStucture)

	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
		return transport.SendError(400, fmt.Sprintf("Couldn't build expression for update. Here's why: %v\n", err))
	}

	log.Printf("Entity Updated")

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: resp.DeviceId}, err
}

func main() {
	lambda.Start(handler)
}

// func initClient(cont context.Context) dynamodb.Client {
// 	sdkConfig, err := config.LoadDefaultConfig(cont)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	dynamodbClient := *dynamodb.NewFromConfig(sdkConfig)
// 	return dynamodbClient
// }
