package bot

import (
	"fmt"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
)

// Function to display help message
func (k *KramerBot) Help(chat *tgbotapi.Chat) {
	// Send kramer's photo
	fpath, _ := filepath.Abs("./static/kramer_drnostrand.jpg")
	k.SendPhoto(chat.ID, fpath)

	// Show the help banner
	helpText := fmt.Sprintf("Hi %s! Welcome to @kramerbot\n\n"+
		"Your ChatID is %d\n\n"+
		"Available commands:\n"+
		"/start - Register or view your status\n"+
		"/help - Show this help message\n"+
		"/preferences - Show your current notification preferences\n"+
		"/toggle_ozbgood - Toggle OzBargain 'Good' deals (25+ votes)\n"+
		"/toggle_ozbsuper - Toggle OzBargain 'Super' deals (50+ votes)\n"+
		"/toggle_amzdaily - Toggle Amazon Daily deals\n"+
		"/toggle_amzweekly - Toggle Amazon Weekly deals\n"+
		"/addkeyword <keyword> - Add a keyword to watch\n"+
		"/removekeyword <keyword> - Remove a keyword\n"+
		"/listkeywords - List your watched keywords\n"+
		"/test - Send a test notification",
		chat.FirstName, chat.ID)
	k.SendMessage(chat.ID, helpText)
}

// RegisterUser adds a new user or shows status for an existing user
func (k *KramerBot) RegisterUser(chat *tgbotapi.Chat) {
	// Check if user exists
	user, err := k.DataWriter.GetUser(chat.ID)
	if err != nil || user == nil {
		// User doesn't exist, create a new one
		newUser := &models.UserData{
			ChatID:   chat.ID,
			Username: chat.UserName,
			OzbGood:  true, // Default settings
			OzbSuper: false,
			AmzDaily: false,
			AmzWeekly: false,
			Keywords: []string{},
			OzbSent:  []string{},
			AmzSent:  []string{},
		}
		err := k.DataWriter.AddUser(newUser)
		if err != nil {
			k.Logger.Error("Failed to add new user", zap.Int64("chatID", chat.ID), zap.Error(err))
			k.SendMessage(chat.ID, "Sorry, there was an error registering you. Please try again later.")
			return
		}
		// Add to in-memory store as well
		k.UserStore.Users[chat.ID] = newUser
		k.Logger.Info("Registered new user", zap.String("username", chat.UserName), zap.Int64("chatID", chat.ID))
		k.SendMessage(chat.ID, fmt.Sprintf("Welcome %s! You are now registered. Use /help to see available commands.", chat.FirstName))
		k.ShowPreferences(chat) // Show current (default) preferences
	} else {
		// User exists, update username if changed and show status
		if user.Username != chat.UserName {
			user.Username = chat.UserName
			k.UpdateUser(user) // Update in DB
			k.UserStore.Users[chat.ID].Username = chat.UserName // Update in memory
			k.Logger.Info("Updated username for existing user", zap.String("username", chat.UserName), zap.Int64("chatID", chat.ID))
		}
		k.Logger.Info("User already registered", zap.String("username", chat.UserName), zap.Int64("chatID", chat.ID))
		k.SendMessage(chat.ID, fmt.Sprintf("Welcome back %s!", chat.FirstName))
		k.ShowPreferences(chat) // Show current preferences
	}
}

// ShowPreferences displays the user's current notification settings
func (k *KramerBot) ShowPreferences(chat *tgbotapi.Chat) {
	user, err := k.DataWriter.GetUser(chat.ID)
	if err != nil || user == nil {
		k.SendMessage(chat.ID, "Could not find your user data. Have you registered using /start ?")
		return
	}

	prefsText := fmt.Sprintf("Your current preferences:\n"+
		"Ozbargain Good Deals (25+): %t\n"+
		"Ozbargain Super Deals (50+): %t\n"+
		"Amazon Daily Deals: %t\n"+
		"Amazon Weekly Deals: %t\n"+
		"Watched Keywords: %d",
		user.OzbGood, user.OzbSuper, user.AmzDaily, user.AmzWeekly, len(user.Keywords))

	k.SendMessage(chat.ID, prefsText)
	k.ListKeywords(chat) // Also list the keywords
}

