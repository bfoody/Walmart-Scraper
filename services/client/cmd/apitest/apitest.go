package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bfoody/Walmart-Scraper/services/client/internal/api"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/api/walmart"
)

// exit prints a message and exits the process.
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		exit("usage: apitest <type> <args...>")
	}

	http := api.NewHTTPClient()
	client := walmart.NewClient(http)

	action := args[0]
	switch action {
	case "get-product":
		if len(args) < 3 {
			exit("usage: apitest get-product <# repetitions> <interval ms>")
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			exit("err: # repetitions not a number")
		}

		interval, err := strconv.Atoi(args[2])
		if err != nil {
			exit("err: interval not a number")
		}

		for i := 0; i < num; i++ {
			item, err := client.GetItemDetails("onn-32-Class-HD-720P-Roku-Smart-LED-TV-100012589", "314022535")
			if err != nil {
				exit(err.Error())
			}

			fmt.Printf("success: %+v\n", item)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	case "get-product-once":
		item, err := client.GetItemDetails("onn-32-Class-HD-720P-Roku-Smart-LED-TV-100012589", "314022535")
		if err != nil {
			exit(err.Error())
		}

		exit(fmt.Sprintf("success: %+v", item))
	default:
		exit("usage: apitest <type> <args...>")
	}
}
