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

const service = "http://localhost:3000"

func main() {
	var bar *pb.ProgressBar
	var file *os.File
	var fileInfo os.FileInfo
	var err error

	if file, err = os.Open("test.txt"); err != nil {
		log.Fatal(err)
	}

	if fileInfo, err = file.Stat(); err != nil {
		log.Fatal(err)
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

	resp, err := http.Post(fmt.Sprintf("%s/send", service), mpw.FormDataContentType(), r)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(ret))
}
