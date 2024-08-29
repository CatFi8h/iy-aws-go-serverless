package repository

import (
	"context"
	"log"
	"os"
	"time"

	"testing"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	containers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDeviceInfo = []model.DeviceInfo{
	{
		DeviceId:   "21e4e1bc-b2f8-4a47-b092-3e0c452462e0",
		DeviceName: "My Device 1",
		Mac:        "mac-mac-mac-mac1",
		DeviceType: "Phone",
		HomeId:     "1",
		CreateAt:   time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	},
	{
		DeviceId:   "21e4e1bc-b2f8-4a47-b092-3e0c452462e1",
		DeviceName: "My Device 2",
		Mac:        "mac-mac-mac-mac2",
		DeviceType: "Phone",
		HomeId:     "1",
		CreateAt:   time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	},
}

// Create the test container and wait for it to be ready
func setupContainer(t *testing.T) (string, func(t *testing.T)) {
	ctx := context.Background()
	req := containers.ContainerRequest{
		Image:        "amazon/dynamodb-local:latest",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	container, err := containers.GenericContainer(ctx, containers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start DynamoDB: %s", err)
	}
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("Could not get DynamoDB endpoint: %s", err)
	}

	return endpoint, func(t *testing.T) {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Could not stop DynamoDB: %s", err)
		}
	}
}

func TestConnectToLocalDBContainer(t *testing.T) {
	ep, tearDown := setupContainer(t)
	defer tearDown(t)

	connect(ep, t)
	log.Println("Connection check Test finished")
}

func TestSaveDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	dynamodbClient = *connect(e, t)

	if err := CreateDeviceInfo(context.TODO(), testDeviceInfo[0]); err != nil {
		t.Errorf("Expected to be able to save item, but received error: %s", err)
	}
}

func TestGetDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	dynamodbClient = *connect(e, t)

	if err := CreateDeviceInfo(context.TODO(), testDeviceInfo[0]); err != nil {
		t.Errorf("Expected to be able to save item, but received error: %s", err)
	}

	result, err := GetDeviceInfo(context.TODO(), testDeviceInfo[0])
	if err != nil {
		t.Errorf("Expected to be able to get item, but received error: %s", err)
	}
	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 0, t)

}

func validateDeviceIdWithDataById(result *model.DeviceInfo, deviceVersionId int, t *testing.T) {
	if result == nil {
		t.Errorf("Expected entry is not nil")
	}
	if result.DeviceId != testDeviceInfo[deviceVersionId].DeviceId {
		t.Errorf("Expected entry Device ID to be '%s' but received: %s", testDeviceInfo[deviceVersionId].DeviceId, result.DeviceId)
	}
}

func validateResultWithDataById(result *model.DeviceInfo, deviceVersionId int, t *testing.T) {

	if result.DeviceName != testDeviceInfo[deviceVersionId].DeviceName {
		t.Errorf("Expected entry Device Name to be '%s' but received: %s", testDeviceInfo[deviceVersionId].DeviceName, result.DeviceName)
	}
	if result.DeviceType != testDeviceInfo[deviceVersionId].DeviceType {
		t.Errorf("Expected entry Device Type to be '%s' but received: %s", testDeviceInfo[deviceVersionId].DeviceType, result.DeviceType)
	}
	if result.HomeId != testDeviceInfo[deviceVersionId].HomeId {
		t.Errorf("Expected entry Device Home ID to be '%s' but received: %s", testDeviceInfo[deviceVersionId].HomeId, result.HomeId)
	}
	if result.Mac != testDeviceInfo[deviceVersionId].Mac {
		t.Errorf("Expected entry Device MAC to be '%s' but received: %s", testDeviceInfo[deviceVersionId].Mac, result.Mac)
	}
}

func TestUpdateDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	dynamodbClient = *connect(e, t)

	if err := CreateDeviceInfo(context.TODO(), testDeviceInfo[0]); err != nil {
		t.Errorf("Expected to be able to save item, but received error: %s", err)
	}

	result, err := GetDeviceInfo(context.TODO(), testDeviceInfo[0])
	if err != nil {
		t.Errorf("Expected to be able to get item, but received error: %s", err)
	}
	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 0, t)
	deviceInfo := testDeviceInfo[1]
	deviceInfo.DeviceId = testDeviceInfo[0].DeviceId
	UpdateDeviceInfo(context.TODO(), deviceInfo)

	result, err = GetDeviceInfo(context.TODO(), testDeviceInfo[0])
	if err != nil {
		t.Errorf("Expected to be able to get item, but received error: %s", err)
	}
	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 1, t)
}

func TestDeleteDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	dynamodbClient = *connect(e, t)

	if err := CreateDeviceInfo(context.TODO(), testDeviceInfo[0]); err != nil {
		t.Errorf("Expected to be able to save item, but received error: %s", err)
	}

	result, err := GetDeviceInfo(context.TODO(), testDeviceInfo[0])
	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 0, t)
	err = DeleteDeviceInfoByDeviceId(context.TODO(), testDeviceInfo[0])
	if err != nil {
		t.Errorf("Expected to be able to get item, but received error: %s", err)
	}
	result, err = GetDeviceInfo(context.TODO(), testDeviceInfo[0])
	if err != nil {
		t.Errorf("Expected to be able to get item, but received error: %s", err)
	}
	if result != nil {
		t.Errorf("Expected entry is not nil")
	}
}

func connect(e string, t *testing.T) *dynamodb.Client {
	t.Setenv("DEVICE_INFO_TABLE", "device-info-table")
	client := createClient("http://" + e)
	if err := createTable(client); err != nil {
		t.Errorf("Expected to be able to create Dy	namoDB table, but received: %s", err)
	}

	return client
}

func createClient(endpoint string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	// , func(o *config.LoadOptions) error {
	// 	o.Region = "eu-east-1"
	// 	return nil
	// })
	if err != nil {
		panic(err)
	}
	if endpoint == "" {
		endpoint = "https://localhost"
	}

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	})
}

func createTable(c *dynamodb.Client) error {
	_, err := c.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName:   aws.String(os.Getenv("DEVICE_INFO_TABLE")), //tableName
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("deviceid"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("deviceid"),
				KeyType:       types.KeyTypeHash,
			},
		},
	})

	return err
}
