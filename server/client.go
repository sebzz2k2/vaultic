package server

import (
	"io"
	"net"

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

func handleClient(client net.Conn) {
	for {
		buff, beof := readBuffer(client)
		if beof {
			break
		}
		tkns := lexer.Tokenize(string(buff))
		val, err := cmd.ProcessCommand(tkns)
		if err != nil {
			utils.WriteToClient(client, err.Error()+"\n")
		} else {
			utils.WriteToClient(client, val+"\n")
		}
		utils.WriteToClient(client, "> ")
	}
}
