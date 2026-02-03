package config

import (
	"DiscordBotAgent/internal/core/startup"
	"fmt"
)

type Config struct {
	BotToken string
	Prefix   string
	AppID    string
	GuildID  string
}

func New() (*Config, error) {
	sConf, err := startup.GetStartupConfig()
	if err != nil {
		return nil, fmt.Errorf("startup: %w", err)
	}

	return &Config{
		BotToken: sConf.Token,
		Prefix:   sConf.Prefix,
		AppID:    sConf.AppID,
		GuildID:  sConf.GuildID,
	}, nil
}
