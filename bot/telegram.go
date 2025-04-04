package bot

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// send message to chat
func (k *KramerBot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := k.BotApi.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// send html message to chat
func (k *KramerBot) SendHTMLMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := k.BotApi.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send HTML message: %w", err)
	}
	return nil
}

// send markdown message to chat
func (k *KramerBot) SendMarkdownMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := k.BotApi.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Markdown message: %w", err)
	}
	return nil
}

// Send a photo to the user
func (k *KramerBot) SendPhoto(chatID int64, fileName string) error {
	// Convert to absolute path if relative path sent
	if !filepath.IsAbs(fileName) {
		fileName, _ = filepath.Abs(fileName)
	}

	filebytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	// Get filename from path
	fname := filepath.Base(fileName)

	photobytes := tgbotapi.FileBytes{
		Name:  fname,
		Bytes: filebytes,
	}
	msg := tgbotapi.NewPhotoUpload(chatID, photobytes)
	_, err = k.BotApi.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send photo: %w", err)
	}
	return nil
}

// Send a video to the user
func (k *KramerBot) SendVideo(chatID int64, fileName string) error {
	// Convert to absolute path if relative path sent
	if !filepath.IsAbs(fileName) {
		fileName, _ = filepath.Abs(fileName)
	}

	filebytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	// Get filename from path
	fname := filepath.Base(fileName)

	photobytes := tgbotapi.FileBytes{
		Name:  fname,
		Bytes: filebytes,
	}
	msg := tgbotapi.NewVideoUpload(chatID, photobytes)
	_, err = k.BotApi.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send video: %w", err)
	}
	return nil
}
