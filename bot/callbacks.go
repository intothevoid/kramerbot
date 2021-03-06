package bot

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

// Function to send latest deals i.e. NUM_DEALS_TO_SEND
func (k *KramerBot) SendLatestDeals(chatID int64, s *scrapers.OzBargainScraper) {
	latestDeals := s.GetLatestDeals(scrapers.NUM_DEALS_TO_SEND)

	// Send latest deals to the user
	for _, deal := range latestDeals {
		shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
		formattedDeal := fmt.Sprintf("π<a href='%s' target='_blank'>%s</a>πΊ%s", deal.Url, shortenedTitle, deal.Upvotes)

		k.SendHTMLMessage(chatID, formattedDeal)

		// Delay for a bit don't send all deals at once
		time.Sleep(1 * time.Second)
	}
}

// Function to display help message
func (k *KramerBot) Help(chat *tgbotapi.Chat) {
	// Send kramer's photo
	fpath, _ := filepath.Abs("./static/kramer_drnostrand.jpg")
	k.SendPhoto(chat.ID, fpath)

	// Show the help banner
	k.SendMessage(chat.ID, fmt.Sprintf("Hi %s! Available commands are: \n\n"+
		"π /help - View this help message \n\n"+
		"π /latest - View the 5 latest deals on OzBargain\n\n"+
		"π  /watchgood - Ozbargain: Watch out for deals with 25+ upvotes within the hour\n\n"+
		"π  /watchsuper - Ozbargain: Watch out for deals with 100+ upvotes within 24 hours\n\n"+
		"π°οΈ /amazondaily - Amazon: Watch out for top daily Amazon deals with price drops greater than 20 percent\n\n"+
		"π°οΈ /amazonweekly - Amazon: Watch out for top weekly Amazon deals with price drops greater than 20 percent\n\n"+
		"π /watchkeyword - Watch deals with specified keywords across π Ozbargain and π°οΈAmazon\n\n"+
		"β /clearkeyword - Clear deals with specified keyword\n\n"+
		"β /clearallkeywords - Clear deals with all watched keywords\n\n"+
		"π¨βπ¦° /status - Get the current user status\n\n"+
		"π /kramerism - Get a Kramer quote from Seinfeld", chat.FirstName))
}

// Send test message
func (k *KramerBot) SendTestMessage(chat *tgbotapi.Chat) {

	shortenedTitle := util.ShortenString("π₯ This is a test deal not a real deal... Beep Boop", 30) + "..."
	dealUrl := "https://news.google.com.au"
	formattedDeal := fmt.Sprintf(`π₯<a href='%s' target='_blank'>%s</a>`, dealUrl, shortenedTitle)

	k.Logger.Debug(fmt.Sprintf("Sending deal %s to user %s", shortenedTitle, chat.FirstName))
	k.SendHTMLMessage(chat.ID, formattedDeal)
}

// Process keyword watch request
func (k *KramerBot) ProcessKeyword(chat *tgbotapi.Chat, keyword string) {
	var keywords []string

	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]
		userData.Keywords = append(userData.Keywords, keyword)

		// For messaging the user
		keywords = userData.Keywords
	} else {
		// Key does not exist, create new user data
		userData := k.CreateUserData(chat.ID, chat.FirstName, keyword, false, false, false, false)
		k.UserStore.Users[chat.ID] = userData

		// For messaging the user
		keywords = userData.Keywords
	}

	// Save user store
	k.SaveUserStore()

	k.SendMessage(chat.ID, fmt.Sprintf("π Currently watching keywords: %s for user %s", keywords, chat.FirstName))
}

// Process clear keyword request
func (k *KramerBot) ProcessClearKeyword(chat *tgbotapi.Chat, keyword string) {
	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]

		// Delete keyword from userData
		for i, v := range userData.Keywords {
			if v == keyword {
				userData.Keywords = append(userData.Keywords[:i], userData.Keywords[i+1:]...)
			}
		}
	} else {
		// User does not exist, nothing to clear
		k.SendMessage(chat.ID, fmt.Sprintf("User data for %s not found. Nothing to clear", chat.FirstName))
		return
	}

	// Save user store
	k.SaveUserStore()

	k.SendMessage(chat.ID, fmt.Sprintf("π Cleared watched keyword: %s for user %s", keyword, chat.FirstName))
}

