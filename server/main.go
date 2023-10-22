package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

/* type FileServer struct{}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf.Bytes())
		fmt.Printf("received %d bytes over the network\n", n)
	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.Copy(conn, bytes.NewReader(file))
	if err != nil {
		return err
	}

	fmt.Printf("written %d bytes over the network\n", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		sendFile(200000)
	}()
	server := &FileServer{}
	server.start()
} */

const MAX_UPLOAD_SIZE = 1024 * 1024 * 100 // 100 MB

/* func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(
			w,
			"The uploaded file is too big. Please choose an file that's less than 100MB in size",
			http.StatusBadRequest,
		)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	n, err := io.CopyN(buf, file, 400)
} */

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
		filename := generateFileName(part.FileName())
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

		fmt.Fprintln(w, "File uploaded successfully!")
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

func generateFileName(originalFilename string) string {
	currentTime := time.Now()

	parsedTime := currentTime.Format(time.RFC3339)

	filename := fmt.Sprintf(
		"%s(%s)",
		originalFilename,
		parsedTime,
	)
	return filename
}
