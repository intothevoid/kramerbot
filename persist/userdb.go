package persist

import (
	"database/sql"
	"encoding/json"

	"github.com/intothevoid/kramerbot/models"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type UserDB struct {
	DB     *sql.DB
	Name   string
	Logger *zap.Logger
}

// Connect to the database
func CreateDatabaseConnection(dbName string, logger *zap.Logger) *UserDB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	return &UserDB{
		DB:     db,
		Name:   dbName,
		Logger: logger,
	}
}

// Close the database
func (udb *UserDB) Close() {
	err := udb.DB.Close()
	if err != nil {
		udb.Logger.Error("Error closing database", zap.Error(err))
	}
}

// Create *models.UserData table in database
func (udb *UserDB) CreateTable() error {
	_, err := udb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			chat_id INTEGER PRIMARY KEY,
			username TEXT,
			good_deals INTEGER,
			super_deals INTEGER,
			keywords BLOB,
			deals_sent BLOB
			);
		`)

	if err != nil {
		return err
	}

	return nil
}

// Add user to the database
func (udb *UserDB) AddUser(user *models.UserData) error {
	// Convert string array to bytes
	// We do this as sqlite does not allow us to store string slices
	// Instead we convert to JSON bytes and store in the database
	keywords, err := json.Marshal(user.Keywords)
	if err != nil {
		udb.Logger.Error("Error marshalling user keywords", zap.Error(err))
	}

	dealsSent, err := json.Marshal(user.DealsSent)
	if err != nil {
		udb.Logger.Error("Error marshalling user deals sent", zap.Error(err))
	}

	_, err = udb.DB.Exec(`
		INSERT INTO users (
				chat_id, username, good_deals, super_deals, keywords, deals_sent
			) VALUES (?, ?, ?, ?, ?, ?)`,
		user.ChatID, user.Username, user.GoodDeals, user.SuperDeals, keywords, dealsSent,
	)

	if err != nil {
		udb.Logger.Error("Error adding user", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// Update user in the database
func (udb *UserDB) UpdateUser(user *models.UserData) error {
	// Convert string array to bytes
	// We do this as sqlite does not allow us to store string slices
	// Instead we convert to JSON bytes and store in the database
	keywords, err := json.Marshal(user.Keywords)
	if err != nil {
		udb.Logger.Error("Error marshalling user keywords", zap.Error(err))
	}

	dealsSent, err := json.Marshal(user.DealsSent)
	if err != nil {
		udb.Logger.Error("Error marshalling user deals sent", zap.Error(err))
	}

	_, err = udb.DB.Exec(`
		UPDATE users SET
			username = ?, good_deals = ?, super_deals = ?, keywords = ?, deals_sent = ?
		WHERE chat_id = ?`,
		user.Username, user.GoodDeals, user.SuperDeals, keywords, dealsSent, user.ChatID,
	)

	if err != nil {
		udb.Logger.Error("Error updating user", zap.Error(err))
		return err
	}

	return nil
}

// Delete user from the database
func (udb *UserDB) DeleteUser(user *models.UserData) error {
	_, err := udb.DB.Exec(`DELETE FROM users WHERE chat_id = ?`, user.ChatID)

	if err != nil {
		udb.Logger.Error("Error deleting user", zap.Error(err))
		return err
	}

	return nil
}

// Get user from the database by chat_id
func (udb *UserDB) GetUser(chatID int64) (*models.UserData, error) {
	user := &models.UserData{}
	keywords := []byte{}
	dealsSent := []byte{}

	err := udb.DB.QueryRow(`SELECT * FROM users WHERE chat_id = ?`, chatID).Scan(
		&user.ChatID, &user.Username, &user.GoodDeals, &user.SuperDeals, &keywords, &dealsSent,
	)
	if err != nil {
		udb.Logger.Error("Error getting user", zap.Error(err))
		return nil, err
	}

	// Bytes to string array - keywords
	if err := json.Unmarshal([]byte(keywords), &user.Keywords); err != nil {
		udb.Logger.Error("Error unmarshalling user keywords", zap.Error(err))
	}

	// Bytes to string array - dealsSent
	if err := json.Unmarshal([]byte(dealsSent), &user.DealsSent); err != nil {
		udb.Logger.Error("Error unmarshalling user deals sent", zap.Error(err))
	}

	return user, nil
}
