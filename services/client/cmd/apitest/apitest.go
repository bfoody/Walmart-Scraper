package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
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
	case "get-product-parallel":
		if len(args) < 3 {
			exit("usage: apitest get-product-parallel <proxy list txt> <# threads>")
		}

		proxyList := args[1]

		num, err := strconv.Atoi(args[2])
		if err != nil {
			exit("err: # threads not a number")
		}

		file, err := ioutil.ReadFile(proxyList)
		if err != nil {
			exit("err: proxy list file does not exist")
		}

		proxies := strings.Split(string(file), "\n")

		wg := &sync.WaitGroup{}

		for i := 0; i < num; i++ {
			wg.Add(1)

			go func(i int) {
				http := api.NewHTTPClient()
				http.SetProxy(proxies[i%len(proxies)])
				client := walmart.NewClient(http)

				item, err := client.GetItemDetails("onn-32-Class-HD-720P-Roku-Smart-LED-TV-100012589", "314022535")
				if err != nil {
					exit(err.Error())
				}

				fmt.Printf("success on thread %d: %+v\n", i, item)

				wg.Done()
			}(i)
		}

		wg.Wait()
	default:
		exit("usage: apitest <type> <args...>")
	}
}
