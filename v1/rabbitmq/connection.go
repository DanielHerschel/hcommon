package rabbitmq

import (
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQParams struct {
	hostname string
	port int
	username string
	password string
}

func NewRabbitMQ(params RabbitMQParams) (*RabbitMQ, error) {
	url := strings.Join(
		[]string{
			"amqp://", params.username, ":", params.password, 
			"@", params.hostname, ":", fmt.Sprintf("%d", params.port), "/",
		}, 
	"")

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{connection: conn, channel: ch}, nil
}

type RabbitMQ struct {
	connection *amqp.Connection
	channel *amqp.Channel
	queue amqp.Queue
}

func (r *RabbitMQ) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	return r.channel.ExchangeDeclare(name, kind, durable, autoDelete, false, false, nil)
}

func (r *RabbitMQ) DeclareQueue(name string, durable, autoDelete bool) error {
	q, err := r.channel.QueueDeclare(name, durable, autoDelete, false, false, nil)
	if err != nil {
		return err
	}
	r.queue = q
	return nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.connection.Close()
}
