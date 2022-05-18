package models

// This package stores the model for user data

// User model - indexed by chat ID
type UserStore struct {
	Users map[int64]*UserData
}

// User data model
type UserData struct {
	ChatID     int64    `json:"chatID"`     // Telegram chat ID
	Username   string   `json:"username"`   // Telegram username
	GoodDeals  bool     `json:"gooddeals"`  // watch deals with 25+ upvotes in the last 24 hours
	SuperDeals bool     `json:"superdeals"` // watch deals with 50+ upvotes in the last 24 hours
	Keywords   []string `json:"keywords"`   // list of keywords / deals to watch for
	DealsSent  []string `json:"dealssent"`  // comma separated list of deals sent to user
}
