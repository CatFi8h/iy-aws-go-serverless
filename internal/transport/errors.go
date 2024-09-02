package transport

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func SendBadRequestError(message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       message,
	}, nil
}

func SendNotFoundError(message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       message,
	}, nil
}

func SendError(errorCode int16, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: int(errorCode),
		Body:       message,
	}, nil
}
