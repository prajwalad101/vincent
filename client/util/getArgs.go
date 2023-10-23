package util

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	TransferType string
	Filepath     string
	JobId        string
	DownloadPath string
}

// Parse command line subcommands and its flags
func GetArgs() (Args, error) {
	args := Args{}

	// receive subcommand
	receiveCmd := flag.NewFlagSet("receive", flag.ExitOnError)
	jobId := receiveCmd.String(
		"id",
		"",
		"The job id of a pending send.",
	)
	downloadPath := receiveCmd.String(
		"path",
		"",
		"The path to download the received file",
	)

	// send subcommand
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	filepath := sendCmd.String(
		"file",
		"",
		"The path of the file to upload.",
	)

	// if no command line args provided
	if len(os.Args) < 2 {
		return args, fmt.Errorf("Expected 'send' or 'receive'")
	}

	args.TransferType = os.Args[1]

	switch args.TransferType {
	case "send":
		sendCmd.Parse(os.Args[2:])
		args.Filepath = *filepath
	case "receive":
		receiveCmd.Parse(os.Args[2:])
		args.JobId = *jobId
		args.DownloadPath = *downloadPath
	default:
		return args, fmt.Errorf("Expected 'send' or 'receive'")
	}

	return args, nil
}
