package main

import (
	"os"

	"github.com/intothevoid/kramerbot/bot"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/pipup"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap/zapcore"
)

func main() {
	// create a new instance of our bot
	k := new(bot.KramerBot)

	// Setup configuration
	confPath, _ := os.Getwd()
	confPath += "/config.yaml"
	appconf, err := util.SetupConfig(confPath)
	if err != nil {
		panic("Error reading application config")
	}
	k.Config = appconf

	// initialise logger
	k.Logger = util.SetupLogger(zapcore.Level(appconf.GetInt("log_level")), appconf.GetBool("log_to_file"))

	// migration mode
	if appconf.GetBool("mongo.migration_mode") {
		k.Logger.Info("Migration mode enabled")
		k.MigrateSqliteToMongo(
			appconf.GetString("mongo.mongo_uri"),
			appconf.GetString("mongo.mongo_dbname"),
			appconf.GetString("mongo.mongo_collname"),
		)
		return
	}

	// Android TV notifications via Pipup
	k.Pipup = pipup.New(k.Config, k.Logger)

	// Get the token for the telegram bot api
	// it is safer to keep the token in an environment variable
	// and not store it in the config file to avoid security issues
	k.Token = k.GetToken()

	if k.Token == "" {
		k.Logger.Fatal("Cannot proceed without a bot token")
	}

	// Create Ozbargain ozbscraper
	ozbscraper := new(scrapers.OzBargainScraper)
	ozbscraper.SID = scrapers.SID_OZBARGAIN
	ozbscraper.Logger = k.Logger
	ozbscraper.BaseUrl = scrapers.URL_OZBARGAIN
	ozbscraper.Deals = []models.OzBargainDeal{}
	ozbscraper.ScrapeInterval = k.Config.GetInt("scrapers.ozbargain.scrape_interval")   // mins
	ozbscraper.MaxDealsToStore = k.Config.GetInt("scrapers.ozbargain.max_stored_deals") // max. no. of deals to store

	// Create camel camel camel (amazon) scraper
	cccscraper := new(scrapers.CamCamCamScraper)
	cccscraper.SID = scrapers.SID_CCC_AMAZON
	cccscraper.Logger = k.Logger
	cccscraper.BaseUrl = k.Config.GetStringSlice("scrapers.amazon.urls")
	cccscraper.Deals = []models.CamCamCamDeal{}
	cccscraper.ScrapeInterval = k.Config.GetInt("scrapers.amazon.scrape_interval")   // mins
	cccscraper.MaxDealsToStore = k.Config.GetInt("scrapers.amazon.max_stored_deals") // max. no. of deals to store

	// create a new bot
	k.NewBot(ozbscraper, cccscraper)

	// start receiving updates from telegram
	k.StartBot()
}
