package main

import (
	"fmt"
	"os"

	"github.com/mettbox/gotoggr/account"
)

func processArgs(args []string) {
	if len(args) == 2 && args[0] == "--token" {
		account.SaveToken(args[1])
	} else {
		fmt.Println("List your latest Toggl time entries")
		fmt.Println("usage: gotoggr [--token <toggl-api-token>]")
		fmt.Println("  ------- Options -------")
		fmt.Println("  --token <toggl-api-token>\t Stores your Toggl API token in your Home Folder ")
	}
	os.Exit(0)
}

func main() {
	// First value in this slice is the path to the program and os.Args[1:] holds the arguments to the program
	if args := os.Args[1:]; len(args) > 0 {
		processArgs(args)
	}

	account.LatestEntries()
}
