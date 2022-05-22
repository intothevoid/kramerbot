package persist

import (
	"database/sql"

	"github.com/intothevoid/kramerbot/models"
)

// Create sqlite database
func CreateSqliteDB(dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	return db
}

// Write user store to sqlite database
func CreateUserStore(db *sql.DB, userStore *models.UserStore) {

	// Check if user store table exists
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='user_store';")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// create the table for *models.UserData
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT,
		created_at TEXT,
		updated_at TEXT,
		deleted_at TEXT
	)`)

	if err != nil {
		panic(err)
	}
}
