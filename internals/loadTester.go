package internals

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type HttpMethods string

const (
	POST  HttpMethods = "POST"
	GET   HttpMethods = "GET"
	PUT   HttpMethods = "PUT"
	PATCH HttpMethods = "PATCH"
)

var clientWithTimeOut http.Client

type Urls struct {
	Urls     []entities    `json:"urls"`
	Requests int           `json:"requests"`
	Timeout  time.Duration `json:"timeout"`
}

type entities struct {
	Path       string      `json:"path"`
	HttpMethod HttpMethods `json:"method"`
	results    requestResult
}

func (entity entities) GetRequestsResult() (int, int) {
	return entity.results.dropCount, entity.results.failCount
}

type requestResult struct {
	dropCount int
	failCount int
}

func init() {
	clientWithTimeOut = http.Client{
		Timeout: 5 * time.Second,
	}
}

func CallUrl(entity *entities) {
	req, err := http.NewRequest(string(entity.HttpMethod), entity.Path, nil)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	res, err := clientWithTimeOut.Do(req)

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

func ReadConfig(name string) *Urls {
	jsonFile, err := os.Open(name)
	var apis Urls

	if err != nil {
		fmt.Println("File not found.")
	}

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	_ = json.Unmarshal(byteValue, &apis)

	if apis.Requests == 0 {
		apis.Requests = 100
	}

	if apis.Timeout != 0 {
		clientWithTimeOut.Timeout = apis.Timeout * time.Second
	}

	return &apis
}
