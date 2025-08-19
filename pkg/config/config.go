package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type loggingConfig struct {
	Path      string `yaml:"path"`
	ToFile    bool   `yaml:"toFile"`
	Level     string `yaml:"level"`
	ToConsole bool   `yaml:"toConsole"`
}
type Config struct {
	Port    int           `yaml:"port"`
	Logging loggingConfig `yaml:"logging"`
}

func DefaultConfig() Config {
	return Config{
		Logging: loggingConfig{
			Path:      "./logs/root.log",
			ToFile:    true,
			Level:     "trace",
			ToConsole: true,
		},
		Port: 5381,
	}
}

func LoadConfig(path string) (Config, error) {

	cfg := DefaultConfig()
	file, err := os.Open(path)
	if err != nil {
		return cfg, nil
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
