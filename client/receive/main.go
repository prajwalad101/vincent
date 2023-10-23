package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
)

const service = "http://localhost:3000"

func main() {
	downloadFile(fmt.Sprintf("%s/receive", service))
}

func downloadFile(url string) (err error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(contentDisposition)

	filename := params["filename"]

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
