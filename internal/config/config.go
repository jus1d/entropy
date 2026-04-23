package config

import (
	"apigo/pkg/log/sl"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Env string

const (
	EnvLocal       Env = "local"
	EnvDevelopment Env = "dev"
	EnvProduction  Env = "prod"
)

type Config struct {
	Env        Env        `yaml:"env" env-required:"true"`
	Server     Server     `yaml:"http" env-required:"true"`
	ClickHouse ClickHouse `yaml:"clickhouse" env-required:"true"`
}

type Server struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type ClickHouse struct {
	Address  string `yaml:"address" env-required:"true"`
	Database string `yaml:"database" env-default:"default"`
	Username string `yaml:"username" env-default:"default"`
	Password string `yaml:"password" env-default:"default"`
}

// MustLoad loads config to a new Config instance and return it
func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		slog.Error("missed CONFIG_PATH parameter")
		os.Exit(1)
	}

	var err error
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("config file does not exist", slog.String("path", configPath))
		os.Exit(1)
	}

	var config Config

	if err = cleanenv.ReadConfig(configPath, &config); err != nil {
		slog.Error("cannot read config", sl.Err(err))
		os.Exit(1)
	}

	return &config
}

func Empty() *Config {
	return &Config{}
}
