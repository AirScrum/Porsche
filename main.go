package main

import (
	"fmt"
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
	start()
}
