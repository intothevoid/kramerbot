package bot

import (
	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
)

// Create user data from parameters passed in
func (k *KramerBot) CreateUserData(chatID int64, username string, keyword string,
	goodDeals bool, superDeals bool) *models.UserData {

	userData := models.UserData{}
	userData.ChatID = chatID
	userData.Username = username
	userData.Keywords = append(userData.Keywords, keyword)
	userData.GoodDeals = goodDeals
	userData.SuperDeals = superDeals

	return &userData
}

// Function to load user store from file
func (k *KramerBot) LoadUserStore() {
	// Load user store i.e. user data indexed by chat id
	store := persist.DataStore{Logger: k.Logger}
	k.UserStore = store.ReadUserStore()
}

// Function to save user store to file
func (k *KramerBot) SaveUserStore() {
	// Save user store i.e. user data indexed by chat id
	store := persist.DataStore{Logger: k.Logger}
	store.WriteUserStore(k.UserStore)
}

// Check if the deal has already been sent to the user
func DealSent(user *models.UserData, deal *models.OzBargainDeal) bool {
	// Check if deal.Id is in user.DealsSent
	for _, dealId := range user.DealsSent {
		if dealId == deal.Id {
			return true
		}
	}
	return false
}