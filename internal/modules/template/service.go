package template

import (
	"context"

	"DiscordBotAgent/internal/core/zap_logger"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	cfg Config,
) {
	fields := make([]zapcore.Field, 0, 5)

	if cfg.LogDetails.Guild {
		fields = append(fields, zap.String("guild_id", m.GuildID))
	}
	if cfg.LogDetails.Channel {
		fields = append(fields, zap.String("channel_id", m.ChannelID))
	}
	if cfg.LogDetails.Author {
		fields = append(fields, zap.String("author_name", m.Author.Username))
		fields = append(fields, zap.String("author_id", m.Author.ID))
	}
	if cfg.LogDetails.Content {
		fields = append(fields, zap.String("content", m.Content))
	}

	s.log.WithCtx(ctx).Info("message processed by template module", fields...)
}
