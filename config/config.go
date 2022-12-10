package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config defines the service config.
type Config struct {
	Port    int
	Spotify Spotify
	Logger
}

// Spotify defines the Spotify API auth config.
type Spotify struct {
	ClientID     string
	ClientSecret string
}

// New initialises a Config from environment variables.
func New() Config {
	debug := flag.Bool("debug", false, "log level")
	flag.Parse()

	logLevel := zapcore.InfoLevel
	if *debug == true {
		logLevel = zapcore.DebugLevel
	}
	logger := newLogger(logLevel)

	// attempt to get config environment vars, or default them
	return Config{
		Port: getEnvVarInt(logger, "PORT", 8080),
		Spotify: Spotify{
			ClientID:     getEnvVar(logger, "SPOTIFY_CLIENT_ID", ""),
			ClientSecret: getEnvVar(logger, "SPOTIFY_CLIENT_SECRET", ""),
		},
		Logger: logger,
	}
}

// getEnvVar gets a string environment variable or defaults it if unset.
func getEnvVar(logger Logger, key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Warn("environment var undefined, defaulting",
			zap.String("var", key),
			zap.String("default", defaultValue),
		)
		return defaultValue
	}

	logger.Info("environment var found", zap.String("var", key))
	return val
}

// getEnvVarInt gets an integer environment variable or defaults it if unset.
func getEnvVarInt(logger Logger, key string, defaultValue int) int {
	varStr := getEnvVar(logger, key, strconv.Itoa(defaultValue))
	varInt, _ := strconv.Atoi(varStr)
	return varInt
}

// Logger defines the required logger methods.
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zap.Field) *zap.Logger
}

func newLogger(level zapcore.Level) Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.EncoderConfig.TimeKey = "ts"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.Level.SetLevel(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Printf("failed to initialise logger: %s", err)
		os.Exit(1)
	}

	return logger
}
