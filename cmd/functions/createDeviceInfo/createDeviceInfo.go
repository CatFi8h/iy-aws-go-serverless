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

var respository = *repository.NewDeviceInfoRepository()
var deviceInfoService = service.NewDeviceInfoService(&respository)

func handler(ctx context.Context, request Request) (Response, error) {

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{StatusCode: 400}, errors.New("request body is empty")
	}

	var responseStr string = request.Body

	err := deviceInfoService.CreateDeviceInfo(ctx, responseStr)

	if err != nil {
		log.Printf("Couldn't insert an item. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}
	log.Println("Insert success!")
	a := events.APIGatewayProxyResponse{StatusCode: 200}
	a.Body = `{"Success"}`
	return a, nil

}

func main() {
	lambda.Start(handler)
}
