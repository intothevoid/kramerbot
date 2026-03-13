package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/intothevoid/kramerbot/models"
)

// createWebUsersTableSQL creates the core table (no indexes).
const createWebUsersTableSQL = `
CREATE TABLE IF NOT EXISTS web_users (
	id                    TEXT PRIMARY KEY,
	email                 TEXT UNIQUE NOT NULL,
	password_hash         TEXT NOT NULL,
	display_name          TEXT,
	email_verified        INTEGER NOT NULL DEFAULT 0,
	verify_token          TEXT,
	verify_token_expires  DATETIME,
	telegram_chat_id      INTEGER,
	telegram_username     TEXT,
	link_token            TEXT,
	link_token_expires    DATETIME,
	reset_token           TEXT,
	reset_token_expires   DATETIME,
	ozb_good              INTEGER NOT NULL DEFAULT 0,
	ozb_super             INTEGER NOT NULL DEFAULT 0,
	amz_daily             INTEGER NOT NULL DEFAULT 0,
	amz_weekly            INTEGER NOT NULL DEFAULT 0,
	email_summary         INTEGER NOT NULL DEFAULT 0,
	keywords              TEXT NOT NULL DEFAULT '[]',
	created_at            DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at            DATETIME DEFAULT CURRENT_TIMESTAMP
)`

// migrateStmts adds columns or indexes that may be missing on an existing database.
// All statements are executed best-effort (errors ignored) so both old and new
// databases work without a manual migration step.
var migrateStmts = []string{
	// New columns added after initial release.
	`ALTER TABLE web_users ADD COLUMN reset_token TEXT`,
	`ALTER TABLE web_users ADD COLUMN reset_token_expires DATETIME`,
	`ALTER TABLE web_users ADD COLUMN ozb_good INTEGER NOT NULL DEFAULT 0`,
	`ALTER TABLE web_users ADD COLUMN ozb_super INTEGER NOT NULL DEFAULT 0`,
	`ALTER TABLE web_users ADD COLUMN amz_daily INTEGER NOT NULL DEFAULT 0`,
	`ALTER TABLE web_users ADD COLUMN amz_weekly INTEGER NOT NULL DEFAULT 0`,
	`ALTER TABLE web_users ADD COLUMN keywords TEXT NOT NULL DEFAULT '[]'`,
	// Email verification columns.
	`ALTER TABLE web_users ADD COLUMN email_verified INTEGER NOT NULL DEFAULT 0`,
	`ALTER TABLE web_users ADD COLUMN verify_token TEXT`,
	`ALTER TABLE web_users ADD COLUMN verify_token_expires DATETIME`,
	// Daily email summary preference.
	`ALTER TABLE web_users ADD COLUMN email_summary INTEGER NOT NULL DEFAULT 0`,
	// Indexes — created after columns to avoid "no such column" on old schemas.
	`CREATE INDEX IF NOT EXISTS idx_web_users_email ON web_users(email)`,
	`CREATE INDEX IF NOT EXISTS idx_web_users_link_token ON web_users(link_token)`,
	`CREATE INDEX IF NOT EXISTS idx_web_users_reset_token ON web_users(reset_token)`,
	`CREATE INDEX IF NOT EXISTS idx_web_users_verify_token ON web_users(verify_token)`,
}

// CreateWebUsersTable creates the web_users table and migrates any missing columns/indexes.
func (udb *UserStoreDB) CreateWebUsersTable() error {
	if _, err := udb.DB.Exec(createWebUsersTableSQL); err != nil {
		return fmt.Errorf("failed to create web_users table: %w", err)
	}
	// Best-effort: add missing columns/indexes. Errors are intentionally ignored
	// (SQLite returns an error for duplicate columns/existing indexes).
	for _, stmt := range migrateStmts {
		udb.DB.Exec(stmt) //nolint:errcheck
	}
	return nil
}

// webUserColumns is the explicit column list used in all SELECT queries.
const webUserColumns = `
	id, email, password_hash, display_name,
	email_verified, verify_token, verify_token_expires,
	telegram_chat_id, telegram_username,
	link_token, link_token_expires,
	reset_token, reset_token_expires,
	ozb_good, ozb_super, amz_daily, amz_weekly, email_summary, keywords,
	created_at, updated_at`

