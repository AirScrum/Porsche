package queuepackage

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ctx context.Context
var textQueue amqp.Queue
var textChannel amqp.Channel

func start() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	textChannel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer textChannel.Close()

	textQueue, err := textChannel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	fmt.Println(textQueue)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println(ctx)

}

func sendToTextQueue(msg string) {

	err := textChannel.PublishWithContext(ctx,
		"",             // exchange
		textQueue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msg)

}

func sendToUserStoriesQueue(msg string) {
	/*
		TODO Implement  sendToUserStoriesQueue
	*/
}
func failOnError(err error, msg string) {

	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
