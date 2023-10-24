package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func send(jobId, filepath string) error {
	var bar *pb.ProgressBar
	var file *os.File
	var fileInfo os.FileInfo
	var err error

	if jobId == "" {
		return fmt.Errorf("Please provide a job id")
	}

	if filepath == "" {
		return fmt.Errorf("Please provide a filepath")
	}

	if file, err = os.Open(filepath); err != nil {
		return err
	}

	if fileInfo, err = file.Stat(); err != nil {
		return err
	}

	bar = pb.New64(fileInfo.Size()).SetRefreshRate(time.Millisecond * 10)
	bar.Set(pb.Bytes, true)
	bar.SetRefreshRate(time.Millisecond * 10)

	bar.Start()

	r, w := io.Pipe()

	mpw := multipart.NewWriter(w)

	go func() {
		var part io.Writer
		defer w.Close()
		defer file.Close()

		if part, err = mpw.CreateFormFile("file", file.Name()); err != nil {
			log.Fatal(err)
		}

		reader := bar.NewProxyReader(file)

		if _, err = io.Copy(part, reader); err != nil {
			log.Fatal(err)
		}

		if err = mpw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	resp, err := http.Post(
		fmt.Sprintf("%s/send?jobId=%s", API_URL, jobId),
		mpw.FormDataContentType(),
		r,
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(ret))

	return nil
}
