package template

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) OnMessageCreate(
	ctx context.Context,
	payload any,
) {
	event, ok := payload.(*discordgo.MessageCreate)
	if !ok {
		return
	}

	h.service.LogMessageDetails(ctx, event)
}
