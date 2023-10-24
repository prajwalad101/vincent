package main

import (
	"fmt"

	"github.com/prajwalad101/vincent/client/util"
)

func receive(jobId string, downloadPath string) error {
	if jobId == "" {
		return fmt.Errorf("Please provide a job id")
	}

	url := fmt.Sprintf("%s/receive?jobId=%s", API_URL, jobId)
	err := util.DownloadFile(url, downloadPath)
	if err != nil {
		return err
	}

	return nil
}
