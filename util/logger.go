package util

// import zap logger
import "go.uber.org/zap"

// function to setup the logger
func SetupLogger() *zap.Logger {
	// create the logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// return the logger
	return logger
}
