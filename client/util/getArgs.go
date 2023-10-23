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
}

// Parse command line subcommands and its flags
func GetArgs() (Args, error) {
	args := Args{}

	// send subcommand
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	args.Filepath = *sendCmd.String(
		"file",
		"",
		"The path of the file to upload. (Only required when sending)",
	)

	// receive subcommand
	receiveCmd := flag.NewFlagSet("receive", flag.ExitOnError)
	args.JobId = *receiveCmd.String(
		"id",
		"",
		"The job id of a pending send. (Only required when receiving)",
	)

	// if no command line args provided
	if len(os.Args) < 2 {
		return args, fmt.Errorf("Expected 'send' or 'receive'")
	}

	args.TransferType = os.Args[1]

	switch args.TransferType {
	case "send":
		sendCmd.Parse(os.Args[2:])
	case "receive":
		sendCmd.Parse(os.Args[2:])
	default:
		return args, fmt.Errorf("Expected 'send' or 'receive'")
	}

	return args, nil
}
