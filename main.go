package main

import (
	"fmt"
	"sync"

	"github.com/wassafshahzad/go-load-tester/internals"
)

func main() {

	api := internals.ReadConfig("requests.json")

	var waitGroup sync.WaitGroup
	c := 1
	for index := range api.Urls {
		for i := 0; i < api.Requests; i++ {
			waitGroup.Add(1)
			go func() {
				defer func() {
					waitGroup.Done()
					fmt.Printf("Done please %v \n", i)
				}()
				internals.CallUrl(&api.Urls[index])
				fmt.Printf("Calling thread %v \n", i)
			}()
		}
	}
	waitGroup.Wait()
	d, f := api.Urls[0].GetRequestsResult()
	fmt.Printf("Total Requests Dropped %v Failed %v", d, f)
}
