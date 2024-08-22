package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	Deviceid   string `dynamodbav:"deviceid" json:"deviceid"`
	DeviceName string `dynamodbav:"deviceName" json:"deviceName"`
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

	var deviceInfo DeviceInfo

	err = json.Unmarshal([]byte(request.Body), &deviceInfo)

	if err != nil {
		log.Fatalln(fmt.Printf("Could not Unmarshal JSON : [%s]", err.Error()))
	}

	runner := PartiQLRunner{
		DynamoDbClient: dynamodb.NewFromConfig(sdkConfig),
		TableName:      os.Getenv("DEVICE_INFO_TABLE"),
	}

	log.Println("Getting Entity by DeviceId from DB")
	_, err = runner.DynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:        aws.String(runner.TableName),
		Key:              deviceInfo.getKey(),
		UpdateExpression: aws.String("SET deviceName = if_not_exists(deviceName, :deviceName)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":deviceName": &types.AttributeValueMemberS{Value: "I've been updated."},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400}, nil
	}

	log.Println("Entity Updated.")

	a := events.APIGatewayProxyResponse{StatusCode: 200}

	return a, err

}

func main() {
	lambda.Start(handler)
}
