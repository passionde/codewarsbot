package helper

import (
	"flag"
	"github.com/BurntSushi/toml"
	apiserver2 "github.com/Yarik-xxx/CodeWarsRestApi/internal/app/telegram"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/telegram.toml", "path to config file")
}

func InitConfig() *apiserver2.Config {
	flag.Parse()

	config := apiserver2.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
