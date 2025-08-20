package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"listener/pkg/logger"
	// amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler defines how to process incoming messages
type MessageHandler func(ctx context.Context, message Message) error

type Message struct {
	Body       []byte
	Headers    map[string]interface{}
	Topic      string
	RoutingKey string
	MessageID  string
}

func HandleLogMessage(ctx context.Context, message Message) error {
	var payload logger.LogPayload
	json.Unmarshal(message.Body, &payload)

	// // In case defining the requestID inside the LogPayload in both sides, Extract request ID from headers if available
	// if message.Headers != nil {
	// 	if requestID, ok := message.Headers["request_id"].(string); ok {
	// 		payload.RequestID = requestID
	// 	}
	// }

	return logEvent(payload)
}


func logEvent(logPayload logger.LogPayload) error {
	jsonData, _ := json.MarshalIndent(logPayload, "", "\t")
	logServiceURL := "http://logger-service/log"
	
	request, reqErr := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if reqErr != nil {
		log.Println("couldn't create request for sending to logger")
		return reqErr
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, respErr := client.Do(request)
	if respErr != nil {
		log.Println("couldn't send request to logger")
		return respErr
	}
	defer response.Body.Close()

	// in case of error, return correct status code
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("faulty logger service")
	}

	return nil
}
