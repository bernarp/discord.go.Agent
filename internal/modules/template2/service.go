package template2

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
	return &Service{log: log}
}

func (s *Service) ProcessMessage(
	ctx context.Context,
	m *discordgo.MessageCreate,
	cfg Config,
) {
	s.log.WithCtx(ctx).Info(
		"template2 processed message",
		zap.String("prefix", cfg.Prefix),
		zap.String("content", m.Content),
		zap.Int("max_logs", cfg.MaxLogs),
	)
}
