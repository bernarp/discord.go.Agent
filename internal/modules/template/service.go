package template

import (
	"context"
	"fmt"
	"runtime"
	"time"

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

var startTime = time.Now()

func (s *Service) BuildStatusEmbed(latency time.Duration) *discordgo.MessageEmbed {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := time.Since(startTime)
	uptimeStr := fmt.Sprintf(
		"%dd %dh %dm %ds",
		int(uptime.Hours())/24,
		int(uptime.Hours())%24,
		int(uptime.Minutes())%60,
		int(uptime.Seconds())%60,
	)

	color := 0x00FF00
	if latency > 200*time.Millisecond {
		color = 0xFFA500
	}
	if latency > 500*time.Millisecond {
		color = 0xFF0000
	}

	return &discordgo.MessageEmbed{
		Title:       "🤖 System Status",
		Color:       color,
		Description: "Operational metrics for DiscordBotAgent",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "📡 API Latency",
				Value:  fmt.Sprintf("`%d ms`", latency.Milliseconds()),
				Inline: true,
			},
			{
				Name:   "⏱ Uptime",
				Value:  fmt.Sprintf("`%s`", uptimeStr),
				Inline: true,
			},
			{
				Name:   "🧊 Goroutines",
				Value:  fmt.Sprintf("`%d`", runtime.NumGoroutine()),
				Inline: true,
			},
			{
				Name:   "💾 Memory (Alloc)",
				Value:  fmt.Sprintf("`%v MB`", bToMb(memStats.Alloc)),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Updated at %s", time.Now().Format("15:04:05")),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
