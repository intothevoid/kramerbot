package bot

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/intothevoid/kramerbot/util"
)

// welcomeMessage returns the standard welcome text including the web app URL.
func (k *KramerBot) welcomeMessage(firstName string) string {
	webURL := "http://localhost:8080"
	if k.Config != nil && k.Config.API.WebURL != "" {
		webURL = k.Config.API.WebURL
	}
	return fmt.Sprintf("👋 Welcome to KramerBot - Aussie Deals, %s!\n\nManage your deal preferences and subscriptions at:\n%s", firstName, webURL)
}

// Help sends the welcome message with the web app URL.
func (k *KramerBot) Help(chat *tgbotapi.Chat) {
	k.SendMessage(chat.ID, k.welcomeMessage(chat.FirstName))
}

// RegisterUser adds a new user or shows the welcome message for an existing user.
func (k *KramerBot) RegisterUser(chat *tgbotapi.Chat) {
	user, err := k.DataWriter.GetUser(chat.ID)
	if err != nil || user == nil {
		newUser := &models.UserData{
			ChatID:    chat.ID,
			Username:  chat.UserName,
			OzbGood:   false, // no subscriptions until user opts-in via the web UI
			OzbSuper:  false,
			AmzDaily:  false,
			AmzWeekly: false,
			Keywords:  []string{},
			OzbSent:   []string{},
			AmzSent:   []string{},
		}
		if err := k.DataWriter.AddUser(newUser); err != nil {
			k.Logger.Error("Failed to add new user", zap.Int64("chatID", chat.ID), zap.Error(err))
			k.SendMessage(chat.ID, "Sorry, there was an error registering you. Please try again later.")
			return
		}
		k.UserStore.SetUser(chat.ID, newUser)
		k.Logger.Info("Registered new user", zap.String("username", chat.UserName), zap.Int64("chatID", chat.ID))
	} else {
		if user.Username != chat.UserName {
			user.Username = chat.UserName
			k.UpdateUser(user)
			k.UserStore.SetUser(chat.ID, user)
			k.Logger.Info("Updated username", zap.String("username", chat.UserName), zap.Int64("chatID", chat.ID))
		}
	}
	k.SendMessage(chat.ID, k.welcomeMessage(chat.FirstName))
}

