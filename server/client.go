package server

import (
	"io"
	"net"
	"strings"

	"github.com/sebzz2k2/vaultic/cmd"
	"github.com/sebzz2k2/vaultic/lexer"
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
func parse(tokens []lexer.Token) {
}
func handleClient(client net.Conn) {
	for {
		buff, beof := readBuffer(client)
		if beof {
			break
		}
		tkns := lexer.Tokenize(string(buff))
		parse(tkns) // This is a dummy function to show how to use the lexer
		tokens := utils.Tokenize(buff)
		cmd := cmd.CommandFactory(tokens[0])
		if cmd == nil {
			utils.WriteToClient(client, "Invalid command\n")
			utils.WriteToClient(client, "> ")
			continue
		}
		isValidArgCount := cmd.Validate(len(tokens) - 1)
		if !isValidArgCount {
			utils.WriteToClient(client, "Expected syntax is: "+utils.CmdArgsErrors[strings.ToLower(tokens[0])]+"\n")
			utils.WriteToClient(client, "> ")
			continue
		}

		val, err := cmd.Process(tokens[1:])
		if err != nil {
			utils.WriteToClient(client, err.Error())
		} else {

			utils.WriteToClient(client, val+"\n")
		}
		utils.WriteToClient(client, "> ")
	}
}
