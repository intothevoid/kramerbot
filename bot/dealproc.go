package bot

import (
	"strings"
	"time"

	"github.com/intothevoid/kramerbot/scrapers"
	"go.uber.org/zap"
)

// Process deals returned by the scraper, check deal type and notify user
// if they are subscribed to a particular deal type
func (k *KramerBot) StartProcessing() {
	// Load user store i.e. user data indexed by chat id
	k.LoadUserStore()

	// Begin timed processing and scraping
	// tick := time.NewTicker(time.Second * 60)
	tick := time.NewTicker(time.Minute * time.Duration(k.OzbScraper.ScrapeInterval))
	for range tick.C {

		// Process Ozbargain deals
		k.processOzbargainDeals()

		// Process Camel camel camel (Amazon) deals
		k.processCCCDeals()
	}
}

func (k *KramerBot) processOzbargainDeals() {
	err := k.OzbScraper.Scrape()
	if err != nil {
		k.Logger.Error("Error scraping deals", zap.Error(err))
		return
	}

	// Load deals from OzBargain
	deals := k.OzbScraper.GetData()
	userdata := k.UserStore.Users

	for _, deal := range deals {
		k.Logger.Debug("Ozbargain deal found", zap.Any("deal", deal))

		// Check deal type
		dealType := k.OzbScraper.GetDealType(deal)

		// Go through all registered users and check deals they are subscribed to
		for _, user := range userdata {
			if user.OzbGood && dealType == int(scrapers.OZB_GOOD) && !OzbDealSent(user, &deal) {
				// User is subscribed to good deals, notify user
				k.SendOzbGoodDeal(user, &deal)
			}

			if user.OzbSuper && dealType == int(scrapers.OZB_SUPER) && !OzbDealSent(user, &deal) {
				// User is subscribed to super deals, notify user
				k.SendOzbSuperDeal(user, &deal)
			}

			// Check for watched keywords
			for _, keyword := range user.Keywords {
				if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) && !OzbDealSent(user, &deal) {
					// Deal contains keyword, notify user
					k.SendWatchedDeal(user, &deal)

					// Break out of keyword loop
					break
				}
			}
		}
	}
}

func (k *KramerBot) processCCCDeals() {
	err := k.CCCScraper.Scrape()
	if err != nil {
		k.Logger.Error("Error scraping deals", zap.Error(err))
		return
	}

	// Load deals from OzBargain
	deals := k.CCCScraper.GetData()
	userdata := k.UserStore.Users

	for _, deal := range deals {
		k.Logger.Debug("Amazon deal found", zap.Any("deal", deal))

		// Go through all registered users and check deals they are subscribed to
		for _, user := range userdata {
			if user.AmzDaily && !AmzDealSent(user, &deal) {
				// User is subscribed to AMZ daily deals, notify user
				k.SendAmzDailyDeal(user, &deal)
			}

			if user.AmzWeekly && !AmzDealSent(user, &deal) {
				// User is subscribed to AMZ weekly deals, notify user
				k.SendAmzWeeklyDeal(user, &deal)
			}

			// Check for watched keywords
			for _, keyword := range user.Keywords {
				if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) && !AmzDealSent(user, &deal) {
					// Deal contains keyword, notify user
					// k.SendWatchedDeal(user, &deal)

					// Break out of keyword loop
					break
				}
			}
		}
	}
}
