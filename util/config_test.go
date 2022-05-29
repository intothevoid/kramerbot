package util

import (
	"os"
	"path"
	"testing"

	"go.uber.org/zap/zapcore"
)

// Test configuration values from ./config.yaml
//
// Language: go
// Path: util/config_test.go

func TestConfig(t *testing.T) {

	// Setup configuration
	confPath, _ := os.Getwd()
	confPath = path.Join(confPath, "../resources/config_test.yaml")
	appconf, err := SetupConfig(confPath)
	if err != nil {
		t.Error("Error reading application config")
	}

	// Test log level
	if zapcore.Level(appconf.GetInt("log_level")) != zapcore.WarnLevel {
		t.Error("Error reading log level")
	}

	// Test log to file
	if appconf.GetBool("log_to_file") != true {
		t.Error("Error reading log to file")
	}

	// Test scrape interval
	if appconf.GetInt("scrapers.ozbargain.scrape_interval") != 50 {
		t.Error("Error reading scrape interval")
	}

	// Test max deals
	if appconf.GetInt("scrapers.ozbargain.max_stored_deals") != 999 {
		t.Error("Error reading max deals")
	}
}
