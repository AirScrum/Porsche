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
	meeting := models.ModelResponse{
		UserID:"63f772f72db76b133592617c",
		TextID: "64009b4eff8e68d7cd8b7955",
		UserStories: []models.UserStory{
			{
				UserStoryTitle:     "Login",
				UserStoryDescription: "As a user, I want to log in",
			},
			{
				UserStoryTitle:     "Sign up",
				UserStoryDescription: "As a user, I want to sign up",
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
