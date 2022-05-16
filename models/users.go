package models

// This package stores the model for user data
// User data model - indexed by chat ID
type UserData struct {
	ChatID     int64  `json:"chatID"`     // Telegram chat ID
	Username   string `json:"username"`   // Telegram username
	GoodDeals  bool   `json:"gooddeals"`  // watch deals with 25+ upvotes in the last 24 hours
	SuperDeals bool   `json:"superdeals"` // watch deals with 50+ upvotes in the last 24 hours
	Deals100   bool   `json:"deals100"`   // watch deals with 100+ upvotes
	Keywords   string `json:"deals"`      // comma separated list of keywords
}
