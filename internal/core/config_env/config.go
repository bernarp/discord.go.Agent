package config

import (
	"DiscordBotAgent/internal/core/startup"
	"fmt"
)

type Config struct {
	BotToken string
}

func New() (*Config, error) {
	token, err := startup.GetBotToken()
	if err != nil {
		return nil, fmt.Errorf("startup: %w", err)
	}

	return &Config{
		BotToken: token,
	}, nil
}
