package main

import (
	_ "laoinBot/plugin/chat"
	_ "laoinBot/plugin/sendLike"
	_ "laoinBot/plugin/sexyImage"

	zero "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/driver"
)

func main() {
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"laoin"},
		CommandPrefix: "/",
		SuperUsers:    []int64{2908451607},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("ws://123.56.140.27:6101", "laoinNB666"),
		},
	}, nil)
}
