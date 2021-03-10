package main

import (
	"fmt"

	"github.com/bfoody/Walmart-Scraper/logging"
)

func main() {
	// Initialize logging.
	log, err := logging.Configure()
	if err != nil {
		fmt.Println("Error initializing logging: ", err)
	}

	log.Info("Hello world")
}
