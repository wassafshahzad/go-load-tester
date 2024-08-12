package main

import (
	"fmt"
	"sync"

	"github.com/wassafshahzad/go-load-tester/internals"
)

func main() {

	api := internals.ReadConfig("requests.json")

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	count := 1
	cutoff := int(api.Requests / api.Batches)
	batchCounter := 1

	for index := range api.Urls {
		fmt.Printf("Currently requesting url %v\n", api.Urls[index].Path)
		for i := 0; i < api.Requests; i++ {

			if i == cutoff*batchCounter {
				fmt.Printf("dRegistered goroutines in Batch %v \n", batchCounter)
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
				internals.CallUrl(&api.Urls[index])
			}()
		}
	}
	waitGroup.Wait()
	d, f := api.Urls[0].GetRequestsResult()
	fmt.Printf("\nTotal Requests Dropped %v Failed %v", d, f)
}
