package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/intothevoid/kramerbot/models"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type UserStoreDB struct {
	DB     *sql.DB
	Name   string
	Logger *zap.Logger
}

// Connect to the database with connection pooling
func CreateDatabaseConnection(dbName string, logger *zap.Logger) (*UserStoreDB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection
	db.SetConnMaxIdleTime(1 * time.Minute) // Maximum idle time of a connection

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &UserStoreDB{
		DB:     db,
		Name:   dbName,
		Logger: logger,
	}, nil
}

// Close the database
func (udb *UserStoreDB) Close() error {
	if err := udb.DB.Close(); err != nil {
		udb.Logger.Error("Error closing database", zap.Error(err))
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

// Create *models.UserData table in database
func (udb *UserStoreDB) CreateTable() error {
	// Set WAL mode first, outside of transaction
	if _, err := udb.DB.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Set other pragmas
	if _, err := udb.DB.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		return fmt.Errorf("failed to set synchronous mode: %w", err)
	}

	if _, err := udb.DB.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return fmt.Errorf("failed to set busy timeout: %w", err)
	}

	// Start a transaction for table creation
	tx, err := udb.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() is called

	// Create the table
	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			chat_id INTEGER PRIMARY KEY,
			username TEXT,
			ozb_good INTEGER,
			ozb_super INTEGER,
			keywords BLOB,
			ozb_sent BLOB,
			amz_daily INTEGER,
			amz_weekly INTEGER,
			amz_sent BLOB
		);
	`); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Add user to the database with retry mechanism
func (udb *UserStoreDB) AddUser(user *models.UserData) error {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := udb.addUserWithTransaction(user); err != nil {
			lastErr = err
			udb.Logger.Warn("Failed to add user, retrying...",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", maxRetries),
				zap.Error(err))
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to add user after %d retries: %w", maxRetries, lastErr)
}

// Helper function to add user within a transaction
func (udb *UserStoreDB) addUserWithTransaction(user *models.UserData) error {
	// Start a transaction
	tx, err := udb.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() is called

	// Convert string arrays to bytes
	keywords, err := json.Marshal(user.Keywords)
	if err != nil {
		return fmt.Errorf("failed to marshal keywords: %w", err)
	}

	ozbSent, err := json.Marshal(user.OzbSent)
	if err != nil {
		return fmt.Errorf("failed to marshal OZB deals sent: %w", err)
	}

	amzSent, err := json.Marshal(user.AmzSent)
	if err != nil {
		return fmt.Errorf("failed to marshal AMZ deals sent: %w", err)
	}

	// Insert the user
	_, err = tx.Exec(`
		INSERT INTO users (
			chat_id, username, ozb_good, ozb_super, keywords, ozb_sent, amz_daily, amz_weekly, amz_sent
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ChatID, user.Username, user.OzbGood, user.OzbSuper, keywords, ozbSent, user.AmzDaily, user.AmzWeekly, amzSent,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Update user in the database with retry mechanism
func (udb *UserStoreDB) UpdateUser(user *models.UserData) error {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := udb.updateUserWithTransaction(user); err != nil {
			lastErr = err
			udb.Logger.Warn("Failed to update user, retrying...",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", maxRetries),
				zap.Error(err))
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to update user after %d retries: %w", maxRetries, lastErr)
}

// Helper function to update user within a transaction
func (udb *UserStoreDB) updateUserWithTransaction(user *models.UserData) error {
	// Start a transaction
	tx, err := udb.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() is called

	// Convert string arrays to bytes
	keywords, err := json.Marshal(user.Keywords)
	if err != nil {
		return fmt.Errorf("failed to marshal keywords: %w", err)
	}

	ozbSent, err := json.Marshal(user.OzbSent)
	if err != nil {
		return fmt.Errorf("failed to marshal OZB deals sent: %w", err)
	}

	amzSent, err := json.Marshal(user.AmzSent)
	if err != nil {
		return fmt.Errorf("failed to marshal AMZ deals sent: %w", err)
	}

	// Update the user
	result, err := tx.Exec(`
		UPDATE users SET
			username = ?, ozb_good = ?, ozb_super = ?, keywords = ?, ozb_sent = ?, amz_daily = ?, amz_weekly = ?, amz_sent = ?
		WHERE chat_id = ?`,
		user.Username, user.OzbGood, user.OzbSuper, keywords, ozbSent, user.AmzDaily, user.AmzWeekly, amzSent, user.ChatID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated for chat_id %d", user.ChatID)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Delete user from the database with retry mechanism
func (udb *UserStoreDB) DeleteUser(user *models.UserData) error {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := udb.deleteUserWithTransaction(user); err != nil {
			lastErr = err
			udb.Logger.Warn("Failed to delete user, retrying...",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", maxRetries),
				zap.Error(err))
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to delete user after %d retries: %w", maxRetries, lastErr)
}

// Helper function to delete user within a transaction
func (udb *UserStoreDB) deleteUserWithTransaction(user *models.UserData) error {
	// Start a transaction
	tx, err := udb.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() is called

	// Delete the user
	result, err := tx.Exec(`DELETE FROM users WHERE chat_id = ?`, user.ChatID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted for chat_id %d", user.ChatID)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Get user from the database by chat_id
func (udb *UserStoreDB) GetUser(chatID int64) (*models.UserData, error) {
	// Start a transaction
	tx, err := udb.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() is called

	user := &models.UserData{}
	keywords := []byte{}
	ozbSent := []byte{}
	amzSent := []byte{}

	// Get the user
	err = tx.QueryRow(`SELECT * FROM users WHERE chat_id = ?`, chatID).Scan(
		&user.ChatID, &user.Username, &user.OzbGood, &user.OzbSuper, &keywords, &ozbSent, &user.AmzDaily, &user.AmzWeekly, &amzSent,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with chat_id %d", chatID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert bytes back to string arrays
	if err := json.Unmarshal(keywords, &user.Keywords); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keywords: %w", err)
	}
	if err := json.Unmarshal(ozbSent, &user.OzbSent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OZB deals sent: %w", err)
	}
	if err := json.Unmarshal(amzSent, &user.AmzSent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AMZ deals sent: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user, nil
}

// Read all users from the database
func (udb *UserStoreDB) ReadUserStore() (*models.UserStore, error) {
	rows, err := udb.DB.Query(`Select * from users`)
	if err != nil {
		udb.Logger.Error("Error getting all users", zap.Error(err))
		return nil, err
	}

	userStore := &models.UserStore{
		Users: make(map[int64]*models.UserData),
	}

	for rows.Next() {
		user := &models.UserData{}
		keywords := []byte{}
		ozbSent := []byte{}
		amzSent := []byte{}

		err = rows.Scan(
			&user.ChatID, &user.Username, &user.OzbGood, &user.OzbSuper, &keywords, &ozbSent, &user.AmzDaily, &user.AmzWeekly, &amzSent,
		)
		if err != nil {
			udb.Logger.Error("Error getting user", zap.Error(err))
			return nil, err
		}

		// Bytes to string array - keywords
		if err := json.Unmarshal([]byte(keywords), &user.Keywords); err != nil {
			udb.Logger.Error("Error unmarshalling user keywords", zap.Error(err))
		}

		// Bytes to string array - OZB deals sent
		if err := json.Unmarshal([]byte(ozbSent), &user.OzbSent); err != nil {
			udb.Logger.Error("Error unmarshalling OZB deals sent", zap.Error(err))
		}

		// Bytes to string array - AMZ deals sent
		if err := json.Unmarshal([]byte(amzSent), &user.AmzSent); err != nil {
			udb.Logger.Error("Error unmarshalling AMZ deals sent", zap.Error(err))
		}

		userStore.Users[user.ChatID] = user

	}

	return userStore, nil
}

// Write *models.UserStore to the database
func (udb *UserStoreDB) WriteUserStore(userStore *models.UserStore) error {
	for _, user := range userStore.Users {
		// Convert string array to bytes
		// We do this as sqlite does not allow us to store string slices
		// Instead we convert to JSON bytes and store in the database
		keywords, err := json.Marshal(user.Keywords)
		if err != nil {
			udb.Logger.Error("Error marshalling user keywords", zap.Error(err))
		}

		ozbSent, err := json.Marshal(user.OzbSent)
		if err != nil {
			udb.Logger.Error("Error marshalling OZB deals sent", zap.Error(err))
		}

		amzSent, err := json.Marshal(user.AmzSent)
		if err != nil {
			udb.Logger.Error("Error marshalling AMZ deals sent", zap.Error(err))
		}

		_, err = udb.DB.Exec(`
			INSERT INTO users (
				chat_id, username, ozb_good, ozb_super, keywords, ozb_sent, amz_daily, amz_weekly, amz_sent
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(chat_id) DO UPDATE SET
				username = ?, ozb_good = ?, ozb_super = ?, keywords = ?, ozb_sent = ?, amz_daily =?, amz_weekly =?, amz_sent =?
			`,
			user.ChatID, user.Username, user.OzbGood, user.OzbSuper, keywords, ozbSent, user.AmzDaily, user.AmzWeekly, amzSent,
			user.Username, user.OzbGood, user.OzbSuper, keywords, ozbSent, user.AmzDaily, user.AmzWeekly, amzSent,
		)

		if err != nil {
			udb.Logger.Error("Error adding user", zap.String("error", err.Error()))
			return err
		}
	}

	return nil
}
