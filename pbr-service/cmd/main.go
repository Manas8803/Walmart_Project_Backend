package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Manas8803/Walmart_Project_Backend/pbr-service/pkg/models/service"
	"github.com/aws/aws-lambda-go/events"
	lmd "github.com/aws/aws-lambda-go/lambda"cd
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
		log.Printf("Failed to unmarshal request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}
	log.Println("Received beacon_id: " + requestData.BeaconID + ", user_id: " + requestData.UserID)

	// Fetch discount
	discount, err := service.FetchDiscountByBeaconID(requestData.BeaconID)
	if err != nil {
		log.Println("Error in fetching discount : ", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// Fetch user
	user, err := service.FetchUserByID(requestData.UserID)
	if err != nil {
		log.Println("Error in fetching user : ", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// Push the data in websocket connection
	conn, _, err := websocket.DefaultDialer.Dial(os.Getenv("WEBSOCKET_URL"), nil)
	if err != nil {
		log.Printf("Failed to connect to WebSocket: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to connect to WebSocket",
		}, nil
	}
	defer conn.Close()

	// Define the data you want to send
	data := SocketMessage{
		Action: "report", //! CHANGE THIS
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
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to marshal data",
		}, nil
	}

	// Send the message over the WebSocket connection
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Printf("Failed to invoke websocket api: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to invoke websocket api",
		}, nil
	}

	resp := Response{
		Message: "Successfully pushed data to socket api",
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error in marshalling response : ", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(respBytes),
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}

func main() {
	lmd.Start(Handler)
}
