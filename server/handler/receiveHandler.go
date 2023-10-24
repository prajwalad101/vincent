package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func (broker *Broker) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	jobId := r.URL.Query().Get("jobId")
	if jobId == "" {
		http.Error(w, "Job Id is required", http.StatusBadRequest)
		return
	}

	messageChan := make(chan []byte)
	closeChan := make(chan bool)

	client := Client{
		channel: ClientChan{
			messageChan: messageChan,
			closeChan:   closeChan,
		},
		id: jobId,
	}

	// signal to the broker that we have a new connection
	broker.NewClient <- client

	w.Header().Set("Content-Disposition", "attachment; filename=testfilename.txt")

	for {
		select {
		case <-closeChan:
			fmt.Fprint(w, "Transfer complete.")
			return
		case message := <-messageChan:
			byteReader := bytes.NewReader(message)
			io.Copy(w, byteReader)
			flusher.Flush()
		default:
		}
	}
}
