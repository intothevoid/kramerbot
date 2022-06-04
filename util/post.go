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
		zap.L().Error(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error(err.Error())
	}

	zap.L().Debug("Response from pipup", zap.String("body", string(body)))

	defer resp.Body.Close()
}
