package persist

import "github.com/intothevoid/kramerbot/models"

type UserStore interface {
	WriteUserStore(userStore *models.UserStore) error
	ReadUserStore() (*models.UserStore, error)
}

type DatabaseIF interface {
	WriteUserStore(userStore *models.UserStore) error
	ReadUserStore() (*models.UserStore, error)
	GetUser(chatID int64) (*models.UserData, error)
	DeleteUser(user *models.UserData) error
	AddUser(user *models.UserData) error
	UpdateUser(user *models.UserData) error
	Close() error
}
