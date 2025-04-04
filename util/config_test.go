package util

import (
	"os"
	"testing"

	"go.uber.org/zap/zapcore"
)

// Test configuration values from ./config.yaml
//
// Language: go
// Path: util/config_test.go

func TestSetupConfig(t *testing.T) {
	// Create a test logger
	logger := SetupLogger(zapcore.DebugLevel, false)

	// Test case 1: Valid config file
	t.Run("Valid config file", func(t *testing.T) {
		// Create a temporary config file
		configContent := `
log_level: -1
log_to_file: true
test_mode: false
sqlite:
  db_path: "data/users.db"
scrapers:
  ozbargain:
    scrape_interval: 5
    max_stored_deals: 100
  amazon:
    urls: ["https://camelcamelcamel.com"]
    scrape_interval: 10
    max_stored_deals: 50
    target_price_drop: 20
pipup:
  enabled: false
  username: "test"
  base_url: "http://localhost:7979"
  duration: 5
  media_type: "image"
  media_uri: "https://example.com/image.jpg"
  image_width: 200
  position: 2
  title_color: "#FFFFFF"
  title_size: 20
  message_color: "#FFFFFF"
  message_size: 16
  background_color: "#000000"
`
		tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.WriteString(configContent); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		// Test SetupConfig
		config, err := SetupConfig(tmpFile.Name(), logger)
		if err != nil {
			t.Fatalf("SetupConfig failed: %v", err)
		}

		// Test config values
		if config.LogLevel != -1 {
			t.Errorf("Expected LogLevel -1, got %d", config.LogLevel)
		}
		if !config.LogToFile {
			t.Errorf("Expected LogToFile true, got %v", config.LogToFile)
		}
		if config.TestMode {
			t.Errorf("Expected TestMode false, got %v", config.TestMode)
		}
		if config.SQLite.DBPath != "data/users.db" {
			t.Errorf("Expected DBPath 'data/users.db', got '%s'", config.SQLite.DBPath)
		}
		if config.Scrapers.OzBargain.ScrapeInterval != 5 {
			t.Errorf("Expected OzBargain ScrapeInterval 5, got %d", config.Scrapers.OzBargain.ScrapeInterval)
		}
		if config.Scrapers.OzBargain.MaxStoredDeals != 100 {
			t.Errorf("Expected OzBargain MaxStoredDeals 100, got %d", config.Scrapers.OzBargain.MaxStoredDeals)
		}
		if config.Scrapers.Amazon.ScrapeInterval != 10 {
			t.Errorf("Expected Amazon ScrapeInterval 10, got %d", config.Scrapers.Amazon.ScrapeInterval)
		}
		if config.Scrapers.Amazon.MaxStoredDeals != 50 {
			t.Errorf("Expected Amazon MaxStoredDeals 50, got %d", config.Scrapers.Amazon.MaxStoredDeals)
		}
		if config.Scrapers.Amazon.TargetPriceDrop != 20 {
			t.Errorf("Expected Amazon TargetPriceDrop 20, got %d", config.Scrapers.Amazon.TargetPriceDrop)
		}
		if config.Pipup.Enabled {
			t.Errorf("Expected Pipup Enabled false, got %v", config.Pipup.Enabled)
		}
		if config.Pipup.Username != "test" {
			t.Errorf("Expected Pipup Username 'test', got '%s'", config.Pipup.Username)
		}
	})

	// Test case 2: Invalid config file - should fall back to defaults
	t.Run("Invalid config file", func(t *testing.T) {
		config, err := SetupConfig("nonexistent.yaml", logger)
		if err != nil {
			t.Fatalf("SetupConfig should not return error for nonexistent file: %v", err)
		}

		// Verify default values
		if config.LogLevel != -1 {
			t.Errorf("Expected default LogLevel -1, got %d", config.LogLevel)
		}
		if !config.LogToFile {
			t.Errorf("Expected default LogToFile true, got %v", config.LogToFile)
		}
		if config.TestMode {
			t.Errorf("Expected default TestMode false, got %v", config.TestMode)
		}
		if config.SQLite.DBPath != "data/users.db" {
			t.Errorf("Expected default DBPath 'data/users.db', got '%s'", config.SQLite.DBPath)
		}
		if config.Scrapers.OzBargain.ScrapeInterval != 5 {
			t.Errorf("Expected default OzBargain ScrapeInterval 5, got %d", config.Scrapers.OzBargain.ScrapeInterval)
		}
		if config.Scrapers.OzBargain.MaxStoredDeals != 250 {
			t.Errorf("Expected default OzBargain MaxStoredDeals 250, got %d", config.Scrapers.OzBargain.MaxStoredDeals)
		}
		if config.Scrapers.Amazon.ScrapeInterval != 30 {
			t.Errorf("Expected default Amazon ScrapeInterval 30, got %d", config.Scrapers.Amazon.ScrapeInterval)
		}
		if config.Scrapers.Amazon.MaxStoredDeals != 250 {
			t.Errorf("Expected default Amazon MaxStoredDeals 250, got %d", config.Scrapers.Amazon.MaxStoredDeals)
		}
		if config.Scrapers.Amazon.TargetPriceDrop != 20 {
			t.Errorf("Expected default Amazon TargetPriceDrop 20, got %d", config.Scrapers.Amazon.TargetPriceDrop)
		}
		if config.Pipup.Enabled {
			t.Errorf("Expected default Pipup Enabled false, got %v", config.Pipup.Enabled)
		}
	})
}
