package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func (k *KramerBot) BotProc(updates tgbotapi.UpdatesChannel) {
	var announceMode bool = false

	// keep watching updates channel
	for update := range updates {
		if update.Message == nil {
			continue
		}

		k.Logger.Info("Received message", zap.String("text", update.Message.Text), zap.Int64("chatID", update.Message.Chat.ID))

		if announceMode {
			if k.verifyAdminPassword(update.Message.Text) {
				k.MakeAnnouncement(update.Message.Chat, update.Message.Text)
			} else {
				k.SendMessage(update.Message.Chat.ID, "⛔ Admin password incorrect.")
			}

			announceMode = false
			continue
		}

		// User requested admin function - announce
		if strings.Contains(strings.ToLower(update.Message.Text), "announcement") {
			if !announceMode {
				k.SendMessage(update.Message.Chat.ID, "Enter admin password and announcement in format password:announcement")
				announceMode = true
			}
			continue
		}

		// Testing
		if strings.Contains(strings.ToLower(update.Message.Text), "test") {
			k.SendTestMessage(update.Message.Chat)
			continue
		}

		// Help command
		if strings.Contains(strings.ToLower(update.Message.Text), "help") {
			k.Help(update.Message.Chat)
			continue
		}

		// Unknown command - show help banner
		k.Help(update.Message.Chat)
	}
}

// Verify pass for administrative function
func (k *KramerBot) verifyAdminPassword(message string) bool {
	messages := strings.Split(message, ":")
	if len(messages) == 2 {
		if strings.EqualFold(strings.ToLower(messages[0]), strings.ToLower(k.GetAdminPass())) {
			return true
		}
	}
	return false
}
