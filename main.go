package main

import (
	"github.com/intothevoid/kramerbot/bot"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
)

func main() {
	// create a new instance of our bot
	k := new(bot.KramerBot)

	// initialise logger
	k.Logger = util.SetupLogger()

	// get the token for the telegram bot api
	k.Token = k.GetToken()

	if k.Token == "" {
		k.Logger.Fatal("Cannot proceed without a bot token")
	}

	// create a new bot
	k.NewBot()

	// create a scraper
	scraper := new(scrapers.OzBargainScraper)
	scraper.SID = scrapers.SID_OZBARGAIN
	scraper.Logger = k.Logger
	scraper.BaseUrl = scrapers.URL_OZBARGAIN
	scraper.Deals = []models.OzBargainDeal{}
	scraper.ScrapeInterval = 5 // mins

	// Assign scraper
	k.Scraper = scraper

	// start receiving updates from telegram
	k.StartBot(scraper)
}
