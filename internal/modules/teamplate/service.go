package template

import (
	"context"

	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Service struct {
	log *zap_logger.Logger
}

func NewService(log *zap_logger.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) LogMessageDetails(
	ctx context.Context,
	m *discordgo.MessageCreate,
) {
	s.log.WithCtx(ctx).Info(
		"new message received",
		zap.String("guild_id", m.GuildID),
		zap.String("channel_id", m.ChannelID),
		zap.String("author_id", m.Author.ID),
		zap.String("author_name", m.Author.Username),
		zap.String("content", m.Content),
	)
}
