package bot

import (
	"strings"
	"time"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"go.uber.org/zap"
)

// Process deals returned by the scraper, check deal type and notify user
// if they are subscribed to a particular deal type
func (k *KramerBot) StartProcessing() {
	// Begin timed processing and scraping
	go func() {
		ozbTick := time.NewTicker(time.Minute * time.Duration(k.OzbScraper.ScrapeInterval))
		for range ozbTick.C {
			// Process Ozbargain deals
			k.processOzbargainDeals()
		}
	}()

	go func() {
		// amzTick := time.NewTicker(time.Second * 60)
		amzTick := time.NewTicker(time.Minute * time.Duration(k.CCCScraper.ScrapeInterval))
		for range amzTick.C {
			// Process Camel camel camel (Amazon) deals
			k.processCCCDeals()
		}
	}()
}

func (k *KramerBot) processOzbargainDeals() {
	// Add nil checks for k.OzbScraper
	if k.OzbScraper == nil {
		k.Logger.Error("OzbScraper is nil")
		return
	}

	err := k.OzbScraper.Scrape()
	if err != nil {
		k.Logger.Error("Error scraping deals", zap.Error(err))
		return
	}

	// Load deals from OzBargain
	deals := k.OzbScraper.GetData()
	if deals == nil {
		k.Logger.Error("No deals returned from scraper")
	}

	// Strip duplicates by using a map indexed by deal id
	uniqueDeals := make(map[string]models.OzBargainDeal)
	for _, deal := range deals {
		uniqueDeals[deal.Id] = deal
	}

	// Load store
	k.LoadUserStore()

	var userdata map[int64]*models.UserData
	if k.UserStore != nil {
		userdata = k.UserStore.Users
	} else {
		userdata = nil
		k.Logger.Error("No users found in UserStore")
		return
	}

	for _, deal := range uniqueDeals {
		k.Logger.Debug("Ozbargain deal", zap.Any("deal", deal))

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
				// If keyword is empty or only contains spaces
				keyword = strings.TrimSpace(keyword)
				if len(keyword) == 0 {
					continue
				}

				if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) && !OzbDealSent(user, &deal) {
					// Deal contains keyword, notify user
					k.SendOzbWatchedDeal(user, &deal)

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

	// Strip duplicates by using a map indexed by deal id
	uniqueDeals := make(map[string]models.CamCamCamDeal)
	for _, deal := range deals {
		uniqueDeals[deal.Id] = deal
	}

	// Load store
	k.LoadUserStore()
	userdata := k.UserStore.Users

	// Get price drop target from configuration
	priceDropTarget := k.Config.GetInt("scrapers.amazon.target_price_drop")

	for _, deal := range uniqueDeals {
		k.Logger.Debug("Amazon deal", zap.Any("deal", deal))

		// Check if percentage drop meets target
		priceDropTargetMet := k.CCCScraper.IsTargetDropGreater(&deal, priceDropTarget)

		// Go through all registered users and check deals they are subscribed to
		for _, user := range userdata {
			if user.AmzDaily && priceDropTargetMet && deal.DealType == int(scrapers.AMZ_DAILY) && !AmzDealSent(user, &deal) {
				// User is subscribed to AMZ daily deals, notify user
				k.SendAmzDeal(user, &deal)
			}

			if user.AmzWeekly && priceDropTargetMet && deal.DealType == int(scrapers.AMZ_WEEKLY) && !AmzDealSent(user, &deal) {
				// User is subscribed to AMZ weekly deals, notify user
				k.SendAmzDeal(user, &deal)
			}

			// Check for watched keywords
			for _, keyword := range user.Keywords {
				// If keyword is empty or only contains spaces
				keyword = strings.TrimSpace(keyword)
				if len(keyword) == 0 {
					continue
				}

				if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) && !AmzDealSent(user, &deal) {
					// Deal contains keyword, notify user
					k.SendAmzWatchedDeal(user, &deal)

					// Break out of keyword loop
					break
				}
			}
		}
	}
}
