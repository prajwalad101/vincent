package handler

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/prajwalad101/vincent/server/util"
)

func (broker *Broker) SendHandler(w http.ResponseWriter, r *http.Request) {
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

	// store a unique jobId and send it to the client
	jobId := util.GenerateJobId()
	fmt.Fprintf(w, "Job id: %d ", jobId)

	// Create a multipart reader with the request body and boundary
	partReader := multipart.NewReader(r.Body, params["boundary"])

	for {
		// read the next part from the multipart stream
		part, err := partReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Unable to read multipart", http.StatusInternalServerError)
			return
		}

		// copy the uploaded file to the new file on the server
		for {
			// send data to notifier channel every 2 seconds
			buffer := make([]byte, 10000)
			_, err = part.Read(buffer)
			if err != nil {
				break
			}
			broker.Notifier <- buffer
		}
	}
	return
}
