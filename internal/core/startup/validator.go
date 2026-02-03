package startup

import (
	"regexp"
	"strings"
)

var tokenRegex = regexp.MustCompile(TokenRegexPattern)

func ValidateToken(token string) bool {
	return tokenRegex.MatchString(token)
}

func CleanToken(token string) string {
	token = strings.TrimSpace(token)
	token = strings.Trim(token, `"'`)
	return token
}
