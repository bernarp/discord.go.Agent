package startup

const (
	CryptoSalt = "agent-secure-salt-v1"
	AppName    = "DiscordBotAgent"
)

const (
	EnvTokenKey   = "BOT_TOKEN"
	EnvPrefixKey  = "PREFIX"
	EnvAppIDKey   = "APP_ID"
	EnvGuildIDKey = "GUILD_ID"
	EnvPortKey    = "HTTP_PORT"
)

const (
	FileToken   = "token.enc"
	FilePrefix  = "prefix.enc"
	FileAppID   = "appid.enc"
	FileGuildID = "guildid.enc"
	FilePort    = "port.enc"
)

var EnvFiles = []string{".env.dev", ".env"}

const (
	TokenRegexPattern = `^[\w-]{24,28}\.[\w-]{6}\.[\w-]{27,45}$`
	IDRegexPattern    = `^\d{17,20}$`
	PortRegexPattern  = `^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
)
