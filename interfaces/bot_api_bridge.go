package interfaces

import (
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
)

// BotAPIBridge defines the methods the API handlers need to interact with the core bot logic.
// This helps avoid direct import cycles between the 'bot' and 'api' packages.
type BotAPIBridge interface {
	GetUserDataWriter() persist.DatabaseIF
	GetBotToken() string
	GetUserData(chatID int64) (*models.UserData, error) // Added for convenience/direct access if needed
	UpdateUserData(user *models.UserData) error         // Added for updating memory store too
	SendTestMessageToChat(chatID int64, user *models.TelegramUser) error
}
