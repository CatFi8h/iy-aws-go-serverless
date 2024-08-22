package main

import (
	"context"
	"errors"
	"log"

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
	Deviceid   string `dynamodbav:"deviceid" json:"deviceid"`
	Name       string `dynamodbav:"name" json:"name"`
	Mac        string `dynamodbav:"mac" json:"mac"`
	Devicetype string `dynamodbav:"type" json:"type"`
	HomeId     string `dynamodbav:"homeId" json:"homeId"`
	CreatedAt  int    `dynamodbav:"createdAt" json:"createdAt"`
	ModifiedAt int    `dynamodbav:"modifiedAt" json:"modifiedAt"`
}

func (deviceInfo DeviceInfo) getKey() map[string]types.AttributeValue {
	deviceId, err := attributevalue.Marshal(deviceInfo.Deviceid)
	if err != nil {
		panic("DeviceId Not parsed " + err.Error())
	}

	return map[string]types.AttributeValue{"deviceId": deviceId}
}

func handler(request events.LambdaFunctionURLRequest) (events.APIGatewayProxyResponse, error) {

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{StatusCode: 404}, errors.New("request body is empty")
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	tableName := "data-dev"

	runner := PartiQLRunner{
		DynamoDbClient: dynamodb.NewFromConfig(sdkConfig),
		TableName:      tableName,
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

	a := events.APIGatewayProxyResponse{StatusCode: 200}
	a.Body = deviceInfo.Name
	return a, nil

}

func main() {
	lambda.Start(handler)
}
