package rabbitmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func Test_PublishMessage(t *testing.T) {
	t.Run("publishes message", func(t *testing.T) {
		publishPatch := gomonkey.ApplyMethodReturn(rmq.channel, "PublishWithContext", nil)
		defer publishPatch.Reset()

		err := rmq.PublishMessage(context.Background(), Message{}, "")

		assert.Nil(t, err)
	})

	t.Run("publish returns error", func(t *testing.T) {
		publishPatch := gomonkey.ApplyMethodReturn(rmq.channel, "PublishWithContext", fmt.Errorf("publish failed"))
		defer publishPatch.Reset()

		err := rmq.PublishMessage(context.Background(), Message{}, "")

		assert.NotNil(t, err)
	})
}
