package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func connectDatabase() {
	//Connect to database
}

func connectQueueText() (interface{}, interface{}) {
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

	return ch, err

}

func connectQueueUserStories() (interface{}, interface{}) {
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

	return ch, err
}

/*
func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint")
}

func handleRequests(textChannel interface{}, errText interface{}, userChannel interface{}, errUser interface{}) {
	http.HandleFunc("/user", homepage)
	log.Fatal(http.ListenAndServe(":8001", nil))
	// attempt to publish a message to the queue!
	errText = textChannel.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)

	if errText != nil {
		fmt.Println(errText)
	}
	fmt.Println("Successfully Published Message to Queue")
}
*/
/*
func main2() {
	connectDatabase()
	textChannel, errText := connectQueueText()
	userChannel, errUser := connectQueueUserStories()
	handleRequests(textChannel, errText, userChannel, errUser)
}*/
