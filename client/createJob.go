package main

import (
	"fmt"
	"io"
	"net/http"
)

func createJob(_ bool) error {
	resp, err := http.Get(fmt.Sprintf("%s/job", API_URL))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Job Id:", string(data))

	return nil
}
