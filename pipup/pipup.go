package pipup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

type Pipup struct {
	Enabled         bool
	Username        string
	BaseURL         string
	Duration        int
	MediaType       string
	MediaURI        string
	ImageWidth      int
	Position        int
	TitleColor      string
	TitleSize       int
	MessageColor    string
	MessageSize     int
	BackgroundColor string
	Logger          *zap.Logger
}

/* Sample POST request from CURL -
curl -X POST -H "Content-Type: application/json" -d '{"title": "Kramerbot",
"message": "Hi I just found a deal for you!", "position": 2,
"media": {"image": {"uri": "https://i.ebayimg.com/images/g/arUAAOSwMUVfB1Ve/s-l640.jpg", "width": 200}}}'
"http://192.168.1.10:7979/notify"
*/

// Create new pipup instance
func New(config util.PipupConfig, logger *zap.Logger) *Pipup {
	return &Pipup{
		Enabled:         config.Enabled,
		Username:        config.Username,
		BaseURL:         config.BaseURL,
		Duration:        config.Duration,
		MediaType:       config.MediaType,
		MediaURI:        config.MediaURI,
		ImageWidth:      config.ImageWidth,
		Position:        config.Position,
		TitleColor:      config.TitleColor,
		TitleSize:       config.TitleSize,
		MessageColor:    config.MessageColor,
		MessageSize:     config.MessageSize,
		BackgroundColor: config.BackgroundColor,
		Logger:          logger,
	}
}

// SendMediaMessage sends a media message to the user
func (p *Pipup) SendMediaMessage(message string, title string) error {
	if !p.Enabled {
		return fmt.Errorf("pipup is not enabled")
	}

	// Create the URL
	url := fmt.Sprintf("%s/api/message", p.BaseURL)

	// Create the request body
	body := map[string]string{
		"username": p.Username,
		"message":  message,
		"title":    title,
	}

	// Convert body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// SendMessage sends a message to the user
func (p *Pipup) SendMessage(message string) error {
	if !p.Enabled {
		return fmt.Errorf("pipup is not enabled")
	}

	// Create the URL
	url := fmt.Sprintf("%s/api/message", p.BaseURL)

	// Create the request body
	body := map[string]string{
		"username": p.Username,
		"message":  message,
	}

	// Convert body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
