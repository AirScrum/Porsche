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
	Textid string `json:"textid"`
	Text   string `json:"text"`
}

type Request struct {
	Textid string `json:"textid"`
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

	msg := Message{}

	// Get corressponding text from the textid from database
	msg.Textid = request.Textid
	msg.Text = "Testing"

	queuepackage.SendToTextQueue(queues, msg.Textid, msg.Text)
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
