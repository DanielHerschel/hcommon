package rabbitmq

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	amqp "github.com/rabbitmq/amqp091-go"
)


func Test_NewRabbitMQ(t *testing.T) {
	t.Run("test connection successful", func(t *testing.T) {
		conn := &amqp.Connection{}
		ch := &amqp.Channel{}
		connPatch := gomonkey.ApplyFuncReturn(amqp.Dial, conn, nil)
		chanPatch := gomonkey.ApplyMethodReturn(conn, "Channel", ch, nil)
		defer connPatch.Reset()
		defer chanPatch.Reset()

		client, err := NewRabbitMQ(RabbitMQParams{"localhost", 5672, "admin", "admin"})
		assert.Nil(t, err)
		assert.Equal(t, &RabbitMQ{conn, ch}, client)
	})
}
