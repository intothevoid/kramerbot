// package to wrap telegram bot api
package telegram

// Telegram bot api
// https://core.telegram.org/bots/api
// https://core.telegram.org/bots/api#available-methods

// imports
import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// function to read token from environment variable
func GetToken() string {
	return os.Getenv("TELEGRAM_TOKEN")
}

// function to create a new bot
func NewBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(GetToken())
	if err != nil {
		panic(err)
	}
	return bot
}

// send message to chat
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
