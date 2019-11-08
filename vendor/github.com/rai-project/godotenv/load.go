// Package godotenv is a go port of the ruby dotenv library (https://github.com/bkeepers/dotenv)
//
// Examples/readme can be found on the github page at https://github.com/joho/godotenv
//
// The TL;DR is that you make a .env file that looks something like
//
// 		SOME_ENV_VAR=somevalue
//
// and then in your go code you can call
//
// 		godotenv.Load()
//
// and all the env vars declared in .env will be avaiable through os.Getenv("SOME_ENV_VAR")
package godotenv

import (
	"os"
	"strings"
)

// From will read your env s and load them into ENV for this process.
//
// Call this function as close as possible to the start of your program (ideally in main)
//
// It's important to note that it WILL OVERRIDE an env variable that already exists
// - consider the .env file to set dev vars or sensible defaults

func From(data string) (err error) {
	envMap := make(map[string]string)
	lines := strings.Split(data, "\n")

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			key, value, err := parseLine(fullLine)

			if err == nil {
				envMap[key] = value
			}
		}
	}

	for key, value := range envMap {
		os.Setenv(key, value)
	}

	return
}
