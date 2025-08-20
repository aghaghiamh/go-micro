package messagebroker

import "fmt"

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.Username, c.Password, c.Host, c.Port)
}
