package main

import (
	"fmt"

	"github.com/intothevoid/ozbot/scrapers"
	"github.com/intothevoid/ozbot/util"
)

func main() {
	// create a new scraper
	ozbScraper := new(scrapers.OzBargainScraper)

	// initialise logger
	ozbScraper.Logger = util.SetupLogger()
	ozbScraper.BaseUrl = "https://www.ozbargain.com.au/"

	// Start scraping
	ozbScraper.Scrape()

	// Print deals
	fmt.Println(ozbScraper.Deals)
}
