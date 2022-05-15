// package to wrap telegram bot api
package telegram

// Telegram bot api
// https://core.telegram.org/bots/api
// https://core.telegram.org/bots/api#available-methods

// imports
import (
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/cosmobot/scrapers"
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

// send html message to chat
func (k *KramerBot) SendHTMLMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	k.Bot.Send(msg)
}

// send markdown message to chat
func (k *KramerBot) SendMarkdownMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	k.Bot.Send(msg)
}

// start receiving updates from telegram
func (k *KramerBot) StartReceivingUpdates(scraper scrapers.Scraper) {
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

	// create a scraper
	s := scraper.(*scrapers.OzBargainScraper)

	// keep watching updates channel
	for update := range updates {
		if update.Message == nil {
			continue
		}

		k.Logger.Info("Received message", zap.String("text", update.Message.Text), zap.Int64("chatID", update.Message.Chat.ID))

		// User asked for latest deals
		if strings.Contains(strings.ToLower(update.Message.Text), "latest") {
			k.SendLatestDeals(update.Message.Chat.ID, s)
			continue
		}

		// Help command
		if strings.Contains(strings.ToLower(update.Message.Text), "help") {
			k.Help(update.Message.Chat.ID)
			continue
		}

		// Unknown command - show help banner
		k.Help(update.Message.Chat.ID)
	}
}

// Function to send latest deals
func (k *KramerBot) SendLatestDeals(chatID int64, s *scrapers.OzBargainScraper) {
	// Let the scraper go to work
	s.Scrape()
	latestDeals := s.GetLatestDeals()

	// Send latest deals to the user
	for _, deal := range latestDeals {
		formattedDeal := fmt.Sprintf("<a href='%s' target='_blank'>CLICK HERE</a>", deal.Url)

		k.SendHTMLMessage(chatID, formattedDeal)

		// Delay for a bit don't send all deals at once
		time.Sleep(1 * time.Second)
	}
}

// Function to display help message
func (k *KramerBot) Help(chatID int64) {
	// Show the help banner
	k.SendMessage(chatID, "Giddyup! Available commands are: \n\n"+
		"/help - View this help message \n"+
		"/latest - View the 5 latest deals on OzBargain\n"+
		"/watchsuper - Watch out for deals with 50+ upvotes within the hour\n"+
		"/watchgood - Watch out for deals with 25+ upvotes within the hour\n"+
		"/watch100 - Watch out for deals with 100+ upvotes\n"+
		"/watch - Watch deals with specified keyword\n"+
		"/kramerism - Get a Kramer quote from Seinfeld")
}
