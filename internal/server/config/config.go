package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Stores struct {
		UserStoreFile string `env:"USER_STORE_FILE"`
	}
	Settings struct {
		RunAddress string `env:"RUN_ADDRESS"`
	}
}

func GetConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Settings.RunAddress, "a", ":3333", "address for web-server listen")
	flag.StringVar(&cfg.Stores.UserStoreFile, "user-store", "users.json", "user store file path")
	flag.Parse()

	env.Parse(&cfg.Settings)

	return cfg
}
