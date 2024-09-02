package repository

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DeviceInfoRepository struct {
	dbClient  *dynamodb.Client
	tableName string
}

func NewDeviceInfoRepository(ctx context.Context, tableName string) IDeviceInfoRepository {

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config %s", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(sdkConfig)

	return &DeviceInfoRepository{
		dbClient:  dynamodbClient,
		tableName: tableName,
	}
}

func (d *DeviceInfoRepository) CreateDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {
	item, err := attributevalue.MarshalMap(&deviceInfo)
	if err != nil {
		return nil, err
	}
	newItem, err := d.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(d.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(deviceid)"),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			log.Println("failed pre-condition check")
			return nil, errors.New("item with provided deviceID already exists")
		} else {
			return nil, err
		}
	}
	err = attributevalue.UnmarshalMap(newItem.Attributes, &deviceInfo)
	if err != nil {
		return nil, err
	}
	return deviceInfo, nil
}

func (d *DeviceInfoRepository) UpdateDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {

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
	result, err := d.dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(d.tableName),
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
	return deviceInfo, nil
}

func (d *DeviceInfoRepository) GetDeviceInfo(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {

	result, err := d.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
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

func (d *DeviceInfoRepository) DeleteDeviceInfoByDeviceId(ctx context.Context, deviceInfo *model.DeviceInfo) (*model.DeviceInfo, error) {

	resp, err := d.dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key:       deviceInfo.GetKey(),
	})
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalMap(resp.Attributes, &deviceInfo)
	if err != nil {
		return nil, err
	}

	return deviceInfo, nil

}
