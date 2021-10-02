package testfuncs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

// --------  TEST SERVER  --------

// SetupTestServer will setup a test server and respond with the value supplied as `jsonResponse`.
func SetupTestServer(jsonResponse []byte) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(jsonResponse))
	}))

	return server
}

// GetTestJsonResponse will return the contents of `file` after reading it from the `server_responses` directory.
// This directory should include all test responses we need in the application.
func GetTestJsonResponse(file string) []byte {
	path := "./test/server_responses/" + file
	absPath, _ := filepath.Abs(path)
	testJsonResponse, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Println("There was an error attempting to pull in json server response: ", file)
		fmt.Println(err.Error())
	}

	return testJsonResponse
}

// --------  TEST CLIENT  --------

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// GetNewTestClient can be called to return a http.Client to use in tests. Just pass in the response you want back.
func GetNewTestClient(testJsonResponse []byte) *http.Client {
	return NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(string(testJsonResponse))),
			Header:     make(http.Header),
		}
	})
}
