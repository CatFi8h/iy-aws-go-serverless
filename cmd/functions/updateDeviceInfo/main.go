package main

import (
	"context"
	"os"

	"github.com/CatFi8h/iy-aws-go-serverless/handlers"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/repository"
	"github.com/CatFi8h/iy-aws-go-serverless/internal/service"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	tableName, ok := os.LookupEnv("DEVICE_INFO_TABLE")
	if !ok {
		panic("Need DEVICE_INFO_TABLE environment variable")
	}
	dynamodb := repository.NewDeviceInfoRepository(context.TODO(), tableName)

	service := service.NewDeviceInfoService(dynamodb)

	handler := handlers.NewApiGatewayHandler(service)

	lambda.Start(handler.UpdateHandler)

}
