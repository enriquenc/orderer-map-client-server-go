package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	mq "github.com/enriquenc/orderer-map-client-server-go/mq"
)

func main() {
	// Parse command line arguments
	MQURL := flag.String("mq-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ server address")
	queueName := flag.String("queue", "requests", "RabbitMQ queue name")
	fileName := flag.String("file", "", "File name to read actions")
	action := flag.String("action", "", "Action to perform: add, remove, get, or getAll")
	key := flag.String("key", "", "Key to use for the item")
	value := flag.String("value", "", "Value to use for the item")

	flag.Parse()

	testData, err := parseRequestData(*fileName, *action, *key, *value)

	if err != nil {
		log.Fatalf("Failed parse request data: %v", err)
	}

	// Connect to MQ provider
	mq, err := mq.NewMQ(*MQURL, *queueName)

	if err != nil {
		log.Fatalf("Failed to connect to message queue provider: %v", err)
	}
	defer mq.Close()

	// Record start time for performance measurement
	startTime := time.Now()
	// Publish message to queue
	for i := 0; i < len(testData); i++ {
		mq.Publish(testData[i].RequestData)
		if err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}
	}

	duration := time.Since(startTime) // Calculate duration for performance measurement
	fmt.Printf("Message published successfully in %v\n", duration)
}
