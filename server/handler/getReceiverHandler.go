package handler

import (
	"net/http"
	"strconv"
)

// check if the receivers are ready
func (broker *Broker) GetReceiverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobId := r.URL.Query().Get("jobId")
	if jobId == "" {
		http.Error(w, "'jobId' is required", http.StatusBadRequest)
		return
	}

	clients := broker.Clients[jobId]
	if clients == nil {
		w.Write([]byte("0"))
	}
	numReceivers := len(clients)

	w.Write([]byte(strconv.Itoa(numReceivers)))
	return
}
