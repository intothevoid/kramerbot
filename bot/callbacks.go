package bot

import (
	"fmt"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
)

// Function to send latest deals i.e. NUM_DEALS_TO_SEND
func (k *KramerBot) SendLatestDeals(chatID int64, s *scrapers.OzBargainScraper) {
	latestDeals := s.GetLatestDeals(scrapers.NUM_DEALS_TO_SEND)

	// Send latest deals to the user
	for _, deal := range latestDeals {
		shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
		formattedDeal := fmt.Sprintf("🆕<a href='%s' target='_blank'>%s</a>🔺%s", deal.Url, shortenedTitle, deal.Upvotes)

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
		"🙏 /help - View this help message \n"+
		"📈 /latest - View the 5 latest deals on OzBargain\n"+
		"🔥 /watchgood - Watch out for deals with 25+ upvotes within the hour\n"+
		"🔥🔥 /watchsuper - Watch out for deals with 100+ upvotes within 24 hours\n"+
		"👀 /watchkeyword - Watch deals with specified keywords\n"+
		"⛔ /clearkeyword - Clear deals with specified keyword\n"+
		"⛔ /clearallkeywords - Clear deals with all watched keywords\n"+
		"🙃 /kramerism - Get a Kramer quote from Seinfeld", chat.FirstName))
}

// Send test message
func (k *KramerBot) SendTestMessage(chat *tgbotapi.Chat) {

	shortenedTitle := util.ShortenString("This is a test deal not a real deal... Beep Boop", 30) + "..."
	dealUrl := "https://news.google.com.au"
	formattedDeal := fmt.Sprintf(`<a href='%s' target='_blank'>%s</a>`, dealUrl, shortenedTitle)

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
		userData := k.CreateUserData(chat.ID, chat.FirstName, keyword, false, false)
		k.UserStore.Users[chat.ID] = userData

		// For messaging the user
		keywords = userData.Keywords
	}

	// Save user store
	k.SaveUserStore()

	k.SendMessage(chat.ID, fmt.Sprintf("Currently watching keywords: %s for user %s", keywords, chat.FirstName))
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

	k.SendMessage(chat.ID, fmt.Sprintf("Cleared watched keyword: %s for user %s", keyword, chat.FirstName))
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

	k.SendMessage(chat.ID, fmt.Sprintf("Cleared all watched keywords for user %s", chat.FirstName))
}

// Add watch to good deals by chat id
func (k *KramerBot) WatchGoodDeals(chat *tgbotapi.Chat) {
	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]
		userData.GoodDeals = !userData.GoodDeals // toggle

		// Send message to user
		if userData.GoodDeals {
			k.SendMessage(chat.ID, "You have been added to the good deals watchlist.")
		} else {
			k.SendMessage(chat.ID, "You have been removed from the good deals watchlist.")
		}
	} else {
		// Key does not exist, create new user
		userData := k.CreateUserData(chat.ID, chat.FirstName, "", true, false)
		k.UserStore.Users[chat.ID] = userData

		// Send message to user
		k.SendMessage(chat.ID, "You have been added to the good deals watchlist.")
	}

	// Save user store
	k.SaveUserStore()
}

// Add watch to super deals by chat id
func (k *KramerBot) WatchSuperDeals(chat *tgbotapi.Chat) {

	// Check if key exists in user store
	if _, ok := k.UserStore.Users[chat.ID]; ok {
		// Key exists, add to watch list
		userData := k.UserStore.Users[chat.ID]
		userData.SuperDeals = !userData.SuperDeals // toggle

		// Send message to user
		if userData.SuperDeals {
			k.SendMessage(chat.ID, "You have been added to the super deals watchlist.")
		} else {
			k.SendMessage(chat.ID, "You have been removed from the super deals watchlist.")
		}
	} else {
		// Key does not exist, create new user
		userData := k.CreateUserData(chat.ID, chat.FirstName, "", false, true)
		k.UserStore.Users[chat.ID] = userData
		k.SendMessage(chat.ID, "You have been added to the super deals watchlist.")
	}

	// Save user store
	k.SaveUserStore()
}
