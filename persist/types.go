package persist

import "github.com/intothevoid/kramerbot/models"

type DataStoreInterface interface {
	WriteUserStore(userStore *models.UserStore) error
	ReadUserStore() (*models.UserStore, error)
}
