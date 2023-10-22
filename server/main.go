package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/prajwalad101/vincent/server/utils"
)

var jobs = make([]string, 10)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 100 // 100 MB

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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
	jobId := utils.GenerateJobId()
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

		// create a new file on the server to save the uploaded file
		filename := utils.GenerateFileName(part.FileName())
		uploadedFile, err := os.Create(filename)
		if err != nil {
			http.Error(w, "Unable to create file", http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		// copy the uploaded file to the new file on the server
		_, err = io.Copy(uploadedFile, part)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", uploadHandler)

	log.Println("The server is listening on port 3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
