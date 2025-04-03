package bot

import (
	"fmt"
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
			if err := k.processOzbargainDeals(); err != nil {
				k.Logger.Error("Error processing Ozbargain deals", zap.Error(err))
			}
		}
	}()

	go func() {
		amzTick := time.NewTicker(time.Minute * time.Duration(k.CCCScraper.ScrapeInterval))
		for range amzTick.C {
			if err := k.processCCCDeals(); err != nil {
				k.Logger.Error("Error processing CCC deals", zap.Error(err))
			}
		}
	}()
}

func (k *KramerBot) processOzbargainDeals() error {
	// Add nil checks for k.OzbScraper
	if k.OzbScraper == nil {
		return fmt.Errorf("OzbScraper is nil")
	}

	err := k.OzbScraper.Scrape()
	if err != nil {
		return fmt.Errorf("error scraping deals: %w", err)
	}

	// Load deals from OzBargain
	deals := k.OzbScraper.GetData()
	if deals == nil {
		return fmt.Errorf("no deals returned from scraper")
	}

	// Strip duplicates by using a map indexed by deal id
	uniqueDeals := make(map[string]models.OzBargainDeal)
	for i := range deals {
		deal := deals[i]
		if deal.Id == "" {
			k.Logger.Warn("Skipping deal with empty ID")
			continue
		}
		uniqueDeals[deal.Id] = deal
	}

	// Load store
	if err := k.LoadUserStore(); err != nil {
		return fmt.Errorf("error loading user store: %w", err)
	}

	// Get a thread-safe copy of all users
	userdata := k.UserStore.GetAllUsers()
	if userdata == nil {
		return fmt.Errorf("no users found in UserStore")
	}

	// Pre-process user keywords into maps for efficient lookups
	userKeywordMaps := make(map[int64]map[string]bool)
	for chatID, user := range userdata {
		if user == nil {
			k.Logger.Warn("Skipping nil user", zap.Int64("chat_id", chatID))
			continue
		}
		keywordMap := make(map[string]bool)
		for _, keyword := range user.Keywords {
			// Only add non-empty keywords
			if trimmed := strings.TrimSpace(keyword); trimmed != "" {
				keywordMap[strings.ToLower(trimmed)] = true
			}
		}
		userKeywordMaps[chatID] = keywordMap
	}

	for _, deal := range uniqueDeals {
		k.Logger.Debug("Ozbargain deal", zap.Any("deal", deal))

		// Pre-process deal title once
		dealTitleLower := strings.ToLower(deal.Title)

		// Check deal type
		dealType := k.OzbScraper.GetDealType(deal)

		// Go through all registered users and check deals they are subscribed to
		for chatID, user := range userdata {
			if user == nil {
				k.Logger.Warn("Skipping nil user", zap.Int64("chat_id", chatID))
				continue
			}

			// Check deal type subscriptions
			if user.OzbGood && dealType == int(scrapers.OZB_GOOD) && !OzbDealSent(user, &deal) {
				// User is subscribed to good deals, notify user
				if err := k.SendOzbGoodDeal(user, &deal); err != nil {
					k.Logger.Error("Failed to send OZB good deal",
						zap.String("deal_id", deal.Id),
						zap.Int64("user_id", user.ChatID),
						zap.Error(err))
				}
			}

			if user.OzbSuper && dealType == int(scrapers.OZB_SUPER) && !OzbDealSent(user, &deal) {
				// User is subscribed to super deals, notify user
				if err := k.SendOzbSuperDeal(user, &deal); err != nil {
					k.Logger.Error("Failed to send OZB super deal",
						zap.String("deal_id", deal.Id),
						zap.Int64("user_id", user.ChatID),
						zap.Error(err))
				}
			}

			// Check for watched keywords using pre-processed map
			keywordMap := userKeywordMaps[chatID]
			if len(keywordMap) > 0 {
				// Check if any keyword is in the deal title
				for keyword := range keywordMap {
					if strings.Contains(dealTitleLower, keyword) && !OzbDealSent(user, &deal) {
						// Deal contains keyword, notify user
						if err := k.SendOzbWatchedDeal(user, &deal); err != nil {
							k.Logger.Error("Failed to send OZB watched deal",
								zap.String("deal_id", deal.Id),
								zap.Int64("user_id", user.ChatID),
								zap.String("keyword", keyword),
								zap.Error(err))
						}
						break // Break after first match
					}
				}
			}
		}
	}
	return nil
}

