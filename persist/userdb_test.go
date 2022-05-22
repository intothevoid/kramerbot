package persist_test

import (
	"testing"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// Unit tests for the UserDB functions
// Language: go
// Path: persist/userdb_test.go

// Test the persist.New() function
func TestNew(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "./user_test.db"
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
	dbName := "./user_test.db"
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
	dbName := "./user_test.db"
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
