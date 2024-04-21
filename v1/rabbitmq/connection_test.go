package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)


func Test_NewRabbitMQ(t *testing.T) {
	conn := &amqp.Connection{}
	ch := &amqp.Channel{}

	t.Run("test connection successful", func(t *testing.T) {
		connPatch := gomonkey.ApplyFuncReturn(amqp.Dial, conn, nil)
		chanPatch := gomonkey.ApplyMethodReturn(conn, "Channel", ch, nil)
		defer connPatch.Reset()
		defer chanPatch.Reset()

		client, err := NewRabbitMQ(RabbitMQParams{"localhost", 5672, "admin", "admin"})

		assert.Nil(t, err)
		assert.Equal(t, &RabbitMQ{conn, ch}, client)
	})

	t.Run("test invalid url", func(t *testing.T) {
		client, err := NewRabbitMQ(RabbitMQParams{"123", 0, "", ""})

		assert.NotNil(t, err)
		assert.Nil(t, client)
	})

	t.Run("test could not create channel", func(t *testing.T) {
		connPatch := gomonkey.ApplyFuncReturn(amqp.Dial, conn, nil)
		chanPatch := gomonkey.ApplyMethodReturn(conn, "Channel", nil, fmt.Errorf("channel error"))
		defer connPatch.Reset()
		defer chanPatch.Reset()
		
		client, err := NewRabbitMQ(RabbitMQParams{"123", 0, "", ""})

		assert.NotNil(t, err)
		assert.Nil(t, client)
	})
}

func Test_DeclareExchange(t *testing.T) {
	rmq := &RabbitMQ{&amqp.Connection{}, &amqp.Channel{}}
	t.Run("test queue created", func(t *testing.T) {
		qDeclarePatch := gomonkey.ApplyMethodReturn(rmq.channel, "ExchangeDeclare", nil)
		defer qDeclarePatch.Reset()

		err := rmq.DeclareExchange("", "", true, true)

		assert.Nil(t, err)
	})

	t.Run("test error creating queue", func(t *testing.T) {
		qDeclarePatch := gomonkey.ApplyMethodReturn(rmq.channel, "ExchangeDeclare", fmt.Errorf("exchange declare error"))
		defer qDeclarePatch.Reset()

		err := rmq.DeclareExchange("", "", true, true)

		assert.NotNil(t, err)
	})
}

func Test_DeclareQueue(t *testing.T) {
	rmq := &RabbitMQ{&amqp.Connection{}, &amqp.Channel{}}
	t.Run("test queue created", func(t *testing.T) {
		qDeclarePatch := gomonkey.ApplyMethodReturn(rmq.channel, "QueueDeclare", amqp.Queue{}, nil)
		defer qDeclarePatch.Reset()

		err := rmq.DeclareQueue("test", true, true)

		assert.Nil(t, err)
	})

	t.Run("test error creating queue", func(t *testing.T) {
		qDeclarePatch := gomonkey.ApplyMethodReturn(rmq.channel, "QueueDeclare", nil, fmt.Errorf("queue declare error"))
		defer qDeclarePatch.Reset()

		err := rmq.DeclareQueue("test", true, true)

		assert.NotNil(t, err)
	})
}

func Test_Close(t *testing.T) {
	rmq := &RabbitMQ{&amqp.Connection{}, &amqp.Channel{}}
	t.Run("test doesn't panic", func(t *testing.T) {
		chanClosePatch := gomonkey.ApplyMethodReturn(rmq.channel, "Close", nil)
		connClosePatch := gomonkey.ApplyMethodReturn(rmq.connection, "Close", nil)
		defer chanClosePatch.Reset()
		defer connClosePatch.Reset()

		assert.NotPanics(t, rmq.Close)
	})
}
