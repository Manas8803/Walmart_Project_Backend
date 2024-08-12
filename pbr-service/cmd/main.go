package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Manas8803/Walmart_Project_Backend/pbr-service/pkg/helper"
	"github.com/Manas8803/Walmart_Project_Backend/pbr-service/pkg/models/service"
	"github.com/aws/aws-lambda-go/events"
	lmd "github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/websocket"
)

type RequestData struct {
	BeaconID string `json:"beacon_id"`
	UserID   string `json:"user_id"`
}

type Response struct {
	Message string `json:"message"`
}

type SocketMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestData RequestData

	// Receive Beacon id and user id
	err := json.Unmarshal([]byte(req.Body), &requestData)
	if err != nil {
		return helper.RespondWithError(http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
	}

	// Fetch discount
	discount, err := service.FetchDiscountByBeaconID(requestData.BeaconID)
	if err != nil {
		log.Printf("Error in fetching discount: %v", err)
		return helper.RespondWithError(http.StatusNotFound, fmt.Errorf("failed to fetch discount: %v", err))
	}

	// Fetch user
	user, err := service.FetchUserByID(requestData.UserID)
	if err != nil {
		log.Printf("Error in fetching user: %v", err)
		return helper.RespondWithError(http.StatusNotFound, fmt.Errorf("failed to fetch user: %v", err))
	}

	// Push the data in websocket connection
	conn, _, err := websocket.DefaultDialer.Dial(os.Getenv("NOTIFY_WEBSOCKET_URL"), nil)
	if err != nil {
		log.Printf("Failed to connect to WebSocket: %v", err)
		return helper.RespondWithError(http.StatusInternalServerError, fmt.Errorf("failed to connect to WebSocket: %v", err))
	}
	defer conn.Close()

	// Define the data you want to send
	data := SocketMessage{
		Action: "notify",
		Data: map[string]interface{}{
			"user_id":        requestData.UserID,
			"discount_offer": discount.DiscountOffer,
			"connection_ids": user.ConnectionIDs,
		},
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data to JSON: %v", err)
		return helper.RespondWithError(http.StatusInternalServerError, fmt.Errorf("failed to marshal data: %v", err))
	}

	// Send the message over the WebSocket connection
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Printf("Failed to invoke websocket api: %v", err)
		return helper.RespondWithError(http.StatusInternalServerError, fmt.Errorf("failed to invoke websocket api: %v", err))
	}

	resp := Response{
		Message: "Successfully pushed data to socket api",
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error in marshalling response: %v", err)
		return helper.RespondWithError(http.StatusInternalServerError, fmt.Errorf("failed to marshal response: %v", err))
	}

	return events.APIGatewayProxyResponse{
		Body:       string(respBytes),
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}

func main() {
	lmd.Start(Handler)
}
