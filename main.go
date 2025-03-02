package main

import (
	"github.com/sebzz2k2/vaultic/logger"
	"github.com/sebzz2k2/vaultic/server"
)

func main() {
	logger.Infof("Starting Vaultic server")
	server.Start(":5381")
}
