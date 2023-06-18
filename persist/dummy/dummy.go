package persist

import (
	"github.com/intothevoid/kramerbot/models"
	"go.uber.org/zap"
)

type DummyDB struct {
	Logger *zap.Logger
}

// Connect to the dummy database
func New(dummyUri string, dbName string, collName string, logger *zap.Logger) (*DummyDB, error) {
	// Log all parameters
	logger.Info("dummydb: Connecting to the database",
		zap.String("dummyUri", dummyUri),
		zap.String("dbName", dbName),
		zap.String("collName", collName))

	// Finally, return the UserDB value.
	return &DummyDB{
		Logger: logger,
	}, nil
}

// Close the database
func (ddb *DummyDB) Close() error {
	// To close the dummyDB connection, you can use the Disconnect function.
	ddb.Logger.Info("dummydb: Closing the database")
	return nil
}

// Add user to the database
func (ddb *DummyDB) AddUser(user *models.UserData) error {
	ddb.Logger.Info("dummydb: Adding user to the database")

	// The user has been successfully inserted into the collection.
	return nil
}

// Update user in the database
func (ddb *DummyDB) UpdateUser(user *models.UserData) error {
	ddb.Logger.Info("dummydb: Updating user in the database")

	// The user has been successfully updated in the collection.
	return nil
}

// Delete user from the database
func (ddb *DummyDB) DeleteUser(user *models.UserData) error {
	ddb.Logger.Info("dummydb: Deleting user from the database")

	// The user has been successfully deleted from the collection.
	return nil
}

// Get user from the database by chat_id
func (ddb *DummyDB) GetUser(chatID int64) (*models.UserData, error) {
	ddb.Logger.Info("dummydb: Getting user from the database")

	// The user has been successfully retrieved from the collection.
	// Return empty user
	return &models.UserData{}, nil
}

// Read all users from the database
func (ddb *DummyDB) ReadUserStore() (*models.UserStore, error) {
	ddb.Logger.Info("dummydb: Reading all users from the database")

	return &models.UserStore{}, nil
}

// Write *models.UserStore to the database
func (ddb *DummyDB) WriteUserStore(userStore *models.UserStore) error {
	ddb.Logger.Info("dummydb: Writing all users to the database")

	// The users have been successfully inserted into the collection.
	return nil
}
