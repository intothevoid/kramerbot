package pipup

import (
	"encoding/json"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Pipup struct {
	Config *viper.Viper
	Logger *zap.Logger
}

/* Sample POST request from CURL -
curl -X POST -H "Content-Type: application/json" -d '{"title": "Kramerbot",
"message": "Hi I just found a deal for you!", "position": 2,
"media": {"image": {"uri": "https://i.ebayimg.com/images/g/arUAAOSwMUVfB1Ve/s-l640.jpg", "width": 200}}}'
"http://192.168.1.10:7979/notify"
*/

// Create and send sample notification via post.go
func (p *Pipup) SendMediaMessage(message string, title string) {
	duration := p.Config.GetInt("pipup.duration")
	position := p.Config.GetInt("pipup.position")
	mediaUri := p.Config.GetString("pipup.media_uri")
	imageWidth := p.Config.GetInt("pipup.image_width")
	baseUrl := p.Config.GetString("pipup.base_url")
	title_color := p.Config.GetString("pipup.title_color")
	message_color := p.Config.GetString("pipup.message_color")
	message_size := p.Config.GetInt("pipup.message_size")
	background_color := p.Config.GetString("pipup.background_color")
	title_size := p.Config.GetInt("pipup.title_size")

	toast := &models.PipupToast{
		Title:    title,
		Message:  message,
		Duration: duration,
		Position: position,
		Media: &models.PipupMedia{
			Image: &models.PipupImage{
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
