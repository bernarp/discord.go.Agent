package client

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Button interface {
	ID() string
	Execute(
		ctx context.Context,
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) error
}
