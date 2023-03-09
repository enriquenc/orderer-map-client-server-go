package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	orderermap "server/orderer-map"

	types "github.com/enriquenc/orderer-map-client-server-go/shared"

	"github.com/streadway/amqp"
)

func main() {
	// Parse command line arguments
	rabbitMQURL := flag.String("rabbitmq-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
	queueName := flag.String("queue", "requests", "RabbitMQ queue name")
	logFile := flag.String("log-file", "server.log", "Log file name")
	flag.Parse()

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(*rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	dataStorage := orderermap.NewOrderedMap()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
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

	requestProcessingChannel := make(chan types.Request)

	go func() {
		chWorker, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
		}
		defer chWorker.Close()
		msgs, err := chWorker.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}

		for msg := range msgs {
			var req types.Request
			if err := json.Unmarshal(msg.Body, &req); err != nil {
				fmt.Printf("Failed to decode message: %s", err)
				continue
			}
			requestProcessingChannel <- req
			// log.Printf("received message: %s", msg.Body)
			// time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for req := range requestProcessingChannel {
			// Processing of write requests
			switch req.Action {
			case "add":
				dataStorage.Add(req.Key, req.Value)
				fmt.Fprintf(f, "[add] Added key %s with value %s\n", req.Key, req.Value)
			case "remove":
				dataStorage.Remove(req.Key)
				fmt.Fprintf(f, "[remove] key %s\n", req.Key)
			case "get":
				value, exists := dataStorage.Get(req.Key)
				if exists {
					fmt.Fprintf(f, "[get] Got key %s with value %s\n", req.Key, value)
				} else {
					fmt.Fprintf(f, "[get] Key %s doesn't exist\n", req.Key)
				}
			case "getAll":
				items := dataStorage.GetAll()
				b, _ := json.Marshal(items)
				fmt.Fprintf(f, "[getAll] All values %s\n", string(b))
			}
		}
	}()

	// Set up signal handler to gracefully exit the program on interrupt signal
	interruptSignalChannel := make(chan os.Signal, 1)
	signal.Notify(interruptSignalChannel, os.Interrupt)

	select {
	case <-interruptSignalChannel:
		fmt.Println("Interrupt signal received. Exiting the program...")
		close(requestProcessingChannel)
		return
	}
}
