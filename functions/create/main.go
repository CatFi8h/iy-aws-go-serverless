package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PartiQLRunner struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type DeviceInfo struct {
	Deviceid   string `json:"deviceid"`
	Name       string `json:"name"`
	Mac        string `json:"mac"`
	Devicetype string `json:"type"`
	HomeId     string `json:"homeId"`
	CreatedAt  int    `json:"createdAt"`
	ModifiedAt int    `json:"modifiedAt"`
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

	var responseStr string = request.Body

	var detailsStucture DeviceInfo

	err1 := json.Unmarshal([]byte(responseStr), &detailsStucture)

	if err1 != nil {
		log.Fatalln(fmt.Printf("Could not Unmarshal JSON : [%s]", err1.Error()))
	}

	if detailsStucture.Mac == "" {
		log.Println("Printing Details : ", responseStr)
		log.Println("Printing Structure : ", detailsStucture)
		log.Printf("Can not read JSON, %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	params, err := attributevalue.MarshalList([]interface{}{detailsStucture.Deviceid,
		detailsStucture.Name, detailsStucture.Mac, detailsStucture.Devicetype,
		detailsStucture.HomeId, detailsStucture.CreatedAt, detailsStucture.ModifiedAt})

	if err != nil {
		log.Fatalf("Can not get Attributes, %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	_, err = runner.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("INSERT INTO \"%v\" VALUE {'deviceId': ?, 'name': ?, 'mac': ?, 'type': ?, 'homeId': ?, 'createdAt': ?, 'modifiedAt': ?}", runner.TableName)),
		Parameters: params,
	})
	if err != nil {
		log.Printf("Couldn't insert an item with PartiQL. Here's why: %v\n", err)
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
