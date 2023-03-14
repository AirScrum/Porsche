package main

import (
	"fmt"
	"goserver/models"
	"github.com/streadway/amqp"
	"encoding/json"

)

func main() {
	fmt.Println("Go RabbitMQ Tutorial")
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	/*q, err := ch.QueueDeclare(
		"userStoriesQueue",
		false,
		false,
		false,
		false,
		nil,
	)*/
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	//fmt.Println(q)
	// Handle any errors if we were unable to create the queue
	//if err != nil {
	//	fmt.Println(err)
	//}

	// attempt to publish a message to the queue!
	meeting := models.Meeting{
		MeetingID:    "1283182",
		MeetingTitle: "Hello World",
		MeetingUserStories: []models.UserStory{
			{
				UserStoryID:        "12831823",
				UserStoryTitle:     "Baby",
				UserStoryDescription: "Its U that i need",
			},
			{
				UserStoryID:        "12831823",
				UserStoryTitle:     "Baby",
				UserStoryDescription: "Its u that i need",
			},
		},
	}
	meetingBytes, err := json.Marshal(meeting)
	if err != nil {
		fmt.Println("Error marshaling Person:", err)
		return
	}
fmt.Println(meeting)
	err = ch.Publish(
		"",
		"userStoriesQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(meetingBytes),
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")
}
