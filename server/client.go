package server

import (
	"io"
	"net"
	"strings"

	"github.com/sebzz2k2/vaultic/cmd"
	"github.com/sebzz2k2/vaultic/utils"
)

func readBuffer(reader io.Reader) ([]byte, bool) {
	b := make([]byte, 1024)
	bn, err := reader.Read(b)
	if err != nil {
		if err == io.EOF {
			return nil, true
		}
		panic(err)
	}
	return b[:bn], false
}

func handleClient(client io.Reader) {
	for {
		buff, beof := readBuffer(client)
		if beof {
			break
		}
		tokens := utils.Tokenize(buff)
		cmd := cmd.CommandFactory(tokens[0])
		if cmd == nil {
			writeToClient(client, "Invalid command\n")
		}
		isValidArgCount := cmd.Validate(len(tokens) - 1)
		if !isValidArgCount {
			sendInvalidArgumentResponse(client, tokens[0])
		}
		showPrompt(client)
	}
}

func showPrompt(client io.Reader) {
	if conn, ok := client.(net.Conn); ok {
		writeToClient(conn, "> ")
	}
}

func writeToClient(client io.Reader, message string) {
	if conn, ok := client.(net.Conn); ok {
		conn.Write([]byte(message))
	}
}

func sendInvalidArgumentResponse(client io.Reader, token string) {
	writeToClient(client, "Expected syntax is: "+utils.CmdArgsErrors[strings.ToLower(token)]+"\n")
}
