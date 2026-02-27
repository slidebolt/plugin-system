package main

import (
	"log"

	runner "github.com/slidebolt/sdk-runner"
)

func main() {
	if err := runner.NewRunner(NewSystemPlugin()).Run(); err != nil {
		log.Fatal(err)
	}
}
