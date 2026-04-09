package main

import (
	"laoinBot/config"
	_ "laoinBot/plugin/cangmiao_func"
	_ "laoinBot/plugin/help"
	_ "laoinBot/plugin/sendLike"
	_ "laoinBot/plugin/send_lizi_image"
	_ "laoinBot/plugin/xiaoxiao_func"
	"log"

	zero "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/driver"
)

func main() {
	err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	zero.RunAndBlock(&zero.Config{
		NickName:      config.BotConfig.MainConfig.NickName,
		CommandPrefix: "/",
		SuperUsers:    config.BotConfig.MainConfig.SuperUser,
		Driver: []zero.Driver{
			driver.NewWebSocketClient(config.BotConfig.MainConfig.NapcatUrl, config.BotConfig.MainConfig.NapcatToken),
		},
	}, nil)
}
