package handler

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

func (broker *Broker) SendHandler(w http.ResponseWriter, r *http.Request) {
	// get job id
	jobId := r.URL.Query().Get("jobId")
	if jobId == "" {
		http.Error(w, "'jobId' is required", http.StatusBadRequest)
		return
	}

	receivers := r.URL.Query().Get("receivers")
	if receivers == "" {
		http.Error(w, "'receivers' is required", http.StatusBadRequest)
		return
	}

	numReceivers, err := strconv.Atoi(receivers)
	if err != nil {
		http.Error(w, "Invalid 'receivers' value", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, "Unable to parse Content-Type header", http.StatusBadRequest)
		return
	}

	boundary := params["boundary"]
	if boundary == "" {
		http.Error(w, "Malformed multi-part boundary", http.StatusBadRequest)
		return
	}

	// Listen to connection close and un-register messageChan
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		fmt.Println("Connection closed.")
		broker.ClosingClient <- jobId
	}()

	// close all receivers when sender closes
	defer func() {
		broker.ClosingClient <- jobId
	}()

	for {
		// check for any new clients every 100 ms
		time.Sleep(time.Second * 1)

		// check if there are any registered clients
		clients := broker.Clients[jobId]

		// check if the number if receivers is ready
		if len(clients) >= numReceivers {
			break
		}
	}

	// Create a multipart reader with the request body and boundary
	partReader := multipart.NewReader(r.Body, params["boundary"])

	for {
		// read the next part from the multipart stream
		part, err := partReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Unable to read multipart", http.StatusInternalServerError)
			return
		}

		// copy the uploaded file to the new file on the server
		for {
			// send data to notifier channel every 2 seconds
			buffer := make([]byte, 1)
			_, err := part.Read(buffer)
			if err != nil {
				break
			}
			event := Event{data: buffer, jobId: jobId}
			broker.EventNotifier <- event
		}
	}

	fmt.Fprint(w, "Transfer complete.")
	return
}