// ShowPreferences displays the user's current notification settings
func (k *KramerBot) ShowPreferences(chat *tgbotapi.Chat) {
	user, err := k.DataWriter.GetUser(chat.ID)
	if err != nil || user == nil {
		k.SendMessage(chat.ID, "Could not find your user data. Have you registered using /start ?")
		return
	}

	prefsText := fmt.Sprintf("Your current preferences:\n"+
		"OzBargain Regular Deals (all deals): %t\n"+
		"OzBargain Top Deals (25+ votes in 24h): %t\n"+
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

// Helper function to get user data and handle errors
func (k *KramerBot) getUserData(chatID int64) (*models.UserData, error) {
	user, err := k.DataWriter.GetUser(chatID)
	if err != nil || user == nil {
		k.SendMessage(chatID, "Could not find your user data. Have you registered using /start ?")
		return nil, fmt.Errorf("user not found or error fetching data: %w", err)
	}
	return user, nil
}

// ToggleOzbGood toggles the OzbGood preference
func (k *KramerBot) ToggleOzbGood(chat *tgbotapi.Chat) {
	user, err := k.getUserData(chat.ID)
	if err != nil {
		return // Error message already sent by getUserData
	}

	user.OzbGood = !user.OzbGood
	k.UpdateUser(user)                 // Update DB
	k.UserStore.SetUser(chat.ID, user) // Update memory

	k.SendMessage(chat.ID, fmt.Sprintf("OzBargain Regular Deals (all deals) notifications set to: %t", user.OzbGood))
	k.ShowPreferences(chat)
}

// ToggleOzbSuper toggles the OzbSuper preference
func (k *KramerBot) ToggleOzbSuper(chat *tgbotapi.Chat) {
	user, err := k.getUserData(chat.ID)
	if err != nil {
		return
	}

	user.OzbSuper = !user.OzbSuper
	k.UpdateUser(user)                 // Update DB
	k.UserStore.SetUser(chat.ID, user) // Update memory

	k.SendMessage(chat.ID, fmt.Sprintf("OzBargain Top Deals (25+ votes in 24h) notifications set to: %t", user.OzbSuper))
	k.ShowPreferences(chat)
}

// ToggleAmzDaily toggles the AmzDaily preference
func (k *KramerBot) ToggleAmzDaily(chat *tgbotapi.Chat) {
	user, err := k.getUserData(chat.ID)
	if err != nil {
		k.Logger.Error("Failed to get user data", zap.Error(err))
		return // Error message already sent by getUserData
	}

	k.Logger.Debug("Toggling AmzDaily preference",
		zap.Int64("chatID", chat.ID),
		zap.Bool("currentValue", user.AmzDaily),
		zap.Bool("newValue", !user.AmzDaily))

	user.AmzDaily = !user.AmzDaily
	if err := k.UpdateUser(user); err != nil {
		k.Logger.Error("Failed to update user", zap.Error(err))
		k.SendMessage(chat.ID, "Sorry, there was an error updating your preferences. Please try again later.")
		return
	}
	k.UserStore.SetUser(chat.ID, user) // Update memory

	k.Logger.Debug("Successfully updated AmzDaily preference",
		zap.Int64("chatID", chat.ID),
		zap.Bool("newValue", user.AmzDaily))

	k.SendMessage(chat.ID, fmt.Sprintf("Amazon Daily Deals notifications set to: %t", user.AmzDaily))
	k.ShowPreferences(chat)
}

// ToggleAmzWeekly toggles the AmzWeekly preference
func (k *KramerBot) ToggleAmzWeekly(chat *tgbotapi.Chat) {
	user, err := k.getUserData(chat.ID)
	if err != nil {
		return
	}

	user.AmzWeekly = !user.AmzWeekly
	k.UpdateUser(user)                 // Update DB
	k.UserStore.SetUser(chat.ID, user) // Update memory

	k.SendMessage(chat.ID, fmt.Sprintf("Amazon Weekly Deals notifications set to: %t", user.AmzWeekly))
	k.ShowPreferences(chat)
}

// AddKeyword adds a keyword to the user's watch list
func (k *KramerBot) AddKeyword(chat *tgbotapi.Chat, keyword string) {
	user, err := k.getUserData(chat.ID)
	if err != nil {
		return
	}

	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		k.SendMessage(chat.ID, "Please provide a keyword to add. Usage: /addkeyword <keyword>")
		return
	}

	// Check if keyword already exists
	for _, existingKeyword := range user.Keywords {
		if existingKeyword == keyword {
			k.SendMessage(chat.ID, fmt.Sprintf("Keyword '%s' is already in your watch list.", keyword))
			return
		}
	}

	user.Keywords = append(user.Keywords, keyword)
	k.UpdateUser(user)                 // Update DB
	k.UserStore.SetUser(chat.ID, user) // Update memory

	k.SendMessage(chat.ID, fmt.Sprintf("Keyword '%s' added to your watch list.", keyword))
	k.ListKeywords(chat)
}

// RemoveKeyword removes a keyword from the user's watch list
func (k *KramerBot) RemoveKeyword(chat *tgbotapi.Chat, keywordToRemove string) {
	user, err := k.getUserData(chat.ID)
	if err != nil {
		return
	}

	keywordToRemove = strings.TrimSpace(strings.ToLower(keywordToRemove))
	if keywordToRemove == "" {
		k.SendMessage(chat.ID, "Please provide a keyword to remove. Usage: /removekeyword <keyword>")
		return
	}

	found := false
	var updatedKeywords []string
	for _, existingKeyword := range user.Keywords {
		if existingKeyword != keywordToRemove {
			updatedKeywords = append(updatedKeywords, existingKeyword)
		} else {
			found = true
		}
	}

	if !found {
		k.SendMessage(chat.ID, fmt.Sprintf("Keyword '%s' not found in your watch list.", keywordToRemove))
		return
	}

	user.Keywords = updatedKeywords
	k.UpdateUser(user)                 // Update DB
	k.UserStore.SetUser(chat.ID, user) // Update memory

	k.SendMessage(chat.ID, fmt.Sprintf("Keyword '%s' removed from your watch list.", keywordToRemove))
	k.ListKeywords(chat)
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
func (k *KramerBot) SendOzbGoodDeal(user *models.UserData, deal *models.OzBargainDeal) error {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🟠🔥<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`🟠🔥 %s 🔺%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending good deal %s to user %s", shortenedTitle, user.Username))
	if err := k.SendHTMLMessage(user.ChatID, formattedDeal); err != nil {
		return fmt.Errorf("failed to send HTML message: %w", err)
	}

	// Send android notification if username is set
	if k.Pipup != nil && strings.EqualFold(user.Username, k.Pipup.Username) {
		if err := k.Pipup.SendMediaMessage(textDeal, "Kramerbot"); err != nil {
			return fmt.Errorf("failed to send pipup message: %w", err)
		}
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	if err := k.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Send OZB super deal to user
func (k *KramerBot) SendOzbSuperDeal(user *models.UserData, deal *models.OzBargainDeal) error {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🟠🔥<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`🟠🔥 %s 🔺%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending super deal %s to user %s", shortenedTitle, user.Username))
	if err := k.SendHTMLMessage(user.ChatID, formattedDeal); err != nil {
		return fmt.Errorf("failed to send HTML message: %w", err)
	}

	// Send android notification if username is set
	if k.Pipup != nil && strings.EqualFold(user.Username, k.Pipup.Username) {
		if err := k.Pipup.SendMediaMessage(textDeal, "Kramerbot"); err != nil {
			return fmt.Errorf("failed to send pipup message: %w", err)
		}
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	if err := k.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (k *KramerBot) SendAmzDeal(user *models.UserData, deal *models.CamCamCamDeal) error {
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
	if err := k.SendHTMLMessage(user.ChatID, formattedDeal); err != nil {
		return fmt.Errorf("failed to send HTML message: %w", err)
	}

	// Send android notification if username is set
	if k.Pipup != nil && strings.EqualFold(user.Username, k.Pipup.Username) {
		if err := k.Pipup.SendMediaMessage(textDeal, "Kramerbot"); err != nil {
			return fmt.Errorf("failed to send pipup message: %w", err)
		}
	}

	// Mark deal as sent
	user.AmzSent = append(user.AmzSent, deal.Id)
	if err := k.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Send OZB watched deal to user
func (k *KramerBot) SendOzbWatchedDeal(user *models.UserData, deal *models.OzBargainDeal) error {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🟠👀<a href="%s" target="_blank">%s</a>🔺%s`, deal.Url, shortenedTitle, deal.Upvotes)
	textDeal := fmt.Sprintf(`🟠👀 %s 🔺%s`, shortenedTitle, deal.Upvotes)

	k.Logger.Debug(fmt.Sprintf("Sending watched Ozbargain deal %s to user %s", shortenedTitle, user.Username))
	if err := k.SendHTMLMessage(user.ChatID, formattedDeal); err != nil {
		return fmt.Errorf("failed to send HTML message: %w", err)
	}

	// Send android notification if username is set
	if k.Pipup != nil && strings.EqualFold(user.Username, k.Pipup.Username) {
		if err := k.Pipup.SendMediaMessage(textDeal, "Kramerbot"); err != nil {
			return fmt.Errorf("failed to send pipup message: %w", err)
		}
	}

	// Mark deal as sent
	user.OzbSent = append(user.OzbSent, deal.Id)
	if err := k.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Send AMZ watched deal to user
func (k *KramerBot) SendAmzWatchedDeal(user *models.UserData, deal *models.CamCamCamDeal) error {
	shortenedTitle := util.ShortenString(deal.Title, 30) + "..."
	formattedDeal := fmt.Sprintf(`🅰️👀<a href="%s" target="_blank">%s</a> - %s`, deal.Url, shortenedTitle, k.CCCScraper.GetDealDropString(deal))
	textDeal := fmt.Sprintf(`🅰️👀 %s`, shortenedTitle)

	k.Logger.Debug(fmt.Sprintf("Sending watched Amazon deal %s to user %s", shortenedTitle, user.Username))
	if err := k.SendHTMLMessage(user.ChatID, formattedDeal); err != nil {
		return fmt.Errorf("failed to send HTML message: %w", err)
	}

	// Send android notification if username is set
	if k.Pipup != nil && strings.EqualFold(user.Username, k.Pipup.Username) {
		if err := k.Pipup.SendMediaMessage(textDeal, "Kramerbot"); err != nil {
			return fmt.Errorf("failed to send pipup message: %w", err)
		}
	}

	// Mark deal as sent
	user.AmzSent = append(user.AmzSent, deal.Id)
	if err := k.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
