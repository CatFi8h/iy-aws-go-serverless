package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
	DeviceType string `dynamodbav:"type" json:"type"`
	HomeId     string `dynamodbav:"homeId" json:"homeId"`
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

	var deviceInfo DeviceInfo

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	log.Println(request.Body)
	err = json.Unmarshal([]byte(request.Body), &deviceInfo)

	if err != nil {
		log.Fatalln(fmt.Printf("Could not Unmarshal JSON : [%s]", err.Error()))
	}

	runner := PartiQLRunner{
		DynamoDbClient: dynamodb.NewFromConfig(sdkConfig),
		TableName:      os.Getenv("DEVICE_INFO_TABLE"),
	}
	update := expression.Set(expression.Name("modifiedAt"), expression.Value(time.Now().UnixMilli()))
	if deviceInfo.DeviceName != "" {
		update.Set(expression.Name("deviceName"), expression.Value(deviceInfo.DeviceName))
	}
	if deviceInfo.Mac != "" {
		update.Set(expression.Name("mac"), expression.Value(deviceInfo.Mac))
	}
	if deviceInfo.DeviceType != "" {
		update.Set(expression.Name("deviceType"), expression.Value(deviceInfo.DeviceType))
	}
	if deviceInfo.HomeId != "" {
		update.Set(expression.Name("homeId"), expression.Value(deviceInfo.HomeId))
	}
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400}, nil
	} else {
		_, err = runner.DynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName:                 aws.String(runner.TableName),
			Key:                       deviceInfo.getKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			log.Printf("Couldn't update movie %v. Here's why: %v\n", deviceInfo.DeviceName, err)
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 400}, nil
		}
	}

	log.Printf("Entity Updated")

	a := events.APIGatewayProxyResponse{StatusCode: 200}
	return a, err

}

func main() {
	lambda.Start(handler)
}
