package handler

import (
	"bytes"
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

	// register a new channel to the broker
	messageChan := make(chan []byte)

	client := Client{channel: messageChan, id: jobId}

	// signal to the broker that we have a new connection
	broker.NewClient <- client

	// Remove this client from the map of connected clients
	// when this handler exits
	defer func() {
		broker.ClosingClient <- client
	}()

	// Listen to connection close and un-register messageChan
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		broker.ClosingClient <- client
	}()

	w.Header().Set("Content-Disposition", "attachment; filename=testfilename.txt")

	// block waiting for messages broadcast on this connection's messageChan
	for {
		byteReader := bytes.NewReader(<-messageChan)
		io.Copy(w, byteReader)

		flusher.Flush()
	}
}
