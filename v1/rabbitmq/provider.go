package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	RoutingKey  string
	ContentType string
	Body        []byte
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
			Body:        message.Body,
		},
	)
}

func (r *RabbitMQ) Consume(messagesChan chan Message, queue string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	msgsChan, err := r.channel.ConsumeWithContext(ctx, queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgsChan {
		messagesChan <- Message{Body: msg.Body}
	}
	return nil
}

func (r *RabbitMQ) ConsumeAmount(messagesChan chan Message, queue string, amount uint) error {
	msgsChan, err := r.channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for currAmount := uint(0); currAmount < amount; currAmount++ {
		messagesChan <- Message{Body: (<-msgsChan).Body}
		currAmount++
	}
	return nil
}
