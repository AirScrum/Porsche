package main

import (
	"fmt"
	queuepackage "goserver/queuePackage"
	"log"
	"net/http"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint")
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
