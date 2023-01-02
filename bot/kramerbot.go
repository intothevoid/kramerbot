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
	"github.com/intothevoid/kramerbot/api"
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
	ApiServer  *api.GinServer
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

// function to create a new bot
func (k *KramerBot) NewBot(ozbs *scrapers.OzBargainScraper, cccs *scrapers.CamCamCamScraper) {
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

	// Assign scrapers
	k.OzbScraper = ozbs
	k.CCCScraper = cccs

	// Set up data writer
	dataWriter, _ := mongo_persist.New(
		k.Config.GetString("mongo.mongo_uri"),
		k.Config.GetString("mongo.mongo_dbname"),
		k.Config.GetString("mongo.mongo_collname"),
		k.Logger,
	)
	k.DataWriter = dataWriter

	// Load user store
	k.LoadUserStore()

	// Initialise API server
	k.ApiServer = &api.GinServer{
		UserStoreDB: k.DataWriter,
		OzbScraper:  k.OzbScraper,
		CCCScraper:  k.CCCScraper,
		Config:      k.Config,
	}
	go k.ApiServer.StartServer()
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

// migration function to migrate from sqlite to mongo
func (k *KramerBot) MigrateSqliteToMongo(mongoURI string, mongoDBName string, mongoCollectionName string) {
	// Get working directory
	sqliteDBPath, _ := os.Getwd()
	sqliteDBPath = path.Join(sqliteDBPath, "users.db")

	// Start the conversion
	mongo_persist.SqliteToMongoDB(sqliteDBPath, mongoURI, mongoDBName, mongoCollectionName, k.Logger)
}
