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

// setters and getters for UserData
func (u *UserData) SetChatID(chatID int64) {
	u.ChatID = chatID
}
func (u *UserData) GetChatID() int64 {
	return u.ChatID
}
func (u *UserData) SetUsername(username string) {
	u.Username = username
}
func (u *UserData) GetUsername() string {
	return u.Username
}
func (u *UserData) SetGoodDeals(goodDeals bool) {
	u.GoodDeals = goodDeals
}
func (u *UserData) GetGoodDeals() bool {
	return u.GoodDeals
}
func (u *UserData) SetSuperDeals(superDeals bool) {
	u.SuperDeals = superDeals
}
func (u *UserData) GetSuperDeals() bool {
	return u.SuperDeals
}
func (u *UserData) SetKeywords(keywords []string) {
	u.Keywords = keywords
}
func (u *UserData) GetKeywords() []string {
	return u.Keywords
}
func (u *UserData) SetDealsSent(dealsSent []string) {
	u.DealsSent = dealsSent
}
func (u *UserData) GetDealsSent() []string {
	return u.DealsSent
}
