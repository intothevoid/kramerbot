package bot

import (
	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
)

// Create user data from parameters passed in
func (k *KramerBot) CreateUserData(chatID int64, username string, keyword string,
	ozbGood bool, ozbSuper bool, amzDaily bool, amzWeekly bool) *models.UserData {

	userData := models.UserData{}
	userData.ChatID = chatID
	userData.Username = username
	userData.Keywords = append(userData.Keywords, keyword)
	userData.OzbGood = ozbGood
	userData.OzbSuper = ozbSuper
	userData.AmzDaily = amzDaily
	userData.AmzWeekly = amzWeekly

	return &userData
}

// Function to load user store from file
func (k *KramerBot) LoadUserStore() {
	// Load user store i.e. user data indexed by chat id
	if k.DataWriter != nil {
		var err error
		k.UserStore, err = k.DataWriter.ReadUserStore()
		if err != nil {
			k.Logger.Error("Error loading user store: ", zap.Error(err))
		}
	}
}

// Function to save user store to file
func (k *KramerBot) SaveUserStore() {
	// Save user store i.e. user data indexed by chat id
	if k.DataWriter != nil {
		k.DataWriter.WriteUserStore(k.UserStore)
	}
}

// Update single user record in user store
func (k *KramerBot) UpdateUser(userData *models.UserData) error {
	// Update user store
	if k.DataWriter != nil {
		return k.DataWriter.UpdateUser(userData)
	}
	return nil
}

// OzbDealSent checks if an OzBargain deal has already been sent to the user
// by searching for the deal ID in the user's sent deals list
func OzbDealSent(user *models.UserData, deal *models.OzBargainDeal) bool {
	if user == nil || deal == nil {
		return false
	}

	// Use a map for O(1) lookup instead of slice iteration
	sentDeals := make(map[string]bool)
	for _, id := range user.OzbSent {
		sentDeals[id] = true
	}
	return sentDeals[deal.Id]
}

// AmzDealSent checks if an Amazon deal has already been sent to the user
// by searching for the deal ID in the user's sent deals list
func AmzDealSent(user *models.UserData, deal *models.CamCamCamDeal) bool {
	if user == nil || deal == nil {
		return false
	}

	// Use a map for O(1) lookup instead of slice iteration
	sentDeals := make(map[string]bool)
	for _, id := range user.AmzSent {
		sentDeals[id] = true
	}
	return sentDeals[deal.Id]
}
