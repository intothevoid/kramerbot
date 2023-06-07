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

// Check if the OZB deal has already been sent to the user
func OzbDealSent(user *models.UserData, deal *models.OzBargainDeal) bool {
	// Check if deal.Id is in user.DealsSent
	for _, dealId := range user.OzbSent {
		if dealId == deal.Id {
			return true
		}
	}
	return false
}

// Check if the AMZ deal has already been sent to the user
func AmzDealSent(user *models.UserData, deal *models.CamCamCamDeal) bool {
	// Check if deal.Id is in user.DealsSent
	for _, dealId := range user.AmzSent {
		if dealId == deal.Id {
			return true
		}
	}
	return false
}
