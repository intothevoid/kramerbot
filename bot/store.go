package bot

import (
	"github.com/intothevoid/kramerbot/models"
)

// Create user data from parameters passed in
func (k *KramerBot) CreateUserData(chatID int64, username string, keyword string,
	goodDeals bool, superDeals bool) *models.UserData {

	userData := models.UserData{}
	userData.ChatID = chatID
	userData.Username = username
	userData.Keywords = append(userData.Keywords, keyword)
	userData.OzbGood = goodDeals
	userData.OzbSuper = superDeals

	return &userData
}

// Function to load user store from file
func (k *KramerBot) LoadUserStore() {
	// Load user store i.e. user data indexed by chat id
	if k.DataWriter != nil {
		k.UserStore, _ = k.DataWriter.ReadUserStore()
	}
}

// Function to save user store to file
func (k *KramerBot) SaveUserStore() {
	// Save user store i.e. user data indexed by chat id
	if k.DataWriter != nil {
		k.DataWriter.WriteUserStore(k.UserStore)
	}
}

// Check if the deal has already been sent to the user
func DealSent(user *models.UserData, deal *models.OzBargainDeal) bool {
	// Check if deal.Id is in user.DealsSent
	for _, dealId := range user.OzbSent {
		if dealId == deal.Id {
			return true
		}
	}
	return false
}