// ListKeywords displays the user's watched keywords
func (k *KramerBot) ListKeywords(chat *tgbotapi.Chat) {
	user, err := k.DataWriter.GetUser(chat.ID)
	if err != nil || user == nil {
		k.SendMessage(chat.ID, "Could not find your user data. Have you registered using /start ?")
		return
	}

	if len(user.Keywords) == 0 {
		k.SendMessage(chat.ID, "You are not watching any keywords.")
	} else {
		keywordsText := "Your watched keywords:\n- " + strings.Join(user.Keywords, "\n- ")
		k.SendMessage(chat.ID, keywordsText)
	}
}


// Send test message
func (k *KramerBot) SendTestMessage(chat *tgbotapi.Chat) {

	shortenedTitle := util.ShortenString("🔥 This is a test deal not a real deal... Beep Boop", 30) + "..."
	dealUrl := "https://news.google.com.au"
	formattedDeal := fmt.Sprintf(`🔥<a href='%s' target='_blank'>%s</a>`, dealUrl, shortenedTitle)

	k.Logger.Debug(fmt.Sprintf("Sending deal %s to user %s", shortenedTitle, chat.FirstName))
	k.SendHTMLMessage(chat.ID, formattedDeal)
}

// Make an announcement to all users i.e. important messages, updates etc.
// Note: This is an admin function and will need KRAMERBOT_ADMIN_PASSWORD
func (k *KramerBot) MakeAnnouncement(chat *tgbotapi.Chat, announcement string) {
	// Extract message
	messages := strings.Split(announcement, ":")
	var message string
	if len(messages) == 2 {
		message = messages[1]
	}

	formattedAnnouncement := fmt.Sprintf(`📢 Kramerbot Announcement 📢 %s`, message)

	for _, user := range k.UserStore.Users {
		k.Logger.Debug(fmt.Sprintf("Sending announcement %s to user %s", message, user.Username))
		k.SendMessage(user.ChatID, formattedAnnouncement)
	}

	k.SendMessage(chat.ID, "Announcement was sent to all users.")
}

// Send OZB good deal message to user
func (k *KramerBot) SendOzbGoodDeal(user *models.UserData, deal *models.OzBargainDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🟠🔥<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`🟠🔥 %s 🔺%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending good deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	k.UpdateUser(user)
}

// Send OZB super deal to user
func (k *KramerBot) SendOzbSuperDeal(user *models.UserData, deal *models.OzBargainDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🟠🔥<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`🟠🔥 %s 🔺%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending super deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	k.UpdateUser(user)
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
	formattedDeal := fmt.Sprintf(`🅰️<a href="%s" target="_blank">%s</a> - %s`, deal.Url, shortenedTitle, k.CCCScraper.GetDealDropString(deal))
	textDeal := fmt.Sprintf(`🅰️ %s`, shortenedTitle)

	k.Logger.Debug(fmt.Sprintf("Sending Amazon %s deal %s to user %s", dealType, shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.AmzSent = append(user.AmzSent, deal.Id)
	k.UpdateUser(user)
}

// Send OZB watched deal to user
func (k *KramerBot) SendOzbWatchedDeal(user *models.UserData, deal *models.OzBargainDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🟠👀<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`🟠👀 %s 🔺%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending watched Ozbargain deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	k.UpdateUser(user)
}

// Send AMZ watched deal to user
func (k *KramerBot) SendAmzWatchedDeal(user *models.UserData, deal *models.CamCamCamDeal) {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🅰️👀<a href="%s" target="_blank">%s</a> - %s`, deal.Url, shortenedTitle, k.CCCScraper.GetDealDropString(deal))
	textDeal := fmt.Sprintf(`🅰️👀 %s`, shortenedTitle)

	k.Logger.Debug(fmt.Sprintf("Sending watched Amazon deal %s to user %s", shortenedTitle, user.Username))
	k.SendHTMLMessage(user.ChatID, formattedDeal)

	// Send android notification if username is set
	if strings.EqualFold(user.Username, k.Pipup.Username) {
		k.Pipup.SendMediaMessage(textDeal, "Kramerbot")
	}

	// Mark deal as sent
	user.AmzSent = append(user.AmzSent, deal.Id)
	k.UpdateUser(user)
}
