package event

import (
	"context"
	"fmt"
	"listener/internal/messagebroker"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventHandler struct {
	client  *messagebroker.Client
	channel *amqp.Channel
	handlers map[string]MessageHandler
}

func NewEventHandler(client *messagebroker.Client) (*EventHandler, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}

	channel, err := client.CreateChannel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return &EventHandler{
		client:  client,
		channel: channel,
		handlers: make(map[string]MessageHandler),
	}, nil
}

func (eh *EventHandler) RegisterHandler(topic string, handler MessageHandler) {
	eh.handlers[topic] = handler
}

func (event *EventHandler) SetupExchangeAndQueue(exchangeName, queueName string, topics []string) error {
	exchangeErr := event.channel.ExchangeDeclare(
		exchangeName,
		"topic", // kind of exchange
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if exchangeErr != nil {
		return fmt.Errorf("failed to declare exchange: %w", exchangeErr)
	}

	_, qErr := event.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if qErr != nil {
		return fmt.Errorf("failed to declare queue: %w", qErr)
	}

	for _, topic := range topics {
		err := event.channel.QueueBind(queueName, topic, exchangeName, false, nil)
		if err != nil {
			return fmt.Errorf("failed to bind queue to topic %s: %w", topic, err)
		}
	}

	return nil
}

func (event EventHandler) Listen(ctx context.Context, queueName string) error {
	messages, consErr := event.channel.Consume(queueName, "", false, false, false, false, nil)
	if consErr != nil {
		return fmt.Errorf("failed to start consuming: %w", consErr)
	}

	log.Printf("Started listening on queue: %s", queueName)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-messages:
			if !ok {
				return fmt.Errorf("consumer message channel closed")
			}
			event.proccessMessage(ctx, msg)
		}
	}
}

func (event EventHandler) proccessMessage(ctx context.Context, msg amqp.Delivery) {
	// Convert AMQP delivery to our Message type
	message := Message{
		Body:       msg.Body,
		Headers:    msg.Headers,
		Topic:      msg.Exchange,
		RoutingKey: msg.RoutingKey,
		MessageID:  msg.MessageId,
	}

	log.Printf("successfully received message for %s topic and RoutingKey: %s\n", msg.Exchange, msg.RoutingKey)

	// Find appropriate handler
	handler, exists := event.handlers[msg.RoutingKey]
	if !exists {
		log.Printf("No handler registered for routing key: %s", msg.RoutingKey)
		msg.Nack(false, false) // Don't requeue unhandled messages
		return
	}

	// Process message
	if err := handler(ctx, message); err != nil {
		log.Printf("Handler failed for %s: %v", msg.RoutingKey, err)
		msg.Nack(false, true) // Requeue on handler failure
		return
	}

	// Acknowledge successful processing
	msg.Ack(false)
}

func (eh *EventHandler) Close() error {
	if eh.channel != nil {
		return eh.channel.Close()
	}
	return nil
}
