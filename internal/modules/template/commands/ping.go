package commands

import (
	"context"
	"fmt"

	"DiscordBotAgent/internal/client"
	"github.com/bwmarrin/discordgo"
)

type PingCommand struct {
	Session *discordgo.Session
}

func (c *PingCommand) Info() client.CommandInfo {
	return client.CommandInfo{
		Name:        "ping",
		Description: "Check bot latency",
		Type:        client.CmdGuild,
	}
}

func (c *PingCommand) Execute(
	ctx context.Context,
	i *discordgo.InteractionCreate,
) error {
	latency := c.Session.HeartbeatLatency().Milliseconds()
	content := fmt.Sprintf("🏓 Pong! Latency: `%dms`", latency)
	_, err := c.Session.InteractionResponseEdit(
		i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		},
	)
	return err
}
