## Simple Load Testing application Built using Go.

A simple load testing application. Currently only `Get` request is supported.  

### Running the Application

The application requires only one argument. Which is the path of the config file.
The config file consists of the following keys. 
- `requests` : (int) The total number of requests to send to each endpoint.
- `batches`  : (int) The number of batches in which goroutines will be divided too.
- `urls`     : ([]map[string, string]{path, method}) A list of map with the keys `path` and `method`, The `path` defines the endpoint which to hit and the `method` describes the http method to use.

Example of config file structure.
```
{ 
    "requests" : 1000000,
    "batches":5,
    "urls": [{ "path": "http://localhost:8080", "method": "GET" }]
}

```

