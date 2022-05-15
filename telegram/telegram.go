// package to wrap telegram bot api
package telegram

// Telegram bot api
// https://core.telegram.org/bots/api
// https://core.telegram.org/bots/api#available-methods

// imports
import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type KramerBot struct {
	Token  string
	Logger *zap.Logger
	Bot    *tgbotapi.BotAPI
}

// function to read token from environment variable
func (k *KramerBot) GetToken() string {
	// t.me/kramerbot
	token := os.Getenv("TELEGRAM_BOT_TOKEN") // get token from environment variable
	return token
}

// function to create a new bot
func (k *KramerBot) NewBot() {
	// If user has forgotten to set the token
	if k.Token == "" {
		k.Token = k.GetToken()
	}

	if k.Token == "" {
		k.Logger.Fatal("Cannot proceed without a bot token, is the TELEGRAM_BOT_TOKEN environment variable set?")
	}

	bot, err := tgbotapi.NewBotAPI(k.Token)
	if err != nil {
		k.Logger.Fatal(err.Error())
	}

	k.Logger.Info("Authorized on account", zap.String("username", bot.Self.UserName))

	// Allocate bot
	k.Bot = &tgbotapi.BotAPI{}
	k.Bot = bot
}

// send message to chat
func (k *KramerBot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	k.Bot.Send(msg)
}

// start receiving updates from telegram
func (k *KramerBot) StartReceivingUpdates() {
	// log start receiving updates
	k.Logger.Info("Start receiving updates")

	// setup updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// get updates channel
	updates, err := k.Bot.GetUpdatesChan(u)
	if err != nil {
		k.Logger.Fatal(err.Error())
	}

	// keep watching updates channel
	for update := range updates {
		if update.Message == nil {
			continue
		}

		k.Logger.Info("Received message", zap.String("text", update.Message.Text), zap.Int64("chatID", update.Message.Chat.ID))
	}
}
