package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PartiQLRunner struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type DeviceInfo struct {
	Deviceid   string `dynamodbav:"deviceid"`
	DeviceName string `dynamodbav:"deviceName"`
	Mac        string `dynamodbav:"mac"`
	DeviceType string `dynamodbav:"deviceType"`
	HomeId     string `dynamodbav:"homeId"`
	CreatedAt  int    `dynamodbav:"createdAt"`
	ModifiedAt int    `dynamodbav:"modifiedAt"`
}

func (deviceInfo DeviceInfo) getKey() map[string]types.AttributeValue {
	deviceId, err := attributevalue.Marshal(deviceInfo.Deviceid)
	if err != nil {
		panic("DeviceId Not parsed " + err.Error())
	}

	return map[string]types.AttributeValue{"deviceId": deviceId}
}

func handler(request events.LambdaFunctionURLRequest) (events.APIGatewayProxyResponse, error) {

	log.Println(request.QueryStringParameters)
	an := request.QueryStringParameters
	deviceId := an["id"]
	log.Println("QueryParam : ", deviceId)
	if deviceId == "" {
		return events.APIGatewayProxyResponse{
			Body:       "ID is empty",
			StatusCode: 400}, nil
	}
	deviceInfo := DeviceInfo{Deviceid: deviceId}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	runner := PartiQLRunner{
		DynamoDbClient: dynamodb.NewFromConfig(sdkConfig),
		TableName:      os.Getenv("DEVICE_INFO_TABLE"),
	}
	log.Println("Getting Entity by DeviceId from DB")
	result, err := runner.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(runner.TableName),
		Key:       deviceInfo.getKey(),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			Body:       "Result not found",
			StatusCode: 404}, nil
	}
	log.Println("Entity Found. Starting parsing")
	deviceInfo = DeviceInfo{}
	err = attributevalue.UnmarshalMap(result.Item, &deviceInfo)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Couldn't parse result JSON",
			StatusCode: 400}, nil
	}
	log.Println("Parsing finished")
	deviceInfoBytes, err := json.Marshal(deviceInfo)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Couldn't parse result JSON",
			StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(deviceInfoBytes)}, nil

}

func main() {
	lambda.Start(handler)
}
