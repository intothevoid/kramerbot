// package to wrap telegram bot api
package bot

// Telegram bot api
// https://core.telegram.org/bots/api
// https://core.telegram.org/bots/api#available-methods

// imports
import (
	"os"
	"path"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	persist "github.com/intothevoid/kramerbot/persist"
	mongo_persist "github.com/intothevoid/kramerbot/persist/mongo"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/spf13/viper"
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
	Config     *viper.Viper
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
	testMode := k.Config.GetBool("test_mode")
	return testMode
}

// function to create a new bot
func (k *KramerBot) NewBot(ozbs *scrapers.OzBargainScraper, cccs *scrapers.CamCamCamScraper) {
	// check test mode
	testMode := k.getTestMode()
	if testMode {
		// TEST MODE
		k.Logger.Info("****** TEST MODE IS NOW ACTIVE. Telegram not connected. ******")

		// Make entries to dummy database
		// dataWriter, _ := dummy_persist.New(
		// 	"dummy_uri",
		// 	"dummy_dbname",
		// 	"dummy_collname",
		// 	k.Logger,
		// )
		// k.DataWriter = dataWriter
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

	// Check for environment variables first, then fall back to config values
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = k.Config.GetString("mongo.mongo_uri")
	}

	mongoDBName := os.Getenv("MONGO_DBNAME")
	if mongoDBName == "" {
		mongoDBName = k.Config.GetString("mongo.mongo_dbname")
	}

	mongoCollName := os.Getenv("MONGO_COLLNAME")
	if mongoCollName == "" {
		mongoCollName = k.Config.GetString("mongo.mongo_collname")
	}

	k.Logger.Info("Connecting to MongoDB",
		zap.String("uri", mongoURI),
		zap.String("database", mongoDBName),
		zap.String("collection", mongoCollName))

	// Real mode, make entries to real database
	dataWriter, _ := mongo_persist.New(
		mongoURI,
		mongoDBName,
		mongoCollName,
		k.Logger,
	)
	k.DataWriter = dataWriter

	// Check if the database connection is valid
	if err := k.DataWriter.Ping(); err != nil {
		k.Logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	k.Logger.Info("Successfully connected to MongoDB")

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

// migration function to migrate from sqlite to mongo
func (k *KramerBot) MigrateSqliteToMongo(mongoURI string, mongoDBName string, mongoCollectionName string) {
	// Get working directory
	sqliteDBPath, _ := os.Getwd()
	sqliteDBPath = path.Join(sqliteDBPath, "users.db")

	// Start the conversion
	mongo_persist.SqliteToMongoDB(sqliteDBPath, mongoURI, mongoDBName, mongoCollectionName, k.Logger)
}
