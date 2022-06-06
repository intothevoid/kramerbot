package util

// import zap logger
import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// function to setup the logger
func SetupLogger(logLevel zapcore.Level, fileMode bool) *zap.Logger {
	// configure the logger
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	var core zapcore.Core
	if fileMode {
		// log to file
		logFile, err := os.OpenFile("log.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		writer := zapcore.AddSync(logFile)
		defaultLogLevel := logLevel
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
	} else {
		defaultLogLevel := logLevel
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
	}

	// create the logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// return the logger
	return Logger
}
