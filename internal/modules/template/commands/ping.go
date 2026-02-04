package commands

import (
	"context"

	"DiscordBotAgent/internal/client"
	"DiscordBotAgent/internal/modules/template"
	"DiscordBotAgent/internal/modules/template/buttons"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct {
	Service *template.Service
}

func (c *PingCommand) Info() client.CommandInfo {
	return client.CommandInfo{
		Name:        "status",
		Description: "Show bot system status",
		Type:        client.CmdGuild,
	}
}

func (c *PingCommand) Execute(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) error {
	embed := c.Service.BuildStatusEmbed(s.HeartbeatLatency())
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Refresh",
					Style:    discordgo.PrimaryButton,
					CustomID: buttons.BtnStatusRefresh,
					Emoji:    &discordgo.ComponentEmoji{Name: "🔄"},
				},
				discordgo.Button{
					Label:    "Delete",
					Style:    discordgo.DangerButton,
					CustomID: buttons.BtnStatusDelete,
					Emoji:    &discordgo.ComponentEmoji{Name: "🗑️"},
				},
			},
		},
	}

	_, err := s.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{embed},
			Components: &components,
		},
	)

	return err
}
