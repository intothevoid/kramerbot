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
		// Check deal type
		dealType := k.OzbScraper.GetDealType(deal)

		// Go through all registered users and check deals they are subscribed to
		for _, user := range userdata {
			if user.GoodDeals && dealType == int(scrapers.GOOD_DEAL) && !DealSent(user, &deal) {
				// User is subscribed to good deals, notify user
				k.SendGoodDeal(user, &deal)
			}

			if user.SuperDeals && dealType == int(scrapers.SUPER_DEAL) && !DealSent(user, &deal) {
				// User is subscribed to super deals, notify user
				k.SendSuperDeal(user, &deal)
			}

			// Check for watched keywords
			for _, keyword := range user.Keywords {
				if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) && !DealSent(user, &deal) {
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
	// userdata := k.UserStore.Users

	for _, deal := range deals {
		k.Logger.Debug("Amazon deal found", zap.Any("deal", deal))

		// // Go through all registered users and check deals they are subscribed to
		// for _, user := range userdata {
		// 	if user.GoodDeals && dealType == int(scrapers.GOOD_DEAL) && !DealSent(user, &deal) {
		// 		// User is subscribed to good deals, notify user
		// 		k.SendGoodDeal(user, &deal)
		// 	}

		// 	if user.SuperDeals && dealType == int(scrapers.SUPER_DEAL) && !DealSent(user, &deal) {
		// 		// User is subscribed to super deals, notify user
		// 		k.SendSuperDeal(user, &deal)
		// 	}

		// 	// Check for watched keywords
		// 	for _, keyword := range user.Keywords {
		// 		if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) && !DealSent(user, &deal) {
		// 			// Deal contains keyword, notify user
		// 			k.SendWatchedDeal(user, &deal)

		// 			// Break out of keyword loop
		// 			break
		// 		}
		// 	}
		// }
	}
}
