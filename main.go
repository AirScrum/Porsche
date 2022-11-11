package main

import (
	"fmt"
	"log"
	"net/http"
	"goserver/queuePackage"
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
	queuepackage.Start()
	handleRequests()
}
