package main

import (
	"encoding/json"
	"fmt"
	queuepackage "goserver/queuePackage"
	"io/ioutil"
	"log"
	"net/http"
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

var queues = new(queuepackage.Queue)

func homepage(w http.ResponseWriter, r *http.Request) {
	/*
		Read the request body and parse it from JSON to Message Struct
	*/
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	request := Request{}
	err = json.Unmarshal(buf, &request)
	if err != nil {
		panic(err)
	}

	fmt.Println(request.textID)
	fmt.Println(request.userID)
	msg := Message{}

	// Get corresponding text from the textID from database
	msg.textID = request.textID
	msg.text = "Testing"
	msg.userID = request.userID

	queuepackage.SendToTextQueue(queues, msg.textID, msg.text, msg.userID)
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
	queues = queuepackage.QueueInit()
	queuepackage.ReceiveFromUserStoriesQueue(queues)
	handleRequests()

}
