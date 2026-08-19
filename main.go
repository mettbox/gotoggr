package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mettbox/gotoggr/account"
)

func printUsage() {
	fmt.Println("List your latest Toggl time entries")
	fmt.Println("usage: gotoggr [--set-token <toggl-api-token>] [--show-token]")
	fmt.Println("  ------- Options -------")
	fmt.Println("  --set-token <toggl-api-token>\t Stores your Toggl API token in your Home Folder ")
	fmt.Println("  --show-token\t\t\t Prints the stored Toggl API token")
}

// validateToken returns the token to store, rejecting arguments that hold no token.
func validateToken(argument string) (string, error) {
	token := strings.TrimSpace(argument)
	if token == "" {
		return "", errors.New("--set-token requires a non-empty token")
	}

	return token, nil
}

func setToken(argument string) {
	token, err := validateToken(argument)
	if err != nil {
		fmt.Printf("Error: %v. Generate one at https://track.toggl.com/profile\n", err)
		os.Exit(1)
	}

	account.SaveToken(token)
}

func processArgs(args []string) {
	switch {
	case len(args) == 2 && args[0] == "--set-token":
		setToken(args[1])
	case len(args) == 1 && args[0] == "--show-token":
		fmt.Println(account.GetToken())
	default:
		printUsage()
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
