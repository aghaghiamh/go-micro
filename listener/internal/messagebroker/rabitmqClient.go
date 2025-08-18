package messagebroker

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO: in case you assume this package would become deprecated or you'll use another queue tech., implement ports-and-adapters design pattern

type Client struct {
	conn   *amqp.Connection
	config Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Connect() error {
	connectionURL := c.config.ConnectionString()
	conn, connErr := amqp.Dial(connectionURL)
	if connErr != nil {
		errMsg := "failed to connect to RabbitMQ"
		log.Println(errMsg, ": ", connErr)

		return fmt.Errorf(errMsg)
	}

	log.Println("✅ Successfully connected to Rabbitmq!")
	c.conn = conn
	return nil
}

func (c *Client) IsConnected() bool {
	return c.conn != nil && !c.conn.IsClosed()
}

func (c *Client) Close() error {
	if c.conn != nil && !c.conn.IsClosed() {
		return c.conn.Close()
	}
	return nil
}

// Use this carefully and consider whether you need direct access
func (c *Client) GetConnection() *amqp.Connection {
	return c.conn
}
