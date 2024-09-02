Golang SDK v2 CRUD for Amazon DynamoDB

In this project is implemented 5 lambda functions
CRUD for DeviceInfo entity and Lambda SQS listener for updating DeviceInfo on new message received.

In servierless.yaml is configured to create DynamoDB and SQSQueue, publish all lambda functions as 
REST endpoints described below.
SQSMessage mapping show below.

By default the serverless.yaml configured to run all functions in "us-east-1" region 
and with state: "dev". This can be specified in configuration or as environment variables.

Use commands from Makefile to build, zip and deploy the functions.

Endpoints for CRUD :
Create : POST   {server-name}/device-info
                Body : {    
                            "deviceId"   string, //REQUIRED
                            "deviceName" string,
                            "deviceType" string,
                            "mac"        string,
                            "homeId"     string
                        }
                Return :     diviceId    string
Get :    GET    {server-name}/device-info/{deviceId}
                Return : {    
                            "deviceId"   string, //REQUIRED
                            "deviceName" string,
                            "deviceType" string,
                            "mac"        string,
                            "homeId"     string,
                            "createdAt"  int64,
                            "updateAt"   int64
                        }
Update : PUT    {server-name}/device-info/{deviceId}
                Body : {    
                            "deviceName" string,
                            "deviceType" string,
                            "mac"        string,
                            "homeId"     string
                        }
                Return :     diviceId    string
Delete : DELETE {server-name}/device-info/{deviceId}
                Return :     diviceId    string

SQSMessage:             {    
                            "deviceId"   string,
                            "homeId"     string
                        }
