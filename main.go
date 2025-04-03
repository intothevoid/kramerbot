package main

import (
	"os"

	"github.com/intothevoid/kramerbot/bot"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// create a new instance of our bot
	k := new(bot.KramerBot)

	// Setup configuration
	confPath, _ := os.Getwd()
	confPath += "/config.yaml"

	// Initialize logger first with default settings
	logger := util.SetupLogger(zapcore.DebugLevel, true)

	// Load configuration
	config, err := util.SetupConfig(confPath, logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Update logger with config settings
	logger = util.SetupLogger(zapcore.Level(config.LogLevel), config.LogToFile)
	k.Logger = logger
	k.Config = config

	// Android TV notifications via Pipup
	if config.Pipup.Enabled {
		k.Pipup = pipup.New(config.Pipup, logger)
	}

	// Get the token for the telegram bot api
	// it is safer to keep the token in an environment variable
	// and not store it in the config file to avoid security issues
	k.Token = k.GetToken()

	// Test mode doesn't require a token
	if k.Token == "" && !config.TestMode {
		logger.Fatal("Cannot proceed without a bot token, is the TELEGRAM_BOT_TOKEN environment variable set?")
	}

	// Create Ozbargain ozbscraper
	ozbscraper := new(scrapers.OzBargainScraper)
	ozbscraper.SID = scrapers.SID_OZBARGAIN
	ozbscraper.Logger = logger
	ozbscraper.BaseUrl = scrapers.URL_OZBARGAIN
	ozbscraper.Deals = []models.OzBargainDeal{}
	ozbscraper.ScrapeInterval = config.Scrapers.OzBargain.ScrapeInterval  // mins
	ozbscraper.MaxDealsToStore = config.Scrapers.OzBargain.MaxStoredDeals // max. no. of deals to store

	// Create camel camel camel (amazon) scraper
	cccscraper := new(scrapers.CamCamCamScraper)
	cccscraper.SID = scrapers.SID_CCC_AMAZON
	cccscraper.Logger = logger
	cccscraper.BaseUrl = config.Scrapers.Amazon.URLs
	cccscraper.Deals = []models.CamCamCamDeal{}
	cccscraper.ScrapeInterval = config.Scrapers.Amazon.ScrapeInterval  // mins
	cccscraper.MaxDealsToStore = config.Scrapers.Amazon.MaxStoredDeals // max. no. of deals to store

	// create a new bot
	k.NewBot(ozbscraper, cccscraper)

	// start receiving updates from telegram
	k.StartBot()
}
