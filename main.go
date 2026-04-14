package main

import (
	"laoinBot/config"
	"laoinBot/config/db"
	_ "laoinBot/plugin/cangmiao_func"
	_ "laoinBot/plugin/galgame_recommand"
	_ "laoinBot/plugin/help"

	// _ "laoinBot/plugin/jm"
	_ "laoinBot/plugin/sendLike"
	_ "laoinBot/plugin/send_lizi_image"
	_ "laoinBot/plugin/view_china"
	_ "laoinBot/plugin/xiaoxiao_func"
	"log"

	zero "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/driver"
)

func init() {
	err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	err = db.InitDb()
	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
}
func main() {
	zero.RunAndBlock(&zero.Config{
		NickName:      config.BotConfig.MainConfig.NickName,
		CommandPrefix: "/",
		SuperUsers:    config.BotConfig.MainConfig.SuperUser,
		Driver: []zero.Driver{
			driver.NewWebSocketClient(config.BotConfig.MainConfig.NapcatUrl, config.BotConfig.MainConfig.NapcatToken),
		},
	}, nil)
}
