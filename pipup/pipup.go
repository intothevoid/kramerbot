package pipup

import (
	"encoding/json"
	"strings"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Pipup struct {
	Config   *viper.Viper
	Logger   *zap.Logger
	Username string
}

/* Sample POST request from CURL -
curl -X POST -H "Content-Type: application/json" -d '{"title": "Kramerbot",
"message": "Hi I just found a deal for you!", "position": 2,
"media": {"image": {"uri": "https://i.ebayimg.com/images/g/arUAAOSwMUVfB1Ve/s-l640.jpg", "width": 200}}}'
"http://192.168.1.10:7979/notify"
*/

// Create a new pipup instance
func New(config *viper.Viper, logger *zap.Logger) *Pipup {
	return &Pipup{
		Config:   config,
		Logger:   logger,
		Username: config.GetString("pipup.username"),
	}
}

// Create and send sample notification via post.go
func (p *Pipup) SendMediaMessage(message string, title string) {
	// Do not send message if pipup is disabled
	enabled := p.Config.GetBool("pipup.enabled")
	if !enabled {
		return
	}

	duration := p.Config.GetInt("pipup.duration")
	position := p.Config.GetInt("pipup.position")
	mediaUri := p.Config.GetString("pipup.media_uri")
	mediaType := strings.ToLower(p.Config.GetString("pipup.media_type"))
	imageWidth := p.Config.GetInt("pipup.image_width")
	baseUrl := p.Config.GetString("pipup.base_url")
	title_color := p.Config.GetString("pipup.title_color")
	message_color := p.Config.GetString("pipup.message_color")
	message_size := p.Config.GetInt("pipup.message_size")
	background_color := p.Config.GetString("pipup.background_color")
	title_size := p.Config.GetInt("pipup.title_size")

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
	enabled := p.Config.GetBool("pipup.enabled")
	if !enabled {
		return
	}

	duration := p.Config.GetInt("pipup.duration")
	position := p.Config.GetInt("pipup.position")
	baseUrl := p.Config.GetString("pipup.base_url")

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
