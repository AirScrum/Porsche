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
		Text string `json:"text"`
}
type Bird struct {
	Species string
	Description string
  }

func homepage(w http.ResponseWriter, r *http.Request) {
	/*
		Read the request body and parse it from JSON to Message Struct
	*/
	defer r.Body.Close()
    buf, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
	fmt.Println(buf)
	msg:= Message{} 
	err = json.Unmarshal(buf,&msg)
	fmt.Println(msg.Textid)
	fmt.Println(msg.Text)
	/*
		TODO Read the textID from the database
	*/
	/*
		TODO Send the text fetched from the database to the TextQueue
	*/
}

func handleRequests() {
	http.HandleFunc("/main", homepage)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func main() {
	fmt.Println("Server Started")
	queues := new(queuepackage.Queue)
	queues = queuepackage.QueueInit()
	queuepackage.SendToTextQueue(queues, "Hi")
	//queuepackage.Start()
	//queuepackage.SendToTextQueue("Test")
	handleRequests()
}
