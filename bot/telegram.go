package bot

import (
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		var err error
		fileName, err = filepath.Abs(fileName)
		if err != nil {
			k.Logger.Error("Failed to get absolute path for photo", zap.String("file", fileName), zap.Error(err))
			return
		}
	}

	// Check if file exists before creating FilePath (more robust)
	if _, err := os.Stat(fileName); err != nil {
		k.Logger.Error("Photo file not found or inaccessible", zap.String("path", fileName), zap.Error(err))
		return
	}

	// V5 uses FilePath for local files
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(fileName))
	// Optional: Add caption etc.
	// photo.Caption = "Here is your photo!"

	if _, err := k.BotApi.Send(photo); err != nil {
		k.Logger.Error("Failed to send photo", zap.String("file", fileName), zap.Error(err))
	}
}

// Send a video to the user
func (k *KramerBot) SendVideo(chatID int64, fileName string) {
	// Convert to absolute path if relative path sent
	if !filepath.IsAbs(fileName) {
		var err error
		fileName, err = filepath.Abs(fileName)
		if err != nil {
			k.Logger.Error("Failed to get absolute path for video", zap.String("file", fileName), zap.Error(err))
			return
		}
	}

	// Check if file exists
	if _, err := os.Stat(fileName); err != nil {
		k.Logger.Error("Video file not found or inaccessible", zap.String("path", fileName), zap.Error(err))
		return
	}

	// V5 uses FilePath for local files
	video := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(fileName))
	// Optional: Add caption, duration, dimensions etc.
	// video.Caption = "Here is your video!"

	if _, err := k.BotApi.Send(video); err != nil {
		k.Logger.Error("Failed to send video", zap.String("file", fileName), zap.Error(err))
	}
}
