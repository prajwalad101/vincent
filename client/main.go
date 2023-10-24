package main

import (
	"fmt"

	"github.com/prajwalad101/vincent/client/util"
)

const API_URL = "http://localhost:3000"

func main() {
	args, err := util.GetArgs()
	if err != nil {
		fmt.Println("No parameters provided. Please specify send or receive")
		return
	}

	switch args.Command {
	case "send":
		err = send(args.JobId, args.Filepath)
	case "receive":
		err = receive(args.JobId, args.DownloadPath)
	case "create":
		err = createJob(args.SaveOnClipboard)
	default:
		fmt.Println("No parameters provided. Please specify send or receive")
		return
	}

	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
