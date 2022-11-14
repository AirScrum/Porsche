package queuepackage

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	ctx                     context.Context
	ctxCancel               context.CancelFunc
	textQueue               amqp.Queue
	textChannel             *amqp.Channel
	textChannelError        error
	userStoriesQueue        amqp.Queue
	userStoriesChannel      *amqp.Channel
	userStoriesChannelError error
}

func QueueInit() *Queue {

	myQueues := new(Queue)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	myQueues.textChannel, myQueues.textChannelError = conn.Channel()
	failOnError(myQueues.textChannelError, "[textQueue] - Failed to open text Channel")
	//defer myQueues.textChannel.Close()

	myQueues.userStoriesChannel, myQueues.userStoriesChannelError = conn.Channel()
	failOnError(myQueues.userStoriesChannelError, "[userStoriesQueue] - Failed to open userStories Channel")
	//defer myQueues.userStoriesChannel.Close()

	/*
		Declaring textQueue, and userStoriesQueue
	*/

	myQueues.textQueue, myQueues.textChannelError = myQueues.textChannel.QueueDeclare(
		"textQueue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(myQueues.textChannelError, "[textQueue] - Failed to declare a text queue")
	fmt.Println(myQueues.textQueue)

	myQueues.userStoriesQueue, myQueues.userStoriesChannelError = myQueues.userStoriesChannel.QueueDeclare(
		"userStoriesQueue", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(myQueues.userStoriesChannelError, "[userStoriesQueue] - Failed to declare a userStories queue")
	fmt.Println(myQueues.userStoriesQueue)

	myQueues.ctx, myQueues.ctxCancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer myQueues.ctxCancel()

	fmt.Println(myQueues.ctx)

	return myQueues
}

func SendToTextQueue(myQueue *Queue, textid string, text string, userid string) {

	var msg string = textid + "$" + userid + "$" + text
	err := myQueue.textChannel.PublishWithContext(myQueue.ctx,
		"",                     // exchange
		myQueue.textQueue.Name, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "[textQueue] - Failed to publish a message to textQueue")

	log.Printf("[textQueue] - [x] Sent %s\n", msg)

}

func ReceiveFromUserStoriesQueue(myQueue *Queue) {

	msgs, err := myQueue.userStoriesChannel.Consume(
		myQueue.userStoriesQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	go func() {
		for d := range msgs {
			fmt.Printf("Recieved Message: %s\n", d.Body)
			// Save to database
			// Send ID to Gateway
		}
	}()

}
