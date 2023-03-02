package main

/*
Import important libraries
*/
import (
	"context"
	"encoding/json"
	"fmt"
	dbpackage "goserver/dbPackage"
	queuepackage "goserver/queuePackage"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Define the queues we need
*/
var userStoriesQueue *queuepackage.IQueue
var textQueue *queuepackage.IQueue

var mongoClient *mongo.Client
var mongoContext context.Context
var mongoCancel context.CancelFunc
var mongoError error

/*
This function is called when the is a request in our server, which contains the textid and needed to be sent to the text queue
*/
func homepage(w http.ResponseWriter, r *http.Request) {

	fmt.Println(mongoClient, mongoContext, mongoCancel)

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

	data := queuepackage.Message{}
	data.TextID = request.TextID
	data.Text = "Testing"
	data.UserID = request.UserID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)

	/*
		TODO Read the textID from the database
	*/

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70

	objID, err := primitive.ObjectIDFromHex("64009b4eff8e68d7cd8b7955")
	if err != nil {
		panic(err)
	}
	filter = bson.D{{
		Key:   "_id",
		Value: objID,
	}}

	//  option remove id field from all documents
	//option = bson.D{{"_id", 0}}
	option = bson.D{}

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	cursor, err := dbpackage.Query(mongoClient, mongoContext, "test",
		"text", filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	var results []bson.D

	// to get bson object  from cursor,
	// returns error if any.
	if err := cursor.All(mongoContext, &results); err != nil {

		// handle the error
		panic(err)
	}

	// printing the result of query.
	fmt.Println("Query Result")
	for _, doc := range results {
		fmt.Println(doc)
	}

	// Release resource when the main
	// function is returned.
	defer dbpackage.Close(mongoClient, mongoContext, mongoCancel)

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

	//Load the .env file
	err := godotenv.Load(".env")
	// Print error message when failing
	if err != nil {
		panic(err)
	}

	// Connect to mongoDB
	mongoClient, mongoContext, mongoCancel, mongoError = dbpackage.Connect(os.Getenv("MONGO_DB_URI"))
	if err != nil {
		panic(err)
	}

	fmt.Println(mongoClient, mongoContext, mongoCancel)

	// Define the needed queues
	userStoriesQueue = queuepackage.QueueFactory("userStoriesQueue", "userStories")
	textQueue = queuepackage.QueueFactory("textQueue", "text")

	// Begin listening to the user stories queue to send its content back to the gateway
	queuepackage.ReceiveFromQueueConc(userStoriesQueue)

	// Handle any requests sent to our server
	handleRequests()

}
