package startup

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type StartupConfig struct {
	Token   string
	Prefix  string
	AppID   string
	GuildID string
	Port    string
}

func GetStartupConfig() (*StartupConfig, error) {
	for _, file := range EnvFiles {
		_ = godotenv.Load(file)
	}

	key, err := deriveKey()
	if err != nil {
		return nil, err
	}

	conf := &StartupConfig{}
	conf.Token = loadOrPrompt(EnvTokenKey, FileToken, "BOT_TOKEN", "MzE4...", false, ValidateToken, key)
	conf.Prefix = loadOrPrompt(EnvPrefixKey, FilePrefix, "PREFIX", "!", false, ValidatePrefix, key)
	conf.AppID = loadOrPrompt(EnvAppIDKey, FileAppID, "APP_ID", "1129747578100000000", false, ValidateID, key)
	conf.GuildID = loadOrPrompt(EnvGuildIDKey, FileGuildID, "GUILD_ID", "1074340840000000000", true, ValidateID, key)
	conf.Port = loadOrPrompt(EnvPortKey, FilePort, "HTTP_PORT", "8080", true, ValidatePort, key)

	fmt.Println("--------------------------------------------------")
	return conf, nil
}

func loadOrPrompt(
	envKey, fileName, label, example string,
	optional bool,
	validator func(string) bool,
	key []byte,
) string {
	if val := CleanInput(os.Getenv(envKey)); val != "" {
		if validator(val) {
			return val
		}
		fmt.Printf("(!) %s in environment is invalid.\n", label)
	}

	if _, err := os.Stat(fileName); err == nil {
		if data, err := os.ReadFile(fileName); err == nil {
			if decrypted, err := decrypt(data, key); err == nil {
				val := CleanInput(string(decrypted))
				if val == "" && optional {
					return ""
				}
				if validator(val) {
					return val
				}
			}
		}
		fmt.Printf("(!) Saved %s is invalid or corrupted.\n", label)
	}

	val := promptGeneric(label, example, optional, validator)
	if val != "" {
		_ = saveEncrypted(fileName, val, key)
	}
	return val
}
