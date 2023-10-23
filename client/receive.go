package main

import (
	"fmt"
	"log"

	"github.com/prajwalad101/vincent/client/util"
)

func receive(jobId string, downloadPath string) {
	url := fmt.Sprintf("%s/receive?jobId=%s", API_URL, jobId)
	err := util.DownloadFile(url, downloadPath)
	if err != nil {
		log.Fatal(err)
	}
}
