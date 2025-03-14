package utils

import (
	"io"
	"net"
)

func WriteToClient(client io.Reader, message string) {
	if conn, ok := client.(net.Conn); ok {
		conn.Write([]byte(message))
	}
}
