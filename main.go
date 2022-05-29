package main

import (
	"os"

	"github.com/intothevoid/kramerbot/bot"
	"github.com/intothevoid/kramerbot/models"
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

	// Get the token for the telegram bot api
	// it is safer to keep the token in an environment variable
	// and not store it in the config file to avoid security issues
	k.Token = k.GetToken()

	if k.Token == "" {
		k.Logger.Fatal("Cannot proceed without a bot token")
	}

	// create a scraper
	scraper := new(scrapers.OzBargainScraper)
	scraper.SID = scrapers.SID_OZBARGAIN
	scraper.Logger = k.Logger
	scraper.BaseUrl = scrapers.URL_OZBARGAIN
	scraper.Deals = []models.OzBargainDeal{}
	scraper.ScrapeInterval = k.Config.GetInt("scrapers.ozbargain.scrape_interval")   // mins
	scraper.MaxDealsToStore = k.Config.GetInt("scrapers.ozbargain.max_stored_deals") // max. no. of deals to store

	// create a new bot
	k.NewBot(scraper)

	// start receiving updates from telegram
	k.StartBot()
}
