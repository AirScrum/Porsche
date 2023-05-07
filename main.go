package main

/*
Import important libraries
*/
import (
	"encoding/json"
	"fmt"
	dbpackage "goserver/dbPackage"
	"goserver/models"
	queuepackage "goserver/queuePackage"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
Define the queues we need
*/
var userStoriesQueue *queuepackage.IQueue
var textQueue *queuepackage.IQueue

/*
This function is called when the is a request in our server, which contains the textID needed to be sent to the text queue
The database is queried with the textID received, and construct an object contains the textID, userID, and text needed
to be converted. Then the object is converted to array of bytes and sent to the textQueue.
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
	request := models.Request{}
	err = json.Unmarshal(buf, &request)
	if err != nil {
		panic(err)
	}

	// Define Message that will be sent to the text queue
	msg := models.Message{}

	msg, err = dbpackage.GetMessageFromTextId(request.TextID)

	if err != nil {
		mongoClient, mongoContext, mongoCancel, mongoError := dbpackage.Connect(os.Getenv("MONGO_DB_URI"))
		if mongoError != nil {
			panic(mongoError)
		}
		// Release resource when the main
		// function is returned.
		defer dbpackage.Close(mongoClient, mongoContext, mongoCancel)
		msg, err = dbpackage.GetMessageFromTextId(request.TextID)
	}

	fmt.Println(msg)
	// Send the message text queue
	queuepackage.SendToQueue(textQueue, msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)

}

/*
This function is called when a request is sent to our server, and sent the request to homepage function to handle it
*/
func handleRequests() {
	fmt.Println("Server Started")
	http.HandleFunc("/main", homepage)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

/*
This is our main function
*/
func main() {

	// Connect to mongoDB
	mongoClient, mongoContext, mongoCancel, mongoError := dbpackage.Connect(os.Getenv("MONGO_DB_URI"))
	if mongoError != nil {
		panic(mongoError)
	}

	// Release resource when the main
	// function is returned.
	defer dbpackage.Close(mongoClient, mongoContext, mongoCancel)

	// Define the needed queues
	userStoriesQueue = queuepackage.QueueFactory("userStoriesQueue", "userStories")
	textQueue = queuepackage.QueueFactory("textQueue", "text")

	// Begin listening to the user stories queue to send its content back to the gateway
	queuepackage.ReceiveFromQueueConc(userStoriesQueue)

	// Handle any requests sent to our server
	handleRequests()

}
