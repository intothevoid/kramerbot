package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

func (k *KramerBot) BotProc(updates tgbotapi.UpdatesChannel) {
	// Keyword mode is used for registering keywords to watch
	var keywordMode bool = false
	var keywordClearMode bool = false
	var announceMode bool = false

	// keep watching updates channel
	for update := range updates {
		if update.Message == nil {
			continue
		}

		k.Logger.Info("Received message", zap.String("text", update.Message.Text), zap.Int64("chatID", update.Message.Chat.ID))

		// if keyword mode is on, process keyword
		if keywordMode {
			k.ProcessKeyword(update.Message.Chat, update.Message.Text)
			keywordMode = false
			continue
		}

		// if keyword clear mode is on, process clear keyword
		if keywordClearMode {
			k.ProcessClearKeyword(update.Message.Chat, update.Message.Text)
			keywordClearMode = false
			continue
		}

		if announceMode {
			if k.verifyAdminPassword(update.Message.Text) {
				k.MakeAnnouncement(update.Message.Chat, update.Message.Text)
			} else {
				k.SendMessage(update.Message.Chat.ID, "⛔ Admin password incorrect.")
			}

			announceMode = false
			continue
		}

		// User asked for latest deals
		if strings.Contains(strings.ToLower(update.Message.Text), "latest") {
			k.SendLatestDeals(update.Message.Chat.ID, k.OzbScraper)
			continue
		}

		// User asked to watch super deals i.e. 50+ upvotes within the hour
		if strings.Contains(strings.ToLower(update.Message.Text), "watchsuper") {
			k.WatchOzbSuperDeals(update.Message.Chat)
			continue
		}

		// User asked to watch good OZB deals i.e. 25+ upvotes within the hour
		if strings.Contains(strings.ToLower(update.Message.Text), "watchgood") {
			k.WatchOzbGoodDeals(update.Message.Chat)
			continue
		}

		// User asked to watch daily AMZ deals i.e. top daily deals with price drop > 20%
		if strings.Contains(strings.ToLower(update.Message.Text), "amazondaily") {
			k.WatchAmzDailyDeals(update.Message.Chat)
			continue
		}

		// User asked to watch weekly AMZ deals i.e. top weekly deals with price drop > 20%
		if strings.Contains(strings.ToLower(update.Message.Text), "amazonweekly") {
			k.WatchAmzWeeklyDeals(update.Message.Chat)
			continue
		}

		// User asked to watch specific keyword
		if strings.Contains(strings.ToLower(update.Message.Text), "watchkeyword") {
			if !keywordMode {
				k.SendMessage(update.Message.Chat.ID, "Enter keyword to watch. Example: 'headphone' or 'sennheiser'")
				keywordMode = true
			}
			continue
		}

		// User asked to clear specific watched keyword
		if strings.Contains(strings.ToLower(update.Message.Text), "clearkeyword") {
			if !keywordClearMode {
				k.SendMessage(update.Message.Chat.ID, "Enter keyword to clear. Example: 'headphone' or 'sennheiser'")
				keywordClearMode = true
			}
			continue
		}

		// User asked to clear all watched keywords
		if strings.Contains(strings.ToLower(update.Message.Text), "clearallkeywords") {
			k.ProcessClearAllKeywords(update.Message.Chat)
			continue
		}

		// User requested current status
		if strings.Contains(strings.ToLower(update.Message.Text), "status") {
			k.SendStatus(update.Message.Chat)
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

		// User asked for a kramerism
		if strings.Contains(strings.ToLower(update.Message.Text), "kramerism") {
			kramerism := util.GetKramerism()
			k.SendMessage(update.Message.Chat.ID, kramerism)
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
