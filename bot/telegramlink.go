package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

// HandleTelegramLink checks whether args is a valid web-account link token.
// If it is, the user's Telegram chat ID is written to the web_users record and
// a confirmation message is sent. Otherwise the normal registration flow runs.
func (k *KramerBot) HandleTelegramLink(chat *tgbotapi.Chat, token string) {
	if k.WebUserDB == nil {
		// API not initialised — fall back to normal registration.
		k.RegisterUser(chat)
		return
	}

	webUser, err := k.WebUserDB.GetWebUserByLinkToken(token)
	if err != nil {
		k.Logger.Error("error looking up link token", zap.String("token", token), zap.Error(err))
		k.SendMessage(chat.ID, "❌ Something went wrong. Please try generating a new link from the website.")
		return
	}

	if webUser == nil {
		// Token not found or expired — treat as a normal /start.
		k.Logger.Info("link token not found or expired", zap.String("token", token))
		k.RegisterUser(chat)
		return
	}

	// Write the Telegram chat ID and username into the web user record.
	chatID := chat.ID
	webUser.TelegramChatID = &chatID

	username := chat.UserName
	if username != "" {
		webUser.TelegramUsername = &username
	}

	// Clear the one-time token.
	webUser.LinkToken = nil
	webUser.LinkTokenExpires = nil

	if err := k.WebUserDB.UpdateWebUser(webUser); err != nil {
		k.Logger.Error("failed to save telegram link", zap.Error(err))
		k.SendMessage(chat.ID, "❌ Could not save the link. Please try again from the website.")
		return
	}

	k.Logger.Info("telegram account linked to web user",
		zap.Int64("chat_id", chat.ID),
		zap.String("web_user_id", webUser.ID),
	)

	// Ensure the user exists in the bot's own users table so notifications work.
	k.RegisterUser(chat)

	k.SendMessage(chat.ID,
		"✅ Your Telegram account is now linked to your KramerBot web account!\n\n"+
			"You can manage your deal preferences at any time from the web dashboard.",
	)
}
