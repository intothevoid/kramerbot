package pipup

import "github.com/intothevoid/kramerbot/util"

// Create and send sample notification via post.go
func SendMediaMessage(message string, title string) {
	// Create json message
	jsonBody := `{
        "duration": 30,
        "position": 0,
        "title": "` + title + `",
        "titleColor": "#0066cc",
        "titleSize": 20,
        "message": "` + message + `",
        "messageColor": "#000000",
        "messageSize": 14,
        "backgroundColor": "#ffffff",
        "media": { "image": {
            "uri": "../static/kramer_beach.jpg", "width": 480
        }},
    }`

	// Send post request
	util.SendPostRequest("http://192.168.1.10:7979/notify", jsonBody)
}

// Send a simple message
func SendMessage(message string, title string) {
	// Create json message
	jsonBody := `{
        "duration": 15,
        "title": "` + title + `",
        "message": "` + message + `"
    }`

	// Send post request
	util.SendPostRequest("http://192.168.1.10:7979/notify", jsonBody)
}
