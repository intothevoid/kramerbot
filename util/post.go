package util

import (
	"bytes"
	"net/http"
)

// Send a post request to the webhook
func SendPostRequest(url string, jsonBody string) {
	// Create a request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		panic(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Close body
	defer resp.Body.Close()
}
