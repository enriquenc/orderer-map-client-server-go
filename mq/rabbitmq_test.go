package rabbitmq

import (
	"testing"

	types "github.com/enriquenc/orderer-map-client-server-go/shared"
	"github.com/stretchr/testify/assert"
)

func TestNewMQ(t *testing.T) {
	_, err := NewMQ("invalid_url", "unit_tests")
	assert.Error(t, err)

	_, err = NewMQ("amqp://guest:guest@localhost:5672/", "unit_tests")
	assert.NoError(t, err)
}

func TestRabbitMQ_Close(t *testing.T) {
	r, err := NewMQ("amqp://guest:guest@localhost:5672/", "unit_tests")
	assert.NoError(t, err)

	err = r.Close()
	assert.NoError(t, err)

	err = r.Close()
	assert.Error(t, err)
}

func TestRabbitMQ_Publish(t *testing.T) {
	r, err := NewMQ("amqp://guest:guest@localhost:5672/", "unit_tests")
	assert.NoError(t, err)

	err = r.Publish(types.Request{Key: "k1", Value: "v1", Action: "add"})
	assert.NoError(t, err)
}

func TestRabbitMQ_Consume(t *testing.T) {
	r, err := NewMQ("amqp://guest:guest@localhost:5672/", "unit_tests")
	assert.NoError(t, err)

	ch, err := r.Consume()
	assert.NoError(t, err)

	go func() {
		err = r.Publish(types.Request{Key: "k1", Value: "v1", Action: "add"})
		assert.NoError(t, err)
	}()

	req := <-ch
	assert.Equal(t, types.Request{Key: "k1", Value: "v1", Action: "add"}, req)
}
