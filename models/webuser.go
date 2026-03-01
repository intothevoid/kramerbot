package models

import "time"

// WebUser represents a user account created via the web interface.
type WebUser struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	PasswordHash     string     `json:"-"`
	DisplayName      string     `json:"display_name"`
	TelegramChatID   *int64     `json:"telegram_chat_id,omitempty"`
	TelegramUsername *string    `json:"telegram_username,omitempty"`
	LinkToken        *string    `json:"-"`
	LinkTokenExpires *time.Time `json:"-"`
	ResetToken       *string    `json:"-"`
	ResetTokenExpires *time.Time `json:"-"`
	// Deal notification preferences (synced to bot's UserData when Telegram is linked).
	OzbGood   bool     `json:"ozb_good"`
	OzbSuper  bool     `json:"ozb_super"`
	AmzDaily  bool     `json:"amz_daily"`
	AmzWeekly bool     `json:"amz_weekly"`
	Keywords  []string `json:"keywords"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
