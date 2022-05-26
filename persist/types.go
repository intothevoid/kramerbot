package persist

import "github.com/intothevoid/kramerbot/models"

type UserStore interface {
	WriteUserStore(userStore *models.UserStore) error
	ReadUserStore() (*models.UserStore, error)
}
