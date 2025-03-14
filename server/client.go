package server

import (
	"fmt"
	"io"
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
			utils.WriteToClient(client, "Invalid command\n")
		}
		isValidArgCount := cmd.Validate(len(tokens) - 1)
		if !isValidArgCount {
			fmt.Println(len(tokens) - 1)
			utils.WriteToClient(client, "Expected syntax is: "+utils.CmdArgsErrors[strings.ToLower(tokens[0])]+"\n")
			continue
		}
		fmt.Println("Validated")

		val, err := cmd.Process(tokens[1:])
		if err != nil {
			utils.WriteToClient(client, err.Error())
		} else {

			utils.WriteToClient(client, val+"\n")
		}
		utils.WriteToClient(client, "> ")
	}
}
