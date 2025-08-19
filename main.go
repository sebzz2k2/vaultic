package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sebzz2k2/vaultic/internal/server"
	"github.com/sebzz2k2/vaultic/internal/storage"
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
	server *server.Server
	engine *storage.StorageEngine
}

func (app *Application) Run() error {
	if err := app.initConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	if err := app.initLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	log.Info().
		Str("app", AppName).
		Str("version", Version).
		Msg("Starting Vaultic Key-Value Store")
	if err := app.initStorageEngine(); err != nil {
		return fmt.Errorf("failed to initialize storage engine: %w", err)
	}
	if err := app.initServer(); err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	svrErrCh := make(chan error, 1)
	go func() {
		log.Info().
			Str("address", app.config.Server.Address).
			Int("port", app.config.Server.Port).
			Msg("Server initialized")
		if err := app.server.Start(); err != nil {
			svrErrCh <- fmt.Errorf("server error: %w", err)
		}
	}()
	return app.waitForShutdown(svrErrCh)
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
func (app *Application) initServer() error {
	log.Info().Msg("Initializing server")
	cfg := &server.Config{
		Address:        app.config.Server.Address,
		Port:           app.config.Server.Port,
		MaxConnections: app.config.Server.MaxConnections,
		MaxMessageSize: app.config.Server.MaxMessageSize,
	}

	svr, err := server.New(cfg, *app.engine)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	app.server = svr

	return nil
}

func (app *Application) initStorageEngine() error {
	log.Info().Msg("Initializing storage engine")
	engine, err := storage.NewStorageEngine()
	if err != nil {
		return fmt.Errorf("failed to create storage engine: %w", err)
	}
	app.engine = engine

	return nil
}
func (app *Application) waitForShutdown(serverErrCh <-chan error) error {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		log.Error().Err(err).Msg("Server error occurred")
		return app.shutdown()
	case sig := <-shutdownCh:
		log.Info().
			Str("signal", sig.String()).
			Msg("Shutdown signal received")
		return app.shutdown()
	}
}
func (app *Application) shutdown() error {
	log.Info().Msg("Initiating graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if app.server != nil {
		if err := app.server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Error shutting down server")
		}
	}
	if app.engine != nil {
		log.Info().Msg("Closing storage engine")
		if err := app.engine.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing storage engine")
			return err
		}
	}

	log.Info().Msg("Graceful shutdown completed")
	return nil
}
