package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func connectDatabase() {
	//Connect to database
}

func connectQueueText() {
	/*
	* Queue 1
	 */
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
	}

	statuesQueue, err := ch.QueueDeclare(
		"TextQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	fmt.Println(statuesQueue)
	// Handle any errors if we were unable to create the queue
	if err != nil {
		fmt.Println(err)
	}

}

func connectQueueUserStories() {
	/*
	* Queue 2
	 */
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
	}

	statuesQueue, err := ch.QueueDeclare(
		"UserStoriesQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	fmt.Println(statuesQueue)
	// Handle any errors if we were unable to create the queue
	if err != nil {
		fmt.Println(err)
	}

}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint")
}

func handleRequests() {
	http.HandleFunc("/user", homepage)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func main2() {
	connectDatabase()
	connectQueueText()
	connectQueueUserStories()
	handleRequests()
}
