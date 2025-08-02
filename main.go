package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sebzz2k2/vaultic/pkg/config"
	"github.com/sebzz2k2/vaultic/pkg/logger"
	"github.com/sebzz2k2/vaultic/server"
	"github.com/sebzz2k2/vaultic/utils"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err = logger.Setup(logger.Config{
		Level:     zerolog.DebugLevel,
		Console:   true,
		LogToFile: true,
		FilePath:  config.Global.LogPath,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to set up logger")
	}

	log.Info().Msg("Starting Vaultic server")
	log.Info().Msg("Building index")
	b := server.NewIndexBuilder(utils.FILENAME)
	err = b.BuildIndexes()
	if err != nil {
		log.Error().Err(err).Msg("Error building index")
		return
	}
	log.Info().Msg("Finished building index")
	server.Start(fmt.Sprintf(":%d", config.Global.Port))
}
