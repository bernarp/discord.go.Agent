package client

import (
	"context"
	"sync"

	config "DiscordBotAgent/internal/core/config_env"
	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type CommandWrapper struct {
	Cmd        CommandSlash
	ModuleName string
}

type InteractionHandler interface {
	Handle(
		ctx context.Context,
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
		wrapper CommandWrapper,
	)
}

type Manager struct {
	log        *zap_logger.Logger
	session    *discordgo.Session
	commands   map[string]CommandWrapper
	cmdHandler InteractionHandler
	mu         sync.RWMutex
}

func NewInteraction(
	log *zap_logger.Logger,
	session *discordgo.Session,
	cmdHandler InteractionHandler,
) *Manager {
	return &Manager{
		log:        log,
		session:    session,
		commands:   make(map[string]CommandWrapper),
		cmdHandler: cmdHandler,
	}
}

func (m *Manager) Register(
	cmd CommandSlash,
	moduleName string,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	info := cmd.Info()
	m.commands[info.Name] = CommandWrapper{
		Cmd:        cmd,
		ModuleName: moduleName,
	}
	m.log.Info("registered command", zap.String("name", info.Name), zap.String("module", moduleName))
}

func (m *Manager) SyncCommands(cfg *config.Config) error {
	var cmdsToCreate []*discordgo.ApplicationCommand
	m.mu.RLock()
	for _, w := range m.commands {
		info := w.Cmd.Info()
		cmdsToCreate = append(
			cmdsToCreate, &discordgo.ApplicationCommand{
				Name:        info.Name,
				Description: info.Description,
				Options:     info.Options,
			},
		)
	}
	m.mu.RUnlock()
	_, err := m.session.ApplicationCommandBulkOverwrite(cfg.AppID, cfg.GuildID, cmdsToCreate)
	return err
}

func (m *Manager) HandleInteraction(
	ctx context.Context,
	payload any,
) {
	event, ok := payload.(*discordgo.InteractionCreate)
	if !ok {
		return
	}

	var (
		iType      = event.Type.String()
		targetName = "unknown"
		user       = "unknown"
	)

	if event.Member != nil && event.Member.User != nil {
		user = event.Member.User.Username + "#" + event.Member.User.Discriminator
	} else if event.User != nil {
		user = event.User.Username + "#" + event.User.Discriminator
	}

	switch event.Type {
	case discordgo.InteractionApplicationCommand:
		targetName = event.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		targetName = event.MessageComponentData().CustomID
	case discordgo.InteractionModalSubmit:
		targetName = event.ModalSubmitData().CustomID
	case discordgo.InteractionApplicationCommandAutocomplete:
		targetName = event.ApplicationCommandData().Name
	}

	m.log.WithCtx(ctx).Debug(
		"interaction received",
		zap.String("type", iType),
		zap.String("target", targetName),
		zap.String("user", user),
		zap.String("guild_id", event.GuildID),
		zap.String("channel_id", event.ChannelID),
	)

	switch event.Type {
	case discordgo.InteractionApplicationCommand:
		m.mu.RLock()
		wrapper, exists := m.commands[targetName]
		m.mu.RUnlock()

		if !exists {
			m.log.Debug("command handler not found", zap.String("name", targetName))
			return
		}
		m.cmdHandler.Handle(ctx, m.session, event, wrapper)

	case discordgo.InteractionModalSubmit:
		m.log.Debug("modal routing not implemented", zap.String("id", targetName))
	}
}
