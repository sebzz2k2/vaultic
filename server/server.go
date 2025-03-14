package server

import (
	"net"

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
		utils.WriteToClient(client, "Welcome to vaultic\n")
		utils.WriteToClient(client, "> ")
		if err != nil {
			panic(err)
		}
		go handleClient(client)
	}
}
