package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"

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

// Connect to the database
func CreateDatabaseConnection(dbName string, logger *zap.Logger) *UserStoreDB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	return &UserStoreDB{
		DB:     db,
		Name:   dbName,
		Logger: logger,
	}
}

// Close the database
func (udb *UserStoreDB) Close() {
	err := udb.DB.Close()
	if err != nil {
		udb.Logger.Error("Error closing database", zap.Error(err))
	}
}

// Create *models.UserData table in database
func (udb *UserStoreDB) CreateTable() error {
	// Set WAL mode and other pragmas first
	_, err := udb.DB.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		udb.Logger.Error("Failed to set WAL mode", zap.Error(err))
		return err
	}

	_, err = udb.DB.Exec("PRAGMA synchronous=NORMAL")
	if err != nil {
		udb.Logger.Error("Failed to set synchronous mode", zap.Error(err))
		return err
	}

	_, err = udb.DB.Exec("PRAGMA busy_timeout=5000")
	if err != nil {
		udb.Logger.Error("Failed to set busy timeout", zap.Error(err))
		return err
	}

	// Create the table
	_, err = udb.DB.Exec(`
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
		`)

	if err != nil {
		udb.Logger.Error("Failed to create table", zap.Error(err))
		return err
	}

	return nil
}

// Add user to the database
func (udb *UserStoreDB) AddUser(user *models.UserData) error {
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
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ChatID, user.Username, user.OzbGood, user.OzbSuper, keywords, ozbSent, user.AmzDaily, user.AmzWeekly, amzSent,
	)

	if err != nil {
		udb.Logger.Error("Error adding user", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// Update user in the database
func (udb *UserStoreDB) UpdateUser(user *models.UserData) error {
	udb.Logger.Debug("Starting user update",
		zap.Int64("chatID", user.ChatID),
		zap.String("username", user.Username),
		zap.Bool("amz_daily", user.AmzDaily))

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

	result, err := udb.DB.Exec(`
		UPDATE users SET
			username = ?, ozb_good = ?, ozb_super = ?, keywords = ?, ozb_sent = ?, amz_daily = ?, amz_weekly = ?, amz_sent = ?
		WHERE chat_id = ?`,
		user.Username, user.OzbGood, user.OzbSuper, keywords, ozbSent, user.AmzDaily, user.AmzWeekly, amzSent, user.ChatID,
	)

	if err != nil {
		udb.Logger.Error("Error updating user", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		udb.Logger.Error("Error getting rows affected", zap.Error(err))
		return err
	}

	udb.Logger.Debug("User update completed",
		zap.Int64("chatID", user.ChatID),
		zap.Int64("rowsAffected", rowsAffected))

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated for chat_id %d", user.ChatID)
	}

	return nil
}

// Delete user from the database
func (udb *UserStoreDB) DeleteUser(user *models.UserData) error {
	_, err := udb.DB.Exec(`DELETE FROM users WHERE chat_id = ?`, user.ChatID)

	if err != nil {
		udb.Logger.Error("Error deleting user", zap.Error(err))
		return err
	}

	return nil
}

// Get user from the database by chat_id
func (udb *UserStoreDB) GetUser(chatID int64) (*models.UserData, error) {
	user := &models.UserData{}
	keywords := []byte{}
	ozbSent := []byte{}
	amzSent := []byte{}

	err := udb.DB.QueryRow(`SELECT * FROM users WHERE chat_id = ?`, chatID).Scan(
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
