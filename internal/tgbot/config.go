package tgbot

import (
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/store"
)

type Config struct {
	Token    string `toml:"token"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		LogLevel: "debug",
		Store:    store.NewConfig(),
	}
}
