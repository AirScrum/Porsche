package queuepackage

/*
Import important libraries
*/
import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
This struct is for sending to the text message queue, for the NLP model to take this text and convert it to user stories
*/
type Message struct {
	TextID string `json:"textID"`
	Text   string `json:"text"`
	UserID string `json:"userID"`
}

/*
This struct is used, when we have a request from the gateway with the text id, to get the corresponding text from database, then send it to text queue
*/
type Request struct {
	TextID string `json:"textID"`
}

/*
This struct is used to get the user stories array from the user stories queue, that is received from the NLP model
*/
type UserStory struct {
	UserStories []string `json:"userStories"`
	TextID      string   `json:"textID"`
	UserID      string   `json:"userID"`
}

/*
This struct specifies all the needed attributes for managing queues
*/
type IQueue struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	queue     amqp.Queue
	channel   *amqp.Channel
	err       error
	name      string
	queueType string
}

/*
This function is for initializing the queues
*/
func QueueFactory(queueName string, queueType string) *IQueue {

	myQueue := new(IQueue)
	myQueue.name = queueName
	myQueue.queueType = queueType

	// To connect to the rabbitmq (message queue) server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	myQueue.channel, myQueue.err = conn.Channel()
	failOnError(myQueue.err, "["+myQueue.name+"]"+" - Failed to open"+myQueue.queueType+" Channel")
	//defer myQueue.channel.Close()

	/*
		Declaring the queue
	*/

	myQueue.queue, myQueue.err = myQueue.channel.QueueDeclare(
		myQueue.name, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	// Print error message when failing
	failOnError(myQueue.err, "["+myQueue.name+"]"+" - Failed to declare a "+myQueue.queueType+" queue")
	fmt.Println(myQueue.queue)

	// Define queue context
	myQueue.ctx, myQueue.ctxCancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer myQueue.ctxCancel()

	// Printing the context of the queue
	fmt.Println(myQueue.ctx)

	return myQueue
}

/*
This function is for sending in queue
*/
func SendToQueue(myQueue *IQueue, msg Message) {

	// To encode the msg object into array of bytes to be sent in the queue
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(msg); err != nil {
		log.Fatal(err)
	}

	// Send to the queue
	err := myQueue.channel.PublishWithContext(myQueue.ctx,
		"",                 // exchange
		myQueue.queue.Name, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buf.Bytes(),
		})

	// Print error message when failing
	failOnError(err, "["+myQueue.name+"] - Failed to publish a message to the queue")

	log.Printf("["+myQueue.name+"] - [x] Sent %s\n", msg.TextID)
}

/*
This function is for receiving from the queue
*/
func ReceiveFromQueueConc(myQueue *IQueue) {

	// Here we receive the message
	msgs, err := myQueue.channel.Consume(
		myQueue.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	// Print error message when failing
	if err != nil {
		panic(err)
	}

	// This function runs in the background and is used to always keep listening to the queue for any incoming message
	go func() {
		for d := range msgs {
			fmt.Printf("[%s] Received Message:\n %s\n\n", myQueue.name, d.Body)

			//TODO
			//Decode the message and deserialize it
			//to JSON format to be saved to the MongoDB

			// Save to database
			// Send ID to Gateway
		}
	}()
}

// This function is used to print the error messages
func failOnError(err error, msg string) {

	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
