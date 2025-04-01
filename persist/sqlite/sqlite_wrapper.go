package sqlite

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	persist_if "github.com/intothevoid/kramerbot/persist" // Alias import for the interface package
	"go.uber.org/zap"
)

// SQLiteWrapper wraps a UserStoreDB and implements the DatabaseIF interface.
// Note: We place this in the 'database' package to access UserStoreDB unexported fields if needed in future,
// but we need to be careful with package naming and imports.
type SQLiteWrapper struct {
	*UserStoreDB // Embed the existing SQLite implementation
}

// Ensure SQLiteWrapper implements DatabaseIF at compile time.
var _ persist_if.DatabaseIF = (*SQLiteWrapper)(nil)

// NewSQLiteWrapper creates a new SQLiteWrapper, initializes the database, and creates the table if needed.
func NewSQLiteWrapper(dbPath string, logger *zap.Logger) (*SQLiteWrapper, error) {
	// Default path if not provided
	if dbPath == "" {
		dbPath = "data/users.db"
	}

	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		logger.Info("Database directory does not exist, creating.", zap.String("directory", dbDir))
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			logger.Error("Failed to create database directory", zap.String("directory", dbDir), zap.Error(err))
			return nil, fmt.Errorf("failed to create database directory '%s': %w", dbDir, err)
		}
	} else if err != nil {
		// Handle other potential errors from os.Stat
		logger.Error("Failed to check database directory status", zap.String("directory", dbDir), zap.Error(err))
		return nil, fmt.Errorf("failed to check status of directory '%s': %w", dbDir, err)
	}

	// Create database connection using the function from userdb.go
	db := CreateDatabaseConnection(dbPath, logger)
	if db == nil || db.DB == nil {
		// CreateDatabaseConnection panics on error, but let's be safe
		return nil, fmt.Errorf("failed to create database connection for path '%s'", dbPath)
	}

	// Create table if it doesn't exist
	logger.Info("Ensuring database table exists", zap.String("path", dbPath))
	if err := db.CreateTable(); err != nil {
		logger.Error("Failed to create database table", zap.String("path", dbPath), zap.Error(err))
		db.Close() // Attempt to close the connection if table creation failed
		return nil, fmt.Errorf("failed to create table in database '%s': %w", dbPath, err)
	}

	logger.Info("SQLite database initialized successfully", zap.String("path", dbPath))
	return &SQLiteWrapper{
		UserStoreDB: db,
	}, nil
}

// Ping checks if the database connection is still active.
// For SQLite, we check if a simple query executes without error.
func (sw *SQLiteWrapper) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 second timeout for ping
	defer cancel()

	if sw.DB == nil {
		return fmt.Errorf("sqlite database connection is nil")
	}

	err := sw.DB.PingContext(ctx)
	if err != nil {
		sw.Logger.Error("Failed to ping SQLite database", zap.Error(err))
		return fmt.Errorf("sqlite ping failed: %w", err)
	}
	sw.Logger.Debug("SQLite database ping successful")
	return nil
}

// Close the database connection. This implements the Close method required by DatabaseIF.
// It wraps the existing Close method from UserStoreDB which doesn't return an error.
func (sw *SQLiteWrapper) Close() error {
	if sw.DB != nil {
		sw.UserStoreDB.Close() // Call the embedded Close (which logs errors internally)
	}
	return nil // Match the interface signature
}

// Note: AddUser, UpdateUser, DeleteUser, GetUser, ReadUserStore, WriteUserStore
// are already implemented by the embedded UserStoreDB and satisfy the interface.
