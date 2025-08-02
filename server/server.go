package server

import (
	"net"

	"github.com/sebzz2k2/vaultic/pkg/config"
	"github.com/sebzz2k2/vaultic/utils"
)

func Start(address string) {
	connection, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	for {
		client, err := connection.Accept()
		utils.WriteToClient(client, config.WelcomeMessage)
		utils.WriteToClient(client, config.PromptMessage)
		if err != nil {
			panic(err)
		}
		go handleClient(client)
	}
}
