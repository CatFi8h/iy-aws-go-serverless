package transport

import "github.com/aws/aws-lambda-go/events"

func SendValidationError(errorCode int16, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: int(errorCode),
		Body:       message,
	}, nil
}

func SendError(errorCode int16, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: int(errorCode),
		Body:       message,
	}, nil
}
