package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	types "github.com/enriquenc/orderer-map-client-server-go/shared"
	"github.com/streadway/amqp"
)

func ParseRequestData(fileName, action, key, value string) []types.TestDataAction {
	var file *os.File

	var testData []types.TestDataAction

	if fileName != "" {
		var err error
		file, err = os.Open(fileName)
		if err != nil {
			log.Fatalf("Failed to open %s file: %v", fileName, err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		for decoder.More() {
			var testDataAction types.TestDataAction
			err := decoder.Decode(&testDataAction)
			if err != nil {
				log.Fatalf("Failed to decode test data: %v", err)
			}
			testData = append(testData, testDataAction)
		}
	} else {
		// Validate command line arguments
		if action == "" || !isValidAction(action) {
			log.Fatalf("Invalid action. Must be one of: %s, %s, %s, %s.", types.AddItem, types.GetItem, types.RemoveItem, types.GetAll)
			os.Exit(1)
		}
		if (action == types.AddItem || action == types.GetItem || action == types.RemoveItem) && key == "" {
			log.Fatalf("Key is required for %s, %s, and %s actions.", types.AddItem, types.GetItem, types.RemoveItem)
			os.Exit(1)
		}
		if (action == types.AddItem) && value == "" {
			log.Fatalf("Value is required for %s action.", types.AddItem)
			os.Exit(1)
		}

		testData = append(testData, types.TestDataAction{Key: key, Value: value, Action: action})
	}

	return testData
}

func main() {
	// Parse command line arguments
	rabbitMQURL := flag.String("rabbitmq-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ server address")
	queueName := flag.String("queue", "requests", "RabbitMQ queue name")
	fileName := flag.String("file", "", "File name to read actions")
	action := flag.String("action", "", "Action to perform: add, remove, get, or getAll")
	key := flag.String("key", "", "Key to use for the item")
	value := flag.String("value", "", "Value to use for the item")

	flag.Parse()

	testData := ParseRequestData(*fileName, *action, *key, *value)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(*rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		*queueName, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Record start time for performance measurement
	startTime := time.Now()
	// Publish message to queue
	for i := 0; i < len(testData); i++ {
		req := types.Request{
			Action: testData[i].Action,
			Key:    testData[i].Key,
			Value:  testData[i].Value,
		}
		reqBytes, err := json.Marshal(req)
		if err != nil {
			log.Fatalf("Failed to encode request: %v", err)
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        reqBytes,
			})
		if err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}
	}

	duration := time.Since(startTime) // Calculate duration for performance measurement
	fmt.Printf("Message published successfully in %v\n", duration)
}

func isValidAction(action string) bool {
	switch action {
	case types.AddItem, types.RemoveItem, types.GetItem, types.GetAll:
		return true
	default:
		return false
	}
}
