package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func checkReceivers(jobId string, receivers int) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/receiver?jobId=%s", API_URL, jobId))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	numReceivers, err := strconv.Atoi(string(data))
	if err != nil {
		return false, err
	}

	if receivers == numReceivers {
		return true, nil
	} else {
		return false, nil
	}
}
