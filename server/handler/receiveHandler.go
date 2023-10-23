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

	// register a new channel to the broker
	messageChan := make(chan []byte)

	// signal to the broker that we have a new connection
	broker.NewClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits
	defer func() {
		broker.ClosingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		broker.ClosingClients <- messageChan
	}()

	w.Header().Set("Content-Disposition", "attachment; filename=testfilename.txt")

	// block waiting for messages broadcast on this connection's messageChan
	for {
		/* buffer := make([]byte, 1)
		buffer = <-messageChan */

		byteReader := bytes.NewReader(<-messageChan)
		io.Copy(w, byteReader)

		flusher.Flush()
	}
}
