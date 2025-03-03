package server

import (
	"net"

	"github.com/google/uuid"
)

func genClientID() string {
	return uuid.New().String()
}

func Start(address string) {
	connection, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	for {
		client, err := connection.Accept()
		writeToClient(client, "Welcome to vaultic\n")
		writeToClient(client, "> ")
		if err != nil {
			panic(err)
		}
		go handleClient(client)
	}
}
