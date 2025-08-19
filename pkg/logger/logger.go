package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level     zerolog.Level
	Console   bool
	LogToFile bool
	FilePath  string
}

func getFileNameFromPath(path2 string) string {
	parts := strings.Split(path2, "/")
	return parts[len(parts)-1]
}
func Setup(cfg Config) error {
	var output *os.File
	var err error

	if cfg.LogToFile {
		output, err = os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			if os.IsNotExist(err) {
				// Create the log directory if it doesn't exist
				if err := os.MkdirAll(cfg.FilePath[:len(cfg.FilePath)-len(getFileNameFromPath(cfg.FilePath))], 0755); err != nil {
					return err
				}
				output, err = os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					return err
				}

			}
			return err
		}
		log.Logger = log.Output(output)
	} else if cfg.Console {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		log.Logger = log.Output(os.Stdout)
	}

	zerolog.SetGlobalLevel(cfg.Level)
	return nil
}
