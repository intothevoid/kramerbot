package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
)

// Process deals returned by the scraper, check deal type and notify user
// if they are subscribed to a particular deal type
func (k *KramerBot) StartProcessing() {
	// Load user store i.e. user data indexed by chat id
	k.LoadUserStore()

	// Begin timed processing and scraping
	// tick := time.NewTicker(time.Second * 60)
	tick := time.NewTicker(time.Minute * PROCESSING_INTERVAL)
	for range tick.C {
		// Load deals from OzBargain
		k.Scraper.Scrape()
		deals := k.Scraper.GetData()
		userdata := k.UserStore.Users

		for _, deal := range deals {
			// Check deal type
			dealType := k.Scraper.GetDealType(deal)

			// Go through all registered users and check deals they are subscribed to
			for _, user := range userdata {
				if user.GoodDeals && dealType == int(scrapers.GOOD_DEAL) && !DealSent(user, &deal) {
					// User is subscribed to good deals, notify user
					shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
					formattedDeal := fmt.Sprintf(`🔥<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)

					k.Logger.Debug(fmt.Sprintf("Sending deal %s to user %s", shortenedTitle, user.Username))
					k.SendHTMLMessage(user.ChatID, formattedDeal)

					// Mark deal as sent
					user.DealsSent = append(user.DealsSent, deal.Id)
					k.SaveUserStore()
				}
				if user.SuperDeals && dealType == int(scrapers.SUPER_DEAL) && !DealSent(user, &deal) {
					// User is subscribed to good deals, notify user
					shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
					formattedDeal := fmt.Sprintf(`🔥🔥<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)

					k.Logger.Debug(fmt.Sprintf("Sending deal %s to user %s", shortenedTitle, user.Username))
					k.SendHTMLMessage(user.ChatID, formattedDeal)

					// Mark deal as sent
					user.DealsSent = append(user.DealsSent, deal.Id)
					k.SaveUserStore()
				}

				// Check for watched keywords
				for _, keyword := range user.Keywords {
					if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) && !DealSent(user, &deal) {
						// Deal contains keyword, notify user
						shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
						formattedDeal := fmt.Sprintf(`👀<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)

						k.Logger.Debug(fmt.Sprintf("Sending deal %s to user %s", shortenedTitle, user.Username))
						k.SendHTMLMessage(user.ChatID, formattedDeal)

						// Mark deal as sent
						user.DealsSent = append(user.DealsSent, deal.Id)
						k.SaveUserStore()

						// Break out of keyword loop
						break
					}
				}
			}
		}
	}
}
