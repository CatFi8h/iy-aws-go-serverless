package transport

import "github.com/aws/aws-lambda-go/events"

func Send(responseCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: responseCode,
		Body:       message,
	}, nil
}
