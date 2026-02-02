package client

import (
	"fmt"

	"DiscordBotAgent/internal/core/config_env"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Client struct {
	Session *discordgo.Session
	log     *zap.Logger
}

func New(
	cfg *config.Config,
	log *zap.Logger,
) (*Client, error) {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("discord session: %w", err)
	}

	return &Client{
		Session: session,
		log:     log,
	}, nil
}

func (c *Client) Connect() error {
	c.log.Info("connecting to discord gateway...")
	if err := c.Session.Open(); err != nil {
		return fmt.Errorf("gateway connect: %w", err)
	}
	c.log.Info("connected successfully")
	return nil
}

func (c *Client) Disconnect() error {
	c.log.Info("disconnecting from discord...")
	if err := c.Session.Close(); err != nil {
		return fmt.Errorf("gateway disconnect: %w", err)
	}
	return nil
}
