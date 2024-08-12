package helper

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func RespondWithError(statusCode int, err error) (events.APIGatewayProxyResponse, error) {
	errorResponse := struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	}

	responseBody, marshalErr := json.Marshal(errorResponse)
	if marshalErr != nil {
		log.Printf("Error marshalling error response: %v", marshalErr)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message": "Internal server error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}
