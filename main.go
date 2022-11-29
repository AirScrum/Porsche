package main

import (
	"encoding/json"
	"fmt"
	queuepackage "goserver/queuePackage"
	"io/ioutil"
	"log"
	"net/http"
)

var queues = new(queuepackage.Queue)
var userStoriesQueue *queuepackage.IQueue
var textQueue *queuepackage.IQueue

func homepage(w http.ResponseWriter, r *http.Request) {
	/*
		Read the request body and parse it from JSON to Message Struct
	*/
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	request := queuepackage.Request{}
	err = json.Unmarshal(buf, &request)
	if err != nil {
		panic(err)
	}

	fmt.Println(request.TextID)
	fmt.Println(request.UserID)
	msg := queuepackage.Message{}

	// Get corresponding text from the TextID from database
	msg.TextID = request.TextID
	msg.Text = "Testing"
	msg.UserID = request.UserID

	//queuepackage.SendToTextQueue(queues, msg.TextID, msg.Text, msg.UserID)

	queuepackage.SendToQueue(textQueue, msg)
	/*
		TODO Read the textID from the database
	*/
	/*
		TODO Send the text fetched from the database to the TextQueue
	*/
}

func handleRequests() {
	fmt.Println("Server Started")
	http.HandleFunc("/main", homepage)
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func main() {
	//queues = queuepackage.QueueInit()
	//queuepackage.ReceiveFromUserStoriesQueue(queues)

	userStoriesQueue = queuepackage.QueueFactory("userStoriesQueue", "userStories")
	textQueue = queuepackage.QueueFactory("textQueue", "text")

	queuepackage.ReceiveFromQueueConc(userStoriesQueue)
	handleRequests()

}
