package handler

import (
	"DiscordBotAgent/internal/client"
	"context"
	"runtime/debug"

	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type ButtonsHandler struct {
	log *zap_logger.Logger
}

func NewButtonsHandler(log *zap_logger.Logger) *ButtonsHandler {
	return &ButtonsHandler{
		log: log,
	}
}

func (h *ButtonsHandler) Handle(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	wrapper client.ButtonWrapper,
) {
	defer func() {
		if r := recover(); r != nil {
			h.log.WithCtx(ctx).Error(
				"button handler panicked",
				zap.Any("error", r),
				zap.String("stack", string(debug.Stack())),
				zap.String("module", wrapper.ModuleName),
				zap.String("button_id", wrapper.Btn.ID()),
			)
		}
	}()

	err := wrapper.Btn.Execute(ctx, s, i)
	if err != nil {
		h.log.WithCtx(ctx).Error(
			"button execution failed",
			zap.Error(err),
			zap.String("module", wrapper.ModuleName),
			zap.String("button_id", wrapper.Btn.ID()),
		)
	}
}
