package util

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
)

func DownloadFile(url, downloadPath string) (err error) {
	fmt.Println("Waiting for sender ...")

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
	if filename == "" {
		return fmt.Errorf("Filename is empty")
	}

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s%s", downloadPath, filename))
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Transfer complete.")

	return nil
}
