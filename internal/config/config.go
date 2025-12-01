package config

import (
	"fmt"
	"marketplace-notifications/internal/marketplaces/wb"
	"marketplace-notifications/internal/marketplaces/yandex"
	"marketplace-notifications/internal/utils/env"
	"time"
)

type Config struct {
	Server   ServerConfig
	Monitor  MonitorConfig
	API      APIConfig
	Telegram TelegramConfig
}

type ServerConfig struct {
	ControlToken string
	Port         int
}

type MonitorConfig struct {
	CheckInterval time.Duration
}

type APIConfig struct {
	WB      wb.Config
	Yandex  yandex.Config
	Timeout time.Duration
}

type TelegramConfig struct {
	BotToken string
	ChatIds  []string
	Timeout  time.Duration
	RPS      int
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         env.GetEnvInt("SERVER_PORT", 8080),
			ControlToken: env.GetEnv("CONTROL_TOKEN", ""),
		},
		Monitor: MonitorConfig{
			CheckInterval: env.GetEnvDuration("CHECK_INTERVAL", 2*time.Minute),
		},
		API: APIConfig{
			WB:      wb.GetConfig(env.GetEnv("WB_JWT", ""), env.GetEnvInt("MAX_NEW_QUESTIONS_TO_FETCH", 20), env.GetEnvInt("MAX_NEW_FEEDBACKS_TO_FETCH", 20)),
			Yandex:  yandex.GetConfig(env.GetEnv("YANDEX_TOKEN", "")),
			Timeout: env.GetEnvDuration("MARKETPLACE_API_TIMEOUT", 30*time.Second),
		},
		Telegram: TelegramConfig{
			BotToken: env.GetEnv("TELEGRAM_BOT_TOKEN", ""),
			ChatIds:  env.GetEnvStringSlice("TELEGRAM_CHAT_IDS", nil),
			Timeout:  env.GetEnvDuration("TELEGRAM_API_TIMEOUT", 30*time.Second),
			RPS:      1,
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return config, nil
}

func (config *Config) validate() error {
	if config.Server.ControlToken == "" {
		return fmt.Errorf("missing control token")
	}

	if config.API.WB.JWT == "" {
		return fmt.Errorf("missing WB_JWT")
	}
	if config.API.Yandex.APIToken == "" {
		return fmt.Errorf("missing YANDEX_TOKEN")
	}

	if config.Telegram.BotToken == "" {
		return fmt.Errorf("missing TELEGRAM_BOT_TOKEN")
	}
	if len(config.Telegram.ChatIds) == 0 {
		return fmt.Errorf("missing TELEGRAM_CHAT_IDS")
	}

	return nil
}
