// package to wrap telegram bot api
package bot

// Telegram bot api
// https://core.telegram.org/bots/api
// https://core.telegram.org/bots/api#available-methods

// imports
import (
	"os"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type KramerBot struct {
	Token      string
	Logger     *zap.Logger
	BotApi     *tgbotapi.BotAPI
	Scraper    *scrapers.OzBargainScraper
	UserStore  *models.UserStore
	DataWriter *persist.UserStoreDB
	Pipup      *pipup.Pipup
	Config     *viper.Viper
}

// function to read token from environment variable
func (k *KramerBot) GetToken() string {
	// t.me/kramerbot
	token := os.Getenv("TELEGRAM_BOT_TOKEN") // get token from environment variable
	return token
}

// function to create a new bot
func (k *KramerBot) NewBot(s *scrapers.OzBargainScraper) {
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
	k.BotApi = &tgbotapi.BotAPI{}
	k.BotApi = bot

	// Assign scraper
	k.Scraper = s

	// Get working directory
	dbPath, _ := os.Getwd()
	dbPath = path.Join(dbPath, "users.db")

	// Set up data writer
	k.DataWriter = persist.CreateDatabaseConnection(dbPath, k.Logger)

	// Load user store
	k.LoadUserStore()
}

// start receiving updates from telegram
func (k *KramerBot) StartBot() {
	// log start receiving updates
	k.Logger.Info("Start receiving updates")

	// setup updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// get updates channel
	updates, err := k.BotApi.GetUpdatesChan(u)
	if err != nil {
		k.Logger.Fatal(err.Error())
	}

	// Start processing deals and scraping
	// Run asyncronously to avoid blocking the main thread
	go func() {
		k.StartProcessing()
	}()

	// Start monitoring the bots updates channel
	k.BotProc(updates)
}
