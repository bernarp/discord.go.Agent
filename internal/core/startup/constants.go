package startup

const (
	CryptoSalt = "agent-secure-salt-v1"
	AppName    = "DiscordBotAgent"
)

// Ключи в .env
const (
	EnvTokenKey   = "BOT_TOKEN"
	EnvPrefixKey  = "PREFIX"
	EnvAppIDKey   = "APP_ID"
	EnvGuildIDKey = "GUILD_ID"
)

// Имена зашифрованных файлов
const (
	FileToken   = "token.enc"
	FilePrefix  = "prefix.enc"
	FileAppID   = "appid.enc"
	FileGuildID = "guildid.enc"
)

var EnvFiles = []string{".env.dev", ".env"}

const (
	//nolint:gosec
	TokenRegexPattern = `^[\w-]{24,28}\.[\w-]{6}\.[\w-]{27,45}$`
	IDRegexPattern    = `^\d{17,20}$` // Для AppID и GuildID
)
