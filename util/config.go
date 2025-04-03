package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds the application configuration
type Config struct {
	LogLevel  int  `mapstructure:"log_level"`
	LogToFile bool `mapstructure:"log_to_file"`
	TestMode  bool `mapstructure:"test_mode"`
	SQLite    SQLiteConfig
	Scrapers  ScrapersConfig
	Pipup     PipupConfig
}

// SQLiteConfig holds SQLite database configuration
type SQLiteConfig struct {
	DBPath string `mapstructure:"db_path"`
}

// ScrapersConfig holds configuration for all scrapers
type ScrapersConfig struct {
	OzBargain OzBargainConfig `mapstructure:"ozbargain"`
	Amazon    AmazonConfig    `mapstructure:"amazon"`
}

// OzBargainConfig holds OzBargain scraper configuration
type OzBargainConfig struct {
	ScrapeInterval int `mapstructure:"scrape_interval"`
	MaxStoredDeals int `mapstructure:"max_stored_deals"`
}

// AmazonConfig holds Amazon scraper configuration
type AmazonConfig struct {
	ScrapeInterval  int      `mapstructure:"scrape_interval"`
	MaxStoredDeals  int      `mapstructure:"max_stored_deals"`
	URLs            []string `mapstructure:"urls"`
	TargetPriceDrop int      `mapstructure:"target_price_drop"`
}

