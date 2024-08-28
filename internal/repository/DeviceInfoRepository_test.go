package repository

import (
	"context"
	"fmt"
	"os"

	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// the TestMain functions runs before any test we execute
func TestMain(m *testing.M) {
	// get a context for our request
	ctx := context.Background()
	// create a container request for DynamoDB
	req := testcontainers.ContainerRequest{
		// we're using the latest version of the image provided by Amazon
		Image: "amazon/dynamodb-local:latest",
		// be sure to use the commands as described in the documentation, but
		// an in-memory version is good enough for us
		Cmd: []string{"-jar", "DynamoDBLocal.jar", "-inMemory"},
		// by default, DynamoDB runs on port 8000
		ExposedPorts: []string{"8000/tcp"},
		// testcontainers let's us block until the port is available, i.e.,
		// DynamoDB has started
		WaitingFor: wait.NewHostPortStrategy("8000"),
	}

	// let's start the container!
	d, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	// be sure to stop the container before we're done!
	defer d.Terminate(ctx)

	// now all we need is the IP and port of our DynamoDB instance to connect
	// to the right endpoints
	ip, err := d.Host(ctx)

	if err != nil {
		panic(err)
	}

	port, err := d.MappedPort(ctx, "8000")

	if err != nil {
		panic(err)
	}

	// create a new session with our custom endpoint
	// we need to specify a region, otherwise we get a fatal error
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(fmt.Sprintf("http://%s:%s", ip, port)),
		Region:   aws.String("eu-central-1"),
	}))

	// and now we have our service!
	dynamodb.New(sess)

	// now we just need to tell go-test that we can run the tests
	os.Exit(m.Run())
}

func TestGreetKeys(t *testing.T) {
	// in this test, we first create a new table with a key, then see if the
	// greeting works

	table := "greetingtable"

	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				// by the way: note how aws uses custom strings?!
				AttributeName: aws.String("theresa"),
				AttributeType: aws.String("S"),
			},
		},
		BillingMode:            nil,
		GlobalSecondaryIndexes: nil,
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("theresa"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		LocalSecondaryIndexes: nil,
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		SSESpecification:    nil,
		StreamSpecification: nil,
		TableName:           aws.String(table),
		Tags:                nil,
	})

	if err != nil {
		t.Fatalf("could not create table: %s", err.Error())
	}

	out, err := GreetKeys(table, svc)

	if err != nil {
		t.Fatalf("could not greet table: %s", err.Error())
	}

	if out != "Hello theresa!" {
		t.Fatalf("output \"%s\" is wrong!", out)
	}
}
