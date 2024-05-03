package env

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func GetString(key string, defaultValue string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}
	return value
}

func GetInt(key string, defaultValue int) int {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func GetBool(key string, defaultValue bool) bool {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func GetLogLevel(key string, defaultValue slog.Level) slog.Level {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	var logLevel slog.Level

	switch strings.ToLower(value) {
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = defaultValue
	}

	return logLevel
}
