package messagebroker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublisherConfig holds publisher-specific configuration.
type PublisherConfig struct {
	ExchangeName string
	ConfirmMode  bool
}

// DefaulPublisherConfig returns sensible default configuration
func DefaulPublisherConfig() PublisherConfig {
	return PublisherConfig{
		ExchangeName: "events",
		ConfirmMode:  false,
	}
}

type Publisher struct {
	client  *Client
	channel *amqp.Channel
	config  PublisherConfig
	mutex   sync.Mutex
}

func NewPublisher(client *Client, config PublisherConfig) (*Publisher, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}

	channel, err := client.CreateChannel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return &Publisher{
		client:  client,
		channel: channel,
		config:  config,
	}, nil
}

func (p *Publisher) Setup(config PublisherConfig) error {
	exchangeErr := p.channel.ExchangeDeclare(
		config.ExchangeName,
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

	// Setup Confirmation Mode
	if config.ConfirmMode {
		if err := p.channel.Confirm(false); err != nil {
			p.channel.Close()
			return fmt.Errorf("failed to enable confirm mode: %w", err)
		}
	}

	return nil
}

// TODO: better to hide out the publishing in one other structure or the Publisher
func (p *Publisher) Publish(ctx context.Context, routingKey string, publishing amqp.Publishing) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	select {
	case <-ctx.Done():
		return fmt.Errorf("publish cancelled: %w", ctx.Err())
	default:
		// The actual publish call
		err := p.channel.PublishWithContext(ctx,
			p.config.ExchangeName, // exchange
			routingKey,            // routing key
			false,                 // mandatory
			false,                 // immediate
			publishing,
		)
		if err != nil {
			return fmt.Errorf("failed to publish message: %w", err)
		}

		log.Printf("message with %s Exchange and %s RoutingKey has been published.\n", p.config.ExchangeName, routingKey)
	}

	// Wait for confirmation if confirm mode is enabled.
	if p.config.ConfirmMode {
		confirmation := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		select {
		case confirm := <-confirmation:
			if !confirm.Ack {
				return fmt.Errorf("message not acknowledged by broker")
			}
		case <-time.After(5 * time.Second): // Timeout for confirmation
			return fmt.Errorf("publish confirmation timeout")
		}
	}

	return nil
}

func (p *Publisher) Close() error {
	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}
