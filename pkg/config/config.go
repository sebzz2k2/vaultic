package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Global Config

type Config struct {
	LogPath string `yaml:"logpath"`
	Port    int    `yaml:"port"`
}

func DefaultConfig() Config {
	return Config{
		LogPath: "./logs/root.log",
		Port:    5431,
	}
}

func loadConfig(path string) (Config, error) {

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

func InitConfig() error {
	cfg, err := loadConfig("vaultic_config.yaml")
	if err != nil {
		return err
	}
	Global = cfg
	return nil
}
