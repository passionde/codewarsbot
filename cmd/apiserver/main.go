package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/tgbot"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := tgbot.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := tgbot.New(config)

	defer func() {
		if err := recover(); err != nil {
			main()
		}
	}()
	
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
