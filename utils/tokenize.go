package utils

import (
	"strings"
)

func Tokenize(inp []byte) []string {
	return strings.Fields(string(inp))
}

func ValidateCmd(token string) bool {
	for _, command := range validCommands {
		if strings.ToLower(token) == command {
			return true
		}
	}
	return false
}

func IsValidArgsCount(token string, argsCount int) bool {
	return cmdArgs[strings.ToLower(token)] == argsCount
}
