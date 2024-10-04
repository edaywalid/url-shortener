package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerID  string `mapstructure:"SERVER_ID"`
	Port      string `mapstructure:"PORT"`
	ZkAddr    string `mapstructure:"ZK_ADDR"`
	RedisAddr string `mapstructure:"REDIS_ADDR"`
	BaseURL   string `mapstructure:"BASE_URL"`
}

func LoadConfig(path string) (*Config, error) {

	var cfg Config
	godotenv.Load()
	cfg.ServerID = os.Getenv("SERVER_ID")
	cfg.Port = os.Getenv("PORT")
	cfg.ZkAddr = os.Getenv("ZK_ADDR")
	cfg.RedisAddr = os.Getenv("REDIS_ADDR")
	cfg.BaseURL = os.Getenv("BASE_URL")

	log.Info().Msgf("Loaded the configuration: %+v", cfg)

	return &cfg, nil
}
