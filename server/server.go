package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	logger "server/logger"
	requestmanager "server/request-manager"

	mq "github.com/enriquenc/orderer-map-client-server-go/mq"
)

func main() {
	// Parse command line arguments
	MQURL := flag.String("mq-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
	queueName := flag.String("queue", "requests", "RabbitMQ queue name")
	logFile := flag.String("log-file", "server.log", "Log file name")
	flag.Parse()

	// Connect to MQ
	mq, err := mq.NewMQ(*MQURL, *queueName)
	if err != nil {
		log.Fatalf("Failed to connect to message queue provider: %v", err)
	}
	// Close the RabbitMQ connection before exiting the program
	defer mq.Close()

	// Start consuming messages from the message queue.
	// This method runs the goroutine which reads request data
	// from the message queue and pushes it to the returned channel
	requestProcessingChannel, err := mq.Consume()
	if err != nil {
		log.Fatalf("Failed to consume from message queue. %v", err)
	}

	logger, err := logger.NewLogger(*logFile)
	if err != nil {
		log.Fatalf("Failed to create new logger. %v", err)
	}
	defer logger.Close()
	// Start processing the requests in a separate goroutine in parallel
	go requestmanager.ProcessRequests(requestProcessingChannel, logger)

	// Set up signal handler to gracefully exit the program on interrupt signal
	interruptSignalChannel := make(chan os.Signal, 1)
	signal.Notify(interruptSignalChannel, os.Interrupt)

	// Wait for the interrupt signal to exit the program
	select {
	case <-interruptSignalChannel:
		fmt.Println("Interrupt signal received. Exiting the program...")

		mq.Close()

		return
	}
}
