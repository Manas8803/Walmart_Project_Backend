package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	lmd "github.com/aws/aws-lambda-go/lambda"
)

type RequestData struct {
	BeaconID string `json:"beacon_id"`
	UserID   string `json:"user_id"`
}

type Response struct {
	Message string `json:"message"`
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestData RequestData

	// Receive Beacon id and user id
	err := json.Unmarshal([]byte(req.Body), &requestData)
	if err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	log.Println("Received beacon_id: " + requestData.BeaconID + ", user_id: " + requestData.UserID)
	// Fetch discount from beacon id

	// Fetch connnection id from user id
	// Push the data in websocket connection
	resp := Response{
		Message: "Successfully fetched vehicles from DB",
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error in marshalling response : ", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{Body: string(respBytes), StatusCode: http.StatusOK, Headers: map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}}, nil
}

func main() {
	lmd.Start(Handler)
}
