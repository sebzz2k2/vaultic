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

type serverConfig struct {
	Address        string `yaml:"address"`
	Port           int    `yaml:"port"`
	MaxConnections int    `yaml:"maxConnections"`
	MaxMessageSize int    `yaml:"maxMessageSizeBytes"` // in bytes
}
type Config struct {
	Port    int           `yaml:"port"`
	Logging loggingConfig `yaml:"logging"`
	Server  serverConfig  `yaml:"server"`
}

func DefaultConfig() Config {
	return Config{
		Logging: loggingConfig{
			Path:      "./logs/root.log",
			ToFile:    true,
			Level:     "trace",
			ToConsole: true,
		},
		Server: serverConfig{
			Address:        "localhost",
			Port:           5381,
			MaxConnections: 100,
			MaxMessageSize: 1024 * 1024, // 1 MB
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
