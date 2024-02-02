package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server Conn `yaml:"server"`
		Client Conn `yaml:"client"`
		Hash   Hash `yaml:"hash"`
	}

	Conn struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}

	Hash struct {
		FirstZerosCount int `yaml:"first_zeros_count"`
	}
)

var instance *Config

func InitConfig(pathConfig string) *Config {
	instance = &Config{}
	if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
		slog.Error(fmt.Errorf("Fail read config: %v", err).Error())
		return nil
	}
	return instance
}

func GetConfig() *Config {
	return instance
}