// Process clear all keywords request
func (k *KramerBot) ProcessClearAllKeywords(chat *tgbotapi.Chat) {
	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]

		// Delete keyword from userData
		userData.Keywords = []string{}
	} else {
		// User does not exist, nothing to clear
		k.SendMessage(chat.ID, fmt.Sprintf("User data for %s not found. Nothing to clear", chat.FirstName))
		return
	}

	// Save user store
	k.SaveUserStore()

	k.SendMessage(chat.ID, fmt.Sprintf("π Cleared all watched keywords for user %s", chat.FirstName))
}

// Add watch to OZB good deals by chat id
func (k *KramerBot) WatchOzbGoodDeals(chat *tgbotapi.Chat) {
	k.watchDeal(chat, scrapers.OZB_GOOD)
}

// Add watch to OZB super deals by chat id
func (k *KramerBot) WatchOzbSuperDeals(chat *tgbotapi.Chat) {
	k.watchDeal(chat, scrapers.OZB_SUPER)
}

// Add watch to AMZ daily deals by chat id
func (k *KramerBot) WatchAmzDailyDeals(chat *tgbotapi.Chat) {
	k.watchDeal(chat, scrapers.AMZ_DAILY)
}

// Add watch to AMZ weekly deals by chat id
func (k *KramerBot) WatchAmzWeeklyDeals(chat *tgbotapi.Chat) {
	k.watchDeal(chat, scrapers.AMZ_WEEKLY)
}

// Helper function to watch deal
func (k *KramerBot) watchDeal(chat *tgbotapi.Chat, dealType scrapers.DealType) {
	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]
		var message string
		var added bool

		switch dealType {
		case scrapers.OZB_GOOD:
			userData.OzbGood = !userData.OzbGood // toggle
			added = userData.OzbGood
			message = " π π₯ ozbargain good deals list."
		case scrapers.OZB_SUPER:
			userData.OzbSuper = !userData.OzbSuper
			added = userData.OzbSuper
			message = " π π₯ ozbargain super deals"
		case scrapers.AMZ_DAILY:
			userData.AmzDaily = !userData.AmzDaily
			added = userData.AmzDaily
			message = " π°οΈ amazon daily deals list."
		case scrapers.AMZ_WEEKLY:
			userData.AmzWeekly = !userData.AmzWeekly
			added = userData.AmzWeekly
			message = " π°οΈ amazon weekly deals list."
		default:
			k.Logger.Error("Invalid deal type passed in", zap.Any("dealtype", dealType))
			k.SendMessage(chat.ID, "There was an error adding / deleting you from the list.")
			return
		}

		// Send message to user
		if added {
			k.SendMessage(chat.ID, "You have been added to the "+message)
		} else {
			k.SendMessage(chat.ID, "You have been removed from the"+message)
		}
	} else {
		var message string
		var ozbGood bool
		var ozbSuper bool
		var amzDaily bool
		var amzWeekly bool

		switch dealType {
		case scrapers.OZB_GOOD:
			ozbGood = true
			message = " π π₯ ozbargain good deals list."
		case scrapers.OZB_SUPER:
			ozbSuper = true
			message = " π π₯ ozbargain super deals"
		case scrapers.AMZ_DAILY:
			amzDaily = true
			message = " π°οΈ amazon daily deals list."
		case scrapers.AMZ_WEEKLY:
			amzWeekly = true
			message = " π°οΈ amazon weekly deals list."
		default:
			k.Logger.Error("Invalid deal type passed in", zap.Any("dealtype", dealType))
			k.SendMessage(chat.ID, "There was an error adding / deleting you from the list.")
			return
		}

		// Key does not exist, create new user
		userData := k.CreateUserData(chat.ID, chat.FirstName, "", ozbGood, ozbSuper, amzDaily, amzWeekly)
		k.UserStore.Users[chat.ID] = userData

		// Send message to user
		k.SendMessage(chat.ID, "You have been added to the "+message)
	}

	// Save user store
	k.SaveUserStore()
}

