package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/wassafshahzad/go-load-tester/internal"
)

func main() {

	arguments := os.Args

	if len(arguments) < 2 {
		panic("Config file required")
	}

	location := arguments[1]
	api, err := internal.ReadConfig(location)

	if err != nil {
		panic("Error reading json file.")
	}

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	clientWithTimeOut := &http.Client{
		Timeout: 5 * time.Second,
	}

	count := 1
	cutoff := int(api.Requests / api.Batches)
	batchCounter := 1

	for index := range api.Urls {
		fmt.Printf("Currently requesting url %v\n", api.Urls[index].Path)
		for i := 0; i < api.Requests; i++ {

			if i == cutoff*batchCounter {
				fmt.Printf("Registered goroutines in Batch %v \n", batchCounter)
				batchCounter += 1
				waitGroup.Wait()
			}
			waitGroup.Add(1)
			go func() {
				defer func() {
					mutex.Lock()
					fmt.Printf("Requests completed %v/%v \r", count, api.Requests)
					count += 1
					mutex.Unlock()
					waitGroup.Done()
				}()
				internal.CallUrl(&api.Urls[index], clientWithTimeOut)
			}()
		}
	}
	waitGroup.Wait()
	d, f := api.Urls[0].GetRequestsResult()
	fmt.Printf("\nTotal Requests Dropped %v Failed %v", d, f)
}