// PipupConfig holds Android TV notification configuration
type PipupConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	Username        string `mapstructure:"username"`
	BaseURL         string `mapstructure:"base_url"`
	Duration        int    `mapstructure:"duration"`
	MediaType       string `mapstructure:"media_type"`
	MediaURI        string `mapstructure:"media_uri"`
	ImageWidth      int    `mapstructure:"image_width"`
	Position        int    `mapstructure:"position"`
	TitleColor      string `mapstructure:"title_color"`
	TitleSize       int    `mapstructure:"title_size"`
	MessageColor    string `mapstructure:"message_color"`
	MessageSize     int    `mapstructure:"message_size"`
	BackgroundColor string `mapstructure:"background_color"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		LogLevel:  -1, // debug
		LogToFile: true,
		TestMode:  false,
		SQLite: SQLiteConfig{
			DBPath: "data/users.db",
		},
		Scrapers: ScrapersConfig{
			OzBargain: OzBargainConfig{
				ScrapeInterval: 5,
				MaxStoredDeals: 250,
			},
			Amazon: AmazonConfig{
				ScrapeInterval: 30,
				MaxStoredDeals: 250,
				URLs: []string{
					"https://au.camelcamelcamel.com/top_drops/feed?t=daily&",
					"https://au.camelcamelcamel.com/top_drops/feed?t=weekly&",
				},
				TargetPriceDrop: 20,
			},
		},
		Pipup: PipupConfig{
			Enabled:         false,
			Username:        "",
			BaseURL:         "http://localhost:7979/notify",
			Duration:        10,
			MediaType:       "video",
			MediaURI:        "",
			ImageWidth:      200,
			Position:        2,
			TitleColor:      "#ffffff",
			TitleSize:       14,
			MessageColor:    "#ffffff",
			MessageSize:     12,
			BackgroundColor: "#000000",
		},
	}
}

// ValidateConfig validates the configuration
func ValidateConfig(config *Config) error {
	// Validate log level
	if config.LogLevel < -1 || config.LogLevel > 2 {
		return fmt.Errorf("invalid log level: %d", config.LogLevel)
	}

	// Validate SQLite config
	if config.SQLite.DBPath == "" {
		return fmt.Errorf("sqlite.db_path cannot be empty")
	}

	// Validate OzBargain config
	if config.Scrapers.OzBargain.ScrapeInterval < 1 {
		return fmt.Errorf("ozbargain.scrape_interval must be at least 1 minute")
	}
	if config.Scrapers.OzBargain.MaxStoredDeals < 1 {
		return fmt.Errorf("ozbargain.max_stored_deals must be at least 1")
	}

	// Validate Amazon config
	if config.Scrapers.Amazon.ScrapeInterval < 1 {
		return fmt.Errorf("amazon.scrape_interval must be at least 1 minute")
	}
	if config.Scrapers.Amazon.MaxStoredDeals < 1 {
		return fmt.Errorf("amazon.max_stored_deals must be at least 1")
	}
	if len(config.Scrapers.Amazon.URLs) == 0 {
		return fmt.Errorf("amazon.urls cannot be empty")
	}
	if config.Scrapers.Amazon.TargetPriceDrop < 0 {
		return fmt.Errorf("amazon.target_price_drop cannot be negative")
	}

	// Validate Pipup config if enabled
	if config.Pipup.Enabled {
		if config.Pipup.Username == "" {
			return fmt.Errorf("pipup.username cannot be empty when pipup is enabled")
		}
		if config.Pipup.BaseURL == "" {
			return fmt.Errorf("pipup.base_url cannot be empty when pipup is enabled")
		}
		if config.Pipup.Duration < 1 {
			return fmt.Errorf("pipup.duration must be at least 1 second")
		}
		if !isValidMediaType(config.Pipup.MediaType) {
			return fmt.Errorf("invalid pipup.media_type: %s", config.Pipup.MediaType)
		}
		if config.Pipup.MediaType != "web" && config.Pipup.MediaURI == "" {
			return fmt.Errorf("pipup.media_uri cannot be empty when media_type is not 'web'")
		}
	}

	return nil
}

// isValidMediaType checks if the media type is valid
func isValidMediaType(mediaType string) bool {
	switch mediaType {
	case "video", "image", "web":
		return true
	default:
		return false
	}
}

// SetupConfig initializes and validates the configuration
func SetupConfig(confPath string, logger *zap.Logger) (*Config, error) {
	// Create default config
	config := DefaultConfig()

	// Initialize Viper
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")

	// Set environment variable support
	v.SetEnvPrefix("KRAMERBOT")
	v.AutomaticEnv()

	// Set default values
	v.SetDefault("log_level", config.LogLevel)
	v.SetDefault("log_to_file", config.LogToFile)
	v.SetDefault("test_mode", config.TestMode)
	v.SetDefault("sqlite.db_path", config.SQLite.DBPath)
	v.SetDefault("scrapers.ozbargain.scrape_interval", config.Scrapers.OzBargain.ScrapeInterval)
	v.SetDefault("scrapers.ozbargain.max_stored_deals", config.Scrapers.OzBargain.MaxStoredDeals)
	v.SetDefault("scrapers.amazon.scrape_interval", config.Scrapers.Amazon.ScrapeInterval)
	v.SetDefault("scrapers.amazon.max_stored_deals", config.Scrapers.Amazon.MaxStoredDeals)
	v.SetDefault("scrapers.amazon.urls", config.Scrapers.Amazon.URLs)
	v.SetDefault("scrapers.amazon.target_price_drop", config.Scrapers.Amazon.TargetPriceDrop)
	v.SetDefault("pipup.enabled", config.Pipup.Enabled)
	v.SetDefault("pipup.username", config.Pipup.Username)
	v.SetDefault("pipup.base_url", config.Pipup.BaseURL)
	v.SetDefault("pipup.duration", config.Pipup.Duration)
	v.SetDefault("pipup.media_type", config.Pipup.MediaType)
	v.SetDefault("pipup.media_uri", config.Pipup.MediaURI)
	v.SetDefault("pipup.image_width", config.Pipup.ImageWidth)
	v.SetDefault("pipup.position", config.Pipup.Position)
	v.SetDefault("pipup.title_color", config.Pipup.TitleColor)
	v.SetDefault("pipup.title_size", config.Pipup.TitleSize)
	v.SetDefault("pipup.message_color", config.Pipup.MessageColor)
	v.SetDefault("pipup.message_size", config.Pipup.MessageSize)
	v.SetDefault("pipup.background_color", config.Pipup.BackgroundColor)

	// Check if config file exists
	if confPath != "" {
		if _, err := os.Stat(confPath); os.IsNotExist(err) {
			logger.Warn("Config file not found, using defaults and environment variables",
				zap.String("path", confPath))
		} else {
			v.SetConfigFile(confPath)
			if err := v.ReadInConfig(); err != nil {
				logger.Warn("Error reading config file, using defaults and environment variables",
					zap.String("path", confPath),
					zap.Error(err))
			}
		}
	}

	// Unmarshal config
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate config
	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Override SQLite path with environment variable if set
	if envPath := os.Getenv("SQLITE_DB_PATH"); envPath != "" {
		config.SQLite.DBPath = envPath
	}

	// Ensure database directory exists
	dbDir := filepath.Dir(config.SQLite.DBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	return config, nil
}
