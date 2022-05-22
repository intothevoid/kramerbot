package persist_test

import (
	"testing"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// Unit tests for the UserDB functions
// Language: go
// Path: persist/userdb_test.go

// Test the persist.New() function
func TestNew(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "user_test.db"
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Check database name
	if udb.Name != dbName {
		t.Errorf("Expected database name %s, got %s", dbName, udb.Name)
	}

	// Check database connection
	if udb.DB == nil {
		t.Error("Expected database connection, got nil")
	}
}

// Test UserDB.CreateTable() function
func TestCreateTable(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "user_test.db"
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Create table
	err := udb.CreateTable()
	if err != nil {
		t.Errorf("Error creating 'users' table: %s", err)
	}
}

// Test UserDB.AddUser() function
func TestAddUser(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "user_test.db"
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Create table
	err := udb.CreateTable()
	if err != nil {
		t.Errorf("Error creating 'users' table: %s", err)
	}

	// Add user
	user := &models.UserData{
		ChatID:     123456789,
		Username:   "test_user",
		GoodDeals:  false,
		SuperDeals: false,
		Keywords:   []string{"test", "test2"},
		DealsSent:  []string{"120", "122"},
	}

	err = udb.AddUser(user)
	if err != nil {
		// Don't fail test if error is due to duplicate key
		sqerr := err.(sqlite3.Error)
		if sqerr.Code == 19 && sqerr.ExtendedCode == 1555 {
			logger.Error("Error adding user", zap.String("error", err.Error()))
		} else {
			t.Errorf("Error adding user: %s", err)
		}
	}

	// Check if user was added
	_, err = udb.DB.Query(`SELECT * FROM users WHERE chat_id = ?`, user.ChatID)
	if err != nil {
		t.Errorf("Error querying database. Added user not found: %s", err)
	}
}

// Test to get user from database
func TestGetUser(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "user_test.db"
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Create table
	err := udb.CreateTable()
	if err != nil {
		t.Errorf("Error creating 'users' table: %s", err)
	}

	// Add user2
	user2 := &models.UserData{
		ChatID:     007,
		Username:   "bond",
		GoodDeals:  true,
		SuperDeals: false,
		Keywords:   []string{"james", "bond"},
		DealsSent:  []string{"007", "mi5"},
	}

	err = udb.AddUser(user2)
	if err != nil {
		// Don't fail test if error is due to duplicate key
		sqerr := err.(sqlite3.Error)
		if sqerr.Code == 19 && sqerr.ExtendedCode == 1555 {
			logger.Error("Error adding user", zap.String("error", err.Error()))
		} else {
			t.Errorf("Error adding user: %s", err)
		}
	}

	// Get user
	user, err := udb.GetUser(007)
	if err != nil {
		t.Errorf("Error getting user: %s", err)
	}

	// Check if user was added
	if user.ChatID != 007 {
		t.Errorf("Expected user chat_id to be %d, got %d", 007, user.ChatID)
	}
	if user.Username != "bond" {
		t.Errorf("Expected user username to be %s, got %s", "bond", user.Username)
	}
	if user.GoodDeals != true {
		t.Errorf("Expected user good_deals to be %t, got %t", true, user.GoodDeals)
	}
	if user.SuperDeals != false {
		t.Errorf("Expected user super_deals to be %t, got %t", false, user.SuperDeals)
	}
	if len(user.Keywords) != 2 {
		t.Errorf("Expected user keywords to have length 2, got %d", len(user.Keywords))
	}
	if len(user.DealsSent) != 2 {
		t.Errorf("Expected user deals_sent to have length 2, got %d", len(user.DealsSent))
	}
}

// Test update user in database
func TestUpdateUser(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "user_test.db"
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Create table
	err := udb.CreateTable()
	if err != nil {
		t.Errorf("Error creating 'users' table: %s", err)
	}

	// Add user
	user := &models.UserData{
		ChatID:     123456789,
		Username:   "test_user",
		GoodDeals:  false,
		SuperDeals: false,
		Keywords:   []string{"test", "test2"},
		DealsSent:  []string{"120", "122"},
	}

	err = udb.AddUser(user)
	if err != nil {
		// Don't fail test if error is due to duplicate key
		sqerr := err.(sqlite3.Error)
		if sqerr.Code == 19 && sqerr.ExtendedCode == 1555 {
			logger.Error("Error adding user", zap.String("error", err.Error()))
		} else {
			t.Errorf("Error adding user: %s", err)
		}
	}

	// Update user
	user.GoodDeals = true
	user.SuperDeals = true
	user.Keywords = []string{"test", "test2", "test3"}
	user.DealsSent = []string{"120", "122", "123"}
	user.Username = "test_user_updated"

	err = udb.UpdateUser(user)
	if err != nil {
		t.Errorf("Error updating user: %s", err)
	}

	// Check if user was updated
	_, err = udb.DB.Query(`SELECT * FROM users WHERE chat_id = ?`, user.ChatID)
	if err != nil {
		t.Errorf("Error querying database. Updated user not found: %s", err)
	}
}
