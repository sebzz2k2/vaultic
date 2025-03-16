package server

import (
	"net"

	"github.com/sebzz2k2/vaultic/logger"
	"github.com/sebzz2k2/vaultic/utils"
)

func Start(address string) {
	b := NewIndexBuilder(utils.FILENAME, utils.DELIMITER[0])
	err := b.BuildIndexes()
	if err != nil {
		logger.Errorf("Error building index %s", err.Error())
		return
	}
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
