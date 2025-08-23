package events

import (
	"brocker/internal/entities"
	"brocker/pkg/events"
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(ctx context.Context, routingKey string, publishing amqp.Publishing) error
	Close() error
}

// EventService provides a high-level API for publishing business events.
type EventService struct {
	publisher   Publisher
	serviceName string
}

func NewEventService(publisher Publisher, serviceName string) *EventService {
	return &EventService{
		publisher:   publisher,
		serviceName: serviceName,
	}
}

func (s *EventService) PublishLogEvent(ctx context.Context, logPayload entities.LogPayload) error {
	logEvent := &events.LogEvent{
		Name:   logPayload.Name,
		Data:   logPayload.Data,
		Level:  events.LogLevel(logPayload.Level),
		Source: s.serviceName,
	}

	body, err := logEvent.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal log event: %w", err)
	}

	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		Timestamp:    time.Now(),
		DeliveryMode: amqp.Persistent, // Messages should be durable
	}

	// Use the log level as the routing key.
	routingKey := "log." + string(logEvent.Level)

	return s.publisher.Publish(ctx, routingKey, publishing)
}

func (s *EventService) Close() error {
	return s.publisher.Close()
}
