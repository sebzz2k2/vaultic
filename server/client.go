package server

import (
	"io"
	"net"
	"strings"

	"github.com/sebzz2k2/vaultic/logger"
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

func handleClient(client io.Reader, clientID string) {
	logger.Infof("Client connected with ID: %s", clientID)
	for {
		buff, beof := readBuffer(client)
		if beof {
			logger.Infof("Client disconnected with ID: %s", clientID)
			break
		}
		tokens := utils.Tokenize(buff)
		isValidCmd := utils.ValidateCmd(tokens[0])
		if !isValidCmd {
			writeToClient(client, "Invalid command\n")
		}
		isValidArgCount := utils.IsValidArgsCount(tokens[0], len(tokens)-1)
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
