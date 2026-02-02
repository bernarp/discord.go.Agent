package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func New() (*Config, error) {
	if err := godotenv.Load(".env.dev"); err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("config: BOT_TOKEN is empty")
	}

	return &Config{
		BotToken: token,
	}, nil
}
