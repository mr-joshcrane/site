package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mr-joshcrane/site"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lambdaRequest, err := convertToLambdaFunctionURLRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := site.Handler(nil, *lambdaRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		convertLambdaFunctionResponse(w, response)
	})

	site := &http.Server{

		Addr:    ":8080",
		Handler: mux,
	}
	err := site.ListenAndServe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// Function to convert an HTTP request to a LambdaFunctionURLRequest
func convertToLambdaFunctionURLRequest(r *http.Request) (*events.LambdaFunctionURLRequest, error) {
	// Read the body of the request
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Extract headers
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = values[0] // Assuming single-value headers for simplicity
	}

	// Extract cookies
	var cookies []string
	for _, cookie := range r.Cookies() {
		cookies = append(cookies, cookie.String())
	}

	// Extract query parameters
	queryParams := make(map[string]string)
	for name, values := range r.URL.Query() {
		queryParams[name] = values[0] // Assuming single-value query parameters for simplicity
	}

	// Create the LambdaFunctionURLRequest
	lambdaRequest := &events.LambdaFunctionURLRequest{
		Version:               "2.0",
		RawPath:               r.URL.Path,
		RawQueryString:        r.URL.RawQuery,
		Cookies:               cookies,
		Headers:               headers,
		QueryStringParameters: queryParams,
		RequestContext:        events.LambdaFunctionURLRequestContext{}, // Populate as needed
		Body:                  string(bodyBytes),
		IsBase64Encoded:       false, // Change this if the body is Base64 encoded
	}

	return lambdaRequest, nil
}

func convertLambdaFunctionResponse(w http.ResponseWriter, response events.LambdaFunctionURLResponse) {
	// Set the status code
	w.WriteHeader(response.StatusCode)

	// Set the headers
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}

	// Set the body
	w.Write([]byte(response.Body))
}
