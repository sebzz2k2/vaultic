package main

import (
	"github.com/sebzz2k2/vaultic/logger"
	"github.com/sebzz2k2/vaultic/server"
	"github.com/sebzz2k2/vaultic/utils"
)

func main() {
	logger.Infof("Starting Vaultic server")
	logger.Infof("Building index")
	b := server.NewIndexBuilder(utils.FILENAME)
	err := b.BuildIndexes()
	if err != nil {
		logger.Errorf("Error building index %s", err.Error())
		return
	}
	logger.Infof("Finished building index")
	server.Start(":5381")
}