// Send OZB good deal message to user
func (k *KramerBot) SendOzbGoodDeal(user *models.UserData, deal *models.OzBargainDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`π π₯<a href="%s" target="_blank">%s</a>πΊ%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`π π₯ %s πΊ%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending good deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	k.SaveUserStore()
}

// Send user their current configured settings / status
func (k *KramerBot) SendStatus(chat *tgbotapi.Chat) {
	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		user := k.UserStore.Users[chat.ID]
		getTruth := func(set bool) string {
			if set {
				return "yes"
			}
			return "no"
		}
		prettyPrint := func(words []string) string {
			var retval string
			for _, word := range words {
				retval += word + "\n"
			}
			return retval
		}
		userDetails := fmt.Sprintf("π¨βπ¦°π©βπ¦° %s\n\nπ OZB Good Deals: %s\nπ OZB Super Deals: %s\nπ°οΈAmazon Top Daily Deals: %s\nπ°οΈAmazon Top Weekly Deals: %s\nπWatched Deals:\n %sβ°OZB Deals sent: %d\nβ°AMZ Deals sent: %d", user.GetUsername(),
			getTruth(user.GetOzbGood()), getTruth(user.GetOzbSuper()), getTruth(user.GetAmzDaily()),
			getTruth(user.GetAmzWeekly()), prettyPrint(user.GetKeywords()), len(user.GetOzbSent()),
			len(user.GetAmzSent()))

		k.SendHTMLMessage(user.ChatID, userDetails)
	} else {
		k.SendHTMLMessage(chat.ID, "This is embarassing. I could not find your details.")
	}
}

// Send OZB super deal to user
func (k *KramerBot) SendOzbSuperDeal(user *models.UserData, deal *models.OzBargainDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`π π₯<a href="%s" target="_blank">%s</a>πΊ%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`π π₯ %s πΊ%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending super deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	k.SaveUserStore()
}

func (k *KramerBot) SendAmzDeal(user *models.UserData, deal *models.CamCamCamDeal) {
	dealType := ""

	// Get deal type
	if deal.DealType == int(scrapers.AMZ_DAILY) {
		dealType = "top daily deal"
	}
	if deal.DealType == int(scrapers.AMZ_WEEKLY) {
		dealType = "top weekly deal"
	}

	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`π°οΈ<a href="%s" target="_blank">%s</a> - %s`, deal.Url, shortenedTitle, k.CCCScraper.GetDealDropString(deal))
	textDeal := fmt.Sprintf(`π°οΈ %s`, shortenedTitle)

	k.Logger.Debug(fmt.Sprintf("Sending Amazon %s deal %s to user %s", dealType, shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.AmzSent = append(user.AmzSent, deal.Id)
	k.SaveUserStore()
}

// Send OZB watched deal to user
func (k *KramerBot) SendOzbWatchedDeal(user *models.UserData, deal *models.OzBargainDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`π π<a href="%s" target="_blank">%s</a>πΊ%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`π π %s πΊ%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending watched Ozbargain deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	k.SaveUserStore()
}

// Send AMZ watched deal to user
func (k *KramerBot) SendAmzWatchedDeal(user *models.UserData, deal *models.CamCamCamDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`π°οΈπ<a href="%s" target="_blank">%s</a> - %s`, deal.Url, shortenedTitle, k.CCCScraper.GetDealDropString(deal))
	textDeal := fmt.Sprintf(`π°οΈπ %s`, shortenedTitle)

	k.Logger.Debug(fmt.Sprintf("Sending watched Amazon deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.AmzSent = append(user.AmzSent, deal.Id)
	k.SaveUserStore()
}
