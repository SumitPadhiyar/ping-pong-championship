package client

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// DoRequest executes an HTTP request to the given requestURL using the given
// requestMethod and requestData.
func DoRequest(requestMethod, requestURL,
	requestData string) (*http.Response, error) {
	// These will hold the return value.
	var res *http.Response
	var err error

	// Convert method to uppercase for easier checking.
	upperRequestMethod := strings.ToUpper(requestMethod)
	switch upperRequestMethod {
	case "DELETE", "PATCH", "PUT":
		// All these methods have no shortcuts in Go's HTTP library, so
		// we have to do them manually.
		if len(requestData) == 0 && upperRequestMethod != "DELETE" {
			// All methods (except for DELETE) require data.
			return nil, fmt.Errorf(
				"--request_data must be provided")
		}

		// NewRequest requires a Reader, so we create a byte buffer
		// for our string data.
		contentBuffer := bytes.NewBufferString(requestData)

		req, err := http.NewRequest(upperRequestMethod, requestURL,
			contentBuffer)
		if err != nil {
			// Failed creating HTTP request.
			return nil, fmt.Errorf(
				"error creating new HTTP request")
		}

		// Use the default HTTP client to execute the request.
		res, err = http.DefaultClient.Do(req)
	case "GET":
		// Use the HTTP library Get() method.
		res, err = http.Get(requestURL)
	case "POST":
		// Use the HTTP library Post() method.
		if len(requestData) == 0 {
			// Post requires data.
			return nil, fmt.Errorf(
				"--request_data must be provided")
		}

		// Create Reader for Post data.
		contentBuffer := bytes.NewBufferString(requestData)

		res, err = http.Post(requestURL, "application/json",
			contentBuffer)
	default:
		// We doń't know how to handle this request.
		return nil, fmt.Errorf(
			"invalid --request_method provided : %s",
			requestMethod)
	}

	return res, err
}
