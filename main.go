package main

import (
	"github.com/sebzz2k2/vaultic/logger"
	"github.com/sebzz2k2/vaultic/server"
	"github.com/sebzz2k2/vaultic/utils"
)

func main() {
	// encoded := storage.EncodeData(1, false, 1234567890, "kekjvkgvkgvkhvkhblbljblblhblby", "valueyuy")

	// logger.Infof("Encoded data: %v", encoded)
	// decoded, err := storage.DecodeData(encoded)
	// if err != nil {
	// 	logger.Errorf("Error decoding data: %s", err.Error())
	// 	return
	// }
	// logger.Infof("Decoded data: %v", decoded)
	logger.Infof("Starting Vaultic server")
	logger.Infof("Building index")
	b := server.NewIndexBuilder(utils.FILENAME, utils.DELIMITER[0])
	err := b.BuildIndexes()
	if err != nil {
		logger.Errorf("Error building index %s", err.Error())
		return
	}
	logger.Infof("Finished building index")
	server.Start(":5381")
}
