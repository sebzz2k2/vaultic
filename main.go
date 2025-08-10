package main

import (
	"fmt"
	"os"
	"time"

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
		panic(fmt.Errorf(config.ErrorFailedToLoadConfig, err))
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
		log.Fatal().Err(err).Msg(config.ErrorFailedToSetUpLogger)
	}

	log.Info().Msg(config.InfoStartingServer)
	log.Info().Msg(config.InfoBuildingIndex)
	b := server.NewIndexBuilder(utils.FILENAME)
	// get time it takes to build index
	start := time.Now()
	err = b.BuildIndexes()
	if err != nil {
		log.Error().Err(err).Msg(config.ErrorBuildIndex)
		return
	}
	duration := time.Since(start)
	log.Info().Msgf(config.InfoIndexBuiltTime, duration)
	log.Info().Msg(config.InfoFinishedIndex)
	server.Start(fmt.Sprintf(":%d", config.Global.Port))
}
