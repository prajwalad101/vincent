package main

import (
	"log"
	"net/http"

	"github.com/prajwalad101/vincent/server/handler"
)

var jobs = make([]string, 10)

func main() {
	mux := http.NewServeMux()

	broker := handler.NewBroker()

	mux.HandleFunc("/send", broker.SendHandler)
	mux.HandleFunc("/receive", broker.ReceiveHandler)
	mux.HandleFunc("/job", broker.JobHandler)

	log.Println("The server is listening on port 3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
