package bot

import (
	"io/ioutil"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

// send message to chat
func (k *KramerBot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	k.BotApi.Send(msg)
}

// send html message to chat
func (k *KramerBot) SendHTMLMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	k.BotApi.Send(msg)
}

// send markdown message to chat
func (k *KramerBot) SendMarkdownMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	k.BotApi.Send(msg)
}

// Send a photo to the user
func (k *KramerBot) SendPhoto(chatID int64, fileName string) {
	// Convert to absolute path if relative path sent
	if !filepath.IsAbs(fileName) {
		fileName, _ = filepath.Abs(fileName)
	}

	filebytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		k.Logger.Error("Unable to read file", zap.Error(err))
		return
	}

	// Get filename from path
	fname := filepath.Base(fileName)

	photobytes := tgbotapi.FileBytes{
		Name:  fname,
		Bytes: filebytes,
	}
	msg := tgbotapi.NewPhotoUpload(chatID, photobytes)
	k.BotApi.Send(msg)
}

// Send a video to the user
func (k *KramerBot) SendVideo(chatID int64, fileName string) {
	// Convert to absolute path if relative path sent
	if !filepath.IsAbs(fileName) {
		fileName, _ = filepath.Abs(fileName)
	}

	filebytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		k.Logger.Error("Unable to read file", zap.Error(err))
		return
	}

	// Get filename from path
	fname := filepath.Base(fileName)

	photobytes := tgbotapi.FileBytes{
		Name:  fname,
		Bytes: filebytes,
	}
	msg := tgbotapi.NewVideoUpload(chatID, photobytes)
	k.BotApi.Send(msg)
}
