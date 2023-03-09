// rabbitmq.go

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	types "github.com/enriquenc/orderer-map-client-server-go/shared"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewMQ(rabbitMQURL, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return fmt.Errorf("failed to close RabbitMQ channel: %v", err)
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("failed to close RabbitMQ connection: %v", err)
		}
	}

	return nil
}

func (r *RabbitMQ) Consume() (<-chan types.Request, error) {
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %v", err)
	}

	requests := make(chan types.Request)

	go func() {
		for msg := range msgs {
			var req types.Request
			if err := json.Unmarshal(msg.Body, &req); err != nil {
				log.Printf("failed to decode message: %v", err)
				continue
			}
			requests <- req
		}
	}()

	return requests, nil
}

func (r *RabbitMQ) Publish(req types.Request) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode message: %v", err)
	}

	err = r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}
