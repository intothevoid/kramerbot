package util

import (
	"encoding/json"
	"os"

	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
)

type DataStore struct {
	Logger *zap.Logger
}

// Function to write User store to a file
func (d *DataStore) WriteUserStore(userStore *models.UserStore) {
	d.Logger.Debug("Writing user store to file")

	// create the file
	file, err := os.Create("user_store.json")
	if err != nil {
		d.Logger.Error(err.Error())

		file.Close()
		return
	}

	defer file.Close()

	// write the user store to the file
	encoder := json.NewEncoder(file)
	encoder.Encode(userStore)
}

// Function to read User store from a file
func (d *DataStore) ReadUserStore() *models.UserStore {
	d.Logger.Debug("Reading user store from file")

	file, err := os.Open("user_store.json")
	if err != nil {
		d.Logger.Error(err.Error())

		file.Close()
		return &models.UserStore{}
	}

	defer file.Close()

	// decode the user store from the file
	decoder := json.NewDecoder(file)
	var userStore models.UserStore
	decoder.Decode(&userStore)

	return &userStore
}
