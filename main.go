package main

/*
Import important libraries
*/
import (
	"encoding/json"
	"fmt"
	queuepackage "goserver/queuePackage"
	"io/ioutil"
	"log"
	"net/http"
)

/*
Define the queues we need
*/
var userStoriesQueue *queuepackage.IQueue
var textQueue *queuepackage.IQueue

/*
This function is called when the is a request in our server, which contains the textid and needed to be sent to the text queue
*/
func homepage(w http.ResponseWriter, r *http.Request) {

	// Read the request body and parse it from JSON to Message Struct
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)

	// Print error message when failing
	if err != nil {
		panic(err)
	}

	// Take the request body and put it in an object
	request := queuepackage.Request{}
	err = json.Unmarshal(buf, &request)
	if err != nil {
		panic(err)
	}

	// Define Message that will be sent to the text queue
	msg := queuepackage.Message{}

	// Get corresponding text from the TextID from database
	msg.TextID = request.TextID
	msg.Text = "Testing"
	msg.UserID = request.UserID

	// Send the message text queue
	queuepackage.SendToQueue(textQueue, msg)

	/*
		TODO Read the textID from the database
	*/
	/*
		TODO Send the text fetched from the database to the TextQueue
	*/
}

/*
This function is called when a request is sent to our server, and sent the request to homepage function to handle it
*/
func handleRequests() {
	fmt.Println("Server Started")
	http.HandleFunc("/main", homepage)
	log.Fatal(http.ListenAndServe(":8002", nil))
}

/*
This is our main function
*/
func main() {

	// Define the needed queues
	userStoriesQueue = queuepackage.QueueFactory("userStoriesQueue", "userStories")
	textQueue = queuepackage.QueueFactory("textQueue", "text")

	// Begin listening to the user stories queue to send its content back to the gateway
	queuepackage.ReceiveFromQueueConc(userStoriesQueue)

	// Handle any requests sent to our server
	handleRequests()

}
