package main

import (
	"fmt"
	"github.com/quii/mockingjay"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestItChecksAValidEndpointsJSON(t *testing.T) {
	body := `{"foo":"bar"}`
	realServer := makeFakeDownstreamServer(body, noSleep)
	checker, _ := makeChecker(testYAML(body))

	if !checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be correct")
	}
}

func TestItFlagsDifferentJSONToBeIncompatible(t *testing.T) {
	serverResponseBody := `{"foo": "bar"}`
	fakeResponseBody := `{"baz": "boo"}`

	realServer := makeFakeDownstreamServer(serverResponseBody, noSleep)
	checker, _ := makeChecker(testYAML(fakeResponseBody))

	if checker.CheckCompatability(realServer.URL) {
		t.Error("Checker should've found this endpoint to be incorrect")
	}
}

func TestItIsIncompatibleWhenRealServerIsntReachable(t *testing.T) {
	yaml := testYAML("doesnt matter")
	checker, _ := makeChecker(yaml)

	if checker.CheckCompatability("http://localhost:12344") {
		t.Error("Checker shouldve found this to be an error as the real server isnt reachable")
	}
}

func TestItHandlesBadURLsInConfig(t *testing.T) {
	yaml := fmt.Sprintf(yamlFormat, "not a real url", "foobar")
	fakeEndPoints, _ := mockingjay.NewFakeEndpoints([]byte(yaml))
	checker := NewCompatabilityChecker(fakeEndPoints)

	if checker.CheckCompatability("also not a real url") {
		t.Error("Checker should've found that the URL in the YAML cannot be made into a request")
	}
}

const noSleep = 1

const defaultRequestURI = "/hello"

const yamlFormat = `
---
 - name: Test endpoint
   request:
     uri: %s
     method: GET
   response:
     code: 200
     body: '%s'
`

func makeFakeDownstreamServer(responseBody string, sleepTime time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepTime * time.Millisecond)
		if r.URL.RequestURI() == defaultRequestURI {
			fmt.Fprint(w, responseBody)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
}

func makeChecker(responseBody string) (*CompatabilityChecker, error) {
	fakeEndPoints, err := mockingjay.NewFakeEndpoints([]byte(responseBody))

	if err != nil {
		return nil, err
	}
	return NewCompatabilityChecker(fakeEndPoints), nil
}

func testYAML(responseBody string) string {
	return fmt.Sprintf(yamlFormat, defaultRequestURI, responseBody)
}
