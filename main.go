package main

import (
	"fmt"
	"sync"

	"github.com/wassafshahzad/go-load-tester/internals"
)

func main() {

	api := internals.ReadConfig("requests.json")

	var waitGroup sync.WaitGroup

	for _, val := range api.Urls {
		waitGroup.Add(api.Requests)
		for i := 0; i < api.Requests; i++ {
			go func() {
				defer waitGroup.Done()
				fmt.Printf("Sending request at %v \n", val.Path)
				internals.CallUrl(&val)
			}()

		}
	}
	waitGroup.Wait()
	fmt.Printf("%v", api)
}
