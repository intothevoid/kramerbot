package main

import (
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/telegram"
	"github.com/intothevoid/kramerbot/util"
)

func main() {
	// create a new instance of our bot
	k := new(telegram.KramerBot)

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
	scraper.BaseUrl = "https://www.ozbargain.com.au/"
	scraper.Deals = []scrapers.OzBargainDeal{}

	// start receiving updates from telegram
	k.StartReceivingUpdates(scraper)
}
