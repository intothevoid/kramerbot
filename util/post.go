package util

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

// Send a post request to the webhook
func SendPostRequest(url string, jsonBody []byte) {
	// Create a request
	// Set headers
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		Logger.Error(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	Logger.Debug("Response from pipup", zap.String("body", string(body)))
}
