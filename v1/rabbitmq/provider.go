package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	RoutingKey string
	ContentType string
	Body []byte
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, message Message, entityName string) error {
	return r.channel.PublishWithContext(
		ctx, 
		entityName, 
		message.RoutingKey, 
		false, 
		false,
		amqp.Publishing{
			ContentType: message.ContentType,
			Body: message.Body,
		},
	)
}
