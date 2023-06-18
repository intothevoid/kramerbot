package models

// This package stores the model for user data

// User model - indexed by chat ID
type UserStore struct {
	Users map[int64]*UserData
}

// User data model
type UserData struct {
	ChatID         int64    `bson:"chat_id"`         // Telegram chat ID
	Username       string   `bson:"username"`        // Telegram username
	OzbGood        bool     `bson:"ozb_good"`        // watch deals with 25+ upvotes in the last 24 hours
	OzbSuper       bool     `bson:"ozb_super"`       // watch deals with 50+ upvotes in the last 24 hours
	Keywords       []string `bson:"keywords"`        // list of keywords / deals to watch for
	OzbSent        []string `bson:"ozb_sent"`        // comma separated list of ozb deals sent to user
	AmzDaily       bool     `bson:"amz_daily"`       // watch top daily deals on amazon
	AmzWeekly      bool     `bson:"amz_weekly"`      // watch top weekly deals on amazon
	AmzSent        []string `bson:"amz_sent"`        // comma separated list of amz deals sent to user
	UsernameChosen string   `bson:"username_chosen"` // username chosen by user on website
	Password       string   `bson:"password"`        // password chosen by user on website
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
func (u *UserData) SetOzbGood(ozbGood bool) {
	u.OzbGood = ozbGood
}
func (u *UserData) GetOzbGood() bool {
	return u.OzbGood
}
func (u *UserData) SetOzbSuper(ozbSuper bool) {
	u.OzbSuper = ozbSuper
}
func (u *UserData) GetOzbSuper() bool {
	return u.OzbSuper
}
func (u *UserData) SetKeywords(keywords []string) {
	u.Keywords = keywords
}
func (u *UserData) GetKeywords() []string {
	return u.Keywords
}
func (u *UserData) SetOzbSent(ozbSent []string) {
	u.OzbSent = ozbSent
}
func (u *UserData) GetOzbSent() []string {
	return u.OzbSent
}
func (u *UserData) GetAmzDaily() bool {
	return u.AmzDaily
}
func (u *UserData) SetAmzDaily(amzDaily bool) {
	u.AmzDaily = amzDaily
}
func (u *UserData) GetAmzWeekly() bool {
	return u.AmzWeekly
}
func (u *UserData) SetAmzWeekly(amzWeekly bool) {
	u.AmzWeekly = amzWeekly
}
func (u *UserData) SetAmzSent(amzSent []string) {
	u.AmzSent = amzSent
}
func (u *UserData) GetAmzSent() []string {
	return u.AmzSent
}
func (u *UserData) SetUsernameChosen(usernameChosen string) {
	u.UsernameChosen = usernameChosen
}
func (u *UserData) GetUsernameChosen() string {
	return u.UsernameChosen
}
