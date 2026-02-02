package template2

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	service *Service
	module  *Module
}

func NewHandler(
	service *Service,
	module *Module,
) *Handler {
	return &Handler{
		service: service,
		module:  module,
	}
}

func (h *Handler) OnMessageCreate(
	ctx context.Context,
	payload any,
) {
	cfg := h.module.GetConfig()

	if !cfg.Enabled {
		return
	}

	event, ok := payload.(*discordgo.MessageCreate)
	if !ok {
		return
	}

	h.service.ProcessMessage(ctx, event, cfg)
}
