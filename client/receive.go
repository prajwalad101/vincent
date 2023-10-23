package main

import (
	"fmt"

	"github.com/prajwalad101/vincent/client/util"
)

func receive(jobId string, downloadPath string) {
	url := fmt.Sprintf("%s/receive/jobId=%s", API_URL, jobId)
	util.DownloadFile(url, downloadPath)
}
