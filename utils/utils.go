package utils

import (
	"io"
	"strings"
)

func WriteToClient(client io.Writer, message string) {
	client.Write([]byte(message))
}
func Tokenize(inp []byte) []string {
	return strings.Fields(string(inp))
}
