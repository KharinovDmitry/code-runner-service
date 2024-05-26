package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	localEnv = "local"
	devEnv   = "dev"
)

type Config struct {
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

func Load(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.New("config file open error: " + configPath)
	}
	cfg := Config{}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, errors.New("load config error: " + err.Error())
	}

	return &cfg, nil
}
