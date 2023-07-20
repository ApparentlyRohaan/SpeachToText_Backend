package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var RabbitMQChannel *amqp.Channel

func RabbitMQConnect() {
	// RabbitMQ connection string
	connectionString := "amqp://guest:guest@localhost:5672/"

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	RabbitMQChannel = ch

	// Declare a queue
	queueName := "hello"
	q, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable (messages survive server restarts)
		false,     // Delete when unused
		false,     // Exclusive (for the connection which declares it)
		false,     // No-wait (for queue declaration)
		nil,       // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Message to send
	message := "file.wav"

	// Convert message to []byte
	body := []byte(message)

	// Publish the message to the queue
	err = ch.Publish(
		"",     // Exchange name (default exchange)
		q.Name, // Routing key (queue name)
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf("Sent: %s\n", message)

	// Wait to receive the response from Python
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto-acknowledge (messages are automatically acknowledged)
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	failOnError(err, "Failed to register a consumer")

	// Receive the response
	for msg := range msgs {
		fmt.Printf("Received response from Python: %s\n", msg.Body)
		break // We only need one response, so break the loop after receiving the first one.
	}
}
