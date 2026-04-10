package galgamerecommand

import (
	"laoinBot/config/db"
	"laoinBot/config/db/table"
	"laoinBot/plugin/help"
	"math/rand"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

func init() {
	help.HelpInstance.SetHelper("黄油推荐", "推荐黄油", "黄油推荐")
	nova.OnFullMatch("黄油推荐").Handle(func(ctx *nova.Ctx) {
		var galgame table.Galgame
		var count int64
		db.Db.Model(&table.Galgame{}).Count(&count)
		if count == 0 {
			ctx.SendChain(message.Text("没有黄油推荐"))
			return
		}
		id := rand.Intn(int(count))
		db.Db.Where("id = ?", id).Take(&galgame)
		ctx.SendChain(message.Text(galgame.Title, "\n\n", galgame.Url, "\n\n", galgame.Introduction))
	})
}
