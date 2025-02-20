package utils

import (
	"strings"
)

func Tokenize(inp []byte) []string {
	return strings.Fields(string(inp))
}
