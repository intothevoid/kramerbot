package util

// import zap logger
import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// function to setup the logger
func SetupLogger() *zap.Logger {
	// configure the logger
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// log to file
	logFile, err := os.OpenFile("log.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zap.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	// create the logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// return the logger
	return logger
}
