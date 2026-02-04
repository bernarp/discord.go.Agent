package buttons

import (
	"context"

	"DiscordBotAgent/internal/modules/template"

	"github.com/bwmarrin/discordgo"
)

const (
	BtnStatusRefresh = "btn_status_refresh"
	BtnStatusDelete  = "btn_status_delete"
)

type RefreshButton struct {
	Service *template.Service
}

func (b *RefreshButton) ID() string {
	return BtnStatusRefresh
}

func (b *RefreshButton) Execute(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) error {
	newEmbed := b.Service.BuildStatusEmbed(s.HeartbeatLatency())

	return s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{newEmbed},
			},
		},
	)
}

type DeleteButton struct{}

func (b *DeleteButton) ID() string {
	return BtnStatusDelete
}

func (b *DeleteButton) Execute(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) error {

	err := s.InteractionRespond(
		i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		},
	)
	if err != nil {
		return err
	}
	return s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
}
