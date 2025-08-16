package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(path string, config Configuration) (*zap.Logger, error) {
	// Encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	if config.Debug {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// File sink with rotasi log
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path + time.Now().Format("20060102") + ".log",
		MaxSize:    10, // MB
		MaxBackups: 7,
		MaxAge:     28, // days
		Compress:   true,
	})

	// Stdout sink
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Combine into one core
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, zap.InfoLevel),
		zapcore.NewCore(encoder, consoleWriter, zap.InfoLevel),
	)

	// create logger
	logger := zap.New(core, zap.AddCaller())
	return logger, nil

}
