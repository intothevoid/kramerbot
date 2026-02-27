package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/intothevoid/kramerbot/models"
)

// CreateWebUsersTable creates the web_users table if it does not exist.
func (udb *UserStoreDB) CreateWebUsersTable() error {
	_, err := udb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS web_users (
			id                   TEXT PRIMARY KEY,
			email                TEXT UNIQUE NOT NULL,
			password_hash        TEXT NOT NULL,
			display_name         TEXT,
			telegram_chat_id     INTEGER,
			telegram_username    TEXT,
			link_token           TEXT,
			link_token_expires   DATETIME,
			created_at           DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at           DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_web_users_email ON web_users(email);
		CREATE INDEX IF NOT EXISTS idx_web_users_link_token ON web_users(link_token);
	`)
	if err != nil {
		return fmt.Errorf("failed to create web_users table: %w", err)
	}
	return nil
}

// CreateWebUser inserts a new web user record.
func (udb *UserStoreDB) CreateWebUser(user *models.WebUser) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := udb.DB.Exec(`
		INSERT INTO web_users (id, email, password_hash, display_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash, user.DisplayName, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create web user: %w", err)
	}
	return nil
}

// GetWebUserByEmail retrieves a web user by email address.
func (udb *UserStoreDB) GetWebUserByEmail(email string) (*models.WebUser, error) {
	return udb.scanWebUser(
		udb.DB.QueryRow(`SELECT * FROM web_users WHERE email = ?`, email),
	)
}

// GetWebUserByID retrieves a web user by UUID.
func (udb *UserStoreDB) GetWebUserByID(id string) (*models.WebUser, error) {
	return udb.scanWebUser(
		udb.DB.QueryRow(`SELECT * FROM web_users WHERE id = ?`, id),
	)
}

// GetWebUserByLinkToken retrieves a web user by a non-expired link token.
func (udb *UserStoreDB) GetWebUserByLinkToken(token string) (*models.WebUser, error) {
	return udb.scanWebUser(
		udb.DB.QueryRow(
			`SELECT * FROM web_users WHERE link_token = ? AND link_token_expires > ?`,
			token, time.Now().UTC(),
		),
	)
}

// UpdateWebUser updates an existing web user record.
func (udb *UserStoreDB) UpdateWebUser(user *models.WebUser) error {
	user.UpdatedAt = time.Now().UTC()
	_, err := udb.DB.Exec(`
		UPDATE web_users SET
			email = ?,
			password_hash = ?,
			display_name = ?,
			telegram_chat_id = ?,
			telegram_username = ?,
			link_token = ?,
			link_token_expires = ?,
			updated_at = ?
		WHERE id = ?`,
		user.Email, user.PasswordHash, user.DisplayName,
		user.TelegramChatID, user.TelegramUsername,
		user.LinkToken, user.LinkTokenExpires,
		user.UpdatedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update web user: %w", err)
	}
	return nil
}

// DeleteWebUser removes a web user by UUID.
func (udb *UserStoreDB) DeleteWebUser(id string) error {
	_, err := udb.DB.Exec(`DELETE FROM web_users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete web user: %w", err)
	}
	return nil
}

// scanWebUser reads a single row into a WebUser struct.
func (udb *UserStoreDB) scanWebUser(row *sql.Row) (*models.WebUser, error) {
	u := &models.WebUser{}
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName,
		&u.TelegramChatID, &u.TelegramUsername,
		&u.LinkToken, &u.LinkTokenExpires,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan web user: %w", err)
	}
	return u, nil
}
