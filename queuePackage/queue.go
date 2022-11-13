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
var textChannel *amqp.Channel

var userStories amqp.Queue
var userStoriesChannel *amqp.Channel

func Start() {
	/*
		Connection to RabbitMQ, TextChannel, and UserStoriesChannel
	*/
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	textChannel, textChannelError := conn.Channel()
	failOnError(textChannelError, "[textQueue] - Failed to open text Channel")
	defer textChannel.Close()

	userStoriesChannel, userStoriesChannelError := conn.Channel()
	failOnError(userStoriesChannelError, "[userStoriesQueue] - Failed to open userStories Channel")
	defer userStoriesChannel.Close()

	/*
		Declaring textQueue, and userStoriesQueue
	*/

	textQueue, textChannelError := textChannel.QueueDeclare(
		"textQueue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(textChannelError, "[textQueue] - Failed to declare a text queue")
	fmt.Println(textQueue)

	userStoriesQueue, userStoriesChannelError := userStoriesChannel.QueueDeclare(
		"userStoriesQueue", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(userStoriesChannelError, "[userStoriesQueue] - Failed to declare a userStories queue")
	fmt.Println(userStoriesQueue)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println(ctx)

	err = textChannel.PublishWithContext(ctx,
		"",             // exchange
		textQueue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("msg"),
		})
	failOnError(err, "[textQueue] - Failed to publish a message to textQueue")
	log.Printf("[textQueue] - [x] Sent %s\n", "msg")

}

/*
func SendToTextQueue(msg string) {

	err := textChannel.PublishWithContext(ctx,
		"",             // exchange
		textQueue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "[textQueue] - Failed to publish a message to textQueue")
	log.Printf("[textQueue] - [x] Sent %s\n", msg)

}*/

func sendToUserStoriesQueue(msg string) {
	/*
		TODO Implement  sendToUserStoriesQueue logic
	*/
	err := userStoriesChannel.PublishWithContext(ctx,
		"",               // exchange
		userStories.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "[userStoriesQueue] - Failed to publish a message to userStoriesQueue")
	log.Printf("[userStoriesQueue] - [x] Sent %s\n", msg)
}
func failOnError(err error, msg string) {

	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
