package client

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type CommandType string

const (
	CmdGlobal CommandType = "global"
	CmdGuild  CommandType = "guild"
)

type CommandInfo struct {
	Name        string
	Description string
	Type        CommandType
	Options     []*discordgo.ApplicationCommandOption
}

type CommandSlash interface {
	Info() CommandInfo
	Execute(
		ctx context.Context,
		i *discordgo.InteractionCreate,
	) error
}
