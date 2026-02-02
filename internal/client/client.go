package client

import (
	"fmt"

	"DiscordBotAgent/internal/core/config_env"
	"DiscordBotAgent/internal/core/eventbus"
	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	Session *discordgo.Session
	log     *zap_logger.Logger
	eb      *eventbus.EventBus
}

func New(
	cfg *config.Config,
	log *zap_logger.Logger,
	eb *eventbus.EventBus,
) (*Client, error) {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("discord session: %w", err)
	}

	c := &Client{
		Session: session,
		log:     log,
		eb:      eb,
	}

	c.registerInternalHandlers()

	return c, nil
}

func (c *Client) registerInternalHandlers() {
	c.Session.AddHandler(
		func(
			s *discordgo.Session,
			m *discordgo.MessageCreate,
		) {
			if m.Author.Bot {
				return
			}
			c.eb.Publish(eventbus.MessageCreate, m)
		},
	)

	c.Session.AddHandler(
		func(
			s *discordgo.Session,
			r *discordgo.Ready,
		) {
			c.eb.Publish(eventbus.ReadyDiscordGateway, r)
		},
	)
}

func (c *Client) Connect() error {
	c.log.Info("connecting to discord gateway")
	if err := c.Session.Open(); err != nil {
		return fmt.Errorf("gateway connect: %w", err)
	}
	return nil
}

func (c *Client) Disconnect() error {
	c.log.Info("closing discord connection")
	return c.Session.Close()
}
