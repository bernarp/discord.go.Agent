package handler

import (
	"context"
	"time"

	"DiscordBotAgent/internal/client"
	"DiscordBotAgent/internal/core/module_manager"
	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type CommandsHandler struct {
	log       *zap_logger.Logger
	moduleMgr *module_manager.Manager
}

func NewCommandsHandler(
	log *zap_logger.Logger,
	moduleMgr *module_manager.Manager,
) *CommandsHandler {
	return &CommandsHandler{
		log:       log,
		moduleMgr: moduleMgr,
	}
}

func (h *CommandsHandler) Handle(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	wrapper client.CommandWrapper,
) {
	start := time.Now()
	data := i.ApplicationCommandData()

	err := s.InteractionRespond(
		i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)
	if err != nil {
		h.log.Error("failed to send deferred response", zap.Error(err))
		return
	}

	if wrapper.ModuleName != "" {
		if !h.moduleMgr.IsModuleEnabled(wrapper.ModuleName) {
			h.log.Warn(
				"command blocked: module disabled",
				zap.String("command", data.Name),
				zap.String("module", wrapper.ModuleName),
			)
			h.sendEditError(s, i, "module disabled.")
			return
		}
	}

	if err := wrapper.Cmd.Execute(ctx, i); err != nil {
		h.log.Error(
			"command execution failed",
			zap.String("command", data.Name),
			zap.Error(err),
		)
		h.sendEditError(s, i, "internal error.")
		return
	}

	h.log.WithCtx(ctx).Debug(
		"command executed successfully",
		zap.String("command", data.Name),
		zap.Duration("duration", time.Since(start)),
		zap.String("user", i.Member.User.Username),
	)
}

func (h *CommandsHandler) sendEditError(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	msg string,
) {
	_, _ = s.InteractionResponseEdit(
		i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		},
	)
}
