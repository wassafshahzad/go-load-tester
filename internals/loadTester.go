package internals

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HttpMethods string

const (
	POST  HttpMethods = "POST"
	GET   HttpMethods = "GET"
	PUT   HttpMethods = "PUT"
	PATCH HttpMethods = "PATCH"
)

type Urls struct {
	Urls     []entities `json:"urls"`
	Requests int        `json:"requests"`
}

type entities struct {
	Path       string      `json:"path"`
	HttpMethod HttpMethods `json:"method"`
	results    requestResult
}

type requestResult struct {
	dropCount int
	failCount int
}

func CallUrl(entity *entities) {
	req, _ := http.NewRequest(string(entity.HttpMethod), entity.Path, nil)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		entity.results.dropCount += 1
		return
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		entity.results.failCount += 1
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
	json.Unmarshal(byteValue, &apis)

	if apis.Requests == 0 {
		apis.Requests = 100
	}
	return &apis
}
