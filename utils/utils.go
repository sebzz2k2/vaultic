package utils

import (
	"io"
	"net"
	"strings"
)

func WriteToClient(client io.Reader, message string) {
	if conn, ok := client.(net.Conn); ok {
		conn.Write([]byte(message))
	}
}
func Tokenize(inp []byte) []string {
	return strings.Fields(string(inp))
}
