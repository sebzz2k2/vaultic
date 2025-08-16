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

func Crc32(data string) uint32 {
	var crc uint32 = 0xFFFFFFFF
	poly := uint32(0x04C11DB7)

	for _, b := range []byte(data) {
		crc ^= uint32(b) << 24
		for i := 0; i < 8; i++ {
			if (crc & 0x80000000) != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
	}
	return crc ^ 0xFFFFFFFF
}
