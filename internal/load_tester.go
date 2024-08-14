package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HttpMethods string

type Urls struct {
	Requests int        `json:"requests"`
	Batches  int        `json:"batches"`
	Urls     []entities `json:"urls"`
}

type entities struct {
	Path       string `json:"path"`
	HttpMethod string `json:"method"`
	results    requestResult
}

func (entity entities) GetRequestsResult() (int, int) {
	return entity.results.dropCount, entity.results.failCount
}

type requestResult struct {
	dropCount int
	failCount int
}

// This can be a method
// No need to create a new request every time pass the request as an argument
func CallUrl(entity *entities, client *http.Client) {
	req, err := http.NewRequest(string(entity.HttpMethod), entity.Path, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		entity.results.dropCount += 1
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		entity.results.failCount += 1
	} else {
		fmt.Printf("%v\n", res.StatusCode)
	}
}

func ReadConfig(name string) (*Urls, error) {
	jsonFile, err := os.Open(name)
	if err != nil {
		fmt.Println("File not found.")
		return nil, err
	}

	var apis Urls

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	_ = json.Unmarshal(byteValue, &apis)

	if apis.Requests == 0 {
		apis.Requests = 100
	}

	return &apis, nil
}
