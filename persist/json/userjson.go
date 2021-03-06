package persist

import (
	"encoding/json"
	"os"

	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
)

type UserStoreJson struct {
	Logger *zap.Logger
}

// Function to write User store to a file
func (d *UserStoreJson) WriteUserStore(userStore *models.UserStore) error {
	d.Logger.Debug("Writing user store to file")

	// create the file
	file, err := os.Create("user_store.json")
	if err != nil {
		d.Logger.Error(err.Error())

		file.Close()
		return err
	}

	defer file.Close()

	// write the user store to the file
	encoder := json.NewEncoder(file)
	encoder.Encode(userStore)

	return nil
}

// Function to read User store from a file
func (d *UserStoreJson) ReadUserStore() (*models.UserStore, error) {
	d.Logger.Debug("Reading user store from file")

	file, err := os.Open("user_store.json")
	if err != nil {
		d.Logger.Error(err.Error())

		file.Close()
		return d.CreateEmptyUserStore(), err
	}

	defer file.Close()

	// decode the user store from the file
	decoder := json.NewDecoder(file)
	var userStore models.UserStore
	decoder.Decode(&userStore)

	return &userStore, nil
}

// Create empty user store
func (d *UserStoreJson) CreateEmptyUserStore() *models.UserStore {
	return &models.UserStore{
		Users: make(map[int64]*models.UserData),
	}
}
