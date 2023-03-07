package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

const (
	addItem    = "add"
	removeItem = "remove"
	getItem    = "get"
	getAll     = "getAll"
)

type Request struct {
	Action string
	Key    string
	Value  string
}

type Action struct {
	Command          string
	Key              string
	Value            string
	ExpectedResponse interface{}
}

func main() {
	// Parse command line arguments
	serverAddr := flag.String("server", "localhost:5672", "RabbitMQ server address")
	queueName := flag.String("queue", "requests", "RabbitMQ queue name")
	fileName := flag.String("file", "", "File name to read actions")
	action := flag.String("action", "", "Action to perform: add, remove, get, or getAll")
	key := flag.String("key", "", "Key to use for the item")
	value := flag.String("value", "", "Value to use for the item")

	flag.Parse()

	var file *os.File
	var testData []Action
	if *fileName != "" {
		var err error
		file, err = os.Open(*fileName)
		if err != nil {
			log.Fatalf("failed to open %s file: %v", *fileName, err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		for decoder.More() {
			var action Action
			err := decoder.Decode(&action)
			if err != nil {
				log.Fatalf("failed to decode test data: %v", err)
			}
			testData = append(testData, action)
		}
	} else {
		// Validate command line arguments
		if *action == "" || !isValidAction(*action) {
			fmt.Println("Invalid action. Must be one of: add, remove, get, getAll")
			os.Exit(1)
		}
		if (*action == addItem || *action == getItem) && *key == "" {
			fmt.Println("Key is required for add and get actions")
			os.Exit(1)
		}
		if *action == addItem && *value == "" {
			fmt.Println("Value is required for add action")
			os.Exit(1)
		}

	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s/", *serverAddr))
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
	if *fileName != "" {
		for i := 0; i < len(testData); i++ {
			req := Request{
				Action: testData[i].Command,
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
	} else {
		// Build request message
		req := Request{
			Action: *action,
			Key:    *key,
			Value:  *value,
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
	case addItem, removeItem, getItem, getAll:
		return true
	default:
		return false
	}
}
