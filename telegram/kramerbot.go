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
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

type KramerBot struct {
	Token     string
	Logger    *zap.Logger
	Bot       *tgbotapi.BotAPI
	Scraper   scrapers.Scraper
	UserStore *models.UserStore
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

	// Load user store
	k.LoadUserStore()
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

	// Start scraping
	s := scraper.(*scrapers.OzBargainScraper)
	s.AutoScrape()

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

		// User asked to watch super deals i.e. 50+ upvotes within the hour
		if strings.Contains(strings.ToLower(update.Message.Text), "watchsuper") {
			k.WatchSuperDeals(update.Message.Chat)
			continue
		}

		// User asked to watch good deals i.e. 25+ upvotes within the hour
		if strings.Contains(strings.ToLower(update.Message.Text), "watchgood") {
			k.WatchGoodDeals(update.Message.Chat)
			continue
		}

		// User asked to watch 100+ upvotes deals
		if strings.Contains(strings.ToLower(update.Message.Text), "watch100") {
			k.Watch100Deals(update.Message.Chat)
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

// Function to send latest deals i.e. NUM_DEALS_TO_SEND
func (k *KramerBot) SendLatestDeals(chatID int64, s *scrapers.OzBargainScraper) {
	latestDeals := s.GetLatestDeals(scrapers.NUM_DEALS_TO_SEND)

	// Send latest deals to the user
	for _, deal := range latestDeals {
		shortenedTitle := util.ShortenString(deal.Title, 40) + "..."
		formattedDeal := fmt.Sprintf("<a href='%s' target='_blank'>%s</a>", deal.Url, shortenedTitle)

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
		"/keywordwatch - Watch deals with specified keyword\n"+
		"/keywordclear - Clear deals with specified keyword\n"+
		"/keywordclearall - Clear deals with all watched keywords\n"+
		"/kramerism - Get a Kramer quote from Seinfeld")
}

// Add watch to super deals by chat id
func (k *KramerBot) WatchSuperDeals(chat *tgbotapi.Chat) {

	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]
		userData.SuperDeals = true
	} else {
		// Key does not exist, create new user
		userData := k.CreateUserData(chat.ID, chat.FirstName, "", false, false, true)
		k.UserStore.Users[chat.ID] = userData
	}

	// Send message to user
	k.SendMessage(chat.ID, fmt.Sprintf("%s, you are now added to the super deals watchlist.", chat.FirstName))
}

// Add watch to good deals by chat id
func (k *KramerBot) WatchGoodDeals(chat *tgbotapi.Chat) {

	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]
		userData.GoodDeals = true
	} else {
		// Key does not exist, create new user
		userData := k.CreateUserData(chat.ID, chat.FirstName, "", false, true, false)
		k.UserStore.Users[chat.ID] = userData
	}

	// Send message to user
	k.SendMessage(chat.ID, fmt.Sprintf("%s, you are now added to the good deals watchlist.", chat.FirstName))
}

// Add watch to super deals by chat id
func (k *KramerBot) Watch100Deals(chat *tgbotapi.Chat) {

	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]
		userData.Deals100 = true
	} else {
		// Key does not exist, create new user
		userData := k.CreateUserData(chat.ID, chat.FirstName, "", true, false, false)
		k.UserStore.Users[chat.ID] = userData
	}

	// Send message to user
	k.SendMessage(chat.ID, fmt.Sprintf("%s, you are now added to the 100+ upvotes deals watchlist.", chat.FirstName))
}

// Create user data from parameters passed in
func (k *KramerBot) CreateUserData(chatID int64, username string, keywords string, deals100 bool,
	goodDeals bool, superDeals bool) models.UserData {

	userData := models.UserData{}
	userData.ChatID = chatID
	userData.Username = username
	userData.Keywords = keywords
	userData.Deals100 = deals100
	userData.GoodDeals = goodDeals
	userData.SuperDeals = superDeals

	return userData
}

// Function to load user store from file
func (k *KramerBot) LoadUserStore() {
	// Load user store i.e. user data indexed by chat id
	store := util.DataStore{Logger: k.Logger}
	k.UserStore = store.ReadUserStore()
}

// Function to save user store to file
func (k *KramerBot) SaveUserStore() {
	// Save user store i.e. user data indexed by chat id
	store := util.DataStore{Logger: k.Logger}
	store.WriteUserStore(k.UserStore)
}
