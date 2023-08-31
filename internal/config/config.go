package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	DataBaseUrl string `yaml:"database_url" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost"`
	Port    string `yaml:"port" env-default:"8080"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