func (k *KramerBot) processCCCDeals() error {
	if k.CCCScraper == nil {
		return fmt.Errorf("CCCScraper is nil")
	}

	err := k.CCCScraper.Scrape()
	if err != nil {
		return fmt.Errorf("error scraping deals: %w", err)
	}

	// Load deals from OzBargain
	deals := k.CCCScraper.GetData()
	if deals == nil {
		return fmt.Errorf("no deals returned from scraper")
	}

	// Strip duplicates by using a map indexed by deal id
	uniqueDeals := make(map[string]models.CamCamCamDeal)
	for _, deal := range deals {
		uniqueDeals[deal.Id] = deal
	}

	// Load store
	if err := k.LoadUserStore(); err != nil {
		return fmt.Errorf("error loading user store: %w", err)
	}

	// Get a thread-safe copy of all users
	userdata := k.UserStore.GetAllUsers()
	if userdata == nil {
		return fmt.Errorf("no users found in UserStore")
	}

	// Pre-process user keywords into maps for efficient lookups
	userKeywordMaps := make(map[int64]map[string]bool)
	for chatID, user := range userdata {
		keywordMap := make(map[string]bool)
		for _, keyword := range user.Keywords {
			// Only add non-empty keywords
			if trimmed := strings.TrimSpace(keyword); trimmed != "" {
				keywordMap[strings.ToLower(trimmed)] = true
			}
		}
		userKeywordMaps[chatID] = keywordMap
	}

	// Get price drop target from configuration
	priceDropTarget := k.Config.Scrapers.Amazon.TargetPriceDrop

	for _, deal := range uniqueDeals {
		k.Logger.Debug("Amazon deal", zap.Any("deal", deal))

		// Pre-process deal title once
		dealTitleLower := strings.ToLower(deal.Title)

		// Check if percentage drop meets target
		priceDropTargetMet := k.CCCScraper.IsTargetDropGreater(&deal, priceDropTarget)

		// Go through all registered users and check deals they are subscribed to
		for chatID, user := range userdata {
			if user == nil {
				k.Logger.Warn("Skipping nil user", zap.Int64("chat_id", chatID))
				continue
			}

			if user.AmzDaily && priceDropTargetMet && deal.DealType == int(scrapers.AMZ_DAILY) && !AmzDealSent(user, &deal) {
				// User is subscribed to AMZ daily deals, notify user
				if err := k.SendAmzDeal(user, &deal); err != nil {
					k.Logger.Error("Failed to send AMZ daily deal",
						zap.String("deal_id", deal.Id),
						zap.Int64("user_id", user.ChatID),
						zap.Error(err))
				}
			}

			if user.AmzWeekly && priceDropTargetMet && deal.DealType == int(scrapers.AMZ_WEEKLY) && !AmzDealSent(user, &deal) {
				// User is subscribed to AMZ weekly deals, notify user
				if err := k.SendAmzDeal(user, &deal); err != nil {
					k.Logger.Error("Failed to send AMZ weekly deal",
						zap.String("deal_id", deal.Id),
						zap.Int64("user_id", user.ChatID),
						zap.Error(err))
				}
			}

			// Check for watched keywords using pre-processed map
			keywordMap := userKeywordMaps[chatID]
			if len(keywordMap) > 0 {
				// Check if any keyword is in the deal title
				for keyword := range keywordMap {
					if strings.Contains(dealTitleLower, keyword) && !AmzDealSent(user, &deal) {
						// Deal contains keyword, notify user
						if err := k.SendAmzWatchedDeal(user, &deal); err != nil {
							k.Logger.Error("Failed to send AMZ watched deal",
								zap.String("deal_id", deal.Id),
								zap.Int64("user_id", user.ChatID),
								zap.String("keyword", keyword),
								zap.Error(err))
						}
						break // Break after first match
					}
				}
			}
		}
	}
	return nil
}
