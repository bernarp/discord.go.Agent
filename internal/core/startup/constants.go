package startup

const (
	TokenFileName = "token.enc"
	CryptoSalt    = "agent-secure-salt-v1"
	EnvTokenKey   = "BOT_TOKEN"
	AppName       = "DiscordBotAgent"
)

var EnvFiles = []string{
	".env.dev",
	".env",
}

//nolint:gosec
const TokenRegexPattern = `^[\w-]{24,28}\.[\w-]{6}\.[\w-]{27,45}$`
