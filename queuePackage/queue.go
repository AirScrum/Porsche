package queuepackage

/*
Import important libraries
*/
import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	dbpackage "goserver/dbPackage"
	"goserver/models"
	"log"
	"time"
)

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
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
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
func SendToQueue(myQueue *IQueue, msg models.Message) {

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
	// This function runs in the background and is used to always keep listening to the queue for any incoming message
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
	go func() {
		//Loop on all messages in the queue
		for d := range msgs {
			fmt.Printf("[%s] Received Message:\n %s\n\n", myQueue.name, d.Body)
			var modelResponse models.ModelResponse
			//Converting message from bytes array to meeting format
			err = json.Unmarshal(d.Body, &modelResponse)
			if err != nil {
				fmt.Println("Meeting marshalling error:", err)
				return
			}
			//Converting meeting to meetingModel before adding it
			//to the database
			for i, userStory := range modelResponse.UserStories {
				//Converting textID into ObjectID
				mappedTextID, err := primitive.ObjectIDFromHex(modelResponse.TextID)
				if err != nil {
					panic(err)
				}
				//Converting userID into ObjectID
				mappedUserID, err := primitive.ObjectIDFromHex(modelResponse.UserID)
				if err != nil {
					panic(err)
				}
				//Constructing the userStory model to match the required format in the database
				userStoryModel := models.UserStoryModel{
					TextID:               mappedTextID,
					UserID:               mappedUserID,
					UserStoryTitle:       userStory.UserStoryTitle,
					UserStoryDescription: userStory.UserStoryDescription,
				}
				res, err := dbpackage.InsertUserStory(userStoryModel)
				if err != nil {
					log.Fatal("User Story InsertOne Error: ", err, i)
				}
				fmt.Println("Inserted ", res.InsertedID, i)
			}
			//TODO Create a POST request on a route
			//to notify user that userStories were saved
		}
	}()
}

// This function is used to print the error messages
func failOnError(err error, msg string) {

	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
