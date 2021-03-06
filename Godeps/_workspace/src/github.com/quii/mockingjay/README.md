# mockingjay

[![GoDoc](https://godoc.org/github.com/quii/mockingjay?status.svg)](https://godoc.org/github.com/quii/mockingjay)

Create a server from configuration. Useful for integration tests, hopefully.

## Example

```go
package main

import (
	"github.com/quii/mockingjay"
	"log"
	"net/http"
)

func main() {
	const testYAML = `
---
 - name: Test endpoint
   request:
     uri: /hello
     method: GET
   response:
     code: 200
     body: hello, world
     headers:
       content-type: text/plain

 - name: Test endpoint 2
   request:
     uri: /world
     method: DELETE
   response:
     code: 200
     body: hello, world

 - name: Failing endpoint
   request:
     uri: /card
     method: POST
     body: Greetings
   response:
     code: 500
     body: Oh bugger
 `
	endpoints, err := mockingjay.NewFakeEndpoints([]byte(testYAML))

	if err != nil {
		log.Fatal(err)
	}

	server := mockingjay.NewServer(endpoints)

	// Mount it just like any other server
	http.Handle("/", server)
	http.ListenAndServe(":9090", nil)
}
```
## Building

- Requires Go 1.3+
- godeps

## Todo

- Although it supports request/response headers, it only supports one value per header (http allows you to set multiple values)
- Tests for stuff inside request.go
