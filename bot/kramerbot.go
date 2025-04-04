// package to wrap telegram bot api
package bot

// Telegram bot api
// https://core.telegram.org/bots/api
// https://core.telegram.org/bots/api#available-methods

// imports
import (
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
	sqlite_persist "github.com/intothevoid/kramerbot/persist/sqlite"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

type KramerBot struct {
	Token      string
	Logger     *zap.Logger
	BotApi     *tgbotapi.BotAPI
	OzbScraper *scrapers.OzBargainScraper
	CCCScraper *scrapers.CamCamCamScraper
	UserStore  *models.UserStore
	DataWriter persist.DatabaseIF
	Pipup      *pipup.Pipup
	Config     *util.Config
}

// function to read token from environment variable
func (k *KramerBot) GetToken() string {
	// t.me/kramerbot
	token := os.Getenv("TELEGRAM_BOT_TOKEN") // get token from environment variable
	return token
}

// function to read admin password from environment variable
func (k *KramerBot) GetAdminPass() string {
	adminPass := os.Getenv("KRAMERBOT_ADMIN_PASS") // get the admin password
	return adminPass
}

// get test mode from configuration
func (k *KramerBot) getTestMode() bool {
	return k.Config.TestMode
}

// function to create a new bot
func (k *KramerBot) NewBot(ozbs *scrapers.OzBargainScraper, cccs *scrapers.CamCamCamScraper) {
	// check test mode
	testMode := k.getTestMode()
	if testMode {
		// TEST MODE
		k.Logger.Info("****** TEST MODE IS NOW ACTIVE. Telegram not connected. ******")
	} else {
		// REGULAR MODE
		// If user has forgotten to set the token
		if k.Token == "" {
			k.Token = k.GetToken()
		}

		if k.Token == "" {
			k.Logger.Fatal("Cannot proceed without a bot token, is the TELEGRAM_BOT_TOKEN environment variable set?")
		}

		// Init telegram bot
		bot, err := tgbotapi.NewBotAPI(k.Token)
		if err != nil {
			k.Logger.Fatal(err.Error())
		}

		k.Logger.Info("Authorized on account", zap.String("username", bot.Self.UserName))

		// Allocate bot
		k.BotApi = &tgbotapi.BotAPI{}
		k.BotApi = bot
	}

	// Assign scrapers
	k.OzbScraper = ozbs
	k.CCCScraper = cccs

	// Database Initialization (SQLite)
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = k.Config.SQLite.DBPath
	}

	k.Logger.Info("Initializing SQLite database", zap.String("path", dbPath))

	// Use the NewSQLiteWrapper from the sqlite package
	dataWriter, err := sqlite_persist.NewSQLiteWrapper(dbPath, k.Logger)
	if err != nil {
		k.Logger.Fatal("Failed to initialize SQLite database", zap.String("path", dbPath), zap.Error(err))
	}
	k.DataWriter = dataWriter // Assign the wrapper which implements DatabaseIF

	// Check if the database connection is valid using Ping
	if err := k.DataWriter.Ping(); err != nil {
		k.Logger.Fatal("Failed to connect to SQLite database", zap.String("path", dbPath), zap.Error(err))
	}
	k.Logger.Info("Successfully connected to SQLite database", zap.String("path", dbPath))

	// Load user store
	k.LoadUserStore()
}

// start receiving updates from telegram
func (k *KramerBot) StartBot() {
	// check test mode
	testMode := k.getTestMode()

	// Do not send any updates when test mode is active
	if !testMode {

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
	} else {
		testTick := time.NewTicker(time.Second * time.Duration(10))
		count := 0
		for range testTick.C {
			// Test mode do nothing
			// log tick count
			count++
			k.Logger.Info("test mode active", zap.Int("tick count", count))

		}
	}
}
