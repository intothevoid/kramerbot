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

		// Check if the message is a command
		if update.Message.IsCommand() {
			command := update.Message.Command()
			args := update.Message.CommandArguments()

			k.Logger.Info("Received command", zap.String("command", command), zap.String("args", args), zap.Int64("chatID", update.Message.Chat.ID))

			switch command {
			case "start", "register":
				k.RegisterUser(update.Message.Chat)
				continue
			case "help":
				k.Help(update.Message.Chat)
				continue
			case "preferences", "status":
				k.ShowPreferences(update.Message.Chat)
				continue
			case "listkeywords":
				k.ListKeywords(update.Message.Chat)
				continue
			case "addkeyword":
				k.AddKeyword(update.Message.Chat, args)
				continue
			case "removekeyword":
				k.RemoveKeyword(update.Message.Chat, args)
				continue
			case "ozbgood":
				k.ToggleOzbGood(update.Message.Chat)
				continue
			case "ozbsuper":
				k.ToggleOzbSuper(update.Message.Chat)
				continue
			case "amzdaily":
				k.ToggleAmzDaily(update.Message.Chat)
				continue
			case "amzweekly":
				k.ToggleAmzWeekly(update.Message.Chat)
				continue
			case "test":
				k.SendTestMessage(update.Message.Chat)
				continue
			case "announce": // Admin command
				if !announceMode {
					k.SendMessage(update.Message.Chat.ID, "Enter admin password and announcement in format password:announcement")
				}
				announceMode = !announceMode // Toggle announce mode
				continue
			default:
				// Unknown command - show help banner
				k.Help(update.Message.Chat)
				continue
			}
		}

		// If it's not a command and not in announce mode, treat as potential announcement password/message
		if announceMode {
			if k.verifyAdminPassword(update.Message.Text) {
				k.MakeAnnouncement(update.Message.Chat, update.Message.Text)
			} else {
				k.SendMessage(update.Message.Chat.ID, "⛔ Admin password incorrect or invalid format.")
			}
			announceMode = false // Reset announce mode after attempt
			continue
		}

		// If it's not a command and not announce mode, ignore or show help?
		// For now, let's show help for any non-command text.
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
