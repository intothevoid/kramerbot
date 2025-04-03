package pipup

import (
	"encoding/json"

	"github.com/intothevoid/kramerbot/models"
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

// Create and send sample notification via post.go
func (p *Pipup) SendMediaMessage(message string, title string) {
	// Do not send message if pipup is disabled
	if !p.Enabled {
		return
	}

	duration := p.Duration
	position := p.Position
	mediaUri := p.MediaURI
	mediaType := p.MediaType
	imageWidth := p.ImageWidth
	baseUrl := p.BaseURL
	title_color := p.TitleColor
	message_color := p.MessageColor
	message_size := p.MessageSize
	background_color := p.BackgroundColor
	title_size := p.TitleSize

	// Initialise toast
	toast := &models.PipupToast{}

	if mediaType == "image" {
		toast = &models.PipupToast{
			Title:    title,
			Message:  message,
			Duration: duration,
			Position: position,
			Media: &models.PipupImage{
				Image: &models.PipupUri{
					Uri:   mediaUri,
					Width: imageWidth,
				},
			},
			TitleColor:      title_color,
			TitleSize:       title_size,
			MessageColor:    message_color,
			MessageSize:     message_size,
			BackgroundColor: background_color,
		}
	}

	if mediaType == "video" {
		toast = &models.PipupToast{
			Title:    title,
			Message:  message,
			Duration: duration,
			Position: position,
			Media: &models.PipupVideo{
				Video: &models.PipupUri{
					Uri:   mediaUri,
					Width: imageWidth,
				},
			},
			TitleColor:      title_color,
			TitleSize:       title_size,
			MessageColor:    message_color,
			MessageSize:     message_size,
			BackgroundColor: background_color,
		}
	}

	if mediaType == "web" {
		toast = &models.PipupToast{
			Title:    title,
			Message:  message,
			Duration: duration,
			Position: position,
			Media: &models.PipupWeb{
				Web: &models.PipupUri{
					Uri:   mediaUri,
					Width: imageWidth,
				},
			},
			TitleColor:      title_color,
			TitleSize:       title_size,
			MessageColor:    message_color,
			MessageSize:     message_size,
			BackgroundColor: background_color,
		}
	}

	// Create json message
	jsonBody, err := json.Marshal(toast)

	if err != nil {
		p.Logger.Error("Error marshalling andtoid tv notification json", zap.Error(err))
	}

	// Send post request
	p.Logger.Debug("Sending pipup request", zap.String("url", baseUrl), zap.String("json", string(jsonBody)))
	util.SendPostRequest(baseUrl, jsonBody)
}

// Send a simple message
func (p *Pipup) SendMessage(message string, title string) {
	// Do not send message if pipup is disabled
	if !p.Enabled {
		return
	}

	duration := p.Duration
	position := p.Position
	baseUrl := p.BaseURL

	toast := &models.PipupSimpleToast{
		Title:    title,
		Message:  message,
		Duration: duration,
		Position: position,
	}

	// Create json message
	jsonBody, err := json.Marshal(toast)

	if err != nil {
		p.Logger.Error("Error marshalling andtoid tv notification json", zap.Error(err))
	}

	// Send post request
	p.Logger.Debug("Sending pipup request", zap.String("url", baseUrl), zap.String("json", string(jsonBody)))
	util.SendPostRequest(baseUrl, jsonBody)
}
