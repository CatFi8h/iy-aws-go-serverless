package repository

import (
	"context"
	"time"

	"testing"

	"github.com/CatFi8h/iy-aws-go-serverless/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"

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

var tableName string = "device_info_table"

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
}

func TestSaveDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	repo := createRepositoryWithConnection(e, t)

	resp, err := repo.CreateDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err, "Expected to be able to save item, but received error")
	assert.NotNil(t, resp)
	assert.Equal(t, testDeviceInfo[0].DeviceId, resp.DeviceId)
}

func TestSaveSameDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	repo := createRepositoryWithConnection(e, t)

	resp, err := repo.CreateDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err, "Expected to be able to save item, but received error")
	assert.NotNil(t, resp)
	assert.Equal(t, testDeviceInfo[0].DeviceId, resp.DeviceId)
	_, err = repo.CreateDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.Error(t, err)
	assert.ErrorContains(t, err, "item with provided deviceID already exists")
}

func TestGetDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	repo := *createRepositoryWithConnection(e, t)

	result, err := repo.CreateDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err, "Expected to be able to save item, but received error")
	validateDeviceIdWithDataById(result, 0, t)

	result, err = repo.GetDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err, "Expected to be able to save item, but received error")
	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 0, t)
}

func TestGetDeviceInfoNoItemFound(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	repo := createRepositoryWithConnection(e, t)

	resp, err := repo.GetDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err)
	assert.Nil(t, resp)
}

func validateDeviceIdWithDataById(result *model.DeviceInfo, deviceVersionId int, t *testing.T) {
	assert.NotNil(t, result, "Expected entry is not nil")
	assert.Equal(t, result.DeviceId, testDeviceInfo[deviceVersionId].DeviceId)
}

func validateResultWithDataById(result *model.DeviceInfo, deviceVersionId int, t *testing.T) {
	assert.Equal(t, result.DeviceName, testDeviceInfo[deviceVersionId].DeviceName)
	assert.Equal(t, result.DeviceType, testDeviceInfo[deviceVersionId].DeviceType)
	assert.Equal(t, result.HomeId, testDeviceInfo[deviceVersionId].HomeId)
	assert.Equal(t, result.Mac, testDeviceInfo[deviceVersionId].Mac)
}

func TestUpdateDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	repo := createRepositoryWithConnection(e, t)

	result, err := repo.CreateDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err, "Expected to be able to save item, but received error")

	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 0, t)
	deviceInfo := testDeviceInfo[1]
	deviceInfo.DeviceId = testDeviceInfo[0].DeviceId
	result, err = repo.UpdateDeviceInfo(context.TODO(), &deviceInfo)

	assert.NoError(t, err, "Expected to be able to save item, but received error")
	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 1, t)
}

func TestDeleteDeviceInfo(t *testing.T) {
	e, tearDown := setupContainer(t)
	defer tearDown(t)

	repo := createRepositoryWithConnection(e, t)

	result, err := repo.CreateDeviceInfo(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err)
	validateDeviceIdWithDataById(result, 0, t)
	validateResultWithDataById(result, 0, t)
	result, err = repo.DeleteDeviceInfoByDeviceId(context.TODO(), &testDeviceInfo[0])
	assert.NoError(t, err)
	validateDeviceIdWithDataById(result, 0, t)

}

func createRepositoryWithConnection(e string, t *testing.T) *DeviceInfoRepository {
	dynamodbClient := *connect(e, t)
	return &DeviceInfoRepository{
		dbClient:  &dynamodbClient,
		tableName: tableName,
	}
}

func connect(e string, t *testing.T) *dynamodb.Client {
	client := createClient("http://" + e)
	if err := createTable(client); err != nil {
		t.Errorf("Expected to be able to create Dy	namoDB table, but received: %s", err)
	}

	return client
}

func createClient(endpoint string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
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
		TableName:   &tableName,
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
