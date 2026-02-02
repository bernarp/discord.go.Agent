package eventbus

const (
	MessageCreate EventType = "message.create"
	MessageUpdate EventType = "message.update"
	MessageDelete EventType = "message.delete"

	GuildCreate EventType = "guild.create"
	GuildDelete EventType = "guild.delete"

	ReadyDiscordGateway EventType = "discordapi.bot.ready"
)