// CreateWebUser inserts a new web user record.
func (udb *UserStoreDB) CreateWebUser(user *models.WebUser) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	if user.Keywords == nil {
		user.Keywords = []string{}
	}
	kw, _ := json.Marshal(user.Keywords)
	_, err := udb.DB.Exec(`
		INSERT INTO web_users
			(id, email, password_hash, display_name,
			 ozb_good, ozb_super, amz_daily, amz_weekly, keywords,
			 created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash, user.DisplayName,
		user.OzbGood, user.OzbSuper, user.AmzDaily, user.AmzWeekly, string(kw),
		now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create web user: %w", err)
	}
	return nil
}

// GetWebUserByEmail retrieves a web user by email address.
func (udb *UserStoreDB) GetWebUserByEmail(email string) (*models.WebUser, error) {
	return udb.scanWebUser(udb.DB.QueryRow(
		`SELECT `+webUserColumns+` FROM web_users WHERE email = ?`, email,
	))
}

// GetWebUserByID retrieves a web user by UUID.
func (udb *UserStoreDB) GetWebUserByID(id string) (*models.WebUser, error) {
	return udb.scanWebUser(udb.DB.QueryRow(
		`SELECT `+webUserColumns+` FROM web_users WHERE id = ?`, id,
	))
}

// GetWebUserByLinkToken retrieves a web user by a non-expired Telegram link token.
func (udb *UserStoreDB) GetWebUserByLinkToken(token string) (*models.WebUser, error) {
	return udb.scanWebUser(udb.DB.QueryRow(
		`SELECT `+webUserColumns+` FROM web_users WHERE link_token = ? AND link_token_expires > ?`,
		token, time.Now().UTC(),
	))
}

// GetWebUserByVerifyToken retrieves a web user by a non-expired email verification token.
func (udb *UserStoreDB) GetWebUserByVerifyToken(token string) (*models.WebUser, error) {
	return udb.scanWebUser(udb.DB.QueryRow(
		`SELECT `+webUserColumns+` FROM web_users WHERE verify_token = ? AND verify_token_expires > ?`,
		token, time.Now().UTC(),
	))
}

// GetWebUserByResetToken retrieves a web user by a non-expired password reset token.
func (udb *UserStoreDB) GetWebUserByResetToken(token string) (*models.WebUser, error) {
	return udb.scanWebUser(udb.DB.QueryRow(
		`SELECT `+webUserColumns+` FROM web_users WHERE reset_token = ? AND reset_token_expires > ?`,
		token, time.Now().UTC(),
	))
}

// UpdateWebUser updates all mutable fields of a web user record.
func (udb *UserStoreDB) UpdateWebUser(user *models.WebUser) error {
	user.UpdatedAt = time.Now().UTC()
	if user.Keywords == nil {
		user.Keywords = []string{}
	}
	kw, _ := json.Marshal(user.Keywords)
	_, err := udb.DB.Exec(`
		UPDATE web_users SET
			email = ?,
			password_hash = ?,
			display_name = ?,
			email_verified = ?,
			verify_token = ?,
			verify_token_expires = ?,
			telegram_chat_id = ?,
			telegram_username = ?,
			link_token = ?,
			link_token_expires = ?,
			reset_token = ?,
			reset_token_expires = ?,
			ozb_good = ?,
			ozb_super = ?,
			amz_daily = ?,
			amz_weekly = ?,
			email_summary = ?,
			keywords = ?,
			updated_at = ?
		WHERE id = ?`,
		user.Email, user.PasswordHash, user.DisplayName,
		user.EmailVerified, user.VerifyToken, user.VerifyTokenExpires,
		user.TelegramChatID, user.TelegramUsername,
		user.LinkToken, user.LinkTokenExpires,
		user.ResetToken, user.ResetTokenExpires,
		user.OzbGood, user.OzbSuper, user.AmzDaily, user.AmzWeekly, user.EmailSummary, string(kw),
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

// GetAllVerifiedWebUsers returns all web users whose email has been verified.
func (udb *UserStoreDB) GetAllVerifiedWebUsers() ([]*models.WebUser, error) {
	rows, err := udb.DB.Query(`SELECT ` + webUserColumns + ` FROM web_users WHERE email_verified = 1`)
	if err != nil {
		return nil, fmt.Errorf("failed to query verified web users: %w", err)
	}
	defer rows.Close()

	var users []*models.WebUser
	for rows.Next() {
		u := &models.WebUser{}
		var kwJSON string
		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName,
			&u.EmailVerified, &u.VerifyToken, &u.VerifyTokenExpires,
			&u.TelegramChatID, &u.TelegramUsername,
			&u.LinkToken, &u.LinkTokenExpires,
			&u.ResetToken, &u.ResetTokenExpires,
			&u.OzbGood, &u.OzbSuper, &u.AmzDaily, &u.AmzWeekly, &u.EmailSummary, &kwJSON,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan web user: %w", err)
		}
		if kwJSON == "" || kwJSON == "null" {
			u.Keywords = []string{}
		} else {
			json.Unmarshal([]byte(kwJSON), &u.Keywords) //nolint:errcheck
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// scanWebUser reads an explicit-column row into a WebUser struct.
func (udb *UserStoreDB) scanWebUser(row *sql.Row) (*models.WebUser, error) {
	u := &models.WebUser{}
	var kwJSON string
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName,
		&u.EmailVerified, &u.VerifyToken, &u.VerifyTokenExpires,
		&u.TelegramChatID, &u.TelegramUsername,
		&u.LinkToken, &u.LinkTokenExpires,
		&u.ResetToken, &u.ResetTokenExpires,
		&u.OzbGood, &u.OzbSuper, &u.AmzDaily, &u.AmzWeekly, &u.EmailSummary, &kwJSON,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan web user: %w", err)
	}
	if kwJSON == "" || kwJSON == "null" {
		u.Keywords = []string{}
	} else {
		json.Unmarshal([]byte(kwJSON), &u.Keywords) //nolint:errcheck
	}
	return u, nil
}
