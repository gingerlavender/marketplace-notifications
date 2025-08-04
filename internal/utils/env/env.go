package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		var chatIds []string

		for chatId := range strings.SplitSeq(value, ",") {
			chatId = strings.TrimSpace(chatId)
			if chatId != "" {
				chatIds = append(chatIds, chatId)
			}
		}

		return chatIds
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
