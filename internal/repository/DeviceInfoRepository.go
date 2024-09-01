package repository

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type deviceInfoRepository struct{}

var (
	dynamodbClient dynamodb.Client
)

func NewDeviceInfoRepository() IDeviceInfoRepository {
	return &deviceInfoRepository{}
}

// func NewDeviceInfoRepository(client dynamodb.Client) IDeviceInfoRepository {
// 	dynamodbClient = client
// 	return &deviceInfoRepository{}
// }

func (*deviceInfoRepository) CreateDeviceInfo(ctx context.Context, deviceInfo model.DeviceInfo) (*model.DeviceInfo, error) {
	item, err := attributevalue.MarshalMap(deviceInfo)
	if err != nil {
		return nil, err
	}
	newItem, err := dynamodbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DEVICE_INFO_TABLE")),
		Item:      item,
	})
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalMap(newItem.Attributes, &deviceInfo)
	if err != nil {
		return nil, err
	}
	return &deviceInfo, nil
}

func (*deviceInfoRepository) UpdateDeviceInfo(ctx context.Context, deviceInfo model.DeviceInfo) (*model.DeviceInfo, error) {

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
		return nil, err
	}
	result, err := dynamodbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(os.Getenv("DEVICE_INFO_TABLE")),
		Key:                       deviceInfo.GetKey(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return nil, err
	}
	if len(result.Attributes) < 1 {
		log.Printf("Device Info with ID: %v Not found", deviceInfo.DeviceId)
	}
	err = attributevalue.UnmarshalMap(result.Attributes, &deviceInfo)
	if err != nil {
		return nil, err
	}
	return &deviceInfo, nil
}

func (*deviceInfoRepository) GetDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {

	result, err := dynamodbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DEVICE_INFO_TABLE")),
		Key:       deviceInfo.GetKey(),
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, &deviceInfo)
	if err != nil {
		return nil, err
	}
	return deviceInfo, nil
}

func (*deviceInfoRepository) DeleteDeviceInfoByDeviceId(ctx context.Context, deviceInfo model.DeviceInfo) error {

	_, err := dynamodbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("DEVICE_INFO_TABLE")),
		Key:       deviceInfo.GetKey(),
	})
	if err != nil {
		return err
	}

	return nil

}

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	dynamodbClient = *dynamodb.NewFromConfig(sdkConfig)
}
