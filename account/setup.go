package account

import (
	"fmt"
	"os"

	"github.com/mettbox/gotoggr/util"
)

var setupFile = util.UserHomeDir() + "/.gotoggr"

func SaveToken(token string) {
	err := os.WriteFile(setupFile, []byte(token), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Store Toggle API Token in: %s \n", setupFile)
}

func GetToken() string {
	fileExists, err := util.FileExists(setupFile)
	if err != nil {
		panic(err)
	}

	if !fileExists {
		fmt.Println("Error: Setup file does not exists. Please run gotoggr --token <your-toggle-api-token>")
		os.Exit(0)
	}

	token, err := os.ReadFile(setupFile)
	if err != nil {
		panic(err)
	}

	return string(token)
}
