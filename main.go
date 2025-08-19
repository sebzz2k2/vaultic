package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sebzz2k2/vaultic/pkg/config"
	"github.com/sebzz2k2/vaultic/pkg/logger"
)

const (
	AppName = "Vaultic"
	Version = "0.1.0"

	shutdownTimeout = 30 * time.Second
)

func main() {
	app := &Application{}

	if err := app.Run(); err != nil {
		log.Fatal().Err(err).Msg("Application failed to run")
	}
}

type Application struct {
	config *config.Config
	// server *server.Server
	// engine storage.StorageEngine
}

func (app *Application) Run() error {
	if err := app.initConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	if err := app.initLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	fmt.Println(app.config.Logging.ToConsole, app.config.Logging.ToFile, "49")

	log.Info().
		Str("app", AppName).
		Str("version", Version).
		Msg("Starting Vaultic Key-Value Store")
	return nil
}

func (app *Application) initConfig() error {
	cfg, err := config.LoadConfig("vaultic_config.yaml")
	if err != nil {
		return err
	}
	app.config = &cfg
	return nil
}

func (app *Application) initLogger() error {
	level, err := zerolog.ParseLevel(app.config.Logging.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	err = logger.Setup(logger.Config{
		Level:     level,
		Console:   app.config.Logging.ToConsole,
		LogToFile: app.config.Logging.ToFile,
		FilePath:  app.config.Logging.Path,
	})
	if err != nil {
		return err
	}
	return nil
}

// func maintemp() {
// 	err := config.InitConfig()
// 	if err != nil {
// 		panic(fmt.Errorf(config.ErrorFailedToLoadConfig, err))
// 	}

// 	zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// 	err = logger.Setup(logger.Config{
// 		Level:     zerolog.DebugLevel,
// 		Console:   true,
// 		LogToFile: true,
// 		FilePath:  config.Global.Logging.Path,
// 	})
// 	if err != nil {
// 		log.Fatal().Err(err).Msg(config.ErrorFailedToSetUpLogger)
// 	}

// 	log.Info().Msg(config.InfoStartingServer)
// 	log.Info().Msg(config.InfoBuildingIndex)
// 	b := index.NewIndexBuilder(utils.FILENAME)
// 	// get time it takes to build index
// 	start := time.Now()
// 	err = b.BuildIndexes()
// 	if err != nil {
// 		log.Error().Err(err).Msg(config.ErrorBuildIndex)
// 		return
// 	}
// 	duration := time.Since(start)
// 	log.Info().Msgf(config.InfoIndexBuiltTime, duration)
// 	log.Info().Msg(config.InfoFinishedIndex)
// 	// creates a task to monitor memtable in the background
// 	// go b.monitorMemtable(size in bytes)

// 	server.Start(fmt.Sprintf(":%d", config.Global.Port))
// }
