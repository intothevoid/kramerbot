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
	k.SendMessage(chat.ID, fmt.Sprintf("Hi %s! Welcome to @kramerbot\n\n"+
		"To configure your bot, please sign up at <TBD>\n"+
		"You ChatID is %d\n", chat.FirstName, chat.ID))
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
