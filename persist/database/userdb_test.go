package persist_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/intothevoid/kramerbot/models"
	persist "github.com/intothevoid/kramerbot/persist/database"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// Unit tests for the UserDB functions
// Language: go
// Path: persist/userdb_test.go

// Test the persist.New() function
func TestNew(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "users_test.db"
	defer DeleteDBFile(dbName)
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
	dbName := "users_test.db"
	defer DeleteDBFile(dbName)
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
	dbName := "users_test.db"
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Create table
	err := udb.CreateTable()
	if err != nil {
		t.Errorf("Error creating 'users' table: %s", err)
	}

	// Add user
	user := &models.UserData{
		ChatID:    123456789,
		Username:  "test_user",
		OzbGood:   false,
		OzbSuper:  false,
		Keywords:  []string{"test", "test2"},
		OzbSent:   []string{"120", "122"},
		AmzDaily:  false,
		AmzWeekly: false,
		AmzSent:   []string{"222", "333"},
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
	dbName := "users_test.db"
	defer DeleteDBFile(dbName)
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Create table
	err := udb.CreateTable()
	if err != nil {
		t.Errorf("Error creating 'users' table: %s", err)
	}

	// Add user2
	user2 := &models.UserData{
		ChatID:    007,
		Username:  "bond",
		OzbGood:   true,
		OzbSuper:  false,
		Keywords:  []string{"james", "bond"},
		OzbSent:   []string{"007", "mi5"},
		AmzDaily:  false,
		AmzWeekly: true,
		AmzSent:   []string{"222", "333", "444"},
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
	if user.OzbGood != true {
		t.Errorf("Expected user ozb_good to be %t, got %t", true, user.OzbGood)
	}
	if user.OzbSuper != false {
		t.Errorf("Expected user ozb_super to be %t, got %t", false, user.OzbSuper)
	}
	if len(user.Keywords) != 2 {
		t.Errorf("Expected user keywords to have length 2, got %d", len(user.Keywords))
	}
	if len(user.OzbSent) != 2 {
		t.Errorf("Expected user ozb_sent to have length 2, got %d", len(user.OzbSent))
	}
	if user.AmzDaily != false {
		t.Errorf("Expected user amz_daily to be %t, got %t", true, user.OzbGood)
	}
	if user.AmzWeekly != true {
		t.Errorf("Expected user amz_weekly to be %t, got %t", false, user.OzbSuper)
	}
	if len(user.AmzSent) != 3 {
		t.Errorf("Expected user amz_sent to have length 3, got %d", len(user.OzbSent))
	}
}

// Test update user in database
func TestUpdateUser(t *testing.T) {
	var logger = *zap.NewExample()
	dbName := "users_test.db"
	defer DeleteDBFile(dbName)
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Create table
	err := udb.CreateTable()
	if err != nil {
		t.Errorf("Error creating 'users' table: %s", err)
	}

	// Add user
	user := &models.UserData{
		ChatID:    123456789,
		Username:  "test_user",
		OzbGood:   false,
		OzbSuper:  false,
		Keywords:  []string{"test", "test2"},
		OzbSent:   []string{"120", "122"},
		AmzDaily:  false,
		AmzWeekly: true,
		AmzSent:   []string{"222", "333", "444"},
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
	user.OzbGood = true
	user.OzbSuper = true
	user.Keywords = []string{"test", "test2", "test3"}
	user.OzbSent = []string{"120", "122", "123"}
	user.Username = "test_user_updated"
	user.AmzDaily = true
	user.AmzWeekly = false
	user.OzbSent = []string{"555", "666", "777"}

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

// Delete database file
func DeleteDBFile(dbName string) {
	err := os.Remove(dbName)
	if err != nil {
		fmt.Printf("Error deleting database file: %s", err)
	}
}
