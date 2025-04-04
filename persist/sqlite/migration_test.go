package sqlite_test

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/persist/sqlite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test migrateUserStoreFromJsonToDatabase() function
func TestMigrateUserStoreFromJsonToDatabase(t *testing.T) {
	// Create logger
	logger := zap.NewExample()

	// Create test data file
	testUser := models.UserData{
		ChatID:    123456789,
		Username:  "testuser",
		Keywords:  []string{"test"},
		OzbSent:   []string{},
		OzbGood:   true,
		OzbSuper:  true,
		AmzDaily:  true,
		AmzWeekly: true,
		AmzSent:   []string{},
	}

	testStore := models.UserStore{
		Users: map[int64]*models.UserData{
			testUser.ChatID: &testUser,
		},
	}

	// Write test data to file
	file, err := os.Create("user_store_test.json")
	if err != nil {
		t.Fatalf("Failed to create test data file: %v", err)
	}
	defer os.Remove("user_store_test.json") // Clean up after test

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(testStore); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	file.Close()

	// Run migration
	if err := sqlite.MigrateUserStoreFromJsonToDatabase(logger); err != nil {
		t.Fatalf("Failed to migrate user store: %v", err)
	}

	// Open database file
	dbName := "users_test.db"
	defer os.Remove(dbName) // Clean up after test
	udb, err := sqlite.CreateDatabaseConnection(dbName, logger)
	if err != nil {
		t.Fatalf("Failed to create database connection: %v", err)
	}
	assert.NoError(t, err)

	// Create table
	if err := udb.CreateTable(); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	assert.NoError(t, err)

	// Get count of users in database
	count, err := udb.DB.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		t.Fatalf("Error getting count of users in database: %v", err)
	}
	defer count.Close()

	actualCount := checkCount(count)

	if actualCount != 1 {
		t.Errorf("Expected 1 user in database, got %d", actualCount)
	}
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}
	return count
}
