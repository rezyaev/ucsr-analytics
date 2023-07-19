package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	command := os.Args[1]

	switch command {
	case "download":
		if error := FetchDemos(); error != nil {
			log.Panic(error)
		}
	case "analyze":
		if error := AnalyzeDemos(); error != nil {
			log.Panic(error)
		}
	default:
		fmt.Printf("Wrong command line argument: %v", command)
	}
}
