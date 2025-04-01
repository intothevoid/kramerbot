// package to wrap telegram bot api
package bot

// Telegram bot api
// https://core.telegram.org/bots/api
// https://core.telegram.org/bots/api#available-methods

// imports
import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/intothevoid/kramerbot/api"
	"github.com/intothevoid/kramerbot/interfaces"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
	mongo_persist "github.com/intothevoid/kramerbot/persist/mongo"
	sqlite_persist "github.com/intothevoid/kramerbot/persist/sqlite"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
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
	WebAppURL  string
}

var _ interfaces.BotAPIBridge = (*KramerBot)(nil) // Compile-time check

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

	// Database Initialization (SQLite)
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = k.Config.GetString("sqlite.db_path")
		if dbPath == "" {
			dbPath = "data/users.db" // Default path if neither env var nor config is set
		}
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

	// Set Web App URL from config
	k.WebAppURL = k.Config.GetString("web_app.base_url")
	// If base_url is not set, try to infer from listen_address (for local testing)
	if k.WebAppURL == "" {
		listenAddr := k.Config.GetString("web_app.listen_address")
		if listenAddr != "" {
			// Basic attempt to make it http://localhost:PORT
			parts := strings.Split(listenAddr, ":")
			if len(parts) == 2 {
				k.WebAppURL = "http://localhost:" + parts[1]
			} else if len(parts) == 1 {
				// Assume default http port if only address is given? Risky.
				k.Logger.Warn("Cannot reliably determine Web App URL from listen_address without port", zap.String("listen_address", listenAddr))
			}
		}
	}
	if k.WebAppURL == "" {
		k.Logger.Warn("web_app.base_url is not configured, Web App feature might not be accessible correctly.")
	}
}

// start receiving updates from telegram
func (k *KramerBot) StartBot() {
	// check test mode
	testMode := k.getTestMode()

	// Start Web App server if enabled
	if k.Config.GetBool("web_app.enabled") && !testMode {
		go k.StartWebServer() // Run web server in a separate goroutine
	}

	// Do not send any updates when test mode is active
	if !testMode {

		// log start receiving updates
		k.Logger.Info("Start receiving updates")

		// setup updates
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		// get updates channel (v5 API change: no error returned here)
		updates := k.BotApi.GetUpdatesChan(u)
		// Note: Error handling might be needed later when reading from the channel
		// if err != nil {
		// 	k.Logger.Fatal(err.Error())
		// }

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

// StartWebServer sets up and starts the HTTP server for the web app
func (k *KramerBot) StartWebServer() {
	listenAddr := k.Config.GetString("web_app.listen_address")
	if listenAddr == "" {
		k.Logger.Error("Web App is enabled but web_app.listen_address is not configured.")
		return
	}

	// Get current working directory to locate webapp files
	workDir, err := os.Getwd()
	if err != nil {
		k.Logger.Error("Failed to get working directory", zap.Error(err))
		return
	}
	webappDir := path.Join(workDir, "webapp")

	// Check if webapp directory exists
	if _, err := os.Stat(webappDir); os.IsNotExist(err) {
		k.Logger.Error("Web App directory not found", zap.String("path", webappDir))
		return
	}

	// Create API handler instance
	devMode := k.Config.GetBool("web_app.development_mode") // Get development mode setting
	apiHandler := api.NewAPI(k, k.Logger, devMode)

	// Log startup mode
	if devMode {
		k.Logger.Info("Starting web app in DEVELOPMENT mode - Authentication bypassed, dummy data will be used")
	} else {
		k.Logger.Info("Starting web app in PRODUCTION mode - Full authentication required")
	}

	// Setup HTTP routes
	mux := http.NewServeMux()

	// API routes (with auth middleware)
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("/preferences", apiHandler.HandlePreferences)
	apiRouter.HandleFunc("/keywords/add", apiHandler.HandleAddKeyword)
	apiRouter.HandleFunc("/keywords/remove", apiHandler.HandleRemoveKeyword)
	apiRouter.HandleFunc("/test", apiHandler.HandleTestNotification)
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler.ValidateAuthMiddleware(apiRouter)))

	// Static file server for the web app frontend
	fs := http.FileServer(http.Dir(webappDir))
	mux.Handle("/", fs) // Serve index.html and other assets

	k.Logger.Info("Starting Web App server", zap.String("address", listenAddr), zap.String("serving", webappDir))

	// Start the server
	err = http.ListenAndServe(listenAddr, mux)
	if err != nil {
		k.Logger.Error("Web App server failed", zap.Error(err))
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

// GetUserDataWriter returns the database interface
func (k *KramerBot) GetUserDataWriter() persist.DatabaseIF {
	return k.DataWriter
}

// GetBotToken returns the bot's Telegram token
func (k *KramerBot) GetBotToken() string {
	return k.Token
}

// GetUserData retrieves user data (can potentially just use DataWriter directly via interface)
func (k *KramerBot) GetUserData(chatID int64) (*models.UserData, error) {
	return k.DataWriter.GetUser(chatID)
}

// UpdateUserData updates user data in DB and in-memory store
func (k *KramerBot) UpdateUserData(user *models.UserData) error {
	// Update in DB
	if err := k.DataWriter.UpdateUser(user); err != nil {
		k.Logger.Error("Failed to update user in DB via bridge", zap.Int64("chatID", user.ChatID), zap.Error(err))
		return err
	}
	// Update in-memory store
	k.UserStore.Users[user.ChatID] = user
	k.Logger.Debug("Updated user in DB and memory via bridge", zap.Int64("chatID", user.ChatID))
	return nil
}

// SendTestMessageToChat sends a test message given a chat ID and user info
func (k *KramerBot) SendTestMessageToChat(chatID int64, userInfo *models.TelegramUser) error {
	// Construct a simplified message or use existing logic
	shortenedTitle := util.ShortenString("🔥 This is a test deal from the Web App!", 30) + "..."
	dealUrl := "https://news.google.com.au"
	formattedDeal := fmt.Sprintf(`🔥<a href='%s' target='_blank'>%s</a>`, dealUrl, shortenedTitle)

	userName := "User"
	if userInfo != nil {
		userName = userInfo.FirstName
	}

	k.Logger.Debug(fmt.Sprintf("Sending test deal %s to user %s (%d) via Web App", shortenedTitle, userName, chatID))
	// Directly use SendHTMLMessage (assuming it's safe to call from this context)
	k.SendHTMLMessage(chatID, formattedDeal)
	return nil // Assuming SendHTMLMessage doesn't return an error we need to propagate
}
