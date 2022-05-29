package util

import (
	"os"

	"github.com/spf13/viper"
)

// Setup viper config
func SetupConfig(conf_path string) (*viper.Viper, error) {
	viper := viper.New()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	configPath := conf_path

	// Check if path exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return viper, err
	}

	// Load the config file
	viper.SetConfigFile(configPath)

	// Read the config file
	viper.ReadInConfig()

	return viper, nil
}
