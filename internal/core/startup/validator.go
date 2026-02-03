package startup

import (
	"regexp"
	"strings"
)

var (
	tokenRegex = regexp.MustCompile(TokenRegexPattern)
	idRegex    = regexp.MustCompile(IDRegexPattern)
	portRegex  = regexp.MustCompile(PortRegexPattern)
)

func ValidateToken(token string) bool { return tokenRegex.MatchString(token) }
func ValidateID(id string) bool       { return idRegex.MatchString(id) }
func ValidatePrefix(p string) bool    { return len(p) > 0 && len(p) < 5 }
func ValidatePort(p string) bool      { return portRegex.MatchString(p) }

func CleanInput(in string) string {
	in = strings.TrimSpace(in)
	in = strings.Trim(in, `"'`)
	return in
}
