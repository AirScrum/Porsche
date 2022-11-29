package queuepackage

import (
	"context"
	"fmt"
	"log"
	"time"
	    "bytes"
    "encoding/gob"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	textID string `json:"textID"`
	text   string `json:"text"`
	userID string `json:"userID"`
}

type Request struct {
	textID string `json:"textID"`
	userID string `json:"userID"`
}

type IQueue struct {
	ctx                     context.Context
	ctxCancel               context.CancelFunc
	queue           	    amqp.Queue
	channel             	*amqp.Channel
	err 			      	 error
	name					string
	queueType				string
}

func QueueFactory(queueName string, queueType string) *IQueue {

	myQueue := new(IQueue)
	myQueue.name=queueName
	myQueue.queueType=queueType

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	myQueue.channel, myQueue.err = conn.Channel()
	failOnError(myQueue.err, "[" + myQueue.name +"]" + " - Failed to open"+ myQueue.queueType+ " Channel")
	//defer myQueue.channel.Close()

	/*
		Declaring the queue
	*/

	myQueue.queue, myQueue.err = myQueue.channel.QueueDeclare(
		myQueue.name, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(myQueue.err, "[" + myQueue.name +"]" + " - Failed to declare a "+myQueue.queueType+" queue")
	fmt.Println(myQueue.queue)

	myQueue.ctx, myQueue.ctxCancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer myQueue.ctxCancel()

	fmt.Println(myQueue.ctx)

	return myQueue
}

func SendToQueue(myQueue IQueue, msg Message){
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(msg); err != nil {
		log.Fatal(err)
	}

	err := myQueue.channel.PublishWithContext(myQueue.ctx,
		"",                     // exchange
		myQueue.queue.Name, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buf.Bytes(),
		})
	failOnError(err, "["+myQueue.name+"] - Failed to publish a message to the queue")

	log.Printf("["+myQueue.name+"] - [x] Sent %s\n", msg.textID)
}


func ReceiveFromQueueConc(myQueue *IQueue) {
		msgs, err := myQueue.channel.Consume(
			myQueue.queue.Name,
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
				fmt.Printf("Received Message from: %s\n", d.Body)
				// Save to database
				// Send ID to Gateway
			}
		}()
}
