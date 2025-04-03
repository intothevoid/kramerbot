package sqlite

import (
	"encoding/json"
	"os"

	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
)

// Read user store and convert to sqlite database
func MigrateUserStoreFromJsonToDatabase(logger *zap.Logger) error {
	logger.Debug("Reading user store from file")

	file, err := os.Open("user_store_test.json")

	if err != nil {
		logger.Error(err.Error())

		file.Close()
		return err
	}

	defer file.Close()

	// decode the user store from the file
	decoder := json.NewDecoder(file)
	var userStore models.UserStore
	decoder.Decode(&userStore)

	// create database and table
	dbName := "users_test.db"
	udb, err := CreateDatabaseConnection(dbName, logger)
	if err != nil {
		logger.Error("Failed to create database connection", zap.Error(err))
		return err
	}

	// Create table
	err = udb.CreateTable()
	if err != nil {
		logger.Sugar().Errorf("Error creating 'users' table: %s", err)
	}

	// insert the users into the database
	for _, user := range userStore.Users {
		err = udb.AddUser(user)
		if err != nil {
			logger.Sugar().Errorf("Error adding user: %s", err)
			return err
		}
	}

	return err
}
